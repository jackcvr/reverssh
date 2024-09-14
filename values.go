package main

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Level slog.Level

func (l *Level) String() (s string) {
	return fmt.Sprintf("%d", *l)
}

func (l *Level) Set(value string) error {
	switch value {
	case "debug":
		*l = Level(slog.LevelDebug)
	case "info":
		*l = Level(slog.LevelInfo)
	case "warn":
		*l = Level(slog.LevelWarn)
	case "error":
		*l = Level(slog.LevelError)
	default:
		return errors.New("invalid level")
	}
	return nil
}

type Ports []int

func (p *Ports) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *Ports) Set(s string) error {
	*p = Ports{}
	ps := strings.Split(s, ",")
	for _, v := range ps {
		port, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*p = append(*p, port)
	}
	return nil
}
