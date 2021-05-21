package awair

import (
	"context"
)

type Client interface {
	Devices(ctx context.Context) (*DevicesResponse, error)
	Latest(ctx context.Context, device *Device) (*RawDataPoints, error)
}

