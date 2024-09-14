package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(25)
	flag.CommandLine.SetOutput(os.Stderr)
}

func main() {
	app := App{
		level: Level(slog.LevelInfo),
	}

	flag.BoolVar(&app.quiet, "q", false, "Do not print anything (default false)")
	flag.Var(&app.level, "l", "Log level. Possible values: debug, info, warn, error (default info)")
	flag.StringVar(&app.file, "f", "", "Log file (default stdout)")
	flag.StringVar(&app.bindAddress, "b", "0.0.0.0:22", "Local address to listen on")
	flag.Var(&app.remotePorts, "p", "Remote ports to connect to, e.g. '22,2222'")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	logger, err := NewLogger(app.file, app.level)
	if err != nil {
		app.Error(err.Error())
		return
	}
	slog.SetDefault(logger)

	if err = app.Run(); err != nil {
		app.Error(err.Error())
	}
}

func NewLogger(file string, level Level) (*slog.Logger, error) {
	w := os.Stdout
	if file != "" {
		var err error
		w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.Level(level)})), nil
}
