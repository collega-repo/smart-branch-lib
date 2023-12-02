package stream

import (
	"bufio"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/collega-repo/smart-branch-lib/configs"
	"github.com/goccy/go-json"
	"time"
)

const (
	DeviceIdChan     = `deviceIdChan`
	LogoutChan       = `logoutChan`
	UserIdCustomer   = `userIdCustomer`
	UserIdTeller     = `userIdTeller`
	UserIdCS         = `userIdCS`
	UserIdCustomerCS = `userIdCustomerCS`
)

func StreamStringFlush(w *bufio.Writer, message string) error {
	if _, err := w.Write([]byte(fmt.Sprintf("data: %s\n\n", message))); err != nil {
		return fmt.Errorf(`error write %s: %w`, message, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf(`error flush %s: %w`, message, err)
	}
	return nil
}

func StreamFlush[T any](w *bufio.Writer, response commons.ApiResponse[T]) error {
	if response.Status == commons.SuccessResponse {
		rawJSON, err := json.Marshal(response.Data)
		if err != nil {
			return fmt.Errorf(`error marshall response: %v`, err)
		}
		if _, err := fmt.Fprintf(w, "data: %s\n\n", string(rawJSON)); err != nil {
			return fmt.Errorf(`error write %v`, err)
		}
		if err := w.Flush(); err != nil {
			return fmt.Errorf(`error flush %v`, err)
		}
	}
	return nil
}

func BackgroundStream(w *bufio.Writer, sleepDuration int) {
	valid := sleepDuration != 0
	if valid {
		for {
			if _, err := w.Write([]byte("data: connected\n\n")); err != nil {
				configs.Loggers.Err(err).Msg(`error write connected`)
				return
			}
			if err := w.Flush(); err != nil {
				configs.Loggers.Err(err).Msg(`error flush connected`)
				return
			}
			time.Sleep(time.Duration(sleepDuration) * time.Second)
		}
	}
}
