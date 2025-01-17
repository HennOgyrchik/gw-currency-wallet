FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gw-currency-wallet ./cmd/

FROM scratch
WORKDIR /app
COPY --from=builder ./app/gw-currency-wallet .
COPY --from=builder ./app/config.env .
COPY --from=builder ./app/internal/storages/migrations ./migrations
EXPOSE 80
CMD ["./gw-currency-wallet"]

#docker build -t gw-currency-wallet .
#docker run --name gw-currency-wallet -p 8080:80 -d gw-currency-wallet