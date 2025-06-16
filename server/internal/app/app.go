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
	"sharaga/internal/domain/service/auth"
	"sharaga/internal/domain/service/orders"
	"sharaga/internal/domain/service/partners"
	"sharaga/internal/domain/service/products"
	"sharaga/internal/domain/service/users"
	"sharaga/internal/infrastructure/persistence/mysql"
	"sharaga/internal/infrastructure/persistence/redis"
	"sharaga/internal/server/httpserver"
	"sharaga/pkg/contextx"
	"sharaga/pkg/logx"
	"sharaga/pkg/middlewarex"
	"time"
)

type App struct {
	cfg        *config.Config
	l          *slog.Logger
	httpServer *http.Server
	db         *sqlx.DB
	redis      *redis.Connection
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

	a.redis, err = redis.NewConnection(a.cfg.Redis)

	usersRepo := mysql.NewUsersRepo(a.db)
	productsRepo := mysql.NewProductsRepo(a.db)
	partnersRepo := mysql.NewPartnersRepo(a.db)
	ordersRepo := mysql.NewOrdersRepo(a.db)

	userService := users.NewService(usersRepo)
	productService := products.NewService(productsRepo)
	partnerService := partners.NewService(partnersRepo)
	orderService := orders.NewService(ordersRepo, productsRepo)
	authService := auth.NewService(usersRepo, a.redis)

	a.newHttpServer(
		userService,
		productService,
		partnerService,
		orderService,
		authService,
	)

	return a.httpServer.ListenAndServe()
}

func (a *App) newHttpServer(
	userService *users.Service,
	productService *products.Service,
	partnerService *partners.Service,
	orderService *orders.Service,
	authService *auth.Service,
) {
	restServer := httpserver.NewServer(
		httpserver.NewUserServer(userService),
		httpserver.NewProductServer(productService),
		httpserver.NewPartnerServer(partnerService),
		httpserver.NewOrderServer(orderService),
		httpserver.NewAuthServer(authService),
		httpserver.NewMwAuth(authService, userService),
	)

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
		ReadTimeout:  time.Duration(a.cfg.HttpServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.cfg.HttpServer.WriteTimeout) * time.Second,
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

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.cfg.HttpServer.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.l.Error("shutdown http server", logx.Error(err))
	}
	if err := a.db.Close(); err != nil {
		a.l.Error("close mysql db", logx.Error(err))
	}
	if err := a.redis.Close(); err != nil {
		a.l.Error("close redis db", logx.Error(err))
	}
	if a.logFile != nil {
		if err := a.logFile.Close(); err != nil {
			a.l.Error("close log file", logx.Error(err))
		}
	}
}
