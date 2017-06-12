package hlog

import "runtime"

type Hlog struct {
	isOpen bool
}

func (this *Hlog) Open() {
	this.isOpen = true
}

func (this *Hlog) Close() {
	this.isOpen = false
}

func (this *Hlog) getCaller() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	return file, line
}
