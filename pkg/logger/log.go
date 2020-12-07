package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jinghzhu/goutils/utils"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel string
	console  *log.Logger
)

func init() {
	initConsole()
	initLogrus()
}

func initLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Disable logrun to output logs in local machine
	logrus.SetOutput(ioutil.Discard)
}

func initConsole() {
	console = GetConsoleLogger()
}

// SetLogLevel sets log level for logrus and local console. Accept info and error level.
// By default, it is info level.
func SetLogLevel(level string) {
	if strings.ToUpper(level) == INFO || !isValidLevel(level) {
		// by default, info level
		logrus.SetLevel(logrus.InfoLevel)
		LogLevel = INFO
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
		LogLevel = ERROR
	}
	msg := "The log level is " + LogLevel
	logrus.Infoln(msg)
	console.Println(msg)
}

func isValidLevel(level string) bool {
	level = strings.ToUpper(level)
	if level != INFO && level != ERROR {
		return false
	}
	return true
}

// Print log in simple way
func Info(msg interface{}) {
	if LogLevel == INFO {
		output(msg, INFO)
	}
}

func Error(msg interface{}) {
	output(msg, ERROR)
}

func output(msg interface{}, prefix string) {
	logrus.Println(msg)
	file, line := utils.Locate(3)
	console.Println(
		fmt.Sprintf("[%s] %s Ln%d %+v", prefix, file, line, msg),
	)
}

// Print log with fields
func InfoFields(msg string, fields Fields) {
	if LogLevel == INFO {
		outputFields(msg, fields, INFO)
	}
}

func ErrorFields(msg string, fields Fields) {
	outputFields(msg, fields, ERROR)
}

// ErrorFieldsWithErr accepts error message string, a map and an error. It adds the error into map with
// standard ERROR key.
func ErrorFieldsWithErr(msg string, fields Fields, err error) {
	fields[ERROR] = err
	outputFields(msg, fields, ERROR)
}

func outputFields(msg string, fields Fields, prefix string) {
	e := logrus.WithFields(logrus.Fields(fields))
	e.Time = time.Now()
	e.Println(msg)
	data, err := e.String()
	if err != nil {
		console.Println("Fail to get string representation for logrus.entry. " + err.Error())
		data = fmt.Sprintf("%v", fields)
	}
	file, line := utils.Locate(3)
	console.Println(fmt.Sprintf("[%s] %s Ln%d %s %s", prefix, file, line, msg, data))
}
