package awair

import (
	"context"
	"fmt"
	http2 "github.com/herval/iotcollector/pkg/http"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LocalClient struct {
	host   string
	client *http.Client
}

func NewLocalClient(host string) *LocalClient {
	return &LocalClient{
		host: host,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *LocalClient) Devices(ctx context.Context) (*DevicesResponse, error) {
	res := LocalDeviceConfig{}
	errorRes := ErrorResponse{}
	var err error

	if err = http2.Get(
		ctx,
		c.client,
		fmt.Sprintf("%s/settings/config/data", c.host),
		"",
		&res,
		&errorRes,
	); err != nil {
		return nil, err
	}

	res.DeviceId, err = strconv.Atoi(strings.Split(res.DeviceUuid, "_")[1])
	if err != nil {
		return nil, err
	}

	return &DevicesResponse{
		[]Device{
			toDevice(res),
		},
	}, nil
}

func toDevice(res LocalDeviceConfig) Device {
	return Device{
		DeviceId:   res.DeviceId,
		Name:       res.DeviceUuid,
		DeviceUUID: res.DeviceUuid,
		Timezone:   res.Timezone,
	}
}

func (c *LocalClient) Latest(ctx context.Context, device *Device) (*RawDataPoints, error) {
	res := LocalAirDataResponse{}
	errorRes := ErrorResponse{}
	var err error

	if err = http2.Get(
		ctx,
		c.client,
		fmt.Sprintf("%s/air-data/latest", c.host),
		"",
		&res,
		&errorRes,
	); err != nil {
		return nil, err
	}

	res.DeviceId = device.DeviceId

	return toDataPoints(res), nil
}

func toDataPoints(res LocalAirDataResponse) *RawDataPoints {
	return &RawDataPoints{
		Score:    res.Score,
		DeviceId: res.DeviceId,
		Sensors: []struct {
			Comp  string  `json:"comp"`
			Value float64 `json:"value"`
		}{
			{
				Comp:  "temp",
				Value: res.Temp,
			},
			{
				Comp:  "co2",
				Value: res.Co2,
			},
			{
				Comp:  "voc",
				Value: res.Voc,
			},
			{
				Comp:  "pm25",
				Value: res.Pm25,
			},
			{
				Comp:  "humid",
				Value: res.Humid,
			},
		},
	}
}
