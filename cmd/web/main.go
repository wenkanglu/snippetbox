package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/wenkanglu/snippetbox/internal/metrics"
	"github.com/wenkanglu/snippetbox/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	metrics        *metrics.Metrics
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	appAddrEnv := os.Getenv("APP_ADDRESS")
	metricsAddrEnv := os.Getenv("METRICS_ADDRESS")
	pgConnEnv := os.Getenv("POSTGRES_CONN")

	appAddr := flag.String("app-addr", appAddrEnv, "HTTP network address")
	metricsAddr := flag.String("metrics-addr", metricsAddrEnv, "Metrics network address")
	pgConn := flag.String("db", pgConnEnv, "Postgres database URL")
	flag.Parse()

	pool, err := pgxpool.New(context.Background(), *pgConn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	m, h := metrics.NewMetrics("snippetbox")

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{Pool: pool},
		users:          &models.UserModel{Pool: pool},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		metrics:        m,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	appServer := &http.Server{
		Addr:         *appAddr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// TODO: I think these go routines ruin the panic recovery - must check again
	go func() {
		logger.Info("Starting server", slog.String("addr", *appAddr))
		err = appServer.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		logger.Error(err.Error())
		os.Exit(1)
	}()

	metricsServer := &http.Server{
		Addr:      *metricsAddr,
		Handler:   app.metricsRoutes(h),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
	}

	go func() {
		logger.Info("Starting metrics server", slog.String("addr", *metricsAddr))
		err = metricsServer.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		logger.Error(err.Error())
		os.Exit(1)
	}()

	select {}
}
