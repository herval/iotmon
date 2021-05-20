package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herval/iotcollector/pkg/awair"
	"github.com/herval/iotcollector/pkg/prom"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"time"
)

func main() {

	updates := make(chan *awair.RawDataPoints)

	if host := os.Getenv("AWAIR_LOCAL_HOST"); host != "" {
		cli := awair.NewLocalClient(host)
		go awairPollLocal(
			awair.NeWPollingMonitor(
				context.Background(),
				cli,
			),
			updates,
		)
	}

	if token := os.Getenv("AWAIR_TOKEN"); token != "" {
		go awairPoll(
			awair.NeWPollingMonitor(
				context.Background(),
				awair.NewClient(token),
			),
			updates,
		)
	}

	if port := os.Getenv("PORT"); port != "" {
		go startWebserver(port)
	}

	if promUrl := os.Getenv("PROMETHEUS_URL"); promUrl != "" {
		go promPush(updates, promUrl)
	}

	select {}
}

func awairPollLocal(awairMonitor *awair.Monitor, upd awair.Updates) {
	if err := awairMonitor.Start(upd, time.Second*30); err != nil {
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

func promPush(upd awair.Updates, gatewayUrl string) {
	pushers := PusherBuffer{
		pushers:    map[int]prom.Pusher{},
		gatewayUrl: gatewayUrl,
	}

	for d := range upd {
		fmt.Println("Processing metrics...")
		pusher := pushers.For(d.DeviceId)
		pusher.Update("score", d.Score)
		for _, s := range d.Sensors {
			kind := strings.ToLower(s.Comp)
			if !pusher.Update(kind, s.Value) {
				fmt.Println("Skipping '" + kind + "' measurement")
			}
		}

		fmt.Println(fmt.Sprintf("%+v", d))
		err := pusher.Push()
		if err != nil {
			fmt.Println(err.Error())
			// TODO
		} else {
			fmt.Println("Pushed")
		}
	}
}

type PusherBuffer struct {
	pushers    map[int]prom.Pusher
	gatewayUrl string
}

func (b *PusherBuffer) For(deviceId int) *prom.Pusher {
	if buff, found := b.pushers[deviceId]; found {
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

	b.pushers[deviceId] = *p
	return p
}
