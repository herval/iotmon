package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(ctx context.Context, client *http.Client, url string, bearerToken string, out interface{}, errOut interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	res := out
	if err := sendRequest(client, bearerToken, req, &res, &errOut); err != nil {
		return err
	}

	return nil
}

func sendRequest(client *http.Client, bearerToken string, req *http.Request, v interface{}, errRes interface{}) error {
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
		fmt.Println("Error! " + res.Status)
		//var errRes ErrorResponse
		if errRes != nil {
			if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
				return err
			}
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
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
