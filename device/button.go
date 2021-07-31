package device

import (
	"context"
	"periph.io/x/periph/conn/gpio"
	"time"
)

type Action string

const (
	Push    Action = "push"
	Release Action = "release"
)

func NewButton(pin gpio.PinIO, timeout time.Duration) *Button {
	return &Button{
		pin:     pin,
		timeout: timeout,
	}
}

type Button struct {
	pin     gpio.PinIO
	timeout time.Duration
}

func (t *Button) Run(ctx context.Context, actions chan<- Action) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if t.pin.WaitForEdge(t.timeout) == false {
				continue
			}

			if t.pin.Read() == gpio.High {
				actions <- Push
			} else {
				actions <- Release
			}
		}
	}
}
