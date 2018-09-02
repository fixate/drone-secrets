# Drone Secrets

Will set secrets from a yaml manifest on your drone server.

### Motivation

Currently, there is no simple declaritive way to manage many drone secrets. It's a hassel and error prone to use `drone-cli` to set secrets per repo, per image in bulk. Your manifest file can be kept in source control (should be encrypted and/or only accessible by trusted entities) and is easily understood. Once `drone-secrets` runs you'll know the right secres are set.

### Install:

See [releases](https://github.com/fixate/drone-secrets/releases/latest) or build repo using `make build` or `make builddev`

TODO:

- Use CI for github releases for multiple targets
- Make go get work to fetch the command (I'm new to go so if people start caring
		about this project I'll take things further)

## Configuration

Put this in your shell environment (probably `.bashrc`)

```shell
export DRONE_SERVER=<address to your build server>
export DRONE_TOKEN=<your token (hamburger menu > token)>
```

### Usage:

```shell
drone-secrets apply -f manifest.yml
```

*manifest.yml*

```yaml
---
# Comma delimited or list accepted syntax accepted for repo, value, events and images
- repo: my/repo, my/other-repo
	  secrets:
	  # Set for my/repo and my/other-repo
	  - name: MY_SECRET
	    value: 12345

	  - name: SLACK_WEBHOOK
	    value: abcde
	    # Default events are push, tag, deployment
	    events: push,tag
	    image: plugins/slack

- repo: my/repo
  secrets:
	  # Setting value to a list	  
	  - name: PLUGINS_ENVIRONMENT_VARIABLES
	    events: 
	      - push
	      - tag
	    # List types are converted to a comma delimited string
	    value:
	      - PORT=1234
	      - SECRET_TOKEN=abcd1234
	    # Same as:
	    # value: PORT=1234,SECRET_TOKEN=abcde1234,...
	    images: 
	      - plugins/ecs
	      - plugins/ecs:*
```

## TODO:

- Tests 
- Optionally clean secrets not in manifest
