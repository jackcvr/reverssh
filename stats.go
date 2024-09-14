package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

const StatServerSock = "/var/run/reverssh.sock"

type Stats map[string]*Record

type Record struct {
	StartTime  time.Time
	IsReversed bool
}

func (stats Stats) RunServer(ctx context.Context) {
	_ = os.Remove(StatServerSock)
	defer os.Remove(StatServerSock)

	app := ctx.Value("app").(*App)
	s, err := net.Listen("unix", StatServerSock)
	if err != nil {
		app.LogError("listening", "reason", err)
		return
	}
	defer s.Close()

	app.LogInfo("listening", "addr", StatServerSock)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var c net.Conn
		c, err = s.Accept()
		if err != nil {
			app.LogError("accepting", "reason", err.Error())
			continue
		}
		res := ""
		now := time.Now()
		for addr, record := range stats {
			lifetime := now.Sub(record.StartTime)
			res += fmt.Sprintf("%s lifetime=%d reversed=%t\n", addr, int(lifetime.Seconds()), record.IsReversed)
		}
		if _, err = c.Write([]byte(res)); err != nil {
			app.LogError("writing", "reason", err.Error())
		}
		_ = c.Close()
	}
}

func ReadStats() ([]byte, error) {
	c, err := net.Dial("unix", StatServerSock)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1024)
	var n int
	n, err = c.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}