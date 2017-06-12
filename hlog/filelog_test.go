package hlog_test

import (
	"github.com/hundredlee/wechat_pusher/hlog"
	"testing"
)

func TestFileLog_LogError(t *testing.T) {
	filelog := hlog.LogInstance()
	filelog.LogError("xxxxxxxxx")
}

func TestFileLog_LogInfo(t *testing.T) {
	filelog := hlog.LogInstance()
	filelog.LogError("yyyyyyyyy")
}
