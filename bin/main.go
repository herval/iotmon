package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	aqi2 "github.com/herval/iotcollector/pkg/aqi"
	"github.com/herval/iotcollector/pkg/awair"
	"github.com/herval/iotcollector/pkg/coinmarketcap"
	"github.com/herval/iotcollector/pkg/prom"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"time"
)

func main() {

	awairUpdates := make(chan *awair.RawDataPoints)
	aqiUpdates := make(chan *aqi2.AqiDataResult)
	cryptoUpdates := make(chan *coinmarketcap.QuoteData)

	if key := os.Getenv("COINMARKETCAP_API_KEY"); key != "" {
		coins := strings.Split(os.Getenv("COINMARKETCAP_COINS"), ",")
		cli := coinmarketcap.NewClient(key)

		go cryptoPoll(cli, coins, cryptoUpdates)
	}

	if aqi := os.Getenv("AQICN_API_KEY"); aqi != "" {
		aqiClient := aqi2.NewClient(aqi)
		go aqiPoll(
			aqiClient,
			os.Getenv("AQICN_LOCATIONS"),
			aqiUpdates,
		)
	}

	if host := os.Getenv("AWAIR_LOCAL_HOST"); host != "" {
		cli := awair.NewLocalClient(host)
		go awairPollLocal(
			awair.NeWPollingMonitor(
				context.Background(),
				cli,
			),
			awairUpdates,
		)
	}

	if token := os.Getenv("AWAIR_TOKEN"); token != "" {
		go awairPoll(
			awair.NeWPollingMonitor(
				context.Background(),
				awair.NewClient(token),
			),
			awairUpdates,
		)
	}

	if port := os.Getenv("PORT"); port != "" {
		go startWebserver(port)
	}

	if promUrl := os.Getenv("PROMETHEUS_URL"); promUrl != "" {
		go promPush(awairUpdates, aqiUpdates, cryptoUpdates, promUrl)
	}

	select {}
}

func cryptoPoll(client *coinmarketcap.Client, symbols []string, upd chan *coinmarketcap.QuoteData) {
	for {
		data, err := client.Latest(context.Background(), symbols)
		if err != nil {
			fmt.Println(err.Error())
		}

		//fmt.Println(symbols)
		for _, d := range data {
			upd <- &d
		}

		time.Sleep(time.Hour)
	}
}

func aqiPoll(client *aqi2.AqiCnClient, locations string, upd chan *aqi2.AqiDataResult) {
	latLngRaw := strings.Split(locations, ";")

	latLng := [][]string{}
	for _, ll := range latLngRaw {
		latLng = append(latLng, strings.Split(ll, ","))
	}

	for _, ll := range latLng {
		go func(coords []string) {
			for {
				data, err := client.GetAtLatLng(context.Background(), coords[0], coords[1])
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(data)
				upd <- data

				time.Sleep(time.Hour)
			}
		}(ll)
	}
}

func awairPollLocal(awairMonitor *awair.Monitor, upd awair.Updates) {
	if err := awairMonitor.Start(upd, time.Minute); err != nil {
		panic(err)
	}
}

func startWebserver(port string) {
	r := gin.Default()
	//r.LoadHTMLGlob("./resources/templates/*")

	r.GET("/", func(c *gin.Context) {
		c.Status(200)
	})
	err := r.Run(port)
	if err != nil {
		fmt.Println("Error starting http server - %s", err.Error())
	}
}

func awairPoll(awairMonitor *awair.Monitor, upd awair.Updates) {
	if err := awairMonitor.Start(upd, time.Minute*10); err != nil {
		panic(err)
	}
}

func promPush(awairUpdates awair.Updates, aqiUpdates chan *aqi2.AqiDataResult, cryptoUpdates chan *coinmarketcap.QuoteData, gatewayUrl string) {
	pushers := PusherBuffer{
		pushers:    map[string]prom.Pusher{},
		gatewayUrl: gatewayUrl,
	}

	for {
		select {
		case c := <-cryptoUpdates:
			fmt.Println("Processing coinz....")
			pusher := pushers.ForCrypto(c.Symbol)
			pusher.Update("cmc_rank", float64(c.CmcRank))
			pusher.Update("price_usd", c.Quote["USD"].Price)
			pusher.Update("market_cap_usd", c.Quote["USD"].MarketCap)

			err := pusher.Push()
			if err != nil {
				fmt.Println(err.Error())
				// TODO
			} else {
				fmt.Println("Pushed")
			}

		case u := <-aqiUpdates:
			fmt.Println("Processing aqi metrics...")
			pusher := pushers.ForAqi(u)
			pusher.Update("aqi", u.Data.Aqi)
			pusher.Update("co", u.Data.Iaqi.Co.V)
			pusher.Update("humidity", u.Data.Iaqi.Humidity.V)
			pusher.Update("no2", u.Data.Iaqi.No2.V)
			pusher.Update("ozone", u.Data.Iaqi.Ozone.V)
			pusher.Update("pressure", u.Data.Iaqi.Pressure.V)
			pusher.Update("pm10", u.Data.Iaqi.Pm10.V)
			pusher.Update("pm25", u.Data.Iaqi.Pm25.V)
			pusher.Update("so2", u.Data.Iaqi.So2.V)
			pusher.Update("temperature", u.Data.Iaqi.Temperature.V)
			pusher.Update("wind", u.Data.Iaqi.Wind.V)

			err := pusher.Push()
			if err != nil {
				fmt.Println(err.Error())
				// TODO
			} else {
				fmt.Println("Pushed")
			}

		case d := <-awairUpdates:
			fmt.Println("Processing awair metrics...")
			pusher := pushers.For(d.DeviceId)
			pusher.Update("score", d.Score)
			for _, s := range d.Sensors {
				kind := strings.ToLower(s.Comp)
				if !pusher.Update(kind, s.Value) {
					fmt.Println("Skipping '" + kind + "' measurement")
				}
			}

			//fmt.Println(fmt.Sprintf("%+v", d))
			err := pusher.Push()
			if err != nil {
				fmt.Println(err.Error())
				// TODO
			} else {
				fmt.Println("Pushed")
			}
		}
	}
}

type PusherBuffer struct {
	pushers    map[string]prom.Pusher
	gatewayUrl string
}

func (b *PusherBuffer) For(deviceId int) *prom.Pusher {
	k := fmt.Sprintf("awair_%d", deviceId)
	if buff, found := b.pushers[k]; found {
		return &buff
	}

	gaugeNames := []string{
		"score", "voc", "co2", "pm25", "humid", "temp",
	}

	p := prom.NewPusher(
		b.gatewayUrl,
		"awair",
		map[string]string{
			"sensor_id": fmt.Sprintf("%d", deviceId),
		},
		gaugeNames,
		"awair_sensor")

	b.pushers[k] = *p
	return p
}

func (b *PusherBuffer) ForAqi(u *aqi2.AqiDataResult) *prom.Pusher {
	k := fmt.Sprintf("aqi_%d", u.Data.City.Name)
	if buff, found := b.pushers[k]; found {
		return &buff
	}

	gaugeNames := []string{
		"aqi", "co", "humidity", "no2", "ozone", "pressure", "pm10", "pm25", "wind", "temperature", "so2",
	}

	p := prom.NewPusher(
		b.gatewayUrl,
		"aqicn",
		map[string]string{
			"city": strings.ReplaceAll(u.Data.City.Name, "+", " "),
			"lat":  fmt.Sprintf("%f", u.Data.City.Geo[0]),
			"lng":  fmt.Sprintf("%f", u.Data.City.Geo[1]),
		},
		gaugeNames,
		"aqicn")

	b.pushers[k] = *p
	return p
}

func (b *PusherBuffer) ForCrypto(symbol string) *prom.Pusher {
	k := fmt.Sprintf("crypto_%d", symbol)
	if buff, found := b.pushers[k]; found {
		return &buff
	}

	gaugeNames := []string{
		"price_usd", "market_cap_usd", "cmc_rank",
	}

	p := prom.NewPusher(
		b.gatewayUrl,
		"crypto",
		map[string]string{
			"coin": symbol,
		},
		gaugeNames,
		"coinmarketcap")

	b.pushers[k] = *p
	return p
}
