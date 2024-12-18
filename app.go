package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

type App struct {
	config Config
}

func (app App) Error(format string, args ...any) {
	if !app.config.Quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func (app App) LogInfo(format string, args ...any) {
	if !app.config.Quiet {
		slog.Info(format, args...)
	}
}

func (app App) LogDebug(format string, args ...any) {
	if !app.config.Quiet {
		slog.Debug(format, args...)
	}
}

func (app App) LogError(format string, args ...any) {
	if !app.config.Quiet {
		slog.Error(format, args...)
	}
}

func (app App) Run() error {
	ctx := context.WithValue(context.Background(), "app", &app)
	stats := Stats{}
	go stats.RunServer(ctx)

	if len(app.config.Bind) > 1 {
		var g *errgroup.Group
		g, ctx = errgroup.WithContext(ctx)
		for _, addr := range app.config.Bind {
			g.Go(func() error {
				return app.Listen(ctx, addr, stats)
			})
		}
		return g.Wait()
	}

	return app.Listen(ctx, app.config.Bind[0], stats)
}

func (app App) Listen(ctx context.Context, addr string, stats Stats) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	app.LogInfo("listening", "addr", addr)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		var c net.Conn
		c, err = ln.Accept()
		if err != nil {
			app.LogError("accepting", "reason", err.Error(), "addr", ln.Addr())
			continue
		}
		key := c.RemoteAddr().String()
		if _, ok := stats[key]; ok {
			_ = c.Close()
			app.LogInfo("rejected", "reason", "already connected", "raddr", c.RemoteAddr())
			continue
		}
		info := &ConnInfo{}
		stats[key] = info
		go func() {
			defer func() {
				_ = c.Close()
				delete(stats, key)
			}()
			app.HandleConnection(c, info)
		}()
	}
}

func (app App) HandleConnection(localConn net.Conn, info *ConnInfo) {
	startTime := time.Now()
	info.StartTime = startTime
	laddr := localConn.LocalAddr()
	raddr := localConn.RemoteAddr()

	app.LogInfo("accepted", "laddr", laddr, "raddr", raddr)

	defer func() {
		duration := time.Now().Sub(startTime) - time.Second
		app.LogInfo("closed",
			"laddr", laddr,
			"raddr", raddr,
			"lifetime", int(duration.Seconds()),
			"reversed", info.IsReversed)
	}()

	for _, port := range app.config.RemotePorts {
		if err := app.ConnectRemote(localConn, port, info); err != nil {
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
	defer ticker.Stop()
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

func (app App) ConnectRemote(localConn net.Conn, port int, info *ConnInfo) error {
	remoteAddr, _ := localConn.RemoteAddr().(*net.TCPAddr)
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteAddr.IP, port))
	if err != nil {
		return err
	}
	laddr := remoteConn.LocalAddr()
	raddr := remoteConn.RemoteAddr()
	defer func() {
		remoteConn.Close()
		app.LogInfo("closed", "laddr", laddr, "raddr", raddr)
	}()
	app.LogInfo("connected", "laddr", laddr, "raddr", raddr)
	info.IsReversed = true
	go io.Copy(localConn, remoteConn)
	_, err = io.Copy(remoteConn, localConn)
	return err
}
