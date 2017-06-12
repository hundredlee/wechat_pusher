package hlog

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strconv"
	"github.com/hundredlee/wechat_pusher/config"
	"fmt"
)


type fileLog struct {
	hlog
	logger *log.Logger
}


var (
	hFileLog *fileLog
	conf *config.Config = config.Instance()
)

func LogInstance () *fileLog {
	if hFileLog == nil{
		InitLogFile(true,"")
	}
	return hFileLog
}

func InitLogFile(isOpen bool,fp string) {
	if !isOpen {
		hFileLog = &fileLog{}
		hFileLog.logger = nil
		hFileLog.isOpen = isOpen
		return 
	}

	if fp == "" {
		wd:= os.Getenv("GOPATH")
		if wd =="" {
			file,_ := exec.LookPath(os.Args[0])
			path := filepath.Dir(file)
			wd = path
		}
		if wd == "" {
			 panic("GOPATH is not setted in env or can not get exe path.")
		}
		fp = fmt.Sprintf("%s/%s/",wd,conf.ConMap["Log.LOG_PATH"])
	}
	hFileLog = NewFileLog(isOpen,fp)
}

func NewFileLog(isOpen bool,logPath string) *fileLog {
	year,month,day := time.Now().Date()
	logName := "log." + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	err := os.MkdirAll(logPath,0755)
	if err != nil {
		 panic("logPath error :"+logPath+"\n")
	}
	f, err := os.OpenFile(logPath+"/"+logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("log file open error : " + logPath + "/" + logName + "\n")
	}
	fFileLog := &fileLog{}
	fFileLog.logger = log.New(f, "", log.LstdFlags)
	fFileLog.isOpen = isOpen
	return fFileLog
}


func (this *fileLog) log(lable string, str string) {
	if !this.isOpen {
		return
	}
	file, line := this.getCaller()
	this.logger.Printf("%s:%d: %s %s\n", file, line, lable, str)
}

func (this *fileLog) LogError(str string) {
	this.log("[ERROR]", str)
}

func (this *fileLog) LogInfo(str string) {
	this.log("[INFO]", str)
}

