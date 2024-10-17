package main

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"time"
	_ "time/tzdata"
)

var configPath = "/etc/reverssh/reverssh.toml"

type Config struct {
	TZ          string
	Verbose     bool
	Quiet       bool
	Bind        []string
	RemotePorts []int
}

var config = Config{
	TZ:          "Europe/Vilnius",
	Verbose:     false,
	Quiet:       false,
	Bind:        []string{"0.0.0.0:22"},
	RemotePorts: []int{22},
}

func init() {
	debug.SetGCPercent(25)
	flag.CommandLine.SetOutput(os.Stderr)
}

func main() {
	var showActive bool

	flag.StringVar(&configPath, "c", configPath, "Path to TOML config file")
	flag.BoolVar(&showActive, "active", false, "Show active connections info")
	flag.Parse()

	if showActive {
		data, err := ReadStats()
		if err != nil && err != io.EOF {
			panic(err)
		} else {
			fmt.Print(string(data))
		}
		return
	}

	if data, err := os.ReadFile(configPath); err != nil {
		panic(err)
	} else if err = toml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	if loc, err := time.LoadLocation(config.TZ); err != nil {
		panic(err)
	} else {
		time.Local = loc
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	level := slog.LevelInfo
	if config.Verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	app := App{config: config}
	if err := app.Run(); err != nil {
		app.Error(err.Error())
	}
}
