package main

import (
	"code"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Value:       "stylish",
				DefaultText: "\"stylish\"",
				Usage:       "output format",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			args := c.Args()
			if args.Len() < 2 {
				return fmt.Errorf("usage: gendiff <filepath1> <filepath2>")
			}

			filePath1 := args.Get(0)
			filePath2 := args.Get(1)
			format := c.String("format")

			res, err := code.GenDiff(filePath1, filePath2, format)
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
