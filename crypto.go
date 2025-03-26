package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	var result Result
	coinmarketcapIds := []string{}
	coinQty := make(map[string]float64)
	csvPortfolio := readCsvFile("portfolio.csv")
	coinmarketcapQuoteUrl := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
	apiKey := "<your-apikey>"
	var portfolio Portfolio

	for coins := range csvPortfolio {
		coinmarketcapIds = append(coinmarketcapIds, csvPortfolio[coins][0])
		i, err := strconv.ParseFloat(strings.TrimSpace(csvPortfolio[coins][1]), 64)
		if err != nil {
			printError(&result, err.Error())
			os.Exit(0)
		}
		coinQty[csvPortfolio[coins][0]] = i
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", coinmarketcapQuoteUrl, nil)
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	q := url.Values{}
	q.Add("id", strings.Join(coinmarketcapIds, ","))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println()
		printError(&result, "Error sending request to server")
		os.Exit(0)
	}

	respBody, _ := io.ReadAll(resp.Body)

	var quotes Quotes
	if err := json.Unmarshal(respBody, &quotes); err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

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
		result.Success = true
	}

	marshaledResult, _ := json.MarshalIndent(result, "", "    ")
	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func printError(result *Result, error string) {
	result.Success = false
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
	Success   bool      `json:"success,omitempty"`
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
