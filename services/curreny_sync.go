package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

const currencyUrl = "https://gbi-api.fox.one/currency"

type CnyTicker struct {
	Timestamp   int64  `json:"timestamp"`
	FromSymbol  string `json:"from"`
	ToSymbol    string `json:"to"`
	Price       string `json:"price"`
	ChangeIn24h string `json:"changeIn24h"`
}

type CurrencyMap struct {
	UsdPrice  string `json:"usd"`
	UsdtPrice string `json:"usdt"`
}

type CurrencyResponseData struct {
	CnyTickers []CnyTicker `json:"cnyTickers"`
	Currencies CurrencyMap `json:"currencies"`
}

type CurrencyResponse struct {
	Code int                  `json:"code"`
	Data CurrencyResponseData `json:"data"`
}

func getCurrency(ctx context.Context) error {
	resp, err := http.Get(currencyUrl)
	if err != nil {
		log.Println("getCurrency Error: get currencyUrl", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("getCurrency Error: read body")
		return err
	}

	var currencyResp CurrencyResponse

	err = json.Unmarshal(body, &currencyResp)
	cnyUsdRate, _ := strconv.ParseFloat(currencyResp.Data.Currencies.UsdPrice, 64)
	for _, ticker := range currencyResp.Data.CnyTickers {
		realPrice, _ := strconv.ParseFloat(ticker.Price, 64)
		realPriceUsd := fmt.Sprintf("%.8f", realPrice/cnyUsdRate)
		err = models.UpdateCurrencyRates(ctx, ticker.FromSymbol, ticker.Price, realPriceUsd)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartCurrencySync(name string, db *durable.Database) {
	context := session.WithDatabase(context.Background(), db)
	ctx := session.WithLogger(context, durable.BuildLogger())
	for true {
		getCurrency(ctx)
		time.Sleep(time.Duration(300) * time.Second) // update per 5 mins
	}
}
