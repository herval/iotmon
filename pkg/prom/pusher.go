package prom

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Pusher struct {
	pusher *push.Pusher
	gauges map[string]prometheus.Gauge
}

func NewPusher(url string, jobName string, labels map[string]string, gaugeNames []string, gaugePrefix string) *Pusher {
	gauges := map[string]prometheus.Gauge{}

	for _, g := range gaugeNames {
		gauges[g] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_%s", gaugePrefix, g),
			Help: fmt.Sprintf("%s level", g),
		})
	}

	// TODO
	p := push.New(url, jobName)

	for _, g := range gauges {
		p = p.Collector(g)
	}

	for k, v := range labels {
		p = p.Grouping(k, v)
	}

	return &Pusher{
		pusher: p,
		gauges: gauges,
	}

	//	Push(); err != nil {
	//	fmt.Println("Could not push completion time to Pushgateway:", err)
	//}
}

func (p *Pusher) Update(gauge string, value float64) {
	p.gauges[gauge].Set(value)
}

func (p *Pusher) Push() error {
	return p.pusher.Push()
}
