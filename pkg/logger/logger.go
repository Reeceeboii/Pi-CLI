package logger

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/update"
	"github.com/fatih/color"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"sync"
)

const logFileName = ".picli.log"

// the location of the log file and a sync.Once instance for use in its access function
var (
	logFileLocation     string
	logFileLocationOnce sync.Once
)

/*
	The live PiCLIFileLogger instance used by the program. We can initialise this under the assumption that
	the user doesn't want to log to a file (more often that not this is the case). However, the enabled flag
	can easily be changed later if logging is to be enabled.
*/
var LivePiCLILogger = NewPiCLILogger(false)

type PiCLIFileLogger struct {
	// Is the logger actually enabled or not?
	Enabled bool
	// Logger used to output general information messages to the log file
	infoLog *log.Logger
	// Logger used to output error messages to the log file
	errLog *log.Logger
	// Logger used to output Pi-CLI's status messages to the log file
	statusLog *log.Logger
	// Logger used to output command information to the log file
	commandLog *log.Logger
	// Handle to the log file
	LogFileHandle *os.File
}

// NewPiCLILogger returns a new PiCLIFileLogger instance
func NewPiCLILogger(fileLoggingEnabled bool) *PiCLIFileLogger {
	logFile, err := getLogFileHandle(getLogFileLocation())
	if err != nil {
		color.Red("Failed to open file handle to log file")
		log.Fatal(err)
	}

	l := &PiCLIFileLogger{
		Enabled:       fileLoggingEnabled,
		infoLog:       log.New(logFile, "[info]    ", log.Ltime),
		errLog:        log.New(logFile, "[error]   ", log.Ltime|log.Lshortfile),
		statusLog:     log.New(logFile, "[status]  ", log.Ltime),
		commandLog:    log.New(logFile, "[command] ", log.Ltime),
		LogFileHandle: logFile,
	}

	return l
}

// LogStartupInformation logs various bits of runtime/OS related information
func (fl *PiCLIFileLogger) LogStartupInformation() {
	if fl.Enabled {
		fl.infoLog.Println("")
		fl.infoLog.Println("Pi-CLI Version: v" + update.Version)
		fl.infoLog.Println("SHA short: " + update.GitHash)
		fl.infoLog.Println("GOOS: " + runtime.GOOS)
		fl.infoLog.Println("Runtime v: " + runtime.Version())
		fl.infoLog.Println("GOARCH: " + runtime.GOARCH)
		fl.infoLog.Println()
	}
}

// LogInformation prints information with the info logger
func (fl *PiCLIFileLogger) LogInformation(a interface{}) {
	if fl.Enabled {
		fl.infoLog.Println(a)
	}
}

// LogError prints error information with the error logger
func (fl PiCLIFileLogger) LogError(a interface{}) {
	if fl.Enabled {
		fl.errLog.Println(a)
	}
}

// LogStatus prints status information with the status logger
func (fl PiCLIFileLogger) LogStatus(v interface{}) {
	if fl.Enabled {
		fl.statusLog.Println(v)
	}
}

// LogCommand prints command information with the command logger
func (fl PiCLIFileLogger) LogCommand(a interface{}) {
	if fl.Enabled {
		fl.commandLog.Println(a)
	}
}

// getLogFileHandle returns an open handle to the log file
func getLogFileHandle(logFileLocation string) (*os.File, error) {
	file, err := os.OpenFile(logFileLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// getLogFileLocation will return the expected path of Pi-CLI's log file
func getLogFileLocation() string {
	logFileLocationOnce.Do(func() {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		if runtime.GOOS == "windows" {
			logFileLocation = strings.ReplaceAll(path.Join(usr.HomeDir, logFileName), "/", "\\")
		} else {
			logFileLocation = path.Join(usr.HomeDir, logFileName)
		}
	})
	return logFileLocation
}
