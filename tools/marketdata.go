package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const(
	DEFAULT_MARKET = "kraken"
	DEFAULT_FIAT = "cad"
)

type MarketDataClient struct {
	Market string
	Currency string
}

type candleResponse struct {
	Result candleResult `json:"result"`
}

type candleResult struct {
	Minutes [][]float64 `json:"60"`
}

type candle struct {
	CloseTime time.Time
	OpenPrice float64
	HighPrice float64
	LowPrice float64
	ClosePrice float64
	Volume float64
	QuoteVolume float64
}

func (client MarketDataClient) GetValueAtTime(timeInstance time.Time, coin string) (float64, error) {
	market := DEFAULT_MARKET
	currency := DEFAULT_FIAT
	if client.Market != "" {
		market = client.Market
	}
	if client.Currency != "" {
		currency = client.Currency
	}

	candle, err := getCandle(timeInstance, market, strings.ToLower(coin + currency), 60)
	if err != nil {
		return 0., err
	}
	return candle.ClosePrice, nil
}

func getCandle(before time.Time, market, pair string, periods int) (*candle, error) {
	after := before.Add(-time.Minute)
	url := fmt.Sprintf("https://api.cryptowat.ch/markets/%s/%s/ohlc?after=%d&before=%d&periods=%d", market, pair, after.Unix(), before.Unix(), periods)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var candleResponse candleResponse
	if err := json.NewDecoder(response.Body).Decode(&candleResponse); err != nil {
		return nil, err
	}

	if len(candleResponse.Result.Minutes) == 0 {
		return nil, errors.New("no candles found")
	}
	candleData := candleResponse.Result.Minutes[0]
	return &candle{
		CloseTime: time.Unix(int64(candleData[0]), 0),
		OpenPrice: candleData[1],
		HighPrice: candleData[2],
		LowPrice: candleData[3],
		ClosePrice: candleData[4],
		Volume: candleData[5],
		QuoteVolume: candleData[6],
	}, nil
}
