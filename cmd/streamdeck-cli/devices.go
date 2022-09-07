package main

import (
	"fmt"

	"github.com/Doridian/streamdeck"
	"github.com/muesli/coral"
)

var (
	devicesCmd = &coral.Command{
		Use:   "devices",
		Short: "devices lists all available Stream Deck devices",
		RunE: func(cmd *coral.Command, args []string) error {
			_ = d.Close()

			devs, err := streamdeck.Devices()
			if err != nil {
				return fmt.Errorf("no Stream Deck devices found: %s", err)
			}
			if len(devs) == 0 {
				return fmt.Errorf("no Stream Deck devices found")
			}

			fmt.Printf("Found %d devices:\n", len(devs))

			for _, d := range devs {
				if err := d.Open(); err != nil {
					return fmt.Errorf("can't open device: %s", err)
				}

				ver, err := d.FirmwareVersion()
				if err != nil {
					return fmt.Errorf("can't retrieve device info: %s", err)
				}
				fmt.Printf("Found device with %d keys (firmware %s)\n",
					d.Keys, ver)

				_ = d.Close()
			}

			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(devicesCmd)
}
