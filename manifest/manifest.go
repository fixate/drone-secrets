package manifest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Load(path string) (SecretsManifest, error) {
	data, err := ioutil.ReadFile(path)
	doc := make(SecretsManifest, 0)
	if err != nil {
		return doc, err
	}

	err = yaml.Unmarshal([]byte(data), &doc)
	if err != nil {
		return doc, err
	}

	// TODO: verify manifest

	return doc, nil
}
