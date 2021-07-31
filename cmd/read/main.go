package main

import (
	"context"
	"flag"
	"github.com/toxygene/periphio-gpio-button/device"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"time"
)

func main() {
	pinName := flag.String("pin", "", "pin name for the button")
	help := flag.Bool("help", false, "print help page")

	flag.Parse()

	if *help || *pinName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	button := device.NewButton(gpioreg.ByName(*pinName), 3*time.Second)

	g := new(errgroup.Group)

	ctx, cancel := context.WithCancel(context.Background())

	actions := make(chan device.Action)
	g.Go(func() error {
		err := button.Run(ctx, actions)
		close(actions)
		return err
	})

	g.Go(func() error {
		for a := range actions {
			println(a)
		}

		return nil
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		defer func() {
			signal.Stop(c)
			cancel()
		}()

		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}

	os.Exit(0)
}
