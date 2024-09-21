package main

import (
	"flag"
	"fmt"
	"io"
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
	var verbose bool
	var logFile string
	var showActive bool
	var app App

	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&logFile, "f", "", "Log file (default stdout)")
	flag.BoolVar(&showActive, "active", false, "Show active connections info")
	flag.BoolVar(&app.quiet, "q", false, "Do not print anything")
	flag.Var(&app.bindAddress, "b", "Local address to listen on")
	flag.Var(&app.remotePorts, "p", "Remote ports to connect to, e.g. '22,2222'")
	flag.Parse()

	if len(app.bindAddress) == 0 {
		app.bindAddress = StringList{"0.0.0.0:22"}
	}

	if showActive {
		data, err := ReadStats()
		if err != nil && err != io.EOF {
			app.Error(err.Error())
		} else {
			fmt.Print(string(data))
		}
		return
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	logger, err := NewLogger(logFile, level)
	if err != nil {
		app.Error(err.Error())
		return
	}
	slog.SetDefault(logger)

	if err = app.Run(); err != nil {
		app.Error(err.Error())
	}
}

func NewLogger(file string, level slog.Level) (*slog.Logger, error) {
	w := os.Stdout
	if file != "" {
		var err error
		w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})), nil
}
