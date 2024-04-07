## [1.0.0-beta.2](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.1...v1.0.0-beta.2) (2024-04-07)


### Features

* **generate:** add craft notice in .craft configuration file ([72382bd](https://github.com/kilianpaquier/craft/commit/72382bd5dbf23404637208bb902b42a2d6395c65))
* **generate:** save edited .craft configuration from generate ([a03b11e](https://github.com/kilianpaquier/craft/commit/a03b11e10dc5b90f63e0229cfe0775e75ae98fa2))


### Bug Fixes

* **generic:** remove readme github actions since there's no integration workflow ([26e3a4f](https://github.com/kilianpaquier/craft/commit/26e3a4fcd4c5522d9c1ea06440a765a5ade74a4b))
* **github:** add branches for integration actions workflow (next, beta, alpha) ([6285f1e](https://github.com/kilianpaquier/craft/commit/6285f1e84642fbc6a6985e1c692d3f369b3d8079))
* **github:** bad path on npm-build and hugo-build artifacts for github-pages ([74615f9](https://github.com/kilianpaquier/craft/commit/74615f9010192ee2a63d780ca1bed63d4a81d781))
* **github:** let github actions set the github-pages job's token ([ba97741](https://github.com/kilianpaquier/craft/commit/ba97741c149c24e9583ac109fd321840aa98a42c))
* **hugo:** put specific exclusions of hugo files in gitignore to avoid ignoring hugo configuration files ([9a1d6d7](https://github.com/kilianpaquier/craft/commit/9a1d6d7a24f29c08a65f01edd552b93ba7649cf7))


### Chores

* **github:** update cache action for hugo build ([2119229](https://github.com/kilianpaquier/craft/commit/21192298af59215fce4298d3f49d56e0cf810558))


### Code Refactoring

* **generate:** drop plugin notion in favor of something more flexible to handle languages frameworks ([1f757b8](https://github.com/kilianpaquier/craft/commit/1f757b88f2a7e0d04694ab4370e1bee725f58449))

## 1.0.0-beta.1 (2024-03-31)


### ⚠ BREAKING CHANGES

* **config:** many options moved to substructures, like openapi_version into api.openapi_version or docker_registry into docker.registry

### Features

* add codecov option, rework some things around gitlab-ci cicd integration, rework some things around makefile generation ([a2ba08a](https://github.com/kilianpaquier/craft/commit/a2ba08ac66a56cff4b39cac79732d3e807b42fc1))
* add dependabot configuration for github and remove description from init command ([2f73b05](https://github.com/kilianpaquier/craft/commit/2f73b05ac912bc152d3a7a259596597a6d6b5bf7))
* add github cicd for generic and golang plugins ([a9e5e2d](https://github.com/kilianpaquier/craft/commit/a9e5e2dcddc1771373bbb53ce61d37a87667cd5b))
* add some features around nodejs and golang (includes some refactor around generate engine) ([3c98be8](https://github.com/kilianpaquier/craft/commit/3c98be83d83ff727e1e8fb54f03d4402002f6ac5))
* **ci:** handle only vN as maintenance branch - [#17](https://github.com/kilianpaquier/craft/issues/17) ([e1ef4b7](https://github.com/kilianpaquier/craft/commit/e1ef4b705358add306db01e8e373e9d107cab1ab))
* **ci:** handle properly workflow result even on push events ([fe6af3b](https://github.com/kilianpaquier/craft/commit/fe6af3b0ee6d6fb620b216ce1f56d8bfc4f40ea6))
* **docker:** push snapshot images to specific registry and review codecov projects' configuration ([9b3bd1a](https://github.com/kilianpaquier/craft/commit/9b3bd1a8c4a37fafb926b02975d95e3cac0ce9b2))
* don't handle multiple primary plugins in the same repository yet ([bcba65d](https://github.com/kilianpaquier/craft/commit/bcba65d62e9c9a16dd124464b6f441503552ec48))
* **github:** add docker-hadolint and docker-trivy analysis, remove version on integration, update codecov configuration ([f77eb46](https://github.com/kilianpaquier/craft/commit/f77eb46d72bb4ea87651b45df2f382a31020de3e))
* **gitlab:** handle cicd generation with to-be-continuous - [#18](https://github.com/kilianpaquier/craft/issues/18) ([02f1a97](https://github.com/kilianpaquier/craft/commit/02f1a97621674c71d18a544ea0573e25caaecff2))
* **go:** handle go test with multiple OS ([c257ea3](https://github.com/kilianpaquier/craft/commit/c257ea33ccd26f7e03c3e11e115c8f716e7a40d2))
* **go:** handle golang docker build stage version ([5815d62](https://github.com/kilianpaquier/craft/commit/5815d62556fc3ebfe04cd24c332342424b2cac71))
* **golang:** add docker build in github release workflow and improve Dockerfile labeling ([6381cd1](https://github.com/kilianpaquier/craft/commit/6381cd15de54ae52f69269b833db7fc02262980d))
* import project from gitlab ([36a4f96](https://github.com/kilianpaquier/craft/commit/36a4f969cb9949b93e3751410347b39dcd3a43d2))
* introduce languages property to share more generated files between plugins ([81307f0](https://github.com/kilianpaquier/craft/commit/81307f0caa45ef630b5acc98fcca71d629266a5a))
* **nodejs:** add generation ([742b57c](https://github.com/kilianpaquier/craft/commit/742b57c985cd4597e64a508d8f1ab03af5a0c54b))
* **readme:** add various badges depending on git platform ([82ce9d0](https://github.com/kilianpaquier/craft/commit/82ce9d0155f96be6e9ad5a1697c6d670954ae931))


### Bug Fixes

* add back release specific worfklow ([11bda70](https://github.com/kilianpaquier/craft/commit/11bda70f0bd7e0899361ad6dde888716c43811bf))
* add issues write on release and add exclusions to go build artifacts ([80eaa29](https://github.com/kilianpaquier/craft/commit/80eaa29f8cd998d0cef025fe287eea60ace5d36a))
* bad publish github actions condition ([e0585c8](https://github.com/kilianpaquier/craft/commit/e0585c8a841a7b1144ef26aca3d1a0d8faa70861))
* bad publish github actions condition ([3c7eedc](https://github.com/kilianpaquier/craft/commit/3c7eedc97cfe9b6f06f14e0b66fad16817d63a7e))
* **ci:** add sufficient rights for semantic release comments on issues and pull requests ([c44c717](https://github.com/kilianpaquier/craft/commit/c44c7173b5ddde883b04208b9c76069b0c67fa19))
* **ci:** change os matrix to include to avoid subnames usage ([e75946a](https://github.com/kilianpaquier/craft/commit/e75946a00a3b0ba554b3cbcf5b3e3d0bb6eff50a))
* **ci:** codecov config in subdir really doesn't work ([649583d](https://github.com/kilianpaquier/craft/commit/649583d92130ee3723e8bdc8e8cd1b66c8b73cf3))
* **ci:** ensure releaserc only have necessary artifacts on git plugin for semantic release ([84ca91b](https://github.com/kilianpaquier/craft/commit/84ca91b355cda4d69f1a9112353864d34c1e433d))
* **ci:** executables not being uploaded in release - [#14](https://github.com/kilianpaquier/craft/issues/14) ([dcb3e77](https://github.com/kilianpaquier/craft/commit/dcb3e7760a11dc3a29c4cea0be544f90dbe33f84))
* **ci:** handle correctly dependabot codecov ignore ([bea2c7c](https://github.com/kilianpaquier/craft/commit/bea2c7c37adcd68694ab8562b0835852f5d6818f))
* **ci:** handle correctly push github actions rules with semantic release branches rules ([ca6f7fe](https://github.com/kilianpaquier/craft/commit/ca6f7fe4c6e004e83f84050074ec0d85a7dae3b6))
* **ci:** multiple improvements around github workflow runs ([09924b4](https://github.com/kilianpaquier/craft/commit/09924b4e952f0aab55f513e3c10d5e3ebbcbe97c))
* **github:** handle properly codecov options ([1dca49f](https://github.com/kilianpaquier/craft/commit/1dca49f35176128298a0756b83ca9358c209803f))
* **github:** handle properly release branches for docker build ([5ae920c](https://github.com/kilianpaquier/craft/commit/5ae920c827ce23334b50fc83a77dbf6307f427c3))
* **github:** remove codecov on dependabot branches ([fbedf74](https://github.com/kilianpaquier/craft/commit/fbedf743f5bc9a4c1d0a0dab393e49d471cd4ead))
* **github:** remove dependabot weird prefix ([42c5e1b](https://github.com/kilianpaquier/craft/commit/42c5e1bd4c90ee3d82909be21d1d4c01be31973e))
* **gitlab:** update IMAGE_VERSION to VERSION ([c48c74e](https://github.com/kilianpaquier/craft/commit/c48c74e0549d50c4d379be440c86d526a5aaf735))
* **golang:** invalid order in Dockerfile instructions ([726edc5](https://github.com/kilianpaquier/craft/commit/726edc5e5230969139f1582dd448ced8162349be))
* **go:** remove hadolint pull request comment ([f294d46](https://github.com/kilianpaquier/craft/commit/f294d46d1a1a21f97dabfae76f8a7db0b5e12a43))
* handle correctly windows generation ([c3e8573](https://github.com/kilianpaquier/craft/commit/c3e8573be98b52830f34f3eca7f7351a80ce4b77))
* include publish in base integration github actions workflow ([28b7e04](https://github.com/kilianpaquier/craft/commit/28b7e043ac8e0cb83be9b2022fd8565ececedf8d))
* missing examples update after last github actions feature ([c955bbf](https://github.com/kilianpaquier/craft/commit/c955bbfdf628382a9836b99fedd72a5b8a7bac91))
* **release:** add v prefix on github workflows version output ([85874f8](https://github.com/kilianpaquier/craft/commit/85874f8f0a5d6735157db61ea70e5a025022f43b))
* remove ref_protected condition on publish job (environment constraint should do the job) ([bf664e0](https://github.com/kilianpaquier/craft/commit/bf664e03ce3958209875bcf5d6916a447344b931))
* remove release github actions useless strategy ([e60e901](https://github.com/kilianpaquier/craft/commit/e60e901c72bbdae4d19b4c2acc73e248ca7c4f84))
* try match github actions condition on ref_protected ([24b7603](https://github.com/kilianpaquier/craft/commit/24b76036141cd80b054f274f853c90a264148e5c))


### Documentation

* add newly github actions workflows doc ([f59abe5](https://github.com/kilianpaquier/craft/commit/f59abe500f66d571af07dc514ed12c70de71792b))
* **readme:** add CC of craft commands instead of manul explanation ([867fd60](https://github.com/kilianpaquier/craft/commit/867fd60fd18625b412dcf547d53db8eef0dff02b))
* **readme:** add schema tips in craft file section ([90f425c](https://github.com/kilianpaquier/craft/commit/90f425cb5f03b8d33f0ae094cfd8041a4fbbde1c))
* **readme:** remove code section language to avoid weird colors ([a8f0f91](https://github.com/kilianpaquier/craft/commit/a8f0f913a27643182174e294bcfafbfe1622d806))
* **readme:** update plugins part ([3028000](https://github.com/kilianpaquier/craft/commit/30280000a55219c5583ee78aa9dac1e2cfd28d41))
* **schema:** add root craft schema for easier setup ([2e6d5c5](https://github.com/kilianpaquier/craft/commit/2e6d5c50fdf8c288c004814d9d486ecbf7bbd581))
* **schema:** typo on no_ options ([aae2b6b](https://github.com/kilianpaquier/craft/commit/aae2b6b71b20f3b953379fd30fc6e125f40d65e4))


### Chores

* **ci:** add merge-check ([26641ff](https://github.com/kilianpaquier/craft/commit/26641ff281092ca6d39b2f928b6c51e06913df0a))
* **ci:** remove build/ci directory ([f9e218b](https://github.com/kilianpaquier/craft/commit/f9e218bb7940ac04d412fedd5bf4af99320422d5))
* **ci:** update golangci rules ([e56e047](https://github.com/kilianpaquier/craft/commit/e56e04760e682c8ed560d51a7ffb0d13171d8820))
* **codecov:** add mocks coverage exclusion ([47bb177](https://github.com/kilianpaquier/craft/commit/47bb177887d850be8a6d055e4d404aa335e1dbfa))
* **deps:** update dependencies ([63132ca](https://github.com/kilianpaquier/craft/commit/63132cad9e371dbe13698101d60e5d3b88a87ec9))
* **deps:** update dependencies ([4e8c9a5](https://github.com/kilianpaquier/craft/commit/4e8c9a5bcc98c6f595f9e929a9c7099a02f4c498))
* **deps:** update go-playground/validator ([91ead6b](https://github.com/kilianpaquier/craft/commit/91ead6b174605c901b11805daa520085920c2cbe))
* **release:** update release environment name ([ab50051](https://github.com/kilianpaquier/craft/commit/ab500519a99aa225b6b7a4f2da2683753937f37f))
* **release:** v1.0.0-alpha.1 [skip ci] ([ae7bc77](https://github.com/kilianpaquier/craft/commit/ae7bc77e8c1213c8c4f49d2cb1833327b2b6a357))
* **release:** v1.0.0-alpha.2 [skip ci] ([9d62a52](https://github.com/kilianpaquier/craft/commit/9d62a529b10e035d0c8e6a2a94a23f50568e67db))
* **release:** v1.0.0-alpha.3 [skip ci] ([257069f](https://github.com/kilianpaquier/craft/commit/257069f9e8404cef587415d29a7f8ac908016774))
* **release:** v1.0.0-alpha.4 [skip ci] ([a970772](https://github.com/kilianpaquier/craft/commit/a9707729a12f2f955c63fdab7d31769dd8b551e6)), closes [#14](https://github.com/kilianpaquier/craft/issues/14)


### Code Refactoring

* **config:** rework craft config structure ([83aa207](https://github.com/kilianpaquier/craft/commit/83aa20766cccb0e3041d54fa9b3c6928b1e555ee))

## [1.0.0-alpha.4](https://github.com/kilianpaquier/craft/compare/v1.0.0-alpha.3...v1.0.0-alpha.4) (2024-03-06)


### Bug Fixes

* **ci:** executables not being uploaded in release - [#14](https://github.com/kilianpaquier/craft/issues/14) ([dcb3e77](https://github.com/kilianpaquier/craft/commit/dcb3e7760a11dc3a29c4cea0be544f90dbe33f84))

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
