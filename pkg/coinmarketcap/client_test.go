package coinmarketcap_test

import (
	"context"
	"fmt"
	"github.com/herval/iotcollector/pkg/coinmarketcap"
	"os"
	"testing"
)

func TestLatest(t *testing.T) {
	cli := coinmarketcap.NewClient(os.Getenv("COINMARKETCAP_API_KEY"))

	r, err := cli.Latest(context.Background(), []string{"BTC"})
	if err != nil {
		fmt.Println(err.Error())
	}

	if r != nil {
		for _, q := range r {
			fmt.Println(fmt.Sprintf("%s - %f", q.Name, q.Quote["USD"].Price))
		}
	}
}
