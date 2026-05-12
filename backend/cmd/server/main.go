package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/leno23/ai-api-gateway/internal/config"
	"github.com/leno23/ai-api-gateway/internal/jobs"
	applog "github.com/leno23/ai-api-gateway/internal/log"
	"github.com/leno23/ai-api-gateway/internal/router"
	"github.com/leno23/ai-api-gateway/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := applog.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := store.OpenPostgres(cfg.DatabaseDSN, log)
	if err != nil {
		log.Fatal("db", zap.Error(err))
	}
	rdb, err := store.OpenRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatal("redis", zap.Error(err))
	}

	rootCtx, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()
	if cfg.QuotaReconcileSec > 0 {
		jobs.StartQuotaReconcile(rootCtx, db, rdb, log, time.Duration(cfg.QuotaReconcileSec)*time.Second)
	}
	if cfg.RebateEnabled {
		poll := cfg.RebatePollSec
		if poll <= 0 {
			poll = 10
		}
		jobs.StartRebateWorker(rootCtx, db, rdb, log, time.Duration(poll)*time.Second)
	}

	r := router.New(cfg, db, rdb, log)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("listening", zap.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	cancelRoot()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Warn("shutdown", zap.Error(err))
	}
}
