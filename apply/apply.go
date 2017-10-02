package apply

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/fixate/drone-secrets/client"
	mfst "github.com/fixate/drone-secrets/manifest"
	"github.com/fixate/drone-secrets/utils"
	. "github.com/logrusorgru/aurora"

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

func processManifest(c *cli.Context, client drone.Client, manifest mfst.SecretsManifest) error {
	for _, manifestSecret := range manifest {
		if err := ensureSecrets(client, &manifestSecret); err != nil {
			return err
		}
	}

	return nil
}

func ensureSecrets(client drone.Client, manifestSecret *mfst.SecretDef) error {
	for _, repo := range manifestSecret.Repo {
		if err := ensureSecretsForRepo(client, manifestSecret, repo); err != nil {
			return err
		}
	}
	return nil
}

func ensureSecretsForRepo(client drone.Client, manifestSecret *mfst.SecretDef, repo string) error {
	owner, name, err := utils.ParseRepo(repo)
	if err != nil {
		return err
	}

	fmt.Printf("Creating secrets for repository '%s'.\n\n", Bold(repo))

	secrets := manifestSecret.Secrets
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
			_, uerr := client.SecretUpdate(owner, name, secret)
			if uerr != nil {
				fmt.Printf("%s %s\n", Green("✓"), Bold(secret.Name))
			} else {
				fmt.Printf("%s %s.\n", Red("✕"), Bold(secret.Name))
				return uerr
			}
		} else {
			_, uerr := client.SecretCreate(owner, name, secret)
			if uerr != nil {
				fmt.Printf("%s %s\n", Green("✓"), Bold(secret.Name))
			} else {
				fmt.Printf("%s %s.\n", Red("✕"), Bold(secret.Name))
				return uerr
			}
		}
	}

	fmt.Print("\n")

	return nil
}
