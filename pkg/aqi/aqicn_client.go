package aqi

import (
	"context"
	"fmt"
	http2 "github.com/herval/iotcollector/pkg/http"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.waqi.info"
)

type AqiCnClient struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

type AqiDataResult struct {
	Status string `json:"status"`
	Data   struct {
		Aqi          float64 `json:"aqi"`
		Idx          float64 `json:"idx"`
		Attributions []struct {
			Url  string `json:"url"`
			Name string `json:"name"`
		} `json:"attributions"`
		City struct {
			Geo  []float64 `json:"geo"`
			Name string    `json:"name"`
			Url  string    `json:"url"`
		} `json:"city"`
		Dominentpol string `json:"dominentpol"`
		Iaqi        struct {
			Co struct {
				V float64 `json:"v"`
			} `json:"co"`
			Humidity struct {
				V float64 `json:"v"`
			} `json:"h"`
			No2 struct {
				V float64 `json:"v"`
			} `json:"no2"`
			Ozone struct {
				V float64 `json:"v"`
			} `json:"o3"`
			Pressure struct {
				V float64 `json:"v"`
			} `json:"p"`
			Pm10 struct {
				V float64 `json:"v"`
			} `json:"pm10"`
			Pm25 struct {
				V float64 `json:"v"`
			} `json:"pm25"`
			So2 struct {
				V float64 `json:"v"`
			} `json:"so2"`
			Temperature struct {
				V float64 `json:"v"`
			} `json:"t"`
			Wind struct {
				V float64 `json:"v"`
			} `json:"w"`
		} `json:"iaqi"`
		//Time struct {
		//	Local    string    `json:"s"`
		//	Timezone string    `json:"tz"`
		//	V        float64       `json:"v"`
		//	Iso      time.Time `json:"iso"`
		//} `json:"time"`
		//Forecast struct {
		//	Daily struct {
		//		Ozone []struct {
		//			Avg float64    `json:"avg"`
		//			Day string `json:"day"`
		//			Max float64    `json:"max"`
		//			Min float64    `json:"min"`
		//		} `json:"o3"`
		//		Pm10 []struct {
		//			Avg float64    `json:"avg"`
		//			Day string `json:"day"`
		//			Max float64    `json:"max"`
		//			Min float64    `json:"min"`
		//		} `json:"pm10"`
		//		Pm25 []struct {
		//			Avg float64    `json:"avg"`
		//			Day string `json:"day"`
		//			Max float64    `json:"max"`
		//			Min float64    `json:"min"`
		//		} `json:"pm25"`
		//		UVIndex []struct {
		//			Avg float64    `json:"avg"`
		//			Day string `json:"day"`
		//			Max float64    `json:"max"`
		//			Min float64    `json:"min"`
		//		} `json:"uvi"`
		//	} `json:"daily"`
		//} `json:"forecast"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}

func NewClient(apiKey string) *AqiCnClient {
	return &AqiCnClient{
		apiKey:  apiKey,
		baseUrl: BaseURL,
		client:  &http.Client{},
	}
}

func (c *AqiCnClient) GetAtLatLng(ctx context.Context, lat string, lng string) (*AqiDataResult, error) {
	///feed/geo::lat;:lng/?token=:token
	var d AqiDataResult

	if err := http2.Get(
		ctx,
		c.client,
		fmt.Sprintf("%s/feed/geo:%s;%s/?token=%s", c.baseUrl, lat, lng, c.apiKey),
		"",
		&d,
		nil,
	); err != nil {
		return nil, err
	}

	return &d, nil
}

func (c *AqiCnClient) GetHere(ctx context.Context) (*AqiDataResult, error) {
	var d AqiDataResult

	if err := http2.Get(
		ctx,
		c.client,
		fmt.Sprintf("%s/feed/here/?token=%s", c.baseUrl, c.apiKey),
		"",
		&d,
		nil,
	); err != nil {
		return nil, err
	}

	return &d, nil
}
