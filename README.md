# craft <!-- omit in toc -->

- [How to use ?](#how-to-use-)
- [Commands](#commands)
- [Craft file](#craft-file)
- [Plugins](#plugins)
  - [Generic plugin](#generic-plugin)
  - [Golang plugin](#golang-plugin)
    - [With API layout](#with-api-layout)
    - [With Docker layout](#with-docker-layout)
  - [Helm plugin](#helm-plugin)
  - [License plugin](#license-plugin)
- [Examples](#examples)

## How to use ?

```sh
go install github.com/kilianpaquier/craft/cmd/craft@latest
```

## Commands

CLI commands:

- `init`: initializes a new craft projects. A few questions will be asked in the terminal to tune the generated layout:
  - project's description
  - project's main maintainer (can be anything, a group name, a person's name, alias, etc.)
  - will the project expose an api layer (golang based)
  - which openapi version the api layer will be (only asked if api layer is true)
  - will the project have a helm chart
- `generate`: generates the craft layout. Only works if the project was initialized (it has a `.craft` file at project's root). Available options are:
  - `--force`: forces the generation of a list of already generated files (`.gitlab-ci.yml`, `sonar.properties`, etc.)
  - `--force-all`: forces the generation of all generated files even if they exist (you may stage or stash your changes before running `craft generate` with this option).

## Craft file

Craft project generation is based on root's `.craft` file, it can contain the following configurations:

```yaml
# project's description (optional)
# used in various places like helm Chart.yml description
# Dockerfile description label
# api.yml description
description: some useful description

# project's maintainers (at least one must be provided)
# the first maintainer will be referenced in various places like in goreleaser configuration
# Dockerfile maintainer / authors label
# sonar.properties organization and project key prefix
# helm values.yml for images owner (e.g ghcr.io/maintainer/app_name)
# api.yml main contact
# all maintainers will be referenced in dependabot assignees and reviewers
# helm Chart.yml maintainers
maintainers:
  - name: maintainer
    email: maintainer@example.com
    url: maintainer.example.com

# project's license (optional)
# providing it will download the appropriate license
# used in various places like api.yml license
# goreleaser executables license
# github release workflow license addition to releases 
license: agpl-3.0 | apache-2.0 | bsd-2-clause | bsd-3-clause | bsl-1.0 | cc0-1.0 | epl-2.0 | gpl-2.0 | gpl-3.0 | lgpl-2.1 | mit | mpl-2.0 | unlicense

# project's CI (optional)
# providing it will create the appropriate ci files (.gitlab-ci.yml, .github/workflows/...)
ci:
  # ci name - self-explaining what each value will generate - (required when ci section is given)
  name: github | gitlab
  # ci options, providing one or multiple options with tune the ci generation (optional)
  options: [codecov, dependabot, sonar]

# project's api configuration
# providing it will create an api layer with golang
api:
  # project's api openapi version
  # not provided or provided as 'v2', the api layer will be generated with the help of go-swagger
  # provided as 'v3', the api layer will not be generated (not yet implemented)
  openapi_version: v2 | v3

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
no_makefile: true | false
```

## Plugins

Craft generation is based on plugins. Each plugin detects from `.craft` configuration and project's files if it needs to generate its part (or not).

Overall, the following plugins are implemented:
- [`generic`](#generic-plugin)
- [`golang`](#golang-plugin)
- [`helm`](#helm-plugin)
- [`license`](#license-plugin)
- [`openapi_v2`](#with-api-layout)

### Generic plugin

Craft project generation for anything but golang (because it's the only coding language implemented for now) will be generated with the generic plugin.

The following layout will be created:

```tree
├── .gitlab
│   ├── workflows
│   │   ├── .gitlab-ci.yml
├── .github
│   ├── workflows
│   │   ├── integration.yml
├── .craft (craft configuration file)
├── .gitlab-ci.yml
├── Makefile
└── README.md
```

It's a very simple generation with little features.

### Golang plugin

Craft project generation for golang is following the present layout: https://github.com/golang-standards/project-layout.

```tree
├── .gitlab
│   ├── workflows
│   │   ├── .gitlab-ci.yml
├── .github
│   ├── workflows
│   │   ├── dependencies.yml
│   │   ├── integration.yml
│   │   ├── publish.yml
├── cmd (executable binaries, prefix is useful for kubernetes identification)
│   ├── xyz (as many as desired CLIs)
│   │   ├── main.go
│   ├── cron-xyz (as many as desired cronjobs)
│   │   ├── main.go
│   ├── job-xyz (as many as desired jobs)
│   │   ├── main.go
│   ├── worker-xyz (as many as desired workers)
│   │   ├── main.go
├── internal
├── pkg
├── .craft
├── .gitignore
├── .gitlab-ci.yml
├── .golangci.yml
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── sonar.properties
```

#### With API layout

When the API option is present, then the following layout will be generated.

```tree
├── cmd
│   ├── <project_name>-api
│   │   ├── main.go (main associated to API layer)
├── internal
│   ├── api
│   │   ├── manual implementation with business layer for API layer
├── models
│   ├── generated models files by go-swagger for the API layer
├── pkg
│   ├── api
│   │   ├── generated client files by go-swagger for the API layer (consumers)
├── restapi
│   ├── generated server files by go-swagger for the API layer
```

#### With Docker layout

When the docker option is present and there's at least one executable, then the following files: `Dockerfile`, `.dockerignore` and `launcher.sh` will be generated.

```tree
├── cmd (executable binaries, prefix is useful for kubernetes identification)
│   ├── xyz (as many as desired CLIs)
│   │   ├── main.go
│   ├── cron-xyz (as many as desired cronjobs)
│   │   ├── main.go
│   ├── job-xyz (as many as desired jobs)
│   │   ├── main.go
│   ├── worker-xyz (as many as desired workers)
│   │   ├── main.go
├── .dockerignore
├── Dockerfile
└── launcher.sh (only when there's at least two main.go in cmd folder, parses the BINARY_NAME environment variable to run the right executable)
```

### Helm plugin

The helm plugin is in charge of generating the helm chart for the project. Depending on implemented coding languages, `values.yaml` file will contain values for `cronjobs`, `jobs` or `workers`.

For instance, associated with [golang plugin](#golang-plugin), kubernetes executables will be parsed from `cmd` folder.

The following layout will be created:

```tree
├── chart
│   ├── .craft (override values for helm chart)
└── .craft
```

### License plugin

The license plugin is only in charge of retrieving the appropriate `LICENSE` file according to `.craft > license` value.

## Examples

You may consult the `examples` for more information and details on generated files.