package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	var result Result
	coinmarketcapIds := []string{}
	coinQty := make(map[string]float64)
	yamlConfig := readYamlFile("crypto.yaml")
	var portfolio Portfolio

	for coinId, quantity := range yamlConfig.Positions {
		coinmarketcapIds = append(coinmarketcapIds, coinId)
		coinQty[coinId] = quantity
	}

	quotes := callCoinMarketCap(coinmarketcapIds, yamlConfig, &result)

	portfolio.Position = make(map[string]Position)
	for v, k := range quotes.Data {
		var position Position
		position.Value = k.Quote.USD.Price * coinQty[v]
		position.Quantity = coinQty[v]

		portfolio.Position[k.Symbol] = position
		portfolio.TotalValue += (k.Quote.USD.Price * coinQty[v])
	}

	if len(result.Error) == 0 {
		result.Portfolio = portfolio
		success := new(bool)
		*success = true
		result.Success = success
	}

	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func callCoinMarketCap(coinmarketcapIds []string, userData Yml, result *Result) Quotes {
	coinmarketcapQuoteUrl := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
	apiKey := userData.ApiKey

	client := &http.Client{}
	req, err := http.NewRequest("GET", coinmarketcapQuoteUrl, nil)
	if err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	q := url.Values{}
	q.Add("id", strings.Join(coinmarketcapIds, ","))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		printError(result, "Error sending request to server")
		os.Exit(0)
	}

	respBody, _ := io.ReadAll(resp.Body)

	var quotes Quotes
	if err := json.Unmarshal(respBody, &quotes); err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	return quotes
}

func readYamlFile(filePath string) Yml {
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	var yamlFile Yml

	err = yaml.Unmarshal([]byte(b), &yamlFile)
	if err != nil {
		panic(err)
	}

	return yamlFile
}

func printError(result *Result, error string) {
	success := new(bool)
	*success = false

	result.Success = success
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}

type Quotes struct {
	Data map[string]Coin `json:"data"`
}

type Coin struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	Quote  struct {
		USD struct {
			Price float64 `json:"price"`
		} `json:"USD"`
	} `json:"quote"`
}

type Result struct {
	Portfolio Portfolio `json:"portfolio,omitempty"`
	Success   *bool     `json:"success,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type Portfolio struct {
	Position   map[string]Position `json:"position"`
	TotalValue float64             `json:"totalValue"`
}

type Position struct {
	Quantity float64 `json:"quantity"`
	Value    float64 `json:"value"`
}

type Yml struct {
	ApiKey    string             `yaml:"apiKey"`
	Positions map[string]float64 `yaml:"positions"`
}
