# craft <!-- omit in toc -->

<p align="center">
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
  - [Init](#init)
  - [Generate](#generate)
- [Craft file](#craft-file)
  - [VSCode association and schema](#vscode-association-and-schema)
- [Generations](#generations)
- [Who is using craft ?](#who-is-using-craft-)
- [Craft as an SDK](#craft-as-an-sdk)

## How to use ?

### Go

```sh
go install github.com/kilianpaquier/craft/cmd/craft@latest
```

### Linux

```sh
OS="linux" # change it depending on your case
ARCH="amd64" # change it depending on your case
INSTALL_DIR="$HOME/.local/bin" # change it depending on your case

new_version=$(curl -fsSL "https://api.github.com/repos/kilianpaquier/craft/releases/latest" | jq -r '.tag_name')
url="https://github.com/kilianpaquier/craft/releases/download/$new_version/craft_${OS}_${ARCH}.tar.gz"
curl -fsSL "$url" | (mkdir -p "/tmp/craft/$new_version" && cd "/tmp/craft/$new_version" && tar -xz)
cp "/tmp/craft/$new_version/craft" "$INSTALL_DIR/craft"
```

## Commands

```
Craft initializes or generates craft projects. Craft projects are only defined by a .craft file 
and multiple files automatically generated to avoid multiple hours to setup Continuous Integration, coverage, security analyses, helm chart, etc.

Craft generation can be done with 'craft' command or 'craft generate' command.
Additional generation command are available to generate only subparts of craft layout (like 'craft chart').

Usage:
  craft [flags]
  craft [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate project layout
  help        Help about any command
  init        Initialize craft project
  version     Show current craft version

Flags:
  -h, --help                help for craft
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")

Use "craft [command] --help" for more information about a command.
```

### Init

```
Initialize new craft project

Usage:
  craft init [flags]

Flags:
  -h, --help   help for init

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```

### Generate

```
Generate project layout

Usage:
  craft generate [flags]

Flags:
  -h, --help   help for generate

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```

#### Chart

```
Generate project layout's helm chart

Usage:
  craft chart [flags]

Flags:
  -h, --help   help for chart

Global Flags:
      --log-format string   set logging format (either "text" or "json") (default "text")
      --log-level string    set logging level (default "info")
```

## Craft file

Craft project generation is based on root's `.craft` file, it can contain the following configurations:

```yaml
# bot in charge of keeping dependencies up to date
bot: dependabot | renovate

# project's CI (optional)
# providing it will create the appropriate ci files (.gitlab-ci.yml, .github/actions, .github/workflows)
ci:
  # auth configurations for various features in CI
  auth:
    # maintenance auth strategy for the specified maintenance bot (just above)
    maintenance: github-app | github-token | mend.io | personal-token

    # release auth for github only (how should the release token be retrieved)
    # will stay empty when using gitlab CICD
    release: github-app | github-token | personal-token

  # ci name - self-explaining what each value will generate - (required when ci section is given)
  name: github | gitlab

  # ci global options, providing one or multiple options with tune the ci generation (optional)
  options:
    - codecov
    - codeql
    - labeler
    - sonar

  # release specific options
  release:
    # whether the release should run automatically
    auto: true | false
    # whether backmerging should be configured for main, staging and develop branches
    backmerge: true | false
    # whether releasing should be disabled
    disable: true | false

  # static deployment configuration
  static:
    # static deployment name
    auto: true | false
    # static deployment automatisation (on main branches for github and on protected branches for gitlab)
    name: netlify | pages

# project's description (optional)
# used in various places like helm Chart.yml description
# Dockerfile description label
description: some useful description

docker:
  # specific docker registry to push images on (optional, default is none - docker.io)
  # used in various places like helm values.yml images registry
  # github release workflow to push images
  registry: ghcr.io
  # specific exposed port (optional, default is 3000)
  # used in various places like helm values.yml service port
  # Dockerfile exposed port
  port: 3000

# project's license (optional)
# providing it will download the appropriate license
# used in various places like goreleaser executables license
# github release workflow license addition to releases
license: agpl-3.0 | apache-2.0 | bsd-2-clause | bsd-3-clause | bsl-1.0 | cc0-1.0 | epl-2.0 | gpl-2.0 | gpl-3.0 | lgpl-2.1 | mit | mpl-2.0 | unlicense

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

# whether to generate an helm chart or not (optional)
no_chart: true | false

# whether to use goreleaser or not, it's only useful on golang based projects (optional)
no_goreleaser: true | false

# whether to generate a Makefile with useful commands (optional)
# this option is automatically disabled when working with nodejs generation
no_makefile: true | false

# whether to generate a README.md with initial badges and informations (optional)
no_readme: true | false

# platform override in case of gitlab on premise, bitbucket on premise, etc.
# by default, an on premise gitlab will be matched if the host contains "gitlab"
# by default, an on premise bitbucket will be matched if the host contains "bitbucket" or "stash"
# when not overridden, the platform is matched based on "git config --get remote.origin.url" on the returned host (github.com, gitlab.com, ...)
platform: bitbucket | gitea | github | gitlab
```

### VSCode association and schema

When working on **vscode**, feel free to use craft's schemas to help setup your project:

```json
{
    "files.associations": {
        ".craft": "yaml"
    }
}
```

It's only creating the association between yaml files and `.craft`, however combined with **vscode** extension **redhat.vscode-yaml**, 
it will load the schema fine since a header is added in all `.craft` when written.

## Generations

Craft generation is based on separated parsers.
Each parser checks from `.craft` configuration and project's files to add specific behaviors in a shared structure.
Once all parsers are executed, generation iterates over all templates files and generates the right one needed depending on shared structure information.

Multiple examples:
- A `go.mod` is detected with `Golang` parser, combined with `ci` configuration, then the appropriate CI will be generated.
- A `go.mod` is detected with `Golang` parser and a `hugo.(toml|yaml|...)` or `theme.(toml|yaml|...)` is detected too, combined with the `ci` and `static` options, 
  then the appropriate **Netlify** or **Pages** (it can be **GitLab** or **GitHub**) deployment will be generated in CI files.
- If `no_chart` is not given, a custom craft helm chart will be generated. 
  This helm chart can deploy cronjobs, jobs and workers easily from `values.yaml` file.
- A `package.json` is detected with `Node` parser, combined with `ci` configuration, then the appropriate CI will be generated
  (codecov analysis, sonar analysis, lint, tests, build if needed).

## Who is using craft ?

- https://github.com/kilianpaquier/compare (Golang library)
- https://github.com/kilianpaquier/craft (Golang CLI with executables as artifacts in releases)
- https://github.com/kilianpaquier/gitlab-storage-cleaner (Golang CLI with Docker deployment and executables as artifacts in releases)
- https://github.com/kilianpaquier/go-builder-generator (Golang CLI with executables as artifacts in releases)
- https://github.com/kilianpaquier/kilianpaquier.github.io (Hugo static website deployed with **Netlify**)
- https://github.com/kilianpaquier/pooling (Golang library)
- https://github.com/kilianpaquier/semantic-release-backmerge (**semantic-release** plugin with static build deployed in npmjs.org)
- https://gitlab.com/nath7098/personal-website (Node static website deployed with Docker)

## Craft as an SDK

Craft can also be used as an SDK, for that you may check the official documentation on [pkg.go.dev](https://pkg.go.dev/github.com/kilianpaquier/craft).