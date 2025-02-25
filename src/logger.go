package main

import (
	"log"
	"os"
	"strings"
)

var logger Logger

// map of log levels to their integer values
var logMap = map[string]int{
	"DEBUG":   0,
	"INFO":    1,
	"WARNING": 2,
	"ERROR":   3,
}

type Logger struct {
	logLevel string
}

func (l *Logger) genericPrint(level string, message string) {
	log.Println(level + " - " + message)
}

func (l *Logger) Debug(message string) {
	if logMap[strings.ToUpper(l.logLevel)] <= logMap["DEBUG"] {
		l.genericPrint("DEBUG", message)
	}
}

func (l *Logger) Info(message string) {
	if logMap[strings.ToUpper(l.logLevel)] <= logMap["INFO"] {
		l.genericPrint("INFO", message)
	}
}

func (l *Logger) Error(message string) {
	if logMap[strings.ToUpper(l.logLevel)] <= logMap["ERROR"] {
		l.genericPrint("ERROR", message)
		notifyer.Error(message)
	}
}

func (l *Logger) Warning(message string) {
	if logMap[strings.ToUpper(l.logLevel)] <= logMap["WARNNING"] {
		l.genericPrint("WARNING", message)
		notifyer.Error(message)
	}
}

func (l *Logger) Fatal(message string) {
	l.genericPrint("FATAL", message)
	notifyer.Error(message)
	os.Exit(1)
}
