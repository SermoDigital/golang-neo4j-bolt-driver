package log

import (
	"fmt"
	l "log"
	"os"
	"strings"
)

// Level is the logging level
type Level int

const (
	NoneLevel  Level = iota
	ErrorLevel Level = iota
	InfoLevel  Level = iota
	TraceLevel Level = iota
)

var (
	level    = NoneLevel
	TraceLog = l.New(os.Stderr, "[BOLT][TRACE]", l.LstdFlags)
	InfoLog  = l.New(os.Stderr, "[BOLT][INFO]", l.LstdFlags)
	ErrorLog = l.New(os.Stderr, "[BOLT][ERROR]", l.LstdFlags)
)

// SetLevel sets the logging level of this package
func SetLevel(levelStr string) {
	switch strings.ToLower(levelStr) {
	case "trace":
		level = TraceLevel
	case "info":
		level = InfoLevel
	case "error":
		level = ErrorLevel
	default:
		level = NoneLevel
	}
}

// GetLevel gets the logging level
func GetLevel() Level {
	return level
}

// Trace writes a trace log in the format of Println
func Trace(args ...interface{}) {
	if level >= TraceLevel {
		TraceLog.Println(args...)
	}
}

// Tracef writes a trace log in the format of Printf
func Tracef(msg string, args ...interface{}) {
	if level >= TraceLevel {
		TraceLog.Printf(msg, args...)
	}
}

// Info writes an info log in the format of Println
func Info(args ...interface{}) {
	if level >= InfoLevel {
		InfoLog.Println(args...)
	}
}

// Infof writes an info log in the format of Printf
func Infof(msg string, args ...interface{}) {
	if level >= InfoLevel {
		InfoLog.Printf(msg, args...)
	}
}

// Error writes an error log in the format of Println
func Error(args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Println(args...)
	}
}

// Errorf writes an error log in the format of Printf
func Errorf(msg string, args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Printf(msg, args...)
	}
}

// Fatal writes an error log in the format of Fatalln
func Fatal(args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Println(args...)
		os.Exit(1)
	}
}

// Fatalf writes an error log in the format of Fatalf
func Fatalf(msg string, args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Printf(msg, args...)
		os.Exit(1)
	}
}

// Panic writes an error log in the format of Panicln
func Panic(args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Println(args...)
		panic(fmt.Sprint(args...))
	}
}

// Panicf writes an error log in the format of Panicf
func Panicf(msg string, args ...interface{}) {
	if level >= ErrorLevel {
		ErrorLog.Printf(msg, args...)
		panic(fmt.Sprintf(msg, args...))
	}
}