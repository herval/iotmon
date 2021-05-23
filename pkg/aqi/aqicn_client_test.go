package aqi_test

import (
	"context"
	"fmt"
	"github.com/herval/iotcollector/pkg/aqi"
	"os"
	"testing"
)

func TestAqiCnClient_GetHere(t *testing.T) {
	c := aqi.NewClient(os.Getenv("AQICN_API_KEY"))

	data, err := c.GetHere(context.Background())
	if err != nil {
		t.Fatal()
	}

	fmt.Println(data)
	fmt.Println(aqi.AqiForLevel(data.Data.Aqi))
}
