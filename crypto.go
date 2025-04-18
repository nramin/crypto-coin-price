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

	"github.com/nramin/crypto-coin-price/structs"
	"gopkg.in/yaml.v3"
)

func main() {
	var result structs.CryptoPrices
	coinmarketcapIds := []string{}
	coinQty := make(map[string]float64)
	yamlConfig := readYamlFile("crypto.yaml")

	for coinId, quantity := range yamlConfig.Positions {
		coinmarketcapIds = append(coinmarketcapIds, coinId)
		coinQty[coinId] = quantity
	}

	quotes := callCoinMarketCap(coinmarketcapIds, yamlConfig, &result)

	result.Position = []structs.Position{}
	for v, k := range quotes.Data {
		position := structs.Position{
			Symbol:   k.Symbol,
			Value:    k.Quote.USD.Price * coinQty[v],
			Quantity: coinQty[v],
		}
		result.Position = append(result.Position, position)
		result.TotalValue += (k.Quote.USD.Price * coinQty[v])
	}

	if len(result.Error) == 0 {
		success := new(bool)
		*success = true
		result.Success = success
	}

	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func callCoinMarketCap(coinmarketcapIds []string, userData structs.Yml, result *structs.CryptoPrices) structs.Quotes {
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

	var quotes structs.Quotes
	if err := json.Unmarshal(respBody, &quotes); err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	return quotes
}

func readYamlFile(filePath string) structs.Yml {
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	var yamlFile structs.Yml

	err = yaml.Unmarshal([]byte(b), &yamlFile)
	if err != nil {
		panic(err)
	}

	return yamlFile
}

func printError(result *structs.CryptoPrices, error string) {
	success := new(bool)
	*success = false

	result.Success = success
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}
