# Drone Secrets

Will set secrets from a yaml manifest on your drone server.

### Motivation

Currently, there is no simple declaritive way to manage many drone secrets. It's a hassel and error prone to use `drone-cli` to set secrets per repo, per image in bulk. A shell script quickly becomes unwieldy. With `drone-secrets` you can create a yaml manifest to clearly define your configuration for your repos, per image. Common configuration can be set for multiple repos to keep things DRY. If you need to restore your secrets (say if you mistakenly deleted your repo), you have a way to quickly set all the required configuration. The manifest(s) can be kept in source control.

### Install:

See [releases](https://github.com/fixate/drone-secrets/releases/latest) or build repo using `make build` or `make builddev`

TODO:

- Use CI for github releases for multiple targets
- Make go get work to fetch the command (If people start caring
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
