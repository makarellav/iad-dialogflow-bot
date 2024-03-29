package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Coin struct {
	Id                string  `json:"id"`
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

type CoinData struct {
	Timestamp int64 `json:"timestamp"`
	Data      Coin  `json:"data"`
}

type intent struct {
	DisplayName string `json:"displayName"`
}

type queryResult struct {
	Intent     intent            `json:"intent"`
	Action     string            `json:"action,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type text struct {
	Text []string `json:"text"`
}

type message struct {
	Text text `json:"text"`
}

// webhookRequest is used to unmarshal a WebhookRequest JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by
// the Dialogflow protocol buffers:
// https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2#WebhookRequest
type webhookRequest struct {
	Session     string      `json:"session"`
	ResponseID  string      `json:"responseId"`
	QueryResult queryResult `json:"queryResult"`
}

// webhookResponse is used to marshal a WebhookResponse JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by
// the Dialogflow protocol buffers:
// https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2#WebhookResponse
type webhookResponse struct {
	FulfillmentMessages []message `json:"fulfillmentMessages"`
}

func getPriceResponse(coin CoinData) (webhookResponse, error) {
	coinPrice, err := strconv.ParseFloat(coin.Data.PriceUsd, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	tz, err := time.LoadLocation("Europe/Kyiv")

	if err != nil {
		return webhookResponse{}, err
	}

	timestamp := time.Unix(coin.Timestamp/1000, 0).In(tz).Format(time.TimeOnly)

	return webhookResponse{
		FulfillmentMessages: []message{
			{
				Text: text{
					Text: []string{fmt.Sprintf("Ціна %s на %s\n%.4f USD", coin.Data.Name, timestamp, coinPrice)},
				},
			},
		},
	}, nil
}

func getInfoResponse(coin CoinData) (webhookResponse, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Ось, що я знайшов про %s\n\n", coin.Data.Name))
	sb.WriteString(fmt.Sprintf("Місце в рейтингу: %s\n", coin.Data.Rank))
	sb.WriteString(fmt.Sprintf("Cимвол: %s\n", coin.Data.Symbol))

	supply, err := strconv.ParseFloat(coin.Data.Supply, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	sb.WriteString(fmt.Sprintf("Загальна пропозиція: %.4f\n", supply))

	if coin.Data.MaxSupply != nil {
		maxSupply, err := strconv.ParseFloat(*coin.Data.MaxSupply, 64)

		if err != nil {
			return webhookResponse{}, err
		}

		sb.WriteString(fmt.Sprintf("Максимальна кількість монет: %.4f\n", maxSupply))
	} else {
		sb.WriteString("Максимальна кількість монет: невідомо\n")
	}

	marketCap, err := strconv.ParseFloat(coin.Data.MarketCapUsd, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	sb.WriteString(fmt.Sprintf("Ринкова капіталізація: %.4f USD\n", marketCap))

	volume, err := strconv.ParseFloat(coin.Data.VolumeUsd24Hr, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	sb.WriteString(fmt.Sprintf("Обсяг (24г): %.4f USD\n", volume))

	coinPrice, err := strconv.ParseFloat(coin.Data.PriceUsd, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	sb.WriteString(fmt.Sprintf("Ціна: %.4f USD\n", coinPrice))

	change, err := strconv.ParseFloat(coin.Data.ChangePercent24Hr, 64)

	if err != nil {
		return webhookResponse{}, err
	}

	sb.WriteString(fmt.Sprintf("Зміна ціни (24г): %.2f%%\n", change))

	return webhookResponse{
		FulfillmentMessages: []message{
			{
				Text: text{
					Text: []string{sb.String()},
				},
			},
		},
	}, nil
}

var intentResponses = map[string]func(coin CoinData) (webhookResponse, error){
	"price": getPriceResponse,
	"info":  getInfoResponse,
}

// handleError handles internal errors.
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Println(err.Error())
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var request webhookRequest
	var response webhookResponse

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		fmt.Println("here", err.Error())

		handleError(w, err)

		return
	}

	log.Printf("Request: %+v", request)

	currency, ok := request.QueryResult.Parameters["currency"]
	property, ok := request.QueryResult.Parameters["property"]

	resp, err := http.Get("https://api.coincap.io/v2/assets/" + currency)

	var coin CoinData

	err = json.NewDecoder(resp.Body).Decode(&coin)

	if err != nil {
		handleError(w, err)

		return
	}

	getResponse, ok := intentResponses[property]

	if !ok {
		fmt.Println("nope")

		return
	}

	response, err = getResponse(coin)

	if err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Response: %+v", response)

	// Send response
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		handleError(w, err)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", HandleWebhookRequest)

	log.Fatal(http.ListenAndServe(":7777", mux))
}
