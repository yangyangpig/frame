package timer

import (
	"errors"
	"putil/log"
	systime "time"
)

const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)

/**
面向函数的回调函数类型
*/
type TimeNotify func(param interface{})

/**
面向对象的回调接口
*/
type TimerCallBack interface {
	TimeNotify()
}

/**
面向对象的定时器对象，这里只有关闭接口可使用
*/
type TimerOO interface {
	Close()
}

type Timer interface {

	/**
	调用此接口，是以C风格的方式传入回调函数(类型是TimeNotify)的方式,以定时器id的方式传入，关闭定时器需要调用Close(timeid int)
	*/
	AddEvent(timeid int, timeout systime.Duration, tcb TimeNotify, param interface{}) (err error)
	/**
	调用此接口，必须使用NewTimerC来创建Timer，param将在定时器到时传入NewTimerC中的通道，但定时器期间不容许close原来
	创建时的Channel，否则发生奔溃，使用时如果要关闭NewTimerC中传递的通道参数，必须关闭所有的定时器才可以
	*/
	AddEventC(timeid int, timeout systime.Duration, param interface{})
	/**
	调用此接口会返回TimerOO定时器接口,面向对象的定时器接口，定时器的关闭得使用者自己调用TimerOO的Close
	*/
	StartTimer(interval systime.Duration, tcb TimerCallBack) TimerOO
	/**
	仅仅对调用AddEvent函数中传入的timeid进行有效，对StartTimer是无效的，如果你在既用AddEvent，又在调用StartTimer启动定时器，
	注意，此接口不会对StartTimer启动的定时器进行关闭。
	*/
	Close(timeid int) (err error)
	/**
	关闭所有定时器事件
	*/
	Clear()
}

//==============================================
func NewTimer() Timer {
	t := new(timer)
	t.init()
	return t
}

/**
以通道的方式进行计时器到时的通知。
*/
func NewTimerC(notify chan<- interface{}) Timer {
	t := new(timer)
	t.notify = notify
	t.init()
	return t
}

type timer struct {
	tm     map[int]*timeElem
	notify chan<- interface{}
}

func (t *timer) init() {
	t.tm = make(map[int]*timeElem)
}

func (t *timer) AddEvent(timeid int, timeout systime.Duration, tcb TimeNotify, param interface{}) (err error) {
	te, isfind := t.tm[timeid]
	if isfind {
		plog.Debug("find a timer by timeid", timeid)
		te.Close()
	}
	te = new(timeElem)
	te.init(timeout, tcb, param)
	t.tm[timeid] = te
	te.startTimer()
	return
}
func (t *timer) AddEventC(timeid int, timeout systime.Duration, param interface{}) {
	te, isfind := t.tm[timeid]
	if isfind {
		plog.Debug("find a timer by timeid", timeid)
		te.Close()
	}
	te = new(timeElem)
	te.initchannel(timeout, t.notify, param)
	t.tm[timeid] = te
	te.startTimer()
	return
}
func (t *timer) Close(timeid int) (err error) {
	te, isfind := t.tm[timeid]
	if !isfind {
		plog.Debug("donn't find the timer by timeid", timeid)
		err = errors.New("donn't find the timer by timeid")
		return
	}
	te.Close()
	return
}

func (t *timer) StartTimer(timeout systime.Duration, tcb TimerCallBack) TimerOO {
	oot := new(timeElem)
	oot.initoo(timeout, tcb)
	oot.startTimer()

	return oot
}

func (t *timer) Clear() {
	if t.tm == nil {
		return
	}
	for _, v := range t.tm {
		v.Close()
	}
}

//=================================================

type timeElem struct {
	tcb     TimeNotify
	cboo    TimerCallBack
	param   interface{}
	tm      *systime.Timer
	timeout systime.Duration
	notify  chan<- interface{}
	isclose bool
}

func (t *timeElem) init(timeout systime.Duration, tcb TimeNotify, param interface{}) {
	t.param = param
	t.tcb = tcb
	t.timeout = timeout
	t.isclose = false
}
func (t *timeElem) initoo(timeout systime.Duration, tcb TimerCallBack) {
	t.param = nil
	t.tcb = nil
	t.cboo = tcb
	t.timeout = timeout
}
func (t *timeElem) initchannel(timeout systime.Duration, notify chan<- interface{}, param interface{}) {
	t.param = param
	t.notify = notify
	t.tcb = nil
	t.cboo = nil
	t.timeout = timeout
}

//可以重入
func (t *timeElem) startTimer() {
	t.tm = systime.NewTimer(t.timeout)
	t.isclose = false
	go func() {
		select {
		case time := <-t.tm.C:
			plog.Info("timer is reach!", time.Format(RFC3339))

			if t.isclose {
				return
			} else {
				t.isclose = true
			}

			if t.tcb != nil { //回调函数模式
				t.tcb(t.param)
			} else if t.cboo != nil { //回调接口模式
				t.cboo.TimeNotify()
			} else if t.notify != nil { //通道传递模式
				//通道的有效性判定，玩家关闭通道怎么判断，？？？
				t.notify <- t.param
			}
		}
	}()
	return
}

func (t *timeElem) Close() {
	if t.isclose == false {
		t.isclose = true
		t.tm.Reset(0 * systime.Millisecond)

	}
}
