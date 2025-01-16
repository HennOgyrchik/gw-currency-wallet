package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"time"

	_ "gw-currency-wallet/internal/docs"
)

type Gin struct {
	srv *http.Server
}

func New(url string, handler Handler) Gin {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.POST("/api/v1/register", handler.Register)
	router.POST("/api/v1/login", handler.Login)
	router.GET("/api/v1/wallet/balance", handler.Balance)
	router.POST("/api/v1/wallet/deposit", handler.Deposit)
	router.POST("/api/v1/wallet/withdraw", handler.Withdraw)
	router.GET("/api/v1/exchange/rates", handler.Rates)
	router.POST("/api/v1/exchange", handler.Exchange)
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return Gin{srv: &http.Server{Addr: url, Handler: router.Handler()}}
}

func (g *Gin) Start() error {
	const op = "Web Start"

	if err := g.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (g *Gin) Stop() error {
	const op = "Web Stop"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	if err = g.srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return err
}
