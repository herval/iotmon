package aqi_test

import (
	"fmt"
	"github.com/herval/iotcollector/pkg/aqi"
	"os"
	"testing"
)

func TestAqiCnClient_GetHere(t *testing.T) {
	c := aqi.NewClient(os.Getenv("AQICN_API_KEY"))

	data, err := c.GetHere()
	if err != nil {
		t.Fatal()
	}

	fmt.Println(data)
	fmt.Println(aqi.AqiForLevel(data.Aqi))
}
