package configs

import (
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

var Loggers zerolog.Logger

func NewLogger() {
	logFile, err := os.OpenFile(commons.Configs.Log.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	Loggers = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"},
		logFile,
	)).With().Timestamp().Logger()
}

func NewLoggerReqRes(name string, date time.Time, multiOutput bool) zerolog.Logger {
	path := commons.Configs.Log.Path
	pathSplit := strings.Split(path, `.`)
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
