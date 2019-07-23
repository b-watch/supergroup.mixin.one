package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

const currency_rate_DDL = `
CREATE TABLE IF NOT EXISTS currency_rates (
    symbol            VARCHAR(36) PRIMARY KEY,
    price_usd         VARCHAR(32) NOT NULL,
    price_cny         VARCHAR(32) NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX currency_rates_pkey ON currency_rates(symbol);

`

type CurrencyRate struct {
	Symbol    string    `json:"symbol"`
	PriceUsd  string    `json:"price_usd"`
	PriceCny  string    `json:"price_cny"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var rateColumns = []string{"symbol", "price_usd", "price_cny", "created_at", "updated_at"}

func UpdateCurrencyRates(ctx context.Context, symbol, priceCny, priceUsd string) error {
	var err error
	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		var rate *CurrencyRate
		row := tx.QueryRowContext(ctx, "SELECT * from currency_rates WHERE symbol=$1", symbol)
		rate, err = rateFromRow(row)
		if err == sql.ErrNoRows {
			_, err := tx.ExecContext(ctx, "INSERT INTO currency_rates (symbol, price_cny, price_usd) VALUES ($1, $2, $3)", symbol, priceCny, priceUsd)
			if err != nil {
				log.Println("UpdateCurrencyRates: Failed to Insert", err)
				return err
			}
		} else if err == nil {
			if priceChangedTooMuch(rate.PriceCny, priceCny) || priceChangedTooMuch(rate.PriceUsd, priceUsd) {
				log.Printf("UpdateCurrencyRates: %s Price change too much\n", symbol)
				return nil
			}
			_, err = tx.ExecContext(ctx, "UPDATE currency_rates SET price_cny=$1, price_usd=$2, updated_at=$3 WHERE symbol=$4", priceCny, priceUsd, time.Now(), symbol)
			if err != nil {
				log.Println("UpdateCurrencyRates: Failed to Update", err)
				return err
			}
		} else {
			return session.TransactionError(ctx, err)
		}
		return nil
	})
	if err != nil {
		log.Println("getCurrency: Failed in TX", err)
		return err
	}
	return nil
}

func GetCurrencyRate(ctx context.Context, symbol string) (*CurrencyRate, error) {
	row := session.Database(ctx).QueryRowContext(ctx, "SELECT * from currency_rates WHERE symbol=$1", symbol)
	rate, err := rateFromRow(row)
	return rate, err
}

func GetCurrencyRates(ctx context.Context) ([]*CurrencyRate, error) {
	query := fmt.Sprintf("SELECT %s FROM currency_rates LIMIT 100", strings.Join(rateColumns, ","))
	rows, err := session.Database(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var rates []*CurrencyRate
	for rows.Next() {
		rate, err := rateFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

func priceChangedTooMuch(priceA, priceB string) bool {
	var realPriceA float64
	var realPriceB float64
	var err error
	realPriceA, err = strconv.ParseFloat(priceA, 64)
	if err != nil {
		return true
	}
	realPriceB, err = strconv.ParseFloat(priceB, 64)
	if err != nil {
		return true
	}
	diff := math.Abs(realPriceA - realPriceB)
	if diff/realPriceA > 0.2 {
		return true
	}
	return false
}

func rateFromRow(row durable.Row) (*CurrencyRate, error) {
	var c CurrencyRate
	err := row.Scan(&c.Symbol, &c.PriceUsd, &c.PriceCny, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}
