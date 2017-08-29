package manifest

import (
	"fmt"

	"github.com/drone/drone-go/drone"
)

type SecretDef struct {
	Name: string `yaml:name`,
	Value: string `yaml:value`,
	Images: []string `yaml:images`,
	Events: []string `yaml:events`,
}

type Secret struct {
	Repo: string `yaml:repo`,
	Items: []SecretDef `yaml:items`,
}

type SecretsManifest struct {
	Secrets: []Secret `yaml:secrets`,
}

func toDroneEvents(events []string) []string error {
	var result []string;
	for evt := range events {
		var event string;
		switch evt {
		case "pr":
		case "pull-request":
			event = drone.EventPullRequest
			break
		case "push":
			event = drone.EventPush
			break
		case "tag":
			event = drone.EventTag
			break
		case "deployment":
			event = drone.EventDeploy
			break

		default:
			return nil, fmt.Errorf("manifest: Invalid event type '%s'", evt)
		}

		result = append(result, event)
	}
	return result, nil;
}

func (inst *SecretDef) ToDroneSecret() []drone.Secret error {
	converted := &drone.Secret{}
	converted.Name = inst.Name
	converted.Value = inst.Value
	converted.Images = inst.Images
	newEvents, err := toEvents(inst.Events)
	if err != nil {
		return err
	}
	converted.Events = newEvents 
	return converted, nil
}



