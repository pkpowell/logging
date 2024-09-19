package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
)

type Display func(string, ...any)
type Create func(string, error) error

const (
	YYYYMMDD            = "2006-01-02"
	HHMMSS24h           = "15:04:05"
	HHMMSS12h           = "3:04:05 PM"
	TextDate            = "January 2, 2006"
	TextDateWithWeekday = "Monday, January 2, 2006"
	AbbrTextDate        = "Jan 2 Mon"
)

var (
	logger, errLogger           *slog.Logger
	handler, errHandler         *slog.HandlerOptions
	tintHandler, errTintHandler *tint.Options
	logLevel                    = &slog.LevelVar{} // INFO
)

var (
	Errorf, Warnf, Infof, Debugf Display
	Errorw                       Create
)

func Init(verbose *bool, jsonLogs *bool) {
	handler = &slog.HandlerOptions{
		Level: logLevel,
	}

	errHandler = &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	tintHandler = &tint.Options{
		Level: logLevel,
	}

	errTintHandler = &tint.Options{
		Level:     logLevel,
		AddSource: true,
	}

	if *verbose {
		fmt.Println("Verbose logging")
		logLevel.Set(slog.LevelDebug)
		handler.AddSource = true
		tintHandler.AddSource = true
	} else {
		tintHandler.TimeFormat = HHMMSS24h
		errTintHandler.TimeFormat = HHMMSS24h
	}

	if *jsonLogs {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, handler))
		errLogger = slog.New(slog.NewJSONHandler(os.Stderr, errHandler))
	} else {
		logger = slog.New(tint.NewHandler(os.Stderr, tintHandler))
		errLogger = slog.New(tint.NewHandler(os.Stderr, errTintHandler))
	}

	slog.SetDefault(logger)

	Infof = func(format string, args ...any) {
		if !logger.Enabled(context.Background(), slog.LevelInfo) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
		_ = logger.Handler().Handle(context.Background(), r)
	}

	Warnf = func(format string, args ...any) {
		if !logger.Enabled(context.Background(), slog.LevelWarn) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(format, args...), pcs[0])
		_ = logger.Handler().Handle(context.Background(), r)
	}

	Debugf = func(format string, args ...any) {
		if !logger.Enabled(context.Background(), slog.LevelDebug) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(format, args...), pcs[0])
		_ = logger.Handler().Handle(context.Background(), r)
	}

	Errorf = func(format string, args ...any) {
		if !errLogger.Enabled(context.Background(), slog.LevelError) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
		_ = errLogger.Handler().Handle(context.Background(), r)
	}

	////////////////////////////////
	////////////////////////////////

	Errorw = func(text string, err error) error {
		if !errLogger.Enabled(context.Background(), slog.LevelError) {
			return nil
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Errorf(text+" %w", err).Error(), pcs[0])
		return errLogger.Handler().Handle(context.Background(), r)
	}
}
