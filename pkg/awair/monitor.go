package awair

import (
	"context"
	"fmt"
	"time"
)

type Monitor struct {
	client *Client
	ctx    context.Context
}

type Updates chan *RawDataPoints

func NeWPollingMonitor(ctx context.Context, token string) *Monitor {
	cli := NewClient(token)

	return &Monitor{
		ctx:    ctx,
		client: cli,
	}
}

// fail startup if unable to fetch devices
func (m *Monitor) Start(updatesChann Updates) error {
	dev, err := m.fetchDevices()
	if err != nil {
		return err
	}

	go func(dev []Device) {
		for {
			for _, d := range dev {
				data, err := m.fetchRawData(d)
				if err != nil {
					// TODO !??
				}

				// TODO buffer already posted?
				for _, r := range data {
					updatesChann <- r
				}
			}

			time.Sleep(time.Second * 10)

			dev, err = m.fetchDevices()
			if err != nil {
				// TODO??!?
			}
		}
	}(dev)

	return nil
}

func (m *Monitor) fetchDevices() ([]Device, error) {
	fmt.Println("Fetching Awair devices...")
	res, err := m.client.Devices(m.ctx)
	println(fmt.Sprintf("%+v", res))
	if err != nil {
		return nil, err
	}

	return res.Devices, nil
}

func (m *Monitor) fetchRawData(device Device) ([]*RawDataPoints, error) {
	data, err := m.client.Latest(m.ctx, &device)
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintf("%+v", data))

	return []*RawDataPoints{data}, nil
}
