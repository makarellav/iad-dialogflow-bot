package main

import (
	"encoding/json"
	"errors"
	"github.com/makarellav/iad-dialogflow-bot/internal/models"
	"net/http"
	"time"
	_ "time/tzdata"
)

func (b *bot) getCoin(wr *webhookRequest) (*models.Coin, error) {
	id, err := b.readID(wr)

	if err != nil {
		return nil, err
	}

	return b.coin.Get(id)
}

func (b *bot) writeResponse(w http.ResponseWriter, data string) error {
	response := webhookResponse{FulfillmentMessages: []message{
		{
			Text: text{
				Text: []string{data},
			},
		},
	}}

	resp, err := json.Marshal(response)

	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")

	w.Write(resp)

	return nil
}

func (b *bot) readRequest(r *http.Request) (*webhookRequest, error) {
	var request webhookRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (b *bot) readID(wr *webhookRequest) (string, error) {
	id, ok := wr.QueryResult.Parameters["currency"]

	if !ok {
		return "", errors.New("invalid parameters")
	}

	return id, nil
}

func (b *bot) readInterval(wr *webhookRequest) (string, error) {
	interval, ok := wr.QueryResult.Parameters["history"]

	if !ok {
		return "", errors.New("invalid parameters")
	}

	return interval, nil
}

func (b *bot) dateFromTimestamp(ts int64) (string, error) {
	return b.formatDate(time.Unix(ts/1000, 0))
}

func (b *bot) formatDate(t time.Time) (string, error) {
	tz, err := time.LoadLocation("Europe/Kyiv")

	if err != nil {
		return "", err
	}

	return t.In(tz).Format("02.01.2006 15:04:05"), nil
}
