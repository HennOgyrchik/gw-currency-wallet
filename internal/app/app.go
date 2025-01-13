package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gw-currency-wallet/internal/cache"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logs"
	"net/http"
	"regexp"
	"strings"
)

const patternToken = "[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+"

func New(ctx context.Context, storage storages.Storage, cache cache.Cache, exchanger exchange.Exchanger, logger *logs.Log) *App {
	return &App{ctx: ctx,
		storage:   storage,
		cache:     cache,
		exchanger: exchanger,
		logger:    logger,
	}
}

func (a *App) Register(c *gin.Context) {
	const op = "App Register"

	var user User

	if err := c.BindJSON(&user); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request")
		return
	}

	//TODO тут надо отправить на регистрацию и получить id пользователя
	id := "MEGA_USER"

	//if err != nil {
	//	a.logger.Err(op, err)
	//	sendError(c, http.StatusInternalServerError, "Failed registration")
	//	return
	//}
	///////////////////

	if err := a.storage.NewWallet(a.ctx, id); err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed create new wallet")
		return
	}

	c.JSON(http.StatusCreated, struct {
		Message string
	}{Message: "User registered successfully"})
}

func (a *App) Login(c *gin.Context) {
	const op = "App Login"

	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request")
		return
	}

	//TODO тут надо отправить креды и получить токен
	token := "Тут будет ваш токен"
	//
	//	• Ошибка: ```401 Unauthorized```
	//	```json
	//{
	//  "error": "Invalid username or password"
	//}
	/////

	c.JSON(http.StatusOK, struct {
		Token string
	}{token})

}

func (a *App) Balance(c *gin.Context) {
	const op = "App Balance"

	user, err := authorization(c, a.logger)
	if err != nil {
		return
	}

	balance, err := a.storage.GetBalance(a.ctx, user)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to get user balance")
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (a *App) Deposit(c *gin.Context) {
	a.DepositWithdrawHandler(c, 1)
}

func (a *App) Withdraw(c *gin.Context) {
	a.DepositWithdrawHandler(c, -1)
}

func (a *App) DepositWithdrawHandler(c *gin.Context, multiplier float32) {
	const op = "App Deposit"

	user, err := authorization(c, a.logger)
	if err != nil {
		return
	}

	var cash Cash

	if err = c.BindJSON(&cash); err != nil {
		sendError(c, http.StatusBadRequest, "Invalid amount or currency")
		return
	}

	if cash.Amount < 0 {
		sendError(c, http.StatusBadRequest, "the amount is less than 0")
		return
	}

	cash.Currency = strings.ToUpper(cash.Currency)

	balance, err := a.storage.GetBalance(a.ctx, user)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to get user balance")
		return
	}

	switch cash.Currency {
	case "USD":
		balance.USD, err = adder(balance.USD, cash.Amount, multiplier)
	case "RUB":
		balance.RUB, err = adder(balance.RUB, cash.Amount, multiplier)
	case "EUR":
		balance.EUR, err = adder(balance.EUR, cash.Amount, multiplier)
	default:
		sendError(c, http.StatusBadRequest, "unknown currency")
		return
	}
	if err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = a.storage.UpdateWallet(a.ctx, user, balance); err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed update balance")
		return
	}

	c.JSON(http.StatusOK, struct {
		Message    string
		NewBalance storages.Balance
	}{"successful", balance})

}

func (a *App) Rates(c *gin.Context) {
	const op = "App Rates"

	_, err := authorization(c, a.logger)
	if err != nil {
		return
	}

	res, err := a.exchanger.GetExchangeRates(context.Background())
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to retrieve exchange rates")
		return
	}

	c.JSONP(http.StatusOK, res)
}

func (a *App) Exchange(c *gin.Context) {
	// Брать из кэша
	fmt.Println("Implement Exchange")
}

func sendError(c *gin.Context, code int, message string) {
	c.JSONP(code, struct {
		Error string
	}{message})
}

func verifyToken(token string) (bool, error) {
	const op = "App verifyToken"

	//TODO !!!!!
	fmt.Println("implement verifyToken")
	return true, nil
}

func getTokenFromString(raw string) (string, error) {
	const op = "App getTokenFromString"

	r, err := regexp.Compile(fmt.Sprintf(patternToken))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return r.FindStringSubmatch(raw)[0], nil
}

func authorization(c *gin.Context, logger *logs.Log) (string, error) {
	const op = "App authorization"

	authStr := c.GetHeader("Authorization")
	fmt.Println(authStr)
	if authStr == "" {
		sendError(c, http.StatusUnauthorized, "Access deny")
		return "", fmt.Errorf("%s: %s", op, "Access deny")
	}

	token, err := getTokenFromString(authStr)
	if err != nil {
		logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed token processing")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	ok, err := verifyToken(token)
	if err != nil {
		logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to verify token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if !ok {
		sendError(c, http.StatusUnauthorized, "Access deny")
		return "", fmt.Errorf("%s: %s", op, "Access deny")
	}

	t, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to parse token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		sendError(c, http.StatusBadRequest, "payload is absent")
		return "", fmt.Errorf("payload is absent")
	}

	return fmt.Sprint(claims["user"]), nil
}

func adder(before, amount, multiplier float32) (float32, error) {

	if multiplier == -1 && before < amount {
		return 0, fmt.Errorf("insufficient funds or invalid amount")
	}

	return before + (amount * multiplier), nil
}
