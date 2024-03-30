package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CoinData struct {
	ID                string  `json:"id"`
	Rank              string  `json:"rank"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Supply            string  `json:"supply"`
	MaxSupply         *string `json:"maxSupply"`
	MarketCapUsd      string  `json:"marketCapUsd"`
	VolumeUsd24Hr     string  `json:"volumeUsd24Hr"`
	PriceUsd          string  `json:"priceUsd"`
	ChangePercent24Hr string  `json:"changePercent24Hr"`
	Vwap24Hr          string  `json:"vwap24Hr"`
}

type Coin struct {
	Timestamp int64    `json:"timestamp"`
	Data      CoinData `json:"data"`
}

type HistoryData struct {
	PriceUsd string    `json:"priceUsd"`
	Time     int64     `json:"time"`
	Date     time.Time `json:"date"`
}

type History struct {
	Data      []HistoryData `json:"data"`
	Timestamp int64         `json:"timestamp"`
}

type CoinModel struct {
	BaseURL string
}

func (cm *CoinModel) Get(id string) (*Coin, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", cm.BaseURL, id))

	if err != nil {
		return nil, err
	}

	var coin Coin

	err = json.NewDecoder(resp.Body).Decode(&coin)

	if err != nil {
		return nil, err
	}

	return &coin, nil
}

func (cm *CoinModel) History(id, interval string) (*HistoryData, *HistoryData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/history?interval=%s", cm.BaseURL, id, interval))

	if err != nil {
		return nil, nil, err
	}

	var history History

	err = json.NewDecoder(resp.Body).Decode(&history)

	if err != nil {
		return nil, nil, err
	}

	current := history.Data[len(history.Data)-1]
	prev := history.Data[len(history.Data)-2]

	return &current, &prev, nil
}
