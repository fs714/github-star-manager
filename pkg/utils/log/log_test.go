package log

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestDefaultLog(t *testing.T) {
	Debug(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Debugf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Debugw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Info(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Infof("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Infow("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Error(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Errorf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Errorw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())

	ZapInfo("message", Int("code", 200), String("Src", "192.168.10.64"), String("Method", "GET"),
		String("Url", "/api/v1/health"), Time("Time", time.Now()))
}

func TestCustomizedLog(t *testing.T) {
	logger := New(os.Stderr, JsonFormat, DebugLevel, true)
	ResetCurrentLog(logger)

	Debug(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Debugf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Debugw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Info(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Infof("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Infow("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Error(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Errorf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Errorw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())

	ZapInfo("message", Int("code", 200), String("Src", "192.168.10.64"), String("Method", "GET"),
		String("Url", "/api/v1/health"), Time("Time", time.Now()))
}

func TestCustomizedTeeLog(t *testing.T) {
	file, err := os.OpenFile("/tmp/go_proj_boot.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var tops = []TeeOption{
		{
			W: os.Stderr,
			Lef: func(lvl Level) bool {
				return lvl <= InfoLevel
			},
			F: ConsoleFormat,
		},
		{
			W: file,
			Lef: func(lvl Level) bool {
				return lvl > InfoLevel
			},
			F: JsonFormat,
		},
	}

	logger := NewTee(tops, true)
	ResetCurrentLog(logger)

	Debug(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Debugf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Debugw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Info(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Infof("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Infow("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Error(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Errorf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Errorw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())

	ZapInfo("message", Int("code", 200), String("Src", "192.168.10.64"), String("Method", "GET"),
		String("Url", "/api/v1/health"), Time("Time", time.Now()))
}

func TestCustomizedTeeWithRotateLog(t *testing.T) {
	var tops = []TeeWithRotateOption{
		{
			Filename:   "/tmp/go_proj_boot/go_proj_boot.log",
			MaxSize:    10,
			MaxAge:     16,
			MaxBackups: 16,
			Compress:   true,
			Lef: func(lvl Level) bool {
				return lvl >= InfoLevel
			},
			F: ConsoleFormat,
		},
	}

	logger := NewTeeWithRotate(tops, true)
	ResetCurrentLog(logger)

	Debug(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Debugf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Debugw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Info(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Infof("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Infow("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())
	Error(200, ", 192.168.10.64", ", GET", ", /api/v1/health, ", time.Now())
	Errorf("%d, %s, %s, %s, %s", 200, "192.168.10.64", "GET", "/api/v1/health", time.Now())
	Errorw("message", "code", 200, "Src", "192.168.10.64", "Method", "GET", "Url", "/api/v1/health", "Time", time.Now())

	ZapInfo("message", Int("code", 200), String("Src", "192.168.10.64"), String("Method", "GET"),
		String("Url", "/api/v1/health"), Time("Time", time.Now()))
}

func parseStrToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse %s to int", s)
	}

	return i, nil
}

func middleFunc(s string) error {
	_, err := parseStrToInt(s)
	if err != nil {
		return errors.WithMessagef(err, "failed to parseStrToInt for %s", s)
	}

	return nil
}

func TestConsoleLogWithError(t *testing.T) {
	err := middleFunc("aaa")
	if err != nil {
		Error("middleFunc failed", ", err: ", err)
		Errorf("middleFunc failed with err: %+v", err)
		Errorw("middleFunc failed", "err", err)
	}
}

func TestJsonLogWithError(t *testing.T) {
	logger := New(os.Stderr, JsonFormat, InfoLevel, true)
	ResetCurrentLog(logger)

	err := middleFunc("aaa")
	if err != nil {
		Error("middleFunc failed", ", err: ", err)
		Errorf("middleFunc failed with err: %+v", err)
		Errorw("middleFunc failed", "err", err)
	}
}
