package manifest

import (
	"strings"

	"github.com/drone/drone-go/drone"
)

type Secret struct {
	Name   string      `yaml:name`
	Value  StringArray `yaml:value`
    EnablePullRequests bool `yaml:enable_pull_requests`
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

func (inst *Secret) ToDroneSecret() (*drone.Secret, error) {
	converted := &drone.Secret{}
	converted.Name = inst.Name

	if len(inst.Value) != 0 {
		converted.Data = strings.Join(inst.Value, ",")
	}

	converted.PullRequestPush = inst.EnablePullRequests;
	converted.PullRequest = inst.EnablePullRequests;

	return converted, nil
}
