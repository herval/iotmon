package main

import (
	"context"
	"fmt"
	"github.com/herval/iotcollector/pkg/awair"
	"github.com/herval/iotcollector/pkg/prom"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"time"
)

func main() {

	updates := make(chan *awair.RawDataPoints)

	awairMonitor := awair.NeWPollingMonitor(
		context.Background(),
		os.Getenv("AWAIR_TOKEN"),
	)

	go func(upd awair.Updates) {
		if err := awairMonitor.Start(upd, time.Minute*10); err != nil {
			panic(err)
		}
	}(updates)

	go func(upd awair.Updates) {
		pushers := PusherBuffer{
			pushers:    map[int]prom.Pusher{},
			gatewayUrl: os.Getenv("PROMETHEUS_URL"),
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
			}
			fmt.Println("Pushed")
		}

	}(updates)

	select {}
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
