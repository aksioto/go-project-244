package main

import (
	"code/internal/parser"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	reg := parser.NewRegistry()
	reg.Register(&parser.JSONParser{}, ".json")
	reg.Register(&parser.YAMLParser{}, ".yaml", ".yml")

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

			data1, err := reg.ParseFile(filePath1)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %w", filePath1, err)
			}

			data2, err := reg.ParseFile(filePath2)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %w", filePath2, err)
			}

			fmt.Printf("File 1: %+v\n", data1)
			fmt.Printf("File 2: %+v\n", data2)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
