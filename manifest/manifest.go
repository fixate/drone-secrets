package manifest

import (
	"ioutil"

	"gopkg.in/yaml.v2"
)

func Load(file string) SecretsManifest, error {
	data, ferr := ioutil.ReadFile(path)
	if ferr != nil {
		return nil, ferr
	}

	doc := SecretsManifest{}
	err := yaml.Unmarshal(byte[](data), &doc)
	if err != nil {
		return nil, err
	}

	// TODO: verify manifest

	return doc, nil
}
