package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "wingman",
		Version: "1.0.0",
		Usage:   "Superbird manager",
		Commands: []*cli.Command{
			{
				Name:  "ab",
				Usage: "manipulate a/b data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "boot-result",
						Usage: "set the boot result. 0 for failure, 1 for success",
						Action: func(cCtx *cli.Context, s string) error {
							info, err := OpenAndLoadABData()
							if err != nil {
								return err
							}

							result, err := strconv.Atoi(s)
							if err != nil {
								return fmt.Errorf("invalid boot result: %v", err)
							}

							slot := info.GetActiveSlot()
							if result == 0 {
								// info.Failover()
							} else {
								info.SetSuccessfulBoot(slot)
							}
							return info.Save()
						},
					},
					&cli.StringFlag{
						Name:  "slot",
						Usage: "set the active boot slot. 0 for A, 1 for B",
						Action: func(cCtx *cli.Context, s string) error {
							info, err := OpenAndLoadABData()
							if err != nil {
								return err
							}

							slot, err := strconv.Atoi(s)
							if err != nil || (slot != 0 && slot != 1) {
								return fmt.Errorf("invalid slot number: must be 0 or 1")
							}

							info.SetActiveSlot(slot)
							return info.Save()
						},
					},
					&cli.BoolFlag{
						Name:  "reset",
						Usage: "reset all boot data and switch back to slot A",
						Action: func(cCtx *cli.Context, b bool) error {
							info, err := OpenAndLoadABData()
							if err != nil {
								return err
							}

							info.Reset()
							return info.Save()
						},
					},
				},
				Action: func(cCtx *cli.Context) error {
					info, err := OpenAndLoadABData()
					if err != nil {
						return err
					}
					info.DumpInfo()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
