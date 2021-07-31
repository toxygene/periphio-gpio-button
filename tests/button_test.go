package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/toxygene/periphio-gpio-button/device"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpiotest"
	"testing"
	"time"
)

func TestButton(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		button := device.NewButton(pin, 1*time.Second)

		ctx, cancel := context.WithCancel(context.Background())

		actions := make(chan device.Action)
		defer close(actions)

		go func() {
			pin.EdgesChan <- gpio.High
			assert.Equal(t, <-actions, device.Push)

			pin.EdgesChan <- gpio.Low
			assert.Equal(t, <-actions, device.Release)

			cancel()
		}()

		err := button.Run(ctx, actions)

		assert.Errorf(t, err, "context canceled")
	})
}
