package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Writer(prefix string, suffix string) io.Writer
}

// dualLogger implements Logger, while being able to write to both STDOUT and a file
type dualLogger struct {
	config LoggerConfig // preserve for internal debugging purpose

	mutex     *sync.Mutex
	waitGroup *sync.WaitGroup

	fileLogger *log.Logger

	exitFunc *func()

	_debug   func(string, ...interface{})
	_info    func(string, ...interface{})
	_warning func(string, ...interface{})
	_error   func(string, ...interface{})
	_fatal   func(string, ...interface{})
}

func DualLogger(conf LoggerConfig, wg *sync.WaitGroup, exFn *func()) *dualLogger {
	var dl dualLogger = dualLogger{
		config:    conf,
		waitGroup: wg,
		exitFunc:  exFn,
	}

	if loggerMutex != nil {
		dl.mutex = loggerMutex
	} else {
		dl.mutex = &sync.Mutex{}
	}

	if dl.config.Filepath != "" {
		f, err := os.OpenFile(dl.config.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // skipcq: GSC-G302
		if err != nil {
			return nil // ErrBadOpenFIle
		}
		dl.fileLogger = log.New(f, "", log.LstdFlags)
	}

	switch dl.config.Level {
	case LvlDebug:
		dl._debug = func(format string, v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Debug: ", fmt.Sprintf(format, v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Debug: ", fmt.Sprintf(format, v...), "\n")
			}
		}
		fallthrough
	case LvlInfo:
		dl._info = func(format string, v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Info: ", fmt.Sprintf(format, v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Info: ", fmt.Sprintf(format, v...), "\n")
			}
		}
		fallthrough
	case LvlWarning:
		dl._warning = func(format string, v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Warning: ", fmt.Sprintf(format, v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Warning: ", fmt.Sprintf(format, v...), "\n")
			}
		}
		fallthrough
	case LvlError:
		dl._error = func(format string, v ...interface{}) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print("Error: ", fmt.Sprintf(format, v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Error: ", fmt.Sprintf(format, v...), "\n")
			}
		}
		fallthrough
	case LvlFatal:
		dl._fatal = func(format string, v ...interface{}) {
			dl.mutex.Lock()

			if dl.config.Verbose {
				fmt.Print("Fatal: ", fmt.Sprintf(format, v...), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print("Fatal: ", fmt.Sprintf(format, v...), "\n")
			}

			dl.mutex.Unlock()

			if dl.exitFunc != nil {
				(*dl.exitFunc)()
			}

			if dl.waitGroup != nil {
				dl.waitGroup.Wait()
			}
			os.Exit(1)
		}
	default:
		return nil
	}

	return &dl
}

func (dl *dualLogger) Debug(format string, v ...interface{}) {
	dl._debug(format, v...)
}

func (dl *dualLogger) Info(format string, v ...interface{}) {
	dl._info(format, v...)
}

func (dl *dualLogger) Warning(format string, v ...interface{}) {
	dl._warning(format, v...)
}

func (dl *dualLogger) Error(format string, v ...interface{}) {
	dl._error(format, v...)
}

func (dl *dualLogger) Fatal(format string, v ...interface{}) {
	dl._fatal(format, v...)
}

func (dl *dualLogger) Writer(prefix, suffix string) io.Writer {
	return &dualLoggerWriter{
		_write: func(p []byte) (n int, err error) {
			dl.mutex.Lock()
			defer dl.mutex.Unlock()
			if dl.waitGroup != nil {
				dl.waitGroup.Add(1)
				defer dl.waitGroup.Done()
			}
			if dl.config.Verbose {
				fmt.Print(" ", fmt.Sprint(p), "\n")
			}
			if dl.fileLogger != nil {
				dl.fileLogger.Print(prefix, fmt.Sprint(p), suffix)
			}
			return len(p), nil
		},
	}
}
