package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wlMalk/gorator/generate"

	"github.com/codegangsta/cli"
)

func main() {
	var err error

	gorator := cli.NewApp()
	gorator.Name = "gorator"
	gorator.Usage = "The Go ORM generator and more"
	gorator.Version = generate.VERSION
	gorator.Author = "Waleed AlMalki (wlMalk)"
	gorator.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Value: "",
					Usage: "path to project folder (must be inside GOPATH)",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "config version to use (optional)",
				},
			},
			Usage: "generate Go source files based on config file",
			Action: func(c *cli.Context) {
				if os.Getenv("GOPATH") == "" {
					fmt.Println("GOPATH is not defined")
					os.Exit(0)
				}

				path := c.String("path")
				if path == "" {
					path, err = os.Getwd()
				}

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}

				if !strings.Contains(path, os.Getenv("GOPATH")) {
					fmt.Println("path is not in $GOPATH")
					os.Exit(0)
				}

				version := c.String("version")

				err = generate.Generate(path, version, "orm")
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "migrate databases based on config file",
			Action: func(c *cli.Context) {
			},
		},
	}
	gorator.Run(os.Args)
}
