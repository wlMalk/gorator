package main

import (
	"fmt"
	"os"

	"github.com/wlMalk/gorator/generate"

	"github.com/codegangsta/cli"
)

func main() {
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
			},
			Usage: "generate Go source files based on config file",
			Action: func(c *cli.Context) {
				path := c.String("path")
				err := generate.Generate(path)
				if err != nil {
					fmt.Println(err.Error())
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
