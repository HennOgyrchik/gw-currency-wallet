# Currency wallets
Представляет собой API для работы с валютным кошельком.

Зависимости:
* [gw-exchanger](https://github.com/HennOgyrchik/gw-exchanger.git)
* [gw-authorizer](https://github.com/HennOgyrchik/gw-authorizer.git)

## Endpoints
### Pегистрация нового пользователя
`POST /api/v1/register`
* Принимает JSON
```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```
* Выполняет gRPC-запрос `gw-authorizer.CreateUser`
* Если пользователь уже существует, то возвращает ошибку `400 BadRequest`
* При успешном выполнении возвращает `201 Created`
---
###  Авторизация пользователя
`POST /api/v1/login`
* Принимает JSON
```json
{
  "username": "string",
  "password": "string"
}
```
* Выполняет gRPC-запрос `gw-authorizer.Login`
* Если авторизация неуспешна, то возвращает ошибку `401 Unauthorized`
* При успешном выполнении возвращает `200 Ok` и JWT-токен
 ```json
{
  "token": "string"
}
```
---
### Получение баланса пользователя
`GET /api/v1/wallet/balance`

Требуется заголовок `Authorization: Bearer JWT_TOKEN`

* Выполняет gRPC-запрос `gw-authorizer.VerifyToken`
* Если авторизация неуспешна, то возвращает `401 Unauthorized`
* При успешной проверки авторизации возвращает `200 Ok` и баланс пользователя
 ```json {
{
    "balance": {
      "USD": "float",
      "RUB": "float",
      "EUR": "float"
   }
}
```
 ---
###  Пополнение баланса
`POST /api/v1/wallet/deposit`

Требуется заголовок `Authorization: Bearer JWT_TOKEN`

* Принимает JSON
 ```json
{
  "amount": "float",
  "currency": "string"
}
```
* Выполняет gRPC-запрос `gw-authorizer.VerifyToken`
* Если авторизация неуспешна, то возвращает `401 Unauthorized`
* Получает баланс пользователя из БД
* Вычисляет изменение баланса с коэффициентом `1`
* Выполняет обновление баланса в БД
* Возвращает `200 Ok` и обновленный баланс
 ```json
{
  "message": "successful",
  "new_balance": {
    "USD": "float",
    "RUB": "float",
    "EUR": "float"
  }
}
```
---
###  Списание средств
`POST /api/v1/wallet/withdraw`

Требуется заголовок `Authorization: Bearer JWT_TOKEN`

* Аналогично пополнению баланса, но с коэффициентом `-1`
* Если средств на балансе недостаточно для списания, то возвращает ошибку `400 BadRequest`
 ---
###  Получение курса обмена валют
`GET /api/v1/exchange/rates`

Требуется заголовок `Authorization: Bearer JWT_TOKEN`

* Выполняет gRPC-запрос `gw-authorizer.VerifyToken`
* Если авторизация неуспешна, то возвращает `401 Unauthorized`
* Выполняет gRPC-запрос `gw-exchanger.GetExchangeRates`
* При успешном выполнении возвращает `200 Ok` и курс валют
```json
{
    "rates": 
    {
      "USD": "float",
      "RUB": "float",
      "EUR": "float"
    }
}
```
---
###  Обмен валют
`POST /api/v1/exchange`

Требуется заголовок `Authorization: Bearer JWT_TOKEN`

* Выполняет gRPC-запрос `gw-authorizer.VerifyToken`
* Если авторизация неуспешна, то возвращает `401 Unauthorized`
* Получает курс валют из кэша. Если запись отсутствует, то выполняет gRPC-запрос `gw-exchanger.GetExchangeRates` и заполняет кэш
* Если средств недостаточно, то возвращает `400 BadRequest`
* Вычисляет изменение баланса по валютам и обновляет запись в БД
* При успешном выполнении возвращает `200 Ok` и обновленный баланс
```json
{
  "message": "Exchange successful",
  "exchanged_amount": "float",
  "new_balance":
  {
   "USD": "float",
   "RUB": "float",
   "EUR": "float"
  }
}
```
---
###  SwaggerUI
`GET swagger/*any`

WEB-интерфейс Swagger

## Конфигурация
Чтение конфигурации происходит из файла, переданного флагом `-c` (по умолчанию - чтение из корня проекта).

Конфигурация подключения к PostgreSQL
*  `PSQL_HOST` - default `localhost`
* `PSQL_PORT` - default `5432`
* `PSQL_DB_NAME` - default `postgres`
* `PSQL_USER` - default `postgres`
* `PSQL_PASSWORD` - default `postgres`
* `PSQL_SSL_MODE` - default `disable`
* `PSQL_CONN_TIMEOUT` - default `60` (в секундах)

Конфигурация web-сервера
* `WEB_HOST` - default `localhost`
* `WEB_PORT` - default `80`

Конфигурация gw-exchanger
* `EXCHANGER_HOST` - default `localhost`
* `EXCHANGER_PORT` - default `9090`

Конфигурация gw-authorizer
* `AUTHORIZER_HOST` - default `localhost`
* `AUTHORIZER_PORT` - default `9090`