package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Coin struct {
	Id                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
}

type CoinData struct {
	Data Coin `json:"data"`
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

// welcome creates a response for the welcome intent.
func welcome(request webhookRequest) (webhookResponse, error) {
	response := webhookResponse{
		FulfillmentMessages: []message{
			{
				Text: text{
					Text: []string{"Welcome from Dialogflow Go Webhook"},
				},
			},
		},
	}
	return response, nil
}

// getAgentName creates a response for the get-agent-name intent.
func getAgentName(request webhookRequest) (webhookResponse, error) {
	response := webhookResponse{
		FulfillmentMessages: []message{
			{
				Text: text{
					Text: []string{"My name is Dialogflow Go Webhook"},
				},
			},
		},
	}
	return response, nil
}

func getCurrencyPrice(request webhookRequest) (webhookResponse, error) {
	currency, ok := request.QueryResult.Parameters["currency"]
	property, ok := request.QueryResult.Parameters["property"]

	if !ok {
		fmt.Println("not found")
		return webhookResponse{}, nil
	}

	req, err := http.NewRequest("GET", "https://api.coincap.io/v2/assets/"+currency, nil)

	if err != nil {
		fmt.Println(err.Error())
		return webhookResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return webhookResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return webhookResponse{}, err
	}

	fmt.Println(string(body))

	var coin CoinData

	err = json.Unmarshal(body, &coin)

	if err != nil {
		fmt.Println(err.Error())
		return webhookResponse{}, err
	}

	var response webhookResponse

	if property == "price" {
		response = webhookResponse{
			FulfillmentMessages: []message{
				{
					Text: text{
						Text: []string{fmt.Sprintf("%s %s", coin.Data.Name, coin.Data.PriceUsd)},
					},
				},
			},
		}
	}

	return response, nil
}

// handleError handles internal errors.
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ERROR: %v", err)
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var request webhookRequest
	var response webhookResponse
	var err error

	// Read input JSON
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Request: %+v", request)

	fmt.Println("?")

	// Call intent handler
	switch intent := request.QueryResult.Intent.DisplayName; intent {
	case "Default Welcome Intent":
		response, err = welcome(request)
	case "get-agent-name":
		response, err = getAgentName(request)
	case "price":
		response, err = getCurrencyPrice(request)
	default:
		err = fmt.Errorf("unknown intent: %s", intent)
	}
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

type UserSaysData struct {
	Text  string `json:"text"`
	Alias string `json:"alias,omitempty"`
	Meta  string `json:"meta,omitempty"`
}

type UserSays struct {
	Data []UserSaysData `json:"data"`
}

type Intent struct {
	//Id                        string        `json:"id"`
	Name string `json:"name"`
	//Auto                      bool          `json:"auto"`
	//Condition                 string        `json:"condition"`
	//ConditionalFollowupEvents []interface{} `json:"conditionalFollowupEvents"`
	//ConditionalResponses      []interface{} `json:"conditionalResponses"`
	//Context                   []interface{} `json:"context"`
	//Contexts                  []interface{} `json:"contexts"`
	//EndInteraction            bool          `json:"endInteraction"`
	//Events                    []interface{} `json:"events"`
	//FallbackIntent            bool          `json:"fallbackIntent"`
	//LiveAgentHandoff          bool          `json:"liveAgentHandoff"`
	//ParentId                  interface{}   `json:"parentId"`
	//FollowUpIntents           []interface{} `json:"followUpIntents"`
	//Priority                  int           `json:"priority"`
	//Responses                 []struct {
	//	Action           string        `json:"action"`
	//	AffectedContexts []interface{} `json:"affectedContexts"`
	//	Parameters       []struct {
	//		NoInputPromptMessages []interface{} `json:"noInputPromptMessages"`
	//		NoMatchPromptMessages []interface{} `json:"noMatchPromptMessages"`
	//		PromptMessages        []interface{} `json:"promptMessages"`
	//		DefaultValue          string        `json:"defaultValue"`
	//		Name                  string        `json:"name"`
	//		DataType              string        `json:"dataType"`
	//		IsList                bool          `json:"isList"`
	//		Required              bool          `json:"required"`
	//		Prompts               []interface{} `json:"prompts"`
	//		Value                 string        `json:"value"`
	//		OutputDialogContexts  []interface{} `json:"outputDialogContexts"`
	//	} `json:"parameters"`
	//	DefaultResponsePlatforms struct {
	//	} `json:"defaultResponsePlatforms"`
	//	Messages []struct {
	//		Type      string   `json:"type"`
	//		Condition string   `json:"condition"`
	//		Speech    []string `json:"speech"`
	//	} `json:"messages"`
	//	ResetContexts bool `json:"resetContexts"`
	//} `json:"responses"`
	//RootParentId interface{}   `json:"rootParentId"`
	//Templates    []interface{} `json:"templates"`
	UserSays []UserSays `json:"userSays"`
	//WebhookForSlotFilling bool `json:"webhookForSlotFilling"`
	//WebhookUsed           bool `json:"webhookUsed"`
}

//func main() {
//resp, err := http.Get("https://api.coincap.io/v2/assets?limit=2000")
//
//if err != nil {
//	fmt.Println("failed to load data", err.Error())
//
//	return
//}
//
//data, err := io.ReadAll(resp.Body)
//
//if err != nil {
//	fmt.Println("failed to read data", err.Error())
//
//	return
//}
//
//fmt.Println(string(data))
//
//var coinsData struct {
//	Data []Coin `json:"data"`
//}
//
//err = json.Unmarshal(data, &coinsData)
//
//if err != nil {
//	fmt.Println("failed unmarshalling", err.Error())
//
//	return
//}
//
////headers := []string{"name", "synonyms"}
//
//f, err := os.Create("data.csv")
//
//defer f.Close()
//
//if err != nil {
//	fmt.Println(err.Error())
//
//	return
//}
//
//w := csv.NewWriter(f)
//defer w.Flush()
//
////w.Write(headers)
//
//for _, coin := range coinsData.Data {
//	w.Write([]string{strconv.Quote(*coin.Id), strconv.Quote(*coin.Id), strconv.Quote(*coin.Name), strconv.Quote(*coin.Symbol)})
//}

//{
//	"isTemplate": false,
//	"data": [
//{
//"text": "bitcoin",
//"userDefined": false,
//"alias": "currency",
//"meta": "@currency"
//},
//{
//"text": " цена",
//"userDefined": false
//}
//],
//"count": 0,
//"id": "03c9cc35-8ed2-4833-9f96-53b56c5bce5e",
//"updated": null
//},

//var currencyPriceIntent Intent
//currencyPriceIntent.Name = "currency.price"
//
//priceValues := []string{"ціна", "цена"}
//
//f, err := os.Open("data.csv")
//defer f.Close()
//
//if err != nil {
//	fmt.Println(err.Error())
//
//	return
//}
//
//r := csv.NewReader(f)
//
//coins, err := r.ReadAll()
//
//if err != nil {
//	fmt.Println(err.Error())
//
//	return
//}
//
//for _, coin := range coins[:len(coins)/2+1] {
//	for _, priceValue := range priceValues {
//		var userSaysData UserSays
//
//		userSaysCoin := UserSaysData{
//			Text:  coin[0],
//			Alias: "currency",
//			Meta:  "@currency",
//		}
//
//		userSaysPrice := UserSaysData{
//			Text: " " + priceValue,
//		}
//
//		userSaysData.Data = append(userSaysData.Data, userSaysCoin, userSaysPrice)
//		currencyPriceIntent.UserSays = append(currencyPriceIntent.UserSays, userSaysData)
//	}
//}
//
//data, err := json.Marshal(currencyPriceIntent)
//
//if err != nil {
//	fmt.Println(err.Error())
//
//	return
//}
//
//fmt.Println(len(currencyPriceIntent.UserSays))
//
//log.Fatal(os.WriteFile("price.json", data, 0644))
//}
