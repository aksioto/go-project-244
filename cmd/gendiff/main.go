package main

import (
	"code"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"slices"
	"strings"
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

			allowedFormats := []string{"stylish"}
			if !slices.Contains(allowedFormats, format) {
				return fmt.Errorf("format \"%s\" is not supported\nallowed formats: %s", format, strings.Join(allowedFormats, ", "))
			}

			res, err := code.GenDiff(filePath1, filePath2, format)
			if err != nil {
				return fmt.Errorf("error generating diff: %w", err)
			}

			fmt.Println(res)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
