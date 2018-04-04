package controllers

import (
	"PGConfig/app/proto"
	"strconv"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"
)

type ConfigController struct {}

//分发控制器
//@param struct rq 请求参数
func (n *ConfigController) Dispatch(rq config.ConfigRequest) config.ConfigResponse {
	res := config.ConfigResponse{}
	if rq.Appid < 0 || rq.Region < 0 {
		return res
	}
	s1 := int64(rq.S1)
	s1Bin := strconv.FormatInt(s1, 2)
	s2 := int64(rq.S2)
	s2Bin := strconv.FormatInt(s2, 2)
	arr := map[string]string{
		"S1" : s1Bin,
		"S2" : s2Bin,
	}
	methods := make([]string, 0)
	for key, bin := range arr {
		length := len(bin)
		for i:=0; i<length; i++ {
			if bin[i] == 49 {
				code := math.Pow(2, float64(length-i-1))
				s := strconv.FormatFloat(code, 'G', -1, 64)
				method := key+"_"+s
				v := reflect.ValueOf(n)
				mv := v.MethodByName(method)
				if mv.IsValid() {
					methods = append(methods, method)
				}
			}
		}
	}
	res = n.goAsync(methods, rq)
	return res
}

//放在不同协程处理
func (n *ConfigController) goAsync(methods []string, rq config.ConfigRequest) config.ConfigResponse {
	res := config.ConfigResponse{}
	count := len(methods)
	if count > 0 {
		var wg sync.WaitGroup
		wg.Add(count)
		for _, act := range methods {
			timer := time.NewTimer(2*time.Second)
			ch := make(chan bool)
			go func(tm *time.Timer, cha chan bool) {
				for {
					select {
					case <-tm.C:
						cha<-true
						return
					}
				}
			}(timer, ch)
			go func(method string, tm *time.Timer, cha chan bool) {
				obj := reflect.ValueOf(n)
				mv := obj.MethodByName(method)
				if mv.IsValid() {
					args := []reflect.Value{reflect.ValueOf(rq)}
					result := mv.Call(args)
					slices := strings.Split(method, "_")
					if len(slices) == 2 {
						rel := reflect.ValueOf(&res).Elem()
						name := rel.FieldByName(slices[0])
						cmd := name.FieldByName("Cmd"+slices[1])
						if cmd.CanSet() {
							cmd.Set(result[0])
						}
					}
					if tm.Stop() {
						cha<-true
					}
				}
			}(act, timer, ch)
			go func(cha chan bool) {
				<-cha
				wg.Done()
			}(ch)
		}
		wg.Wait()
	}
	return res
}


//地区配置测试
func (n *ConfigController) S1_1(rq config.ConfigRequest) config.RegionResp {
	time.Sleep(0*time.Second)
	result := config.RegionResp{
		Status:1, //1成功 0失败
		Data:config.RegionItem{
			Id:2000,
			Content:`这几天心里颇不宁静。今晚在院子里坐着乘凉，忽然想起日日走过的荷塘，在这满月的光里，总该另有一番样子吧。月亮渐渐地升高了，墙外马路上孩子们的欢笑，已经听不见了；妻在屋里拍着闰儿⑴，迷迷糊糊地哼着眠歌。我悄悄地披了大衫，带上门出去。
沿着荷塘，是一条曲折的小煤屑路。这是一条幽僻的路；白天也少人走，夜晚更加寂寞。荷塘四面，长着许多树，蓊蓊郁郁⑵的。路的一旁，是些杨柳，和一些不知道名字的树。没有月光的晚上，这路上阴森森的，有些怕人。今晚却很好，虽然月光也还是淡淡的。
路上只我一个人，背着手踱⑶着。这一片天地好像是我的；我也像超出了平常的自己，到了另一个世界里。我爱热闹，也爱冷静；爱群居，也爱独处。像今晚上，一个人在这苍茫的月下，什么都可以想，什么都可以不想，便觉是个自由的人。白天里一定要做的事，一定要说的话，现 在都可不理。这是独处的妙处，我且受用这无边的荷香月色好了。
曲曲折折的荷塘上面，弥望⑷的是田田⑸的叶子。叶子出水很高，像亭亭的舞女的裙。层层的叶子中间，零星地点缀着些白花，有袅娜⑹地开着的，有羞涩地打着朵儿的；正如一粒粒的明珠，又如碧天里的星星，又如刚出浴的美人。微风过处，送来缕缕清香，仿佛远处高楼上渺茫的歌声似的。这时候叶子与花也有一丝的颤动，像闪电般，霎时传过荷塘的那边去了。叶子本是肩并肩密密地挨着，这便宛然有了一道凝碧的波痕。叶子底下是脉脉⑺的流水，遮住了，不能见一些颜色；而叶子却更见风致⑻了。
月光如流水一般，静静地泻在这一片叶子和花上。薄薄的青雾浮起在荷塘里。叶子和花仿佛在牛乳中洗过一样；又像笼着轻纱的梦。虽然是满月，天上却有一层淡淡的云，所以不能朗照；但我以为这恰是到了好处——酣眠固不可少，小睡也别有风味的。月光是隔了树照过来的，高处丛生的灌木，落下参差的斑驳的黑影，峭楞楞如鬼一般；弯弯的杨柳的稀疏的倩影，却又像是画在荷叶上。塘中的月色并不均匀；但光与影有着和谐的旋律，如梵婀玲⑼上奏着的名曲。
荷塘的四面，远远近近，高高低低都是树，而杨柳最多。这些树将一片荷塘重重围住；只在小路一旁，漏着几段空隙，像是特为月光留下的。树色一例是阴阴的，乍看像一团烟雾；但杨柳的丰姿⑽，便在烟雾里也辨得出。树梢上隐隐约约的是一带远山，只有些大意罢了。树缝里也漏着一两点路灯光，没精打采的，是渴睡人的眼。这时候最热闹的，要数树上的蝉声与水里的蛙声；但热闹是它们的，我什么也没有。
忽然想起采莲的事情来了。采莲是江南的旧俗，似乎很早就有，而六朝时为盛；从诗歌里可以约略知道。采莲的是少年的女子，她们是荡着小船，唱着艳歌去的。采莲人不用说很多，还有看采莲的人。那是一个热闹的季节，也是一个风流的季节。梁元帝《采莲赋》里说得好：
于是妖童媛女⑾，荡舟心许；鷁首⑿徐回，兼传羽杯⒀；棹⒁将移而藻挂，船欲动而萍开。尔其纤腰束素⒂，迁延顾步⒃；夏始春余，叶嫩花初，恐沾裳而浅笑，畏倾船而敛裾⒄。
可见当时嬉游的光景了。这真是有趣的事，可惜我们现 在早已无福消受了。
于是又记起，《西洲曲》里的句子：
采莲南塘秋，莲花过人头；低头弄莲子，莲子清如水。
今晚若有采莲人，这儿的莲花也算得“过人头”了；只不见一些流水的影子，是不行的。这令我到底惦着江南了。——这样想着，猛一抬头，不觉已是自己的门前；轻轻地推门进去，什么声息也没有，妻已睡熟好久了。
一九二七年七月，北京清华园。[1]
`,
		},
	}
	return result
}
//bpid配置测试
func (n *ConfigController) S1_2(rq config.ConfigRequest) config.BpidResp {
	time.Sleep(0*time.Second)
	result := config.BpidResp{
		Status:1,
		Data:config.BpidItem{
			Bid:3000,
			Content:`我与父亲不相见已二年余了，我最不能忘记的是他的背影。
那年冬天，祖母死了，父亲的差使1也交卸了，正是祸不单行的日子。我从北京到徐州，打算跟着父亲奔丧2回家。到徐州见着父亲，看见满院狼藉3的东西，又想起祖母，不禁簌簌地流下眼泪。父亲说：“事已如此，不必难过，好在天无绝人之路！”
回家变卖典质4，父亲还了亏空；又借钱办了丧事。这些日子，家中光景5很是惨澹，一半为了丧事，一半为了父亲赋闲6。丧事完毕，父亲要到南京谋事，我也要回北京念书，我们便同行。
到南京时，有朋友约去游逛，勾留7了一日；第二日上午便须渡江到浦口，下午上车北去。父亲因为事忙，本已说定不送我，叫旅馆里一个熟识的茶房8陪我同去。他再三嘱咐茶房，甚是仔细。但他终于不放心，怕茶房不妥帖9；颇踌躇10了一会。其实我那年已二十岁，北京已来往过两三次，是没有什么要紧的了。他踌躇了一会，终于决定还是自己送我去。我再三劝他不必去；他只说：“不要紧，他们去不好！”
我们过了江，进了车站。我买票，他忙着照看行李。行李太多，得向脚夫11行些小费才可过去。他便又忙着和他们讲价钱。我那时真是聪明过分，总觉他说话不大漂亮，非自己插嘴不可，但他终于讲定了价钱；就送我上车。他给我拣定了靠车门的一张椅子；我将他给我做的紫毛大衣铺好座位。他嘱我路上小心，夜里要警醒些，不要受凉。又嘱托茶房好好照应我。我心里暗笑他的迂；他们只认得钱，托他们只是白托！而且我这样大年纪的人，难道还不能料理自己么？我现在想想，我那时真是太聪明了。
我说道：“爸爸，你走吧。”他望车外看了看，说：“我买几个橘子去。你就在此地，不要走动。”我看那边月台的栅栏外有几个卖东西的等着顾客。走到那边月台，须穿过铁道，须跳下去又爬上去。父亲是一个胖子，走过去自然要费事些。我本来要去的，他不肯，只好让他去。我看见他戴着黑布小帽，穿着黑布大马褂12，深青布棉袍，蹒跚13地走到铁道边，慢慢探身下去，尚不大难。可是他穿过铁道，要爬上那边月台，就不容易了。他用两手攀着上面，两脚再向上缩；他肥胖的身子向左微倾，显出努力的样子。这时我看见他的背影，我的泪很快地流下来了。我赶紧拭干了泪。怕他看见，也怕别人看见。我再向外看时，他已抱了朱红的橘子往回走了。过铁道时，他先将橘子散放在地上，自己慢慢爬下，再抱起橘子走。到这边时，我赶紧去搀他。他和我走到车上，将橘子一股脑儿放在我的皮大衣上。于是扑扑衣上的泥土，心里很轻松似的。过一会儿说：“我走了，到那边来信！”我望着他走出去。他走了几步，回过头看见我，说：“进去吧，里边没人。”等他的背影混入来来往往的人里，再找不着了，我便进来坐下，我的眼泪又来了。
近几年来，父亲和我都是东奔西走，家中光景是一日不如一日。他少年出外谋生，独力支持，做了许多大事。哪知老境却如此颓唐！他触目伤怀，自然情不能自已。情郁于中，自然要发之于外；家庭琐屑便往往触他之怒。他待我渐渐不同往日。但最近两年不见，他终于忘却我的不好，只是惦记着我，惦记着他的儿子。我北来后，他写了一信给我，信中说道：“我身体平安，惟膀子疼痛厉害，举箸14提笔，诸多不便，大约大去之期15不远矣。”我读到此处，在晶莹的泪光中，又看见那肥胖的、青布棉袍黑布马褂的背影。唉！我不知何时再能与他相见！`,
		},
	}
	return result
}