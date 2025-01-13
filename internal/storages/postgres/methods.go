package postgres

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/storages"
)

func (p *PSQL) GetBalance(ctx context.Context, user string) (storages.Balance, error) {
	const op = "PSQL GetBalance"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var result storages.Balance
	err := p.pool.QueryRow(ctxWithTimeout, "select cash from wallets where user_id = $1", user).Scan(&result)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return result, err
}

func (p *PSQL) NewWallet(ctx context.Context, id string) error {
	const op = "PSQL NewWallet"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var newWallet storages.Balance
	_, err := p.pool.Exec(ctxWithTimeout, "insert into wallets (user_id, cash) values($1,$2)", id, newWallet)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (p *PSQL) UpdateWallet(ctx context.Context, user string, balance storages.Balance) error {
	const op = "PSQL UpdateWallet"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctxWithTimeout, "update wallets set cash = $1 where user_id = $2", balance, user)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return err
}
