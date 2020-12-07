package config

import (
	"log/syslog"
	"os"
)

type Logger struct {
	syslog.Writer
}

type LoggingLocalOpts struct {
	Output   *os.File
	Flag     int
	Priority syslog.Priority
	Tag      string
}

func defaultLoggingLocalOpts() (logopts *LoggingLocalOpts) {
	return &LoggingLocalOpts{
		Output:   os.Stdout,
		Flag:     defaultLoggingFlag,
		Priority: defaultLoggingPriority,
	}
}
