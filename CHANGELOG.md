## [1.0.0-alpha.3](https://github.com/kilianpaquier/craft/compare/v1.0.0-alpha.2...v1.0.0-alpha.3) (2024-03-06)


### ⚠ BREAKING CHANGES

* **config:** many options moved to substructures, like openapi_version into api.openapi_version or docker_registry into docker.registry

### Features

* **github:** add docker-hadolint and docker-trivy analysis, remove version on integration, update codecov configuration ([f77eb46](https://github.com/kilianpaquier/craft/commit/f77eb46d72bb4ea87651b45df2f382a31020de3e))
* **go:** handle go test with multiple OS ([c257ea3](https://github.com/kilianpaquier/craft/commit/c257ea33ccd26f7e03c3e11e115c8f716e7a40d2))
* **go:** handle golang docker build stage version ([5815d62](https://github.com/kilianpaquier/craft/commit/5815d62556fc3ebfe04cd24c332342424b2cac71))


### Bug Fixes

* **ci:** handle correctly dependabot codecov ignore ([bea2c7c](https://github.com/kilianpaquier/craft/commit/bea2c7c37adcd68694ab8562b0835852f5d6818f))
* **ci:** handle correctly push github actions rules with semantic release branches rules ([ca6f7fe](https://github.com/kilianpaquier/craft/commit/ca6f7fe4c6e004e83f84050074ec0d85a7dae3b6))
* **go:** remove hadolint pull request comment ([f294d46](https://github.com/kilianpaquier/craft/commit/f294d46d1a1a21f97dabfae76f8a7db0b5e12a43))
* handle correctly windows generation ([c3e8573](https://github.com/kilianpaquier/craft/commit/c3e8573be98b52830f34f3eca7f7351a80ce4b77))


### Documentation

* **readme:** add CC of craft commands instead of manul explanation ([867fd60](https://github.com/kilianpaquier/craft/commit/867fd60fd18625b412dcf547d53db8eef0dff02b))


### Code Refactoring

* **config:** rework craft config structure ([83aa207](https://github.com/kilianpaquier/craft/commit/83aa20766cccb0e3041d54fa9b3c6928b1e555ee))

## [1.0.0-alpha.2](https://github.com/kilianpaquier/craft/compare/v1.0.0-alpha.1...v1.0.0-alpha.2) (2024-03-03)


### Features

* add codecov option, rework some things around gitlab-ci cicd integration, rework some things around makefile generation ([a2ba08a](https://github.com/kilianpaquier/craft/commit/a2ba08ac66a56cff4b39cac79732d3e807b42fc1))
* add dependabot configuration for github and remove description from init command ([2f73b05](https://github.com/kilianpaquier/craft/commit/2f73b05ac912bc152d3a7a259596597a6d6b5bf7))
* don't handle multiple primary plugins in the same repository yet ([bcba65d](https://github.com/kilianpaquier/craft/commit/bcba65d62e9c9a16dd124464b6f441503552ec48))
* **golang:** add docker build in github release workflow and improve Dockerfile labeling ([6381cd1](https://github.com/kilianpaquier/craft/commit/6381cd15de54ae52f69269b833db7fc02262980d))


### Bug Fixes

* add issues write on release and add exclusions to go build artifacts ([80eaa29](https://github.com/kilianpaquier/craft/commit/80eaa29f8cd998d0cef025fe287eea60ace5d36a))
* **github:** handle properly codecov options ([1dca49f](https://github.com/kilianpaquier/craft/commit/1dca49f35176128298a0756b83ca9358c209803f))
* **github:** handle properly release branches for docker build ([5ae920c](https://github.com/kilianpaquier/craft/commit/5ae920c827ce23334b50fc83a77dbf6307f427c3))
* **github:** remove codecov on dependabot branches ([fbedf74](https://github.com/kilianpaquier/craft/commit/fbedf743f5bc9a4c1d0a0dab393e49d471cd4ead))
* **github:** remove dependabot weird prefix ([42c5e1b](https://github.com/kilianpaquier/craft/commit/42c5e1bd4c90ee3d82909be21d1d4c01be31973e))
* **gitlab:** update IMAGE_VERSION to VERSION ([c48c74e](https://github.com/kilianpaquier/craft/commit/c48c74e0549d50c4d379be440c86d526a5aaf735))
* **golang:** invalid order in Dockerfile instructions ([726edc5](https://github.com/kilianpaquier/craft/commit/726edc5e5230969139f1582dd448ced8162349be))
* **release:** add v prefix on github workflows version output ([85874f8](https://github.com/kilianpaquier/craft/commit/85874f8f0a5d6735157db61ea70e5a025022f43b))


### Chores

* **codecov:** add mocks coverage exclusion ([47bb177](https://github.com/kilianpaquier/craft/commit/47bb177887d850be8a6d055e4d404aa335e1dbfa))
* **deps:** update go-playground/validator ([91ead6b](https://github.com/kilianpaquier/craft/commit/91ead6b174605c901b11805daa520085920c2cbe))

## 1.0.0-alpha.1 (2024-02-25)


### Features

* add github cicd for generic and golang plugins ([a9e5e2d](https://github.com/kilianpaquier/craft/commit/a9e5e2dcddc1771373bbb53ce61d37a87667cd5b))
* import project from gitlab ([36a4f96](https://github.com/kilianpaquier/craft/commit/36a4f969cb9949b93e3751410347b39dcd3a43d2))


### Bug Fixes

* add back release specific worfklow ([11bda70](https://github.com/kilianpaquier/craft/commit/11bda70f0bd7e0899361ad6dde888716c43811bf))
* bad publish github actions condition ([e0585c8](https://github.com/kilianpaquier/craft/commit/e0585c8a841a7b1144ef26aca3d1a0d8faa70861))
* bad publish github actions condition ([3c7eedc](https://github.com/kilianpaquier/craft/commit/3c7eedc97cfe9b6f06f14e0b66fad16817d63a7e))
* include publish in base integration github actions workflow ([28b7e04](https://github.com/kilianpaquier/craft/commit/28b7e043ac8e0cb83be9b2022fd8565ececedf8d))
* missing examples update after last github actions feature ([c955bbf](https://github.com/kilianpaquier/craft/commit/c955bbfdf628382a9836b99fedd72a5b8a7bac91))
* remove ref_protected condition on publish job (environment constraint should do the job) ([bf664e0](https://github.com/kilianpaquier/craft/commit/bf664e03ce3958209875bcf5d6916a447344b931))
* remove release github actions useless strategy ([e60e901](https://github.com/kilianpaquier/craft/commit/e60e901c72bbdae4d19b4c2acc73e248ca7c4f84))
* try match github actions condition on ref_protected ([24b7603](https://github.com/kilianpaquier/craft/commit/24b76036141cd80b054f274f853c90a264148e5c))


### Documentation

* add newly github actions workflows doc ([f59abe5](https://github.com/kilianpaquier/craft/commit/f59abe500f66d571af07dc514ed12c70de71792b))


### Chores

* **release:** update release environment name ([ab50051](https://github.com/kilianpaquier/craft/commit/ab500519a99aa225b6b7a4f2da2683753937f37f))
