package apply

import (
	"io/ioutil"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/fixate/drone-secrets/client"

	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:   "apply",
	Usage:  "Apply secret manifest",
	Action: run,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "clean",
			Usage: "Remove secrets that are not included in the manifest",
		},
	},
}

var defaultSecretEvents = []string{
	drone.EventPush,
	drone.EventTag,
	drone.EventDeploy,
}

func run(c *cli.Context) error {
	config, merr := manifest.Load(c.String("manifest"))

	client, err := client.NewClient(c)
	if err != nil {
		return err
	}
	secret := &drone.Secret{
		Name:   c.String("name"),
		Value:  c.String("value"),
		Images: c.StringSlice("image"),
		Events: c.StringSlice("event"),
	}
	if len(secret.Events) == 0 {
		secret.Events = defaultSecretEvents
	}
	if strings.HasPrefix(secret.Value, "@") {
		path := strings.TrimPrefix(secret.Value, "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return ferr
		}
		secret.Value = string(out)
	}
	_, err = client.SecretCreate(owner, name, secret)
	return err
}
