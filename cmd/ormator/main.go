package main

import (
	"fmt"
	"os"

	"github.com/wlMalk/ormator/generate"

	"github.com/codegangsta/cli"
)

func main() {
	ormator := cli.NewApp()
	ormator.Name = "ormator"
	ormator.Usage = "The Go ORM generator and more"
	ormator.Version = generate.VERSION
	ormator.Author = "Waleed AlMalki (wlMalk)"
	ormator.Commands = []cli.Command{
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
				err := generate.GenerateFromFile(path + string(os.PathSeparator) + "config.yml")
				if err != nil {
					fmt.Println(err.Error())
				}
			},
		},
	}
	ormator.Run(os.Args)
}
