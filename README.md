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
- `description`: the project description.
- `docker_registry`: a specific docker registry URL that will be used for gitlab cicd docker build job and for helm docker pull registry.
- `license`: one of `agpl-3.0, apache-2.0, bsd-2-clause, bsd-3-clause, bsl-1.0, cc0-1.0, epl-2.0, gpl-2.0, gpl-3.0, lgpl-2.1, mit, mpl-2.0, unlicense`. When given, the according `LICENSE` file will be generated.
- `maintainers`: the list of project's maintainers. Each maintainer must at least have a `name`, they can also have an `email` and `url`.
- `no_api`: when provided, no API layer will be generated.
  - based on `go-swagger` when `openapi_version` is neither given or with `v2` value
- `no_chart`: when provided, no helm chart in `chart` folder will be generated.
- `ci`: one of `gitlab github`.
  - if `gitlab` is provided, continuous integration based on [`kilianpaquier/cicd`](https://gitlab.com/kilianpaquier/cicd) integration templates will be generated depending on project language.
  - if `github` is provided, craft generated github workflows will be generated depending on project language.
- `no_dockerfile`: when provided, no `Dockerfile` will be generated.
- `no_goreleaser`: when provided, no `.goreleaser.yml` file will be generated (in any case, if the project isn't golang based, no file will be generated).
- `no_makefile`: when provided, no `Makefile` will be generated.
- `sonar`: when provided, a `sonar.properties` will be generated. As such, a sonar analysis job will be executed if `ci` is provided.
- `codecov`: when provided, a `codecov.yml` will be generated. If `ci` is valued as `github`, then a codecov step will be added inside appropriate test job.
- `openapi_version`: one of `v2, v3`. When provided and `no_api` is either not provided or `false`, then it will generate the appropriate API layer.
  - Note that `v3` is not implemented.
- `port`: exposed port in docker images (injected as environment variable `BINARY_PORT` in helm values).

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
├── .gitlab (if `ci` is "gitlab")
│   ├── workflows
│   │   ├── .gitlab-ci.yml
├── .github (if `ci` is "github")
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
├── .gitlab (if `ci` is "gitlab")
│   ├── workflows
│   │   ├── .gitlab-ci.yml
├── .github
│   ├── workflows (if `ci` is "github")
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
├── internal (internal code)
├── pkg (public libraries)
├── testdata (tests files - expected generations, etc.)
├── .craft (craft configuration file)
├── .gitignore
├── .gitlab-ci.yml (GitLab CI/CD file)
├── .golangci.yml (golangci-lint configuration file)
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── sonar.properties (sonar analysis properties)
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
└── launcher.sh (only when there's at least one main.go in cmd folder, parses the BINARY_NAME environment variable to run the right executable)
```

### Helm plugin

The helm plugin is in charge of generating the helm chart for the project. Depending on implemented coding languages, `values.yaml` file will contain values for `cronjobs`, `jobs` or `workers`.

For instance, associated with [golang plugin](#golang-plugin), kubernetes executables will be parsed from `cmd` folder.

The following layout will be created:

```tree
├── chart (folder with helm chart templates and values)
│   ├── .craft (override values for helm chart)
└── .craft (craft configuration file)
```

### License plugin

The license plugin is only in charge of retrieving the appropriate `LICENSE` file according to `.craft > license` value.

## Examples

You may consult the `examples` for more information and details on generated files.