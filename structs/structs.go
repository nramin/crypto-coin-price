package structs

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

type CryptoPrices struct {
	Positions  []Position `json:"positions"`
	TotalValue float64    `json:"totalValue,omitempty" bson:"totalValue,omitempty"`
	Success    *bool      `json:"success,omitempty" bson:"success,omitempty"`
	Error      string     `json:"error,omitempty" bson:"error,omitempty"`
}

type Position struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Value    float64 `json:"value"`
}

type Yml struct {
	ApiKey    string             `yaml:"apiKey"`
	Positions map[string]float64 `yaml:"positions"`
}
