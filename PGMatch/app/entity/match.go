package entity

// 单个比赛比赛类型
const (
	_                  = iota
	MatchType_Day       // 日赛:1
	MatchType_Week      // 周循环赛:2
	MatchType_24Hours   // 24小时循环赛:3
	MatchType_Interval  // 间隔循环赛:4
)

// 定时赛比赛信息
type TimingMatchData struct {
	Matchid            int64
	Displaytime        int64
	Begintime          int64
	Looptype           string
	Firstbegintime     int64
	Loopinterval       string
	Loopintervalsecond string
	Loopendtime        string
	Addrobotflag       int
	Configid           string
	Gamename           string
	Matchentrycode     string
	Matchentryinfo     string
	Matchtags          []string
	Vmatchtags         [] string
	AdIcon             string
	ListSort           string
	ScreenDirection    int
	EnableWhiteList    map[string]interface{}
	PartitionFlag      string
	PaxUserCount       string
	MaxPartitions      string
	MaxUserCount       string
	Minversionand      int
	Maxversionand      int
	Minversionios      int
	Maxversionios      int
	Minversion         int
	Endtime            uint64
}
