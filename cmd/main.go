package main

import (
	"context"
	"flag"
	"gw-currency-wallet/internal/app"
	in_mem "gw-currency-wallet/internal/cache/in-mem"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/grpcClient/auth"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/internal/web"
	"gw-currency-wallet/pkg/logs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Wallets API
// @version 1.0
// @description API Server for Wallets Application

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	logger := logs.New(os.Stdout)

	confPath := flag.String("c", "config.env", "path to configuration")
	flag.Parse()

	if err := config.LoadConfig(*confPath); err != nil {
		logger.Err("read configuration", err)
		return
	}

	cfg := config.New()

	dbUrl, err := cfg.Postgres.ConnectionURL()
	if err != nil {
		logger.Err("read db url", err)
		return
	}

	db := postgres.New()

	if err = db.Start(ctx, dbUrl, time.Duration(cfg.Postgres.ConnTimeout)*time.Second, "internal/storages/migrations"); err != nil {
		logger.Err("connection db", err)
		return
	}
	defer db.Stop()

	exchger := exchange.New(cfg.Exchanger.ConnectionURL())
	if err = exchger.Run(); err != nil {
		logger.Err("connection exchange", err)
		return
	}
	defer exchger.Stop()

	authorizer := auth.New(cfg.Auth.ConnectionURL())
	if err = authorizer.Run(); err != nil {
		logger.Err("connection authorizer", err)
		return
	}
	defer exchger.Stop()

	cache := in_mem.New(60 * time.Second)

	srv := app.New(ctx, db, cache, exchger, authorizer, logger)

	webSrv := web.New(cfg.Web.ConnectionURL(), srv)

	go func() {
		<-ctx.Done()
		if err = webSrv.Stop(); err != nil {
			logger.Err("Closing web server", err)
			return
		}
		logger.Info("Closing", logs.Attr{Key: "code", Value: "0"})
	}()

	err = webSrv.Start()
	if err != nil {
		logger.Err("Start web server", err)
		return
	}

}
