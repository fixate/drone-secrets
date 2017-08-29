package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fixate/drone-secrets/apply"

	"github.com/urfave/cli"
)

const version string = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "drone secret"
	app.Version = version
	app.Usage = "Apply secret manifests to your drone server"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "f, manifest",
			Usage:  "File manifest to use for secret creation",
			EnvVar: "DRONE_SECRET_MANIFEST",
		},
		cli.StringFlag{
			Name:   "server",
			Usage:  "drone server address",
			EnvVar: "DRONE_SERVER_ADDRESS",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "drone server token",
			EnvVar: "DRONE_SERVER_TOKEN",
		},
	}
	app.Commands = []cli.Command{
		apply.Command,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	fmt.Println("herheerere")
	return nil
}
