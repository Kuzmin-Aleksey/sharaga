package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sharaga/internal/config"
	"sharaga/internal/infrastructure/persistence/mysql"
	"sharaga/internal/server/httpserver"
	"sharaga/pkg/contextx"
	"sharaga/pkg/logx"
	"sharaga/pkg/middlewarex"
)

type App struct {
	cfg        *config.Config
	l          *slog.Logger
	httpServer *http.Server
	db         *sqlx.DB
	logFile    io.Closer
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() (err error) {
	a.initLogger()

	a.db, err = mysql.Connect(a.cfg.MySQl)
	if err != nil {
		return fmt.Errorf("mysql connect: %w", err)
	}

	// todo

	return nil
}

func (a *App) newHttpServer() {
	restServer := httpserver.NewServer()

	rtr := mux.NewRouter()
	restServer.RegisterRoutes(rtr)

	sensitiveDataMasker := logx.NewSensitiveDataMasker()

	rtr.Use(
		middlewarex.AddTraceId,
		middlewarex.Logger,
		middlewarex.RequestLogging(sensitiveDataMasker, a.cfg.HttpServer.Log.MaxRequestContentLen),
		middlewarex.ResponseLogging(sensitiveDataMasker, a.cfg.HttpServer.Log.MaxResponseContentLen),
	)

	a.httpServer = &http.Server{
		Handler:      rtr,
		Addr:         a.cfg.HttpServer.Addr,
		ReadTimeout:  a.cfg.HttpServer.ReadTimeout,
		WriteTimeout: a.cfg.HttpServer.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(a.l.Handler(), slog.LevelError),
		BaseContext: func(net.Listener) context.Context {
			return contextx.WithLogger(context.Background(), a.l)
		},
	}
}

func (a *App) initLogger() {
	out := os.Stdout

	opt := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}

	if a.cfg.Log.Debug {
		opt.Level = slog.LevelDebug
	}

	var handler slog.Handler

	if a.cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(out, opt)
	} else {
		handler = slog.NewTextHandler(out, opt)
	}

	a.l = slog.New(handler)
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.db.Close(); err != nil {
		a.l.Error("close mysql db", logx.Error(err))
	}
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.l.Error("shutdown http server", logx.Error(err))
	}
}
