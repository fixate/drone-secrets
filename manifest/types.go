package manifest

import (
	"fmt"
	"log"
	"strings"

	"github.com/drone/drone-go/drone"
)

type SecretDef struct {
	Name   string   `yaml:name`
	Value  string   `yaml:value`
	Items  []string `yaml:items`
	Images []string `yaml:images`
	Events []string `yaml:events`
}

type Secret struct {
	Repo  string      `yaml:repo`
	Items []SecretDef `yaml:items`
}

type SecretsManifest struct {
	Secrets []Secret `yaml:secrets`
}

func toDroneEvent(event string) (string, error) {
	switch event {
	case "pr", "pull-request", "pull":
		return drone.EventPull, nil
	case "push":
		return drone.EventPush, nil
	case "tag":
		return drone.EventTag, nil
	case "deployment":
		return drone.EventDeploy, nil

	default:
		return "", fmt.Errorf("manifest: Invalid event type '%s'", event)
	}
}

func (inst *SecretDef) ToDroneSecret() (*drone.Secret, error) {
	converted := &drone.Secret{}
	converted.Name = inst.Name
	if len(inst.Value) != 0 {
		converted.Value = inst.Value
	}

	if len(inst.Items) != 0 {
		converted.Value = strings.Join(inst.Items, ",")
	}

	if len(inst.Value) != 0 && len(inst.Items) != 0 {
		log.Println("WARNING: setting the value and items keys for a secret is invalid. Items will be preferred.")
	}

	converted.Images = inst.Images

	for _, evt := range inst.Events {
		convertedEvt, err := toDroneEvent(evt)
		if err != nil {
			return converted, err
		}

		converted.Events = append(converted.Events, convertedEvt)
	}

	return converted, nil
}
