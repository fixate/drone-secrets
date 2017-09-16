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
			Name:   "t, token",
			Usage:  "server auth token",
			EnvVar: "DRONE_TOKEN",
		},
		cli.StringFlag{
			Name:   "s, server",
			Usage:  "server location",
			EnvVar: "DRONE_SERVER",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip ssl verfification",
			EnvVar: "DRONE_SKIP_VERIFY",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "socks-proxy",
			Usage:  "socks proxy address",
			EnvVar: "SOCKS_PROXY",
			Hidden: true,
		},
		cli.BoolFlag{
			Name:   "socks-proxy-off",
			Usage:  "socks proxy ignored",
			EnvVar: "SOCKS_PROXY_OFF",
			Hidden: true,
		},
	}
	app.Commands = []cli.Command{
		apply.Command,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
