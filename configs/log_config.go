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

type LogEvent struct {
	IdempotencyKey string
	RequestId      string
	UserId         string
	UserAgent      string
	BodyRequest    []byte
	BodyResponse   []byte
	Error          error
	Message        string
	StartTime      time.Time
	Ip             string
	Host           string
	Path           string
	Method         string
	Protocol       string
}

func (l LogEvent) SendRequest(logger zerolog.Logger) {
	var eventLog *zerolog.Event
	if l.Error != nil {
		eventLog = logger.Err(l.Error)
	} else {
		eventLog = logger.Debug()
	}
	if string(l.BodyRequest) != "" {
		eventLog.RawJSON("reqBody", l.BodyRequest)
	}
	if l.IdempotencyKey != "" {
		eventLog.Str(`idempotencyKey`, l.IdempotencyKey)
	}
	if l.RequestId != "" {
		eventLog.Str(`requestId`, l.RequestId)
	}
	if l.UserId != "" {
		eventLog.Str(`userId`, l.UserId)
	}
	if l.UserAgent != "" {
		eventLog.Str(`userAgent`, l.UserAgent)
	}
	if l.Ip != "" {
		eventLog.Str(`ip`, l.Ip)
	}
	if l.Host != "" {
		eventLog.Str(`host`, l.Host)
	}
	if l.Path != "" {
		eventLog.Str(`path`, l.Path)
	}

	eventLog.
		Str(`startTime`, l.StartTime.Format("2006-01-02 15:04:05.000")).
		Str(`method`, l.Method).
		Str(`protocol`, l.Protocol)

	go func() {
		if l.Message != "" {
			eventLog.Msg(l.Message)
		} else {
			eventLog.Send()
		}
	}()
}

func (l LogEvent) SendResponse(logger zerolog.Logger) {
	var eventLog *zerolog.Event
	if l.Error != nil {
		eventLog = logger.Err(l.Error)
	} else {
		eventLog = logger.Info()
	}
	if string(l.BodyResponse) != "" {
		eventLog.RawJSON(`resBody`, l.BodyResponse)
	}
	if string(l.BodyRequest) != "" {
		eventLog.RawJSON("reqBody", l.BodyRequest)
	}
	if l.IdempotencyKey != "" {
		eventLog.Str(`idempotencyKey`, l.IdempotencyKey)
	}
	if l.RequestId != "" {
		eventLog.Str(`requestId`, l.RequestId)
	}
	if l.UserId != "" {
		eventLog.Str(`userId`, l.UserId)
	}
	if l.UserAgent != "" {
		eventLog.Str(`userAgent`, l.UserAgent)
	}
	if l.Ip != "" {
		eventLog.Str(`ip`, l.Ip)
	}
	if l.Host != "" {
		eventLog.Str(`host`, l.Host)
	}
	if l.Path != "" {
		eventLog.Str(`path`, l.Path)
	}

	eventLog.
		Dur(`duration`, time.Since(l.StartTime)).
		Str(`startTime`, l.StartTime.Format("2006-01-02 15:04:05.000")).
		Str(`endTime`, time.Now().Format("2006-01-02 15:04:05.000")).
		Str(`method`, l.Method).
		Str(`protocol`, l.Protocol)

	go func() {
		if l.Message != "" {
			eventLog.Msg(l.Message)
		} else {
			eventLog.Send()
		}
	}()
}
