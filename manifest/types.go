package manifest

import (
	"fmt"
	"strings"

	"github.com/drone/drone-go/drone"
)

type Secret struct {
	Name   string      `yaml:name`
	Value  StringArray `yaml:value`
	Images StringArray `yaml:images`
	Events StringArray `yaml:events`
}

type SecretDef struct {
	Repo    StringArray `yaml:repo`
	Secrets []Secret    `yaml:secrets`
}

type SecretsManifest []SecretDef
type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		split := strings.Split(single, ",")
		for i, part := range split {
			split[i] = strings.Trim(part, " ")
		}
		*a = split
	} else {
		*a = multi
	}
	return nil
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

func (inst *Secret) ToDroneSecret() (*drone.Secret, error) {
	converted := &drone.Secret{}
	converted.Name = inst.Name

	if len(inst.Value) != 0 {
		converted.Value = strings.Join(inst.Value, ",")
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
