package hlog

import "runtime"

type hlog struct {
	isOpen bool
}

func (this *hlog) Open() {
	this.isOpen = true
}

func (this *hlog) Close() {
	this.isOpen = false
}

func (this *hlog) getCaller() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	return file, line
}
