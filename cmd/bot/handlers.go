package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (b *bot) handlers() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", b.webhookRequestHandler)

	return mux
}

func (b *bot) getCoinPriceHandler(w http.ResponseWriter, wr *webhookRequest) {
	coin, err := b.getCoin(wr)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	price, err := strconv.ParseFloat(coin.Data.PriceUsd, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	date, err := b.dateFromTimestamp(coin.Timestamp)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	resp := fmt.Sprintf("Ціна %s на %s\n%.4f USD", coin.Data.Name, date, price)
	err = b.writeResponse(w, resp)

	if err != nil {
		b.errorResponse(w, err)
	}
}

func (b *bot) getCoinInfoHandler(w http.ResponseWriter, wr *webhookRequest) {
	coin, err := b.getCoin(wr)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Ось, що я знайшов про %s\n\n", coin.Data.Name))
	sb.WriteString(fmt.Sprintf("Місце в рейтингу: %s\n", coin.Data.Rank))
	sb.WriteString(fmt.Sprintf("Cимвол: %s\n", coin.Data.Symbol))

	supply, err := strconv.ParseFloat(coin.Data.Supply, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	sb.WriteString(fmt.Sprintf("Загальна пропозиція: %.4f\n", supply))

	if coin.Data.MaxSupply != nil {
		maxSupply, err := strconv.ParseFloat(*coin.Data.MaxSupply, 64)

		if err != nil {
			b.errorResponse(w, err)

			return
		}

		sb.WriteString(fmt.Sprintf("Максимальна кількість монет: %.4f\n", maxSupply))
	} else {
		sb.WriteString("Максимальна кількість монет: ∞\n")
	}

	marketCap, err := strconv.ParseFloat(coin.Data.MarketCapUsd, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	sb.WriteString(fmt.Sprintf("Ринкова капіталізація: %.4f USD\n", marketCap))

	volume, err := strconv.ParseFloat(coin.Data.VolumeUsd24Hr, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	sb.WriteString(fmt.Sprintf("Обсяг (24г): %.4f USD\n", volume))

	coinPrice, err := strconv.ParseFloat(coin.Data.PriceUsd, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	sb.WriteString(fmt.Sprintf("Ціна: %.4f USD\n", coinPrice))

	change, err := strconv.ParseFloat(coin.Data.ChangePercent24Hr, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	sb.WriteString(fmt.Sprintf("Зміна ціни (24г): %.2f%%\n", change))

	err = b.writeResponse(w, sb.String())

	if err != nil {
		b.errorResponse(w, err)
	}
}

func (b *bot) getCoinHistoryHandler(w http.ResponseWriter, wr *webhookRequest) {
	id, err := b.readID(wr)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	interval, err := b.readInterval(wr)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	current, prev, err := b.coin.History(id, interval)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	currentPrice, err := strconv.ParseFloat(current.PriceUsd, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	prevPrice, err := strconv.ParseFloat(prev.PriceUsd, 64)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	currentDate, err := b.formatDate(current.Date)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	prevDate, err := b.formatDate(prev.Date)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	changePercent := (currentPrice - prevPrice) / currentPrice * 100

	resp := fmt.Sprintf("Ціна %s на %s: %.4f\nЦіна %s на %s: %.4f\nЗміна ціни: %.2f%%", id, prevDate, prevPrice, id, currentDate, currentPrice, changePercent)

	err = b.writeResponse(w, resp)

	if err != nil {
		b.errorResponse(w, err)
	}
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func (b *bot) webhookRequestHandler(w http.ResponseWriter, r *http.Request) {
	wr, err := b.readRequest(r)

	if err != nil {
		b.errorResponse(w, err)

		return
	}

	switch wr.QueryResult.Intent.DisplayName {
	case "price":
		b.getCoinPriceHandler(w, wr)
	case "info":
		b.getCoinInfoHandler(w, wr)
	case "history":
		b.getCoinHistoryHandler(w, wr)
	default:
		b.errorResponse(w, errors.New("unknown intent"))
	}
}
