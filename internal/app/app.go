package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages"
)

func New(ctx context.Context, storage storages.Storage, exchanger exchange.Exchanger) *App {
	return &App{ctx: ctx, storage: storage, exchanger: exchanger}
}

func (a *App) Register(c *gin.Context) {

	fmt.Println("Implement Register")
}

func (a *App) Login(c *gin.Context) {

	fmt.Println("Implement Login")
}

func (a *App) Balance(c *gin.Context) {

	fmt.Println("Implement Balance")
}

func (a *App) Deposit(c *gin.Context) {

	fmt.Println("Implement Deposit")
}

func (a *App) Withdraw(c *gin.Context) {

	fmt.Println("Implement Withdraw")
}

func (a *App) Rates(c *gin.Context) {

	res, err := a.exchanger.GetExchangeRates(context.Background())

	fmt.Println(res, err)
}

func (a *App) Exchange(c *gin.Context) {

	fmt.Println("Implement Exchange")
}
