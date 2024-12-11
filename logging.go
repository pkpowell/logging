package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type Displayf func(string, ...any)
type Display func(...string)
type Create func(string, error) error

const (
	YYYYMMDD            = "2006-01-02"
	HHMMSS24h           = "15:04:05"
	HHMMSS12h           = "3:04:05 PM"
	TextDate            = "January 2, 2006"
	TextDateWithWeekday = "Monday, January 2, 2006"
	AbbrTextDate        = "Jan 2 Mon"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

var (
	Logger, errLogger           *slog.Logger
	handler, errHandler         *slog.HandlerOptions
	tintHandler, errTintHandler *tint.Options
	logLevel                    = &slog.LevelVar{} // INFO

	ctx    context.Context
	record slog.Record
)

var (
	Errorf, Warnf, Infof, Debugf Displayf
	Error, Warn, Info, Debug     Display
	Errorw                       Create
)

func Init(verbose *bool, jsonLogs *bool) {
	ctx = context.Background()

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
		Logger = slog.New(slog.NewJSONHandler(os.Stderr, handler))
		errLogger = slog.New(slog.NewJSONHandler(os.Stderr, errHandler))
	} else {
		Logger = slog.New(tint.NewHandler(os.Stderr, tintHandler))
		errLogger = slog.New(tint.NewHandler(os.Stderr, errTintHandler))
	}

	slog.SetDefault(Logger)

	Info = func(args ...string) {
		if !Logger.Enabled(ctx, slog.LevelInfo) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelInfo, strings.Join(args, " "), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Error = func(args ...string) {
		if !Logger.Enabled(ctx, slog.LevelInfo) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelInfo, strings.Join(args, " "), pcs[0])

		_ = Logger.Handler().Handle(ctx, record)
	}

	Warn = func(args ...string) {
		if !Logger.Enabled(ctx, slog.LevelInfo) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelInfo, strings.Join(args, " "), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Debug = func(args ...string) {
		if !Logger.Enabled(ctx, slog.LevelDebug) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelDebug, strings.Join(args, " "), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Infof = func(format string, args ...any) {
		if !Logger.Enabled(ctx, slog.LevelInfo) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Warnf = func(format string, args ...any) {
		if !Logger.Enabled(ctx, slog.LevelWarn) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(format, args...), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Debugf = func(format string, args ...any) {
		if !Logger.Enabled(ctx, slog.LevelDebug) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(format, args...), pcs[0])
		_ = Logger.Handler().Handle(ctx, record)
	}

	Errorf = func(format string, args ...any) {
		if !errLogger.Enabled(ctx, slog.LevelError) {
			return
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
		_ = errLogger.Handler().Handle(ctx, record)
	}

	////////////////////////////////
	////////////////////////////////

	Errorw = func(text string, err error) error {
		if !errLogger.Enabled(ctx, slog.LevelError) {
			return nil
		}
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		record = slog.NewRecord(time.Now(), slog.LevelError, fmt.Errorf(text+" %w", err).Error(), pcs[0])
		return fmt.Errorf(record.Message)
		// return errLogger.Handler().Handle(context.Background(), r)
	}
}
