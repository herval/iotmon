package awair

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURLV1 = "https://developer-apis.awair.is/v1"
)


type CloudClient struct {
	client *http.Client
	token  string
}

func NewClient(authToken string) *CloudClient {
	return &CloudClient{
		token: authToken,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *CloudClient) Devices(ctx context.Context) (*DevicesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices", BaseURLV1), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := DevicesResponse{}
	if err := sendRequest(c.client, c.token, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *CloudClient) Latest(ctx context.Context, device *Device) (*RawDataPoints, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices/%s/%d/air-data/latest", BaseURLV1, device.DeviceType, device.DeviceId), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := RawDataResponse{}
	if err := sendRequest(c.client, c.token, req, &res); err != nil {
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

func (c *CloudClient) RawData(ctx context.Context, device *Device) (*RawDataResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices/%s/%d/air-data/raw", BaseURLV1, device.DeviceType, device.DeviceId), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := RawDataResponse{}
	if err := sendRequest(c.client, c.token, req, &res); err != nil {
		return nil, err
	}

	for _, d := range res.Data {
		d.DeviceId = device.DeviceId
	}

	return &res, nil
}