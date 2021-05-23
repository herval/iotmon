package coinmarketcap

import (
	"context"
	"fmt"
	http2 "github.com/herval/iotcollector/pkg/http"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://pro-api.coinmarketcap.com"
)

type QuoteData struct {
	Id                int              `json:"id"`
	Name              string           `json:"name"`
	Symbol            string           `json:"symbol"`
	Slug              string           `json:"slug"`
	IsActive          int              `json:"is_active"`
	IsFiat            int              `json:"is_fiat"`
	CirculatingSupply float64          `json:"circulating_supply"`
	TotalSupply       float64          `json:"total_supply"`
	MaxSupply         float64          `json:"max_supply"`
	DateAdded         time.Time        `json:"date_added"`
	NumMarketPairs    int              `json:"num_market_pairs"`
	CmcRank           int              `json:"cmc_rank"`
	LastUpdated       time.Time        `json:"last_updated"`
	Tags              []string         `json:"tags"`
	Platform          interface{}      `json:"platform"`
	Quote             map[string]Price `json:"quote"`
}

type Price struct {
	Price            float64   `json:"price"`
	Volume24H        float64   `json:"volume_24h"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	PercentChange30D float64   `json:"percent_change_30d"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

type LatesQuotesResponse struct {
	Data   map[string]QuoteData `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}

type Client struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

func (c Client) Latest(ctx context.Context, symbols []string) ([]QuoteData, error) {
	res := LatesQuotesResponse{}
	s := strings.Join(symbols, ",")
	err := http2.Get(
		ctx,
		c.client,
		fmt.Sprintf("%s/v1/cryptocurrency/quotes/latest?symbol=%s&CMC_PRO_API_KEY=%s", BaseURL, s, c.apiKey),
		"",
		&res,
		nil,
	)

	if err != nil {
		return nil, err
	}

	//fmt.Println(res.Data)

	data := []QuoteData{}
	for _, v := range res.Data {
		data = append(data, v)
	}

	return data, nil
}

// https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseUrl: BaseURL,
		client:  &http.Client{},
	}
}
