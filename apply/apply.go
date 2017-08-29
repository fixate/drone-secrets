package apply

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/fixate/drone-secrets/client"
	mfst "github.com/fixate/drone-secrets/manifest"
	"github.com/fixate/drone-secrets/utils"

	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:   "apply",
	Usage:  "Apply secret manifest",
	Action: run,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "f, manifest",
			Usage:  "File manifest to use for secret creation",
			EnvVar: "DRONE_SECRET_MANIFEST",
		},
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
	manifest, err := mfst.Load(c.String("manifest"))
	if err != nil {
		return err
	}

	client, err := client.NewClient(c)
	if err != nil {
		return err
	}

	return processManifest(c, client, manifest)
}

func processManifest(c *cli.Context, client drone.Client, manifest *mfst.SecretsManifest) error {
	for _, manifestSecret := range manifest.Secrets {
		if err := ensureSecrets(client, &manifestSecret); err != nil {
			return err
		}
	}

	if c.Bool("clean") {
		fmt.Print("TODO: implement clean")
	}

	return nil
}

func ensureSecrets(client drone.Client, manifestSecret *mfst.Secret) error {
	owner, name, err := utils.ParseRepo(manifestSecret.Repo)
	if err != nil {
		return err
	}

	log.Printf("Creating secrets for repository '%s'.\n", manifestSecret.Repo)

	secrets := manifestSecret.Items
	for _, scrt := range secrets {
		secret, err := scrt.ToDroneSecret()
		if err != nil {
			return err
		}

		if len(secret.Events) == 0 {
			secret.Events = defaultSecretEvents
		}

		if strings.HasPrefix(secret.Value, "@") {
			path := strings.TrimPrefix(secret.Value, "@")
			out, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			secret.Value = string(out)
		}

		_, err = client.Secret(owner, name, secret.Name)
		if err != nil {
			_, err = client.SecretCreate(owner, name, secret)
			if err != nil {
				return err
			}
			log.Printf("'%s' created.\n", secret.Name)
		} else {
			_, err = client.SecretUpdate(owner, name, secret)
			if err != nil {
				return err
			}
			log.Printf("'%s' updated.\n", secret.Name)
		}
	}

	return nil
}
