package log

import (
	"fmt"

	"github.com/fatih/color"
)

type Level int

const (
	D Level = 1
	L Level = 2
	E Level = 3
)

var level Level = L

func SetLevel(l Level) {
	level = l
}

func Debug(base string, args ...interface{}) {
	if level <= D {
		color.Cyan(base, args...)
	}
}

func Info(base string, args ...interface{}) {
	if level <= L {
		fmt.Printf(base+"\n", args...)
	}
}

func Error(base string, args ...interface{}) {
	if level <= E {
		color.Red(base, args...)
	}
}
