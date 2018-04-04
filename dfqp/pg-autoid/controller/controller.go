/**
*基本原理：一次从mysql中获取step（1000）个数值，放在内存中，每次请求获取id时优先从内存中取，取不到的情况下再从mysql中拉取step个数值【其实返回min和max即可】
*双buffer机制：为了避免在step处偶尔的尖峰，需要检查：当一个buffer用到了20%时就异步的执行一次mysql请求，将另一份待用step值存储起来，备用！
*
*业务逻辑处理：接收类型string，从该类型对应的id集合中查找id值，如果能查到就返回，同时异步的填充另外一个buffer
 */
//风险提示：频繁的重启会导致一直从DB申请数据段得不到充分的使用，造成浪费！
package controller

import (
	//"PGAutoIdManager/service"
	"dfqp/pg-autoid/service"
	"errors"
	"fmt"
	"putil/log"
	"sync"
)

//id buffer结构体
type idsbuffer struct {
	min     int64  //可选最小值
	max     int64  //可选最大值
	current int64  //当前选择值(说明：本次请求可以使用这个值)
	nextmin int64  //下一组备选最小值
	nextmax int64  //下一组备选最大值
	btag    string //类型标记
	//singlemux sync.Mutex //更新buffer2的锁
	mux sync.Mutex //mutex(用来控制min、max、current、nextmin、nextmax的更新)
	//dbmux       sync.Mutex //控制DB操作，防止并发访问问题
	isbuffering bool //buffer2是否正在缓存中，如果是的话，本次就不触发
}

//
type ctrl map[string]*idsbuffer //the pointer to idsbuffer

var (
	objmap ctrl       //controller object
	mu     sync.Mutex //不同类型的使用中的锁
	dbmux  sync.Mutex //控制DB操作，防止并发访问问题
)

//init func
func init() {
	var err error
	objmap, err = NewController()
	if err != nil {
		plog.Fatal("new err!")
		return
	}
}

//new obj
func NewController() (om ctrl, err error) {
	om = make(map[string]*idsbuffer)
	if om == nil {
		err = errors.New("New map err!")
	}
	om = ctrl(om) //类型转换

	return
}

//get the auto id
func GetId(str string) (int64, error) {
	plog.Debug("GetId GetId GetId GetId GetId GetId:")
	return objmap.getid(str)
}

//obj get id
func (c ctrl) getid(str string) (int64, error) {
	//TODO:srt 的合法性检测,后续添加
	mu.Lock()
	if c[str] == nil {
		//实例化
		c[str] = new(idsbuffer)
		c[str].btag = str
	}
	mu.Unlock()

	re, err := c[str].outputId() //从本地cache中输出数据
	return re, err
}

//do output
func (b *idsbuffer) outputId() (int64, error) {
	//最快的速度出数据
	b.mux.Lock()
	rt, err := b.outputOne()
	b.mux.Unlock()
	for err != nil {
		b.mux.Lock()
		bmin, bmax, bcurrent, bnextmin, bnextmax := b.readBuffer()
		status, err := b.isBufferUseful(bmin, bmax, bcurrent, bnextmin, bnextmax)

		if err != nil {
			if status == -5 {
				b.swapBuffer() //填充buffer1
				b.mux.Unlock()
			} else if status == -6 {
				b.addBufferFromRemote() //填充buffer1
				b.mux.Unlock()
			} else {
				//plog.Fatal(b)
				plog.Fatal("outputId err, status is:" + fmt.Sprintf("%d", status))
				err = errors.New("outputId err, status is:" + fmt.Sprintf("%d", status))
				b.mux.Unlock()
				return rt, err
			}
		} else {
			b.mux.Unlock()
		}

		b.mux.Lock()
		rt, err = b.outputOne()
		b.mux.Unlock()
		if err == nil {
			return rt, err
		}
	}
	return rt, err
}

//放出一个id【写】
func (b *idsbuffer) outputOne() (id int64, err error) {
	id = b.current
	if id > b.max {
		return 0, errors.New("bigger than the max")
	}
	if id == 0 {
		return 0, errors.New("the number is zero")
	}
	if id == b.min+(b.max-b.min+1)/10 && b.nextmin == 0 && b.nextmax == 0 {

		if b.isbuffering == false {
			go b.updateBuffer2FromRemote() //远程更新一份数据放入buffer2（异步操作）
		}
	}
	b.current = b.current + 1

	return
}

//把buffer2中的数据更新到buffer1【写】
func (b *idsbuffer) swapBuffer() error {
	b.min = b.nextmin
	b.max = b.nextmax
	b.current = b.min
	b.nextmin = 0
	b.nextmax = 0
	return nil
}

//远程获取数据，更新buffer1【写】
func (b *idsbuffer) addBufferFromRemote() error {
	rmin, rmax, rerr := b.getFromRemote() //远程拿一份新的ids
	if rmin > rmax || rerr != nil {       //远程拿回的数据不正确
		err := errors.New("the remote data is illegal")
		return err
	}
	b.min = rmin
	b.max = rmax
	b.current = rmin

	return nil
}

//内部不加锁，其他地方调用时锁控制【读】
func (b *idsbuffer) readBuffer() (bmin, bmax, bcurrent, bnextmin, bnextmax int64) {
	bmin = b.min
	bmax = b.max
	bcurrent = b.current
	bnextmin = b.nextmin
	bnextmax = b.nextmax
	return
}

//check whether the buffer is legal  default status is 0
func (b *idsbuffer) isBufferUseful(bmin, bmax, bcurrent, bnextmin, bnextmax int64) (status int, err error) {
	//未初始化
	if b == nil {
		status = -1
		err = errors.New("the buffer is nil")
		return
	}

	//正常情况
	if bcurrent >= bmin && bcurrent <= bmax && bmin > 0 && bmin <= bmax {
		return
	}
	//数据不合法一，报错！
	if bmin < 0 || bmax < bmin || bcurrent < bmin {
		status = -2
		err = errors.New("current buffer is illegal!")
		return
	}
	//数据不合法二，报错
	if bnextmin < 0 || bnextmax < bnextmin {
		status = -3
		err = errors.New("current buffer is illegal s1!")
		return
	}
	//数据不合法三，报错
	if (bnextmin >= bmin && bnextmin <= bmax && bnextmin != 0) || (bnextmax >= bmin && bnextmax <= bmax && bnextmax != 0) {
		status = -4
		err = errors.New("next buffer is illegal s2!")
		return
	}
	//数据初始化了，需要填充第一个buffer
	if bcurrent == 0 && bmin == 0 && bmax == 0 {
		status = -6 //TODO:远程获取数据填充到第一个buffer
		err = errors.New("the current has to the max, do add")
		return
	}
	//当前值已经用到了超过最大值了，需要远程填充第一个buffer或者把第二个buffer换进来！
	if bcurrent > bmax {
		if bnextmin > 0 {
			//第二个buffer换进来
			status = -5
			err = errors.New("the current has to the max , do swap")
			return
		} else {
			//远程获取数据填充到第一个buffer
			status = -6
			err = errors.New("the current has to the max, do add")
			return
		}
	}
	return
}

//远程更新buffer2的数据(在获取远程数据的时候不加锁，取回来后做更新时再加锁 风险：在极快的获取情况下可能会多次获取，导致部分ids浪费)
//注意：该方法在一个buffer中只会触发一次【写】
func (b *idsbuffer) updateBuffer2FromRemote() error {

	b.isbuffering = true
	rmin, rmax, rerr := b.getFromRemote()
	if rmin > rmax || rerr != nil { //远程拿回的数据不合法
		plog.Fatal("远程数据获取异常")
		err := errors.New("the remote data is illegal")
		return err
	}

	//更新数据(异步更新nextmin和nextmax)
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.nextmin == 0 && b.nextmax == 0 {
		b.nextmin = rmin
		b.nextmax = rmax
	} else {
		plog.Fatal(fmt.Sprintf("get data form DB but not use! rmin is %d  rmax is %d ", rmin, rmax))
	}
	b.isbuffering = false

	return nil
}

//从公共存储介质中获取可用的数据
func (b *idsbuffer) getFromRemote() (rmin int64, rmax int64, rerr error) {
	dbmux.Lock()
	defer dbmux.Unlock()
	var step int32
	rmax, step, rerr = service.AutoidService.ModifyAndGet(b.btag)
	if rerr != nil || rmax-int64(step) < 0 {
		plog.Fatal("MOdifyAndGet data error!")
	}
	rmin = rmax - int64(step) + 1
	return
}
