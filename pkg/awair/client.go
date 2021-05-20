package awair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	Devices(ctx context.Context) (*DevicesResponse, error)
	Latest(ctx context.Context, device *Device) (*RawDataPoints, error)
}

func sendRequest(client *http.Client, bearerToken string, req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

	res, err := client.Do(req)
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

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("can't read body: %v", err)
	}

	//fmt.Println(string(b))

	fullResponse := v

	if err = json.Unmarshal(b, &fullResponse); err != nil {
		return err
	}

	return nil
}
