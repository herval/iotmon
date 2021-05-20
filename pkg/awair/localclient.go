package awair

import (
	"context"
	"fmt"
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/settings/config/data", c.host), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := LocalDeviceConfig{}
	if err := sendRequest(c.client, "", req, &res); err != nil {
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/air-data/latest", c.host), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := RawDataResponse{}
	if err := sendRequest(c.client, "", req, &res); err != nil {
		return nil, err
	}

	for _, d := range res.Data {
		d.DeviceId = device.DeviceId
	}

	if len(res.Data) > 0 {
		return res.Data[0], nil
	}

	return nil, nil
}
