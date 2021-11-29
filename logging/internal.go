package logging

import (
	"fmt"
	"os"
)

// Non-block
func _Debug(format string, v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _debug(format, v...)
}

// Non-block
func _Info(format string, v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _info(format, v...)
}

// Non-block
func _Warning(format string, v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _warning(format, v...)
}

// Non-block
func _Error(format string, v ...interface{}) {
	if loggerWaitGroup != nil {
		loggerWaitGroup.Add(1)
	}
	go _error(format, v...)
}

// Block!
func _Fatal(format string, v ...interface{}) {
	_fatal(format, v...) // Not calling as goroutine because non-block.
}

func _debug(format string, v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("DEBUG: ", fmt.Sprintf(format, v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("DEBUG: ", fmt.Sprintf(format, v...), "\n")
	}
}

func _info(format string, v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("INFO: ", fmt.Sprintf(format, v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("INFO: ", fmt.Sprintf(format, v...), "\n")
	}
}

func _warning(format string, v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("WARNING: ", fmt.Sprintf(format, v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("WARNING: ", fmt.Sprintf(format, v...), "\n")
	}
}

func _error(format string, v ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if loggerWaitGroup != nil {
		defer loggerWaitGroup.Done()
	}
	if verboseLogging {
		fmt.Print("ERROR: ", fmt.Sprintf(format, v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("ERROR: ", fmt.Sprintf(format, v...), "\n")
	}
}

func _fatal(format string, v ...interface{}) {
	loggerMutex.Lock()

	if verboseLogging {
		fmt.Print("FATAL: ", fmt.Sprintf(format, v...), "\n")
	}
	if fileLogger != nil {
		fileLogger.Print("FATAL: ", fmt.Sprintf(format, v...), "\n")
	}

	loggerMutex.Unlock()

	if exitingFunc != nil { // if set, call exitingFunc() first to clear goroutines
		(*exitingFunc)()
	}

	if loggerWaitGroup != nil {
		loggerWaitGroup.Wait() // Fatal will exit the system, so make sure the WaitGroup is cleared before that.
	}
	os.Exit(1)
}
