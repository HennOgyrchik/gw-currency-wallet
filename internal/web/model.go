package web

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Balance(ctx *gin.Context)
	Deposit(ctx *gin.Context)
	Withdraw(ctx *gin.Context)
	Rates(ctx *gin.Context)
	Exchange(ctx *gin.Context)
}
