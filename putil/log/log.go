package plog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
	//"errors"
	"framework/rpcclient/szmq"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/natefinch/lumberjack"
)

const (
	INFOLEVEL         = 1 << iota //info日志等级
	DEBUGLEVEL                    //debug日志等级
	WARNLEVEL                     //warn日志等级
	FATALLEVEL                    //fata日志等级
	CORELEVEL                     //core日志等级
	DEFAULTMAXSIZE    = 1         //默认的文件最大容量（单位：M）
	DEFAULTMAXBACKUPS = 10        //默认的保留的历史文件的个数
	DEBAULTMAXAGE     = 0         //默认的保留天数，0表示一直保留
	DEFAULTCOMPRESS   = false     //默认是否压缩
)

type zmqLogData struct {
	ActTime int64  `json:"act_time"`
	Level   string `json:"level"`
	Msg     string `json:"msg"`
}

var (
	l, infohandle, debughandle, warnhandle, fatalhandle, corehandle *log.Logger
	defaultFilePath                                                 = "" //文件路径
	loglevel                                                        = 0  //日志等级初始化
	logRegister                                                     map[string]*log.Logger
	defaultFileName                                                 = "" //日志文件名
	prefix                                                          = "default"
	zmqloglevel                                                     = 0

	zmqCli *szmq.SzmqPushClient
)

func init() {
	//	_, err := os.Create("./" + time.Now().Format("20060102") + ".txt")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	logRegister = make(map[string]*log.Logger)
	iniconf, _ := config.NewConfig("ini", "./conf/app.conf") //config 模块来解析你各种格式的文件
	defaultFilePath = iniconf.String("app.logFilePath")
	args := os.Args
	if args != nil && len(args) > 1 {
		prefix = args[1]
	}
	if defaultFilePath == "" {
		defaultFilePath = "./log"
	} else {
		index := strings.HasSuffix(defaultFilePath, ".log")
		if index {
			//存在，不是路径,抛出错误
			panic("日志文件配置logFilePath不是正确的路径 !")
		}
	}

	//设置默认的writer
	//		writer := &lumberjack.Logger{
	//			Filename:   defaultFileName,
	//			MaxSize:    DEFAULTMAXSIZE, // megabytes
	//			MaxBackups: DEFAULTMAXBACKUPS,
	//			MaxAge:     DEBAULTMAXAGE,   //days
	//			Compress:   DEFAULTCOMPRESS, // disabled by default
	//		}

	//	//writer 设置为标准输出
	//writer := os.Stdout
	//l = log.New(writer, "", log.Ldate|log.Lmicroseconds|log.Lshortfile) //Ldate显示形如 2009/01/23 的日期 //Lmicroseconds显示时分秒日期//Lshortfile显示文件名和行号: d.go:23
	//读取日志等级
	loglevel, _ = iniconf.Int("app.logLevel")
	zmqloglevel, _ = iniconf.Int("app.zmqlogLevel")
	//初始化日志文件句柄
	createWrite()
	if zmqloglevel != 0 {
		zmqCli, _ = szmq.NewSzmqPushClient()
	}

}

func Info(v ...interface{}) {
	if loglevel != 0 && INFOLEVEL&loglevel == INFOLEVEL {
		infohandle.SetPrefix("[INFO]")
		infohandle.Output(2, fmt.Sprint(v...))
	}
	if zmqloglevel != 0 && INFOLEVEL&zmqloglevel == INFOLEVEL {
		data, err := assembleZmqLogData(fmt.Sprint(v...), "info")
		if err != nil {
			return
		}
		//普通日志上报
		zmqCli.WriteNormalLog("program_log", string(data))
	}
}

func Debug(v ...interface{}) {
	if loglevel != 0 && DEBUGLEVEL&loglevel == DEBUGLEVEL {
		debughandle.SetPrefix("[DEBUG]")
		debughandle.Output(2, fmt.Sprint(v...))
	}
	if zmqloglevel != 0 && DEBUGLEVEL&zmqloglevel == DEBUGLEVEL {
		data, err := assembleZmqLogData(fmt.Sprint(v...), "debug")
		if err != nil {
			return
		}
		//调试日志上报
		zmqCli.WriteDebugLog("program_log", string(data))
	}
}

func Warn(v ...interface{}) {
	if loglevel != 0 && WARNLEVEL&loglevel == WARNLEVEL {
		warnhandle.SetPrefix("[WARN]")
		warnhandle.Output(2, fmt.Sprint(v...))
	}
	if zmqloglevel != 0 && WARNLEVEL&zmqloglevel == WARNLEVEL {
		data, err := assembleZmqLogData(fmt.Sprint(v...), "warn")
		if err != nil {
			return
		}
		//调试日志上报
		zmqCli.WriteDebugLog("program_log", string(data))
	}
}
func Fatal(v ...interface{}) {
	if loglevel != 0 && FATALLEVEL&loglevel == FATALLEVEL {
		fatalhandle.SetPrefix("[FATAL]")
		fatalhandle.Output(2, fmt.Sprint(v...))
	}
	if zmqloglevel != 0 && FATALLEVEL&zmqloglevel == FATALLEVEL {
		data, err := assembleZmqLogData(fmt.Sprint(v...), "fatal")
		if err != nil {
			return
		}
		//实时日志上报
		zmqCli.WriteRealLog("program_log", string(data))
	}
}
func Core(v ...interface{}) {
	if loglevel != 0 && CORELEVEL&loglevel == CORELEVEL {
		corehandle.SetPrefix("[CORE]")
		corehandle.Output(2, fmt.Sprint(v...))
	}
	if zmqloglevel != 0 && CORELEVEL&zmqloglevel == CORELEVEL {
		data, err := assembleZmqLogData(fmt.Sprint(v...), "core")
		if err != nil {
			return
		}
		//实时日志上报
		zmqCli.WriteRealLog("program_log", string(data))
	}

}

//在main的开始引入，其他goroutine中也需要引入！
func CatchPanic() {
	//panic异常处理
	if r := recover(); r != nil {
		//fmt.Println("recover happened!", r)
		var s string = string(debug.Stack())
		Core(s)
		//服务退出(经过验证：使用os.Exit(1)时会主动向Net层发送Fin包)
		os.Exit(1)
		//TODO:告警！
	}
}

/**
 * 设置log输出使用轮转文件日志存储
 * @param {[type]} logFilePath string        日志目录
 * @param {[type]} maxSize     日志开始备份大小
 * @param {[type]} maxBackups  日志备份数量(就是多少个文件轮转)
 * @param {[type]} maxAge      int           存留日期
 */
func SetOutPutFileRotate(logFilePath string, maxSize, maxBackups, maxAge int) {
	l.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, //days
		Compress:   false,  // disabled by default
	})
}
func createWrite() {
	switch {
	case loglevel != 0 && INFOLEVEL&loglevel == INFOLEVEL:
		infohandle = LogRegister("info")
		fallthrough
	case loglevel != 0 && DEBUGLEVEL&loglevel == DEBUGLEVEL:
		debughandle = LogRegister("debug")
		fallthrough
	case loglevel != 0 && WARNLEVEL&loglevel == WARNLEVEL:
		warnhandle = LogRegister("warn")
		fallthrough
	case loglevel != 0 && FATALLEVEL&loglevel == FATALLEVEL:
		fatalhandle = LogRegister("fatal")
		fallthrough
	case loglevel != 0 && CORELEVEL&loglevel == CORELEVEL:
		corehandle = LogRegister("core")
	}

}

func LogRegister(logLevelName string) *log.Logger {
	flag := strings.HasSuffix(defaultFilePath, "/")
	if flag {
		defaultFileName = defaultFilePath + logLevelName + "-" + prefix + ".log"
	} else {
		defaultFileName = defaultFilePath + "/" + logLevelName + "-" + prefix + ".log"
	}

	if exist, ok := logRegister[defaultFileName]; ok {
		return exist
	}

	//设置默认的writer
	writer := &lumberjack.Logger{
		Filename:   defaultFileName,
		MaxSize:    DEFAULTMAXSIZE, // megabytes
		MaxBackups: DEFAULTMAXBACKUPS,
		MaxAge:     DEBAULTMAXAGE,   //days
		Compress:   DEFAULTCOMPRESS, // disabled by default
	}

	//writer 设置为标准输出
	l = log.New(writer, "", log.Ldate|log.Lmicroseconds|log.Lshortfile) //Ldate显示形如 2009/01/23 的日期 //Lmicroseconds显示时分秒日期//Lshortfile显示文件名和行号: d.go:23
	logRegister[defaultFileName] = l
	return l
}

func assembleZmqLogData(data string, loglevel string) (realZmqData []byte, err error) {
	//初始化数据zmq数据结构
	var zmqData zmqLogData = zmqLogData{
		ActTime: 0,
		Level:   "",
		Msg:     "",
	}

	zmqData.ActTime = time.Now().UnixNano() / 1e6
	zmqData.Level = loglevel
	zmqData.Msg = data

	realZmqData, err = json.Marshal(zmqData)
	if err != nil {
		Fatal("marshal data err: ", err)
		return
	}
	return

}
