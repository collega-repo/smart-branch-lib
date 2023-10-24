package configs

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

var Loggers zerolog.Logger

func NewLogger(path ...string) {
	if len(path) == 0 {
		path = append(path, "logger.out")
	}
	logFile, err := os.OpenFile(path[0], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	Loggers = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"},
		logFile,
	)).With().Timestamp().Logger()
}

func NewLoggerReqRes(name string, date time.Time, multiOutput bool, path ...string) zerolog.Logger {
	if len(path) == 0 {
		path = append(path, "logger.out")
	}
	pathSplit := strings.Split(path[0], `.`)
	pathName := fmt.Sprintf(`%s_%s_%s`, pathSplit[0], date.Format(`2006-01-02`), name)
	if len(pathSplit) > 1 {
		pathName = fmt.Sprintf(`%s.%s`, pathName, pathSplit[1])
	}
	logFile, err := os.OpenFile(pathName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println(err.Error())
	}
	if multiOutput {
		return zerolog.New(zerolog.MultiLevelWriter(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"},
			logFile,
		)).With().Timestamp().Logger()
	}
	return zerolog.New(logFile).With().Timestamp().Logger()
}
