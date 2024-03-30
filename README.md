<!-- This file is safe to edit. Once it exists it will not be overwritten. -->

# craft <!-- omit in toc -->

<p align="center">
  <img alt="GitHub Actions" src="https://img.shields.io/github/actions/workflow/status/kilianpaquier/craft/integration.yml?branch=main&style=for-the-badge">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/kilianpaquier/craft?include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues-raw/kilianpaquier/craft?style=for-the-badge">
  <img alt="GitHub License" src="https://img.shields.io/github/license/kilianpaquier/craft?style=for-the-badge">
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/kilianpaquier/craft/main?style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/kilianpaquier/craft/main?style=for-the-badge&label=Go+Version">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kilianpaquier/craft?style=for-the-badge">
</p>

---

- [How to use ?](#how-to-use-)
  - [Go](#go)
  - [Linux](#linux)
- [Commands](#commands)
  - [Generate](#generate)
- [Craft file](#craft-file)
  - [VSCode association and schema](#vscode-association-and-schema)
- [Generations](#generations)
  - [Generic](#generic)
  - [Golang](#golang)
  - [Hugo](#hugo)
  - [Helm](#helm)
  - [License](#license)
  - [Nodejs](#nodejs)

## How to use ?

### Go

```sh
go install github.com/kilianpaquier/craft/cmd/craft@latest
```

### Linux

**With dpkg:**

```sh
version=$(curl -fsSL "https://api.github.com/repos/kilianpaquier/craft/releases/latest" | jq -r '.tag_name')
curl -fsSL "https://github.com/kilianpaquier/craft/releases/download/${version}/craft_linux_amd64.deb" -o /tmp/craft.deb
sudo dpkg -i craft.deb && rm /tmp/craft.deb
```

**With the tar archive:**

```sh
version=$(curl -fsSL "https://api.github.com/repos/kilianpaquier/craft/releases/latest" | jq -r '.tag_name')
curl -fsSL "https://github.com/kilianpaquier/craft/releases/download/${version}/craft_linux_amd64.tar.gz" | (cd "${HOME}/.local/craft" && tar -xz)
chmod +x "${HOME}/.local/craft/craft" && ln -sf "${HOME}/.local/craft/craft" "${HOME}/.local/bin/craft"
```

## Commands

```
Craft stands here to generate a similar project layout for all your projects. 
Multiple coding languages are supported and even helm chart can be generated. 
For more information please consult each command specificities.

Usage:
  craft [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate the project layout
  help        Help about any command
  init        Initialize the project layout
  version     Shows current craft version

Flags:
  -h, --help               help for craft
  -l, --log-level string   set logging level

Use "craft [command] --help" for more information about a command.
```

### Generate

```
Generate the project layout

Usage:
  craft generate [flags]

Flags:
  -f, --force strings   force regenerating a list of templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)
      --force-all       force regenerating all templates (.gitlab-ci.yml, sonar.properties, Dockerfile, etc.)
  -h, --help            help for generate

Global Flags:
  -l, --log-level string   set logging level
```

## Craft file

Craft project generation is based on root's `.craft` file, it can contain the following configurations:

```yaml
# project's description (optional)
# used in various places like helm Chart.yml description
# Dockerfile description label
description: some useful description

# project's maintainers (at least one must be provided)
# the first maintainer will be referenced in various places like in goreleaser configuration
# Dockerfile maintainer / authors label
# sonar.properties organization and project key prefix
# helm values.yml for images owner (e.g ghcr.io/maintainer/app_name)
# all maintainers will be referenced in dependabot assignees and reviewers
# helm Chart.yml maintainers
maintainers:
  - name: maintainer
    email: maintainer@example.com
    url: maintainer.example.com

# project's license (optional)
# providing it will download the appropriate license
# used in various places like goreleaser executables license
# github release workflow license addition to releases 
license: agpl-3.0 | apache-2.0 | bsd-2-clause | bsd-3-clause | bsl-1.0 | cc0-1.0 | epl-2.0 | gpl-2.0 | gpl-3.0 | lgpl-2.1 | mit | mpl-2.0 | unlicense

# project's CI (optional)
# providing it will create the appropriate ci files (.gitlab-ci.yml, .github/workflows/...)
ci:
  # ci name - self-explaining what each value will generate - (required when ci section is given)
  name: github | gitlab
  # release specific options
  release:
    # whether the release should run automatically
    auto: true | false
    # whether backmerging should be configured for main, staging and develop branches
    backmerge: true | false
    # whether the release should be disabled
    disable: true | false
    # release tokenization mode for github only (how should the release token be retrieved)
    # will stay empty when using gitlab CICD
    mode: github-apps | personal-token | github-token
  # ci global options, providing one or multiple options with tune the ci generation (optional)
  options:
    - codecov
    - codeql
    - dependabot
    - pages
    - renovate
    - sonar

# platform override in case of gitlab on premise, bitbucket on premise, etc.
# by default, an on premise gitlab will be matched if the host contains "gitlab"
# by default, an on premise bitbucket will be matched if the host contains "bitbucket" or "stash"
# when not overridden, the platform is matched based on "git config --get remote.origin.url" on the returned host (github.com, gitlab.com, ...)
platform: bitbucket | gitea | github | gitlab

docker:
  # specific docker registry to push images on (optional, default is none - docker.io)
  # used in various places like helm values.yml images registry
  # github release workflow to push images
  registry: ghcr.io
  # specific exposed port (optional, default is 3000)
  # used in various places like helm values.yml service port
  # Dockerfile exposed port
  port: 3000

# whether to generate an helm chart or not (optional)
no_chart: true | false

# whether to use goreleaser or not, it's only useful on golang based projects (optional)
no_goreleaser: true | false

# whether to generate a Makefile with useful commands (optional)
# this option is automatically disabled when working with nodejs generation
no_makefile: true | false
```

### VSCode association and schema

When working on vscode, feel free to use craft's schemas to help setup your project:

```json
{
    "files.associations": {
        ".craft": "yaml"
    },
    "yaml.schemas": {
        "https://raw.githubusercontent.com/kilianpaquier/craft/main/.schemas/craft.schema.json": [
            "**/.craft", 
            "!**/chart/.craft"
        ],
        "https://raw.githubusercontent.com/kilianpaquier/craft/main/.schemas/chart.schema.json": [ 
            "**/chart/.craft"
        ]
    }
}
```

## Generations

Craft generation is based on separated detections. Each detection checks from `.craft` configuration and project's files if it needs to generate its part (or not).

Multiple detections are implemented, and generates various files (please consult the [`examples`](./examples/) folder for more information):

### Generic

When no primary generation is detected (golang, nodejs or hugo), then this generation is automatically used. 

It doesn't generate much, just some CI files (in case CI option is provided) for versioning (semantic release), a README.md and makefiles (to easily generate again and clean the git repository).

Feel free to check either [`generic_github`](./examples/generic_github/) or [`generic_gitlab`](./examples/generic_gitlab/) to see what this generation specifically does.

### Golang

When golang generation is detected (a `go.mod` is present at root and is valid), it will generate the appropriate files around golang testing, makefiles, CI integration (in case CI option is given), etc.

Feel free to check either [`golang_github`](./examples/golang_github/) or [`golang_gitlab`](./examples/golang_gitlab/) to see what this generation specifically does.

### Hugo

When hugo generation is detected (a `go.mod` is present at root, is valid and either a `hugo.(toml|yaml|...)` or `theme.(toml|yaml|...)` are present), it will generate the appropriate files around hugo static web build, CI integration (in case CI option is given), etc.

Feel free to check either [`hugo_github`](./examples/hugo_github/) or [`hugo_gitlab`](./examples/hugo_gitlab/) to see what this generation specifically does.

### Helm

When helm generation is detected (`no_chart` is either not provided or false), it will generate a specific helm chart capable of easily deploying cronjobs, jobs, workers or even chart dependencies.

Feel free to check [`helm`](./examples/helm/) to see what this generation specifically does.

### License

When license generation is detected (`license` is provided in `.craft` root file and is valid), it will generate the appropriate license file.

### Nodejs

When nodejs generation is detected (a `package.json` is present at root and is valid), it will the appropriate files around nodejs testing, integration, etc.

Feel free to check either [`nodejs_github`](./examples/nodejs_github/) or [`nodejs_gitlab`](./examples/nodejs_gitlab/) to see what this generation specifically does.
