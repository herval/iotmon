package awair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURLV1 = "https://developer-apis.awair.is/v1"
)

type Client struct {
	client *http.Client
	token  string
}

func NewClient(authToken string) *Client {
	return &Client{
		token: authToken,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) Devices(ctx context.Context) (*DevicesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices", BaseURLV1), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := DevicesResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) Latest(ctx context.Context, device *Device) (*RawDataPoints, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices/%s/%d/air-data/latest", BaseURLV1, device.DeviceType, device.DeviceId), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := RawDataResponse{}
	if err := c.sendRequest(req, &res); err != nil {
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

func (c *Client) RawData(ctx context.Context, device *Device) (*RawDataResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/self/devices/%s/%d/air-data/raw", BaseURLV1, device.DeviceType, device.DeviceId), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := RawDataResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	for _, d := range res.Data {
		d.DeviceId = device.DeviceId
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := v
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
