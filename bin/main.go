package main

import (
	"context"
	"fmt"
	awair2 "github.com/herval/iotcollector/pkg/awair"
	prom2 "github.com/herval/iotcollector/pkg/prom"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {

	updates := make(chan *awair2.RawDataPoints)

	awairMonitor := awair2.NeWPollingMonitor(
		context.Background(),
		os.Getenv("AWAIR_TOKEN"),
	)

	go func(upd awair2.Updates) {
		if err := awairMonitor.Start(upd); err != nil {
			panic(err)
		}
	}(updates)

	go func(upd awair2.Updates) {
		pushers := PusherBuffer{
			pushers: map[int]prom2.Pusher{},
		}

		select {
		case d := <-upd:

			pusher := pushers.For(d.DeviceId)
			pusher.Update("score", d.Score)

			fmt.Println(fmt.Sprintf("%+v", d))
			err := pusher.Push()
			if err != nil {
				fmt.Println(err.Error())
				// TODO
			}
		}

	}(updates)

	select {}
}

type PusherBuffer struct {
	pushers map[int]prom2.Pusher
}

func (b *PusherBuffer) For(deviceId int) *prom2.Pusher {
	if buff, found := b.pushers[deviceId]; found {
		return &buff
	}

	gaugeNames := []string{
		"score", "voc", "co2", "pm25", "humid", "temp",
	}

	p := prom2.NewPusher(
		map[string]string{
			"sensor_id": fmt.Sprintf("%d", deviceId),
		},
		gaugeNames,
		"awair_sensor")

	b.pushers[deviceId] = *p
	return p
}
