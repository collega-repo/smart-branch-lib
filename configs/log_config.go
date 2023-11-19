package configs

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

var Loggers zerolog.Logger

func NewLogger(path, fileName string) {
	pathSplit := strings.Split(path, `.`)

	if len(pathSplit) > 1 {
		path = pathSplit[0]
	}

	if !strings.HasSuffix(path, `/`) {
		path = fmt.Sprintf(`%s/`, path)
	}

	if _, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0777); err != nil {
			panic(err)
		}
	}

	path = fmt.Sprintf(`%s%s.out`, path, fileName)
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	Loggers = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"},
		logFile,
	)).With().Timestamp().Logger()
}

func NewLoggerReqRes(name string, date time.Time, multiOutput bool, path, fileName string) zerolog.Logger {
	pathSplit := strings.Split(path, `.`)
	if len(pathSplit) > 1 {
		path = pathSplit[0]
	}

	if !strings.HasSuffix(path, `/`) {
		path = fmt.Sprintf(`%s/`, path)
	}

	if _, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0777); err != nil {
			fmt.Println(err.Error())
		}
	}

	path = fmt.Sprintf(`%s%s_%s_%s.out`, path, fileName, date.Format(time.DateOnly), name)
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
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

type modeLog string

const (
	Debug modeLog = "DEBUG"
	Info  modeLog = "INFO"
	Trace modeLog = "TRACE"
	Error modeLog = "ERROR"
	Warn  modeLog = "WARNING"
)

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
	Mode           modeLog
}

func (l LogEvent) SendRequest(logger zerolog.Logger) {
	var eventLog *zerolog.Event
	if l.Error != nil {
		eventLog = logger.Err(l.Error)
	} else {
		switch l.Mode {
		case Debug:
			eventLog = logger.Debug()
		case Trace:
			eventLog = logger.Trace()
		case Error:
			eventLog = logger.Error()
		case Warn:
			eventLog = logger.Warn()
		default:
			eventLog = logger.Info()
		}
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
		switch l.Mode {
		case Debug:
			eventLog = logger.Debug()
		case Trace:
			eventLog = logger.Trace()
		case Error:
			eventLog = logger.Error()
		case Warn:
			eventLog = logger.Warn()
		default:
			eventLog = logger.Info()
		}
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
