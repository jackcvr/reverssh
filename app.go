package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

type App struct {
	quiet       bool
	level       Level
	bindAddress string
	remotePorts Ports
}

func (app App) Error(format string, args ...any) {
	if !app.quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func (app App) LogInfo(format string, args ...any) {
	if !app.quiet {
		slog.Info(format, args...)
	}
}

func (app App) LogDebug(format string, args ...any) {
	if !app.quiet {
		slog.Debug(format, args...)
	}
}

func (app App) LogError(format string, args ...any) {
	if !app.quiet {
		slog.Error(format, args...)
	}
}

func (app App) Run() error {
	ln, err := net.Listen("tcp", app.bindAddress)
	if err != nil {
		return err
	}
	defer ln.Close()
	app.LogInfo("listening", "addr", app.bindAddress)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stats := Stats{}
	go stats.RunServer(context.WithValue(ctx, "app", &app))

	var c net.Conn
	for {
		c, err = ln.Accept()
		if err != nil {
			app.LogError("accepting", "reason", err.Error(), "addr", ln.Addr().String())
			continue
		}
		remoteAddr := c.RemoteAddr().String()
		app.LogInfo("accepted", "laddr", c.LocalAddr().String(), "raddr", remoteAddr)
		record := Record{}
		stats[remoteAddr] = &record
		go func() {
			defer func() {
				_ = c.Close()
				delete(stats, remoteAddr)
			}()
			app.Handle(c, &record)
		}()
	}
}

func (app App) Handle(localConn net.Conn, record *Record) {
	startTime := time.Now()
	record.StartTime = startTime
	laddr := localConn.LocalAddr().String()
	raddr := localConn.RemoteAddr().String()

	defer func() {
		duration := time.Now().Sub(startTime) - time.Second
		app.LogInfo("closed", "laddr", laddr, "raddr", raddr, "lifetime", int(duration.Seconds()))
	}()

	for _, port := range app.remotePorts {
		if err := app.ConnectRemote(localConn, port, record); err != nil {
			app.LogDebug("error", "reason", err.Error())
		} else {
			return
		}
	}

	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		buf := make([]byte, 1024)
		if n, err := localConn.Read(buf); err != nil {
			if err != io.EOF {
				app.LogError("reading", "reason", err, "laddr", laddr, "raddr", raddr)
			}
			return
		} else {
			app.LogDebug("received", "laddr", laddr, "raddr", raddr, "payload", string(buf[:n]))
		}
	}

	payload := make([]byte, 1)
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		RandBytes(payload)
		if _, err := localConn.Write(payload); err != nil {
			app.LogDebug("writing", "reason", err.Error(), "laddr", laddr, "raddr", raddr)
			return
		}
		app.LogDebug("sent", "laddr", laddr, "raddr", raddr, "payload", string(payload))
	}
}

func (app App) ConnectRemote(localConn net.Conn, port int, record *Record) error {
	remoteAddr, _ := localConn.RemoteAddr().(*net.TCPAddr)
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteAddr.IP, port))
	if err != nil {
		return err
	}
	laddr := remoteConn.LocalAddr().String()
	raddr := remoteConn.RemoteAddr().String()
	defer func() {
		remoteConn.Close()
		app.LogInfo("closed", "laddr", laddr, "raddr", raddr)
	}()
	app.LogInfo("connected", "laddr", laddr, "raddr", raddr)
	record.IsReversed = true
	return Swap(remoteConn, localConn)
}
