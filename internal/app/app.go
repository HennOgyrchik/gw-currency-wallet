package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gw-currency-wallet/internal/cache"
	"gw-currency-wallet/internal/grpcClient/auth"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logs"
	"net/http"
	"regexp"
	"strings"
)

const patternToken = "[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+"

func New(ctx context.Context, storage storages.Storage, cache cache.Cache, exchanger exchange.Exchanger, authorizer auth.Authorizer, logger *logs.Log) *App {
	return &App{ctx: ctx,
		storage:    storage,
		cache:      cache,
		exchanger:  exchanger,
		authorizer: authorizer,
		logger:     logger,
	}
}

// @Summary Registration
// @Tags Auth
// @Descriotion create account
// @ID register-account
// @Accept json
// @Produce json
// @Param input body User true "user info"
// @Success 201 {object} MessageResponseJSON
// @Failure 400 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/register [post]
func (a *App) Register(c *gin.Context) {
	const op = "App Register"

	var userRequest User

	if err := c.BindJSON(&userRequest); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request")
		return
	}

	userResponse, err := a.authorizer.CreateUser(a.ctx, auth.CreateUserRequest{
		Username: userRequest.Username,
		Password: userRequest.Password,
		Email:    userRequest.Email,
	})

	switch {
	case errors.Is(err, auth.UserAlreadyExistsErr):
		sendError(c, http.StatusBadRequest, "Username or email already exists")
		return
	case err != nil:
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed registration")
		return

	}

	if err := a.storage.NewWallet(a.ctx, userResponse.UserId); err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed create new wallet")
		return
	}

	c.JSON(http.StatusCreated, MessageResponseJSON{"User registered successfully"})
}

// @Summary Login
// @Tags Auth
// @Descriotion login account
// @ID login-account
// @Accept json
// @Produce json
// @Param input body Credentials true "user credentials"
// @Success 200 {object} TokenResponseJSON
// @Failure 400 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/login [post]
func (a *App) Login(c *gin.Context) {
	const op = "App Login"

	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := a.authorizer.Login(a.ctx, auth.LoginCredentials{
		Username: credentials.Username,
		Password: credentials.Password,
	})
	switch {
	case errors.Is(err, auth.InvalidCredentialsErr):
		sendError(c, http.StatusBadRequest, "Invalid username or password")
		return
	case err != nil:
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed login")
		return

	}

	c.JSON(http.StatusOK, TokenResponseJSON{token.Value})

}

// @Summary Balance
// @Security ApiKeyAuth
// @Tags Wallet
// @Descriotion user balance
// @ID user-balance
// @Produce json
// @Success 200 {object} storages.Balance
// @Failure 400 {object} ErrResponseJSON
// @Failure 401 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/wallet/balance [get]
func (a *App) Balance(c *gin.Context) {
	const op = "App Balance"

	user, err := a.authorization(c)
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

// @Summary Deposit
// @Security ApiKeyAuth
// @Tags Wallet
// @Descriotion deposit wallet
// @ID deposit-wallet
// @Accept json
// @Produce json
// @Param input body Cash true "desired currency and amount"
// @Success 200 {object} NewBalanceResponseJSON
// @Failure 400 {object} ErrResponseJSON
// @Failure 401 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/wallet/deposit [post]
func (a *App) Deposit(c *gin.Context) {
	a.DepositWithdrawHandler(c, 1)
}

// @Summary Withdraw
// @Security ApiKeyAuth
// @Tags Wallet
// @Descriotion withdraw wallet
// @ID withdraw-wallet
// @Accept json
// @Produce json
// @Param input body Cash true "desired currency and amount"
// @Success 200 {object} NewBalanceResponseJSON
// @Failure 400 {object} ErrResponseJSON
// @Failure 401 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/wallet/withdraw [post]
func (a *App) Withdraw(c *gin.Context) {
	a.DepositWithdrawHandler(c, -1)
}

func (a *App) DepositWithdrawHandler(c *gin.Context, multiplier float32) {
	const op = "App Deposit"

	user, err := a.authorization(c)
	if err != nil {
		return
	}

	var request Cash

	if err = c.BindJSON(&request); err != nil {
		sendError(c, http.StatusBadRequest, "Invalid amount or currency")
		return
	}

	if request.Amount < 0 {
		sendError(c, http.StatusBadRequest, "the amount is less than 0")
		return
	}

	request.Currency = strings.ToUpper(request.Currency)

	balance, err := a.storage.GetBalance(a.ctx, user)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to get user balance")
		return
	}

	_, err = changeBalance(request.Currency, &balance, request.Amount, multiplier)
	if err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = a.storage.UpdateWallet(a.ctx, user, balance); err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed update balance")
		return
	}

	c.JSON(http.StatusOK, NewBalanceResponseJSON{
		Message:    "successful",
		NewBalance: balance,
	})
}

// @Summary Rates
// @Security ApiKeyAuth
// @Tags Exchange
// @Descriotion rates exchange
// @ID rates-exchange
// @Produce json
// @Success 200 {object} exchange.Rates
// @Failure 400 {object} ErrResponseJSON
// @Failure 401 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/exchange/rates [get]
func (a *App) Rates(c *gin.Context) {
	const op = "App Rates"

	_, err := a.authorization(c)
	if err != nil {
		return
	}

	res, err := a.exchanger.GetExchangeRates(a.ctx)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to retrieve exchange rates")
		return
	}

	c.JSONP(http.StatusOK, res)
}

// @Summary Exchange
// @Security ApiKeyAuth
// @Tags Exchange
// @Descriotion exchange currency
// @ID exchange-wallet
// @Accept json
// @Produce json
// @Param input body ExchangeRequest true "desired currency and amount"
// @Success 200 {object} ExchangeResponseJSON
// @Failure 400 {object} ErrResponseJSON
// @Failure 401 {object} ErrResponseJSON
// @Failure 500 {object} ErrResponseJSON
// @Router /api/v1/exchange [post]
func (a *App) Exchange(c *gin.Context) {
	const op = "App Exchange"

	user, err := a.authorization(c)
	if err != nil {
		return
	}

	var request ExchangeRequest

	if err = c.BindJSON(&request); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request")
		return
	}

	if request.Amount < 0 {
		sendError(c, http.StatusBadRequest, "the amount is less than 0")
		return
	}

	request.FromCurrency, request.ToCurrency = strings.ToUpper(request.FromCurrency), strings.ToUpper(request.ToCurrency)

	balance, err := a.storage.GetBalance(a.ctx, user)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to get user balance")
		return
	}

	var rates exchange.Rates

	valueFromCache, ok := a.cache.Get("rates")
	if !ok {
		rates, err = a.exchanger.GetExchangeRates(a.ctx)
		if err != nil {
			a.logger.Err(op, err)
			sendError(c, http.StatusInternalServerError, "Failed to retrieve exchange rates")
			return
		}
		a.cache.Set("rates", rates)
	} else {
		rates = valueFromCache.(exchange.Rates)
	}

	fromCurrencyRate, ok := rates.Rates[request.FromCurrency]
	if !ok {
		sendError(c, http.StatusBadRequest, "Insufficient funds or invalid currencies")
		return
	}

	toCurrencyRate, ok := rates.Rates[request.ToCurrency]
	if !ok {
		sendError(c, http.StatusBadRequest, "Insufficient funds or invalid currencies")
		return
	}

	_, err = changeBalance(request.FromCurrency, &balance, request.Amount, -1)
	if err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	exchangeAmount, err := changeBalance(request.ToCurrency, &balance, request.Amount, fromCurrencyRate/toCurrencyRate)
	if err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = a.storage.UpdateWallet(a.ctx, user, balance); err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed update balance")
		return
	}

	c.JSON(http.StatusOK, ExchangeResponseJSON{
		Message:        "Exchange successful",
		ExchangeAmount: exchangeAmount,
		NewBalance:     balance,
	})

}

func sendError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrResponseJSON{message})
}

func (a *App) verifyToken(userId, token string) (bool, error) {

	response, err := a.authorizer.VerifyToken(a.ctx, auth.TokenRequest{UserId: userId, Token: token})

	return response.Ok, err
}

func getTokenFromString(raw string) (string, error) {
	const op = "App getTokenFromString"

	r, err := regexp.Compile(fmt.Sprintf(patternToken))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	result := r.FindStringSubmatch(raw)
	if len(result) < 1 {
		return "", fmt.Errorf("%s: %s", op, "invalid token")
	}

	return result[0], nil
}

func (a *App) authorization(c *gin.Context) (string, error) {
	const op = "App authorization"

	authStr := c.GetHeader("Authorization")

	if authStr == "" {
		sendError(c, http.StatusUnauthorized, "Access deny")
		return "", fmt.Errorf("%s: %s", op, "Access deny")
	}

	token, err := getTokenFromString(authStr)
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed token processing")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	jwtParser, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to parse token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := jwtParser.Claims.(jwt.MapClaims)
	if !ok {
		sendError(c, http.StatusBadRequest, "payload is absent")
		return "", fmt.Errorf("payload is absent")
	}

	userId := fmt.Sprint(claims["id"])
	if userId == "" {
		sendError(c, http.StatusBadRequest, "invalid token")
		return "", fmt.Errorf("invalid token")
	}

	ok, err = a.verifyToken(userId, token)

	switch {
	case errors.Is(err, auth.InvalidCredentialsErr):
		sendError(c, http.StatusBadRequest, "Invalid token")
		return "", fmt.Errorf("%s: %s", op, "Invalid token")
	case err != nil:
		a.logger.Err(op, err)
		sendError(c, http.StatusInternalServerError, "Failed to verify token")
		return "", fmt.Errorf("%s: %w", op, err)
	case !ok:
		sendError(c, http.StatusUnauthorized, "Access deny")
		return "", fmt.Errorf("%s: %s", op, "Access deny")
	default:
		return userId, nil
	}

}

func adder(before, amount, multiplier float32) (float32, error) {
	if multiplier < 0 && before < amount*(-multiplier) {
		return 0, fmt.Errorf("insufficient funds or invalid amount")
	}

	return amount * multiplier, nil
}

func changeBalance(currency string, wallet *storages.Balance, amount, multiplier float32) (float32, error) {
	var (
		change float32
		err    error
	)

	switch currency {
	case "USD":
		change, err = adder(wallet.USD, amount, multiplier)
		wallet.USD += change

	case "RUB":
		change, err = adder(wallet.RUB, amount, multiplier)
		wallet.RUB += change
	case "EUR":
		change, err = adder(wallet.EUR, amount, multiplier)
		wallet.EUR += change

	default:
		err = fmt.Errorf("unknown currency")
	}
	return change, err
}
