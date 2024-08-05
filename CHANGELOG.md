## [1.0.0-beta.10](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.9...v1.0.0-beta.10) (2024-08-05)


### Features

* **github:** use package.json packageManager version for bun projects ([28be5cd](https://github.com/kilianpaquier/craft/commit/28be5cd8058dd711e2efce8f34f0f6d07471e535))
* **makefile:** add install.mk for golang and hugo with various installation scripts ([48f9038](https://github.com/kilianpaquier/craft/commit/48f9038c9a36c2f02187c4dd301d984138c3d7b7))
* **netlify:** add github action job and as such option with github actions ([9f3c250](https://github.com/kilianpaquier/craft/commit/9f3c2500b94161ea3d123f4c7b7e58f0e8d6b3b3))
* **sdk:** move and refactor craft to be used also as a SDK - [#45](https://github.com/kilianpaquier/craft/issues/45) ([78b0e4e](https://github.com/kilianpaquier/craft/commit/78b0e4e98f56d15cbf8aa8ba7c61f4e7ee3372ab))


### Bug Fixes

* **go:** rename revive option imports-blocklist ([11f3f44](https://github.com/kilianpaquier/craft/commit/11f3f441a790627573685328da5fb9bf0cee5583))
* **netlify:** add dev folder in gitignore for nodejs and hugo ([1a98fd1](https://github.com/kilianpaquier/craft/commit/1a98fd1e330a9b63a258430a59dcb59d365649ff))
* **netlify:** add netlify.toml file in github actions job ([06c7567](https://github.com/kilianpaquier/craft/commit/06c7567972e881a4dc93dc89c09fff365710dd9f))
* **npm:** add id-token to release for provenance signature ([33ea480](https://github.com/kilianpaquier/craft/commit/33ea480bb3af34bb9c7563f0767a764c07cfe924))


### Chores

* **deps:** bump github.com/samber/lo in the minor-patch group ([f4e0991](https://github.com/kilianpaquier/craft/commit/f4e0991f8f7974057da323b38f36cd2db803a2b6))
* **deps:** bump github.com/samber/lo in the minor-patch group ([f7466b4](https://github.com/kilianpaquier/craft/commit/f7466b4ac0d36d02cd02f5dab5f2d6d233be2275))
* **deps:** bump github.com/xanzy/go-gitlab in the minor-patch group ([1c372be](https://github.com/kilianpaquier/craft/commit/1c372be6396dbc846ac3f753e959e84b5f1e52b0))
* **deps:** bump golang.org/x/mod in the minor-patch group ([b950e15](https://github.com/kilianpaquier/craft/commit/b950e1559297fb241192da2589ae5e113054d519))
* **deps:** upgrade various dependencies ([9917992](https://github.com/kilianpaquier/craft/commit/9917992221119dff45a5e3dceab3806b6d35ea61))
* **generate:** remove SplitSlice unused function ([23d1925](https://github.com/kilianpaquier/craft/commit/23d1925a69b694af2e864786f98d552f66eab967))
* **schema:** add chart schema for craft chart file ([e54f313](https://github.com/kilianpaquier/craft/commit/e54f313c38ac88fede7d4a99b7414df56a2b64d8))

## [1.0.0-beta.9](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.8...v1.0.0-beta.9) (2024-06-24)


### Documentation

* **deps:** upgrade kubernetes in values.yml ([d6999a3](https://github.com/kilianpaquier/craft/commit/d6999a3be04a58d0eaf73dbfd4e833d0de6fdf76))
* **readme:** update some description around ci release options ([d017c47](https://github.com/kilianpaquier/craft/commit/d017c4791ce108864eb506a254ea6eb7d1e1962f))


### Chores

* **dependabot:** setup bot to run every day to avoid batches every week ([eba239f](https://github.com/kilianpaquier/craft/commit/eba239f3b5da11fda97ed76c83b9183790d006d5))
* **deps:** bump docker build push action to v6 ([d17ffa4](https://github.com/kilianpaquier/craft/commit/d17ffa40e28d20aa01d0e8d9dbaea7240067499e))
* **deps:** bump github.com/go-playground/validator/v10 ([76c309b](https://github.com/kilianpaquier/craft/commit/76c309b717f6282ee9f369f5fe081ce0651c05e1))
* **deps:** bump github.com/xanzy/go-gitlab in the minor-patch group ([d8a900f](https://github.com/kilianpaquier/craft/commit/d8a900f94280bd694421b9f229881100d23060ea))
* **deps:** bump setup-bun github actions to v2 ([bf7666b](https://github.com/kilianpaquier/craft/commit/bf7666b28eaed271078295ad4d28ed7969466a5d))
* **deps:** only use major version for bun setup in ci ([5235da8](https://github.com/kilianpaquier/craft/commit/5235da82fc3f300423deaea0e37062d3764f09de))
* **deps:** upgrade goreleaser action to v6 and associated schema usage ([8cf2d52](https://github.com/kilianpaquier/craft/commit/8cf2d527e42715be04c05c65f43daa7d80f37f2b))
* **deps:** upgrade multiple dependencies ([8331b0f](https://github.com/kilianpaquier/craft/commit/8331b0f851acfbd286925cac47f8f0333056c40d))
* **deps:** upgrade pnpm version to 9 in ci templates ([355f0b8](https://github.com/kilianpaquier/craft/commit/355f0b83659af1a54370ca559420e51fc6fdf70c))

## [1.0.0-beta.8](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.7...v1.0.0-beta.8) (2024-05-10)


### Bug Fixes

* **config:** bad json tag for newly added disable key in release ([e795500](https://github.com/kilianpaquier/craft/commit/e795500bee5ac94d6a97ecbbbe7c5596a632a5ae))


### Chores

* **ci:** update golangci-lint action to v6 ([0c45037](https://github.com/kilianpaquier/craft/commit/0c45037d6ef5d20d457ccde5b5e49efa4d502b33))
* **deps:** upgrade toolchain to go1.22.3 ([116315f](https://github.com/kilianpaquier/craft/commit/116315f0036f51002d28ba74afccf2493c1fa2c4))
* **gitlab:** update semrel_ref variable name and prof_ref value since semrel_ref is the variable used for semantic-release ([9bb2f61](https://github.com/kilianpaquier/craft/commit/9bb2f61f302289cf0d1ba941553af8b6b1773dd0))


### Code Refactoring

* **gitlab:** rework version.yml in CICD ([977d1ce](https://github.com/kilianpaquier/craft/commit/977d1ce8b45de29c12a31b4fa844b3bbbec714dd))

## [1.0.0-beta.7](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.6...v1.0.0-beta.7) (2024-05-07)


### Features

* **config:** add release mode for github and migration backmerge / auto_release option in new ci release section ([5d0bc67](https://github.com/kilianpaquier/craft/commit/5d0bc679df78c3afbccba695f1ecfb3ccb810930))
* **config:** add release option to disable release at all ([4f4051e](https://github.com/kilianpaquier/craft/commit/4f4051e7af888c9b01ad8fcef64be54881117f78))


### Bug Fixes

* **readme:** only add CI badges when CI is provided ([561791a](https://github.com/kilianpaquier/craft/commit/561791a90973e9a7b87ff60d63691ae5357f0fe4))


### Documentation

* **readme:** add linux installation section for no go developers ([8a9ae0f](https://github.com/kilianpaquier/craft/commit/8a9ae0f03171cdd42ea3ff2c5584dd381031224c))


### Chores

* **github:** upgrade pnpm action to v4 and remove version from actions in case it's specified in package.json ([e6ae1af](https://github.com/kilianpaquier/craft/commit/e6ae1af2cd1a9671b5680145592ce6a9f0373ec6))
* **gitlab:** update semrel_ref to include all semantic-release to release ([072a546](https://github.com/kilianpaquier/craft/commit/072a5462f1e85a8f74c699157cfe7818dc9b1f7f))

## [1.0.0-beta.6](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.5...v1.0.0-beta.6) (2024-05-06)


### Bug Fixes

* **ci:** ensure for both gitlab and github that a version is computed and fix associated regexp ([0839a1d](https://github.com/kilianpaquier/craft/commit/0839a1db1fd7a9f9d20cde88121773f0260c2bdc))
* **dependabot:** bad commit prefix for code dependencies ([f328ecd](https://github.com/kilianpaquier/craft/commit/f328ecd4afca2c7cac90fefc7de0d8beede703e8))

## [1.0.0-beta.5](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.4...v1.0.0-beta.5) (2024-05-04)


### Bug Fixes

* **backmerge:** only set fetch-depth: 0 for gitlab and github cicd when option is provided ([a44ef0e](https://github.com/kilianpaquier/craft/commit/a44ef0e41cc0a78e289e70d3483b306a99b3a4c3))
* **gitlab:** bad regexp on version job ([0530a6d](https://github.com/kilianpaquier/craft/commit/0530a6db9c97152f65679acd3b814f3d6fdf2a09))


### Documentation

* **ci:** add information about dot not being espaced in github branches regexp ([35517ed](https://github.com/kilianpaquier/craft/commit/35517ed97575c5a26c807af1cbc2a4b0be0f86df))


### Chores

* **gitlab:** update prod_ref and integ_ref in main .gitlab-ci.yml ([1a58a2b](https://github.com/kilianpaquier/craft/commit/1a58a2b98bf3b1b88f988fefd753e1c245018204))


### Code Refactoring

* **github:** rework release workflow in github actions following latest modifications ([557a6d4](https://github.com/kilianpaquier/craft/commit/557a6d4de299b3531dff9e33abcf6abb7f0b3308))

## [1.0.0-beta.4](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.3...v1.0.0-beta.4) (2024-05-04)


### Features

* **ci:** add backmerge option using @kilianpaquier/semantic-release-backmerge (alpha plugin) ([0e4fea1](https://github.com/kilianpaquier/craft/commit/0e4fea1494119a25637541e0a4015f600041d70b))
* **ci:** handle minor maintenance branch ([4a50859](https://github.com/kilianpaquier/craft/commit/4a5085930b6159e552b38e5d5846b040a08477d8))
* **nodejs:** add NPM_TOKEN when package.json is not private for semantic releasing ([7717b43](https://github.com/kilianpaquier/craft/commit/7717b433bea2bd11109a24b6d9fff5a37271c4c7))
* **nodejs:** add reports folder exclusion ([40ede9a](https://github.com/kilianpaquier/craft/commit/40ede9a726014c50e7bf014713e5abc0a3ba0a9b))
* **nodejs:** handle properly package managers ([97349f9](https://github.com/kilianpaquier/craft/commit/97349f9906af881e1ad08e68a4d0127fa97c405a))


### Bug Fixes

* **generic:** enable auto release and backmerge features for empty languages projects (like readme only, etc.) ([49ed7ba](https://github.com/kilianpaquier/craft/commit/49ed7ba0282ee7f9c525de917658bad2050852ea))
* **github:** release action not having the right conditions nor the right rights for version job ([fb0d770](https://github.com/kilianpaquier/craft/commit/fb0d770969d823be52aa1ccebe59ede5cae6f18c))
* **gitlab:** bad semrel plugin in plugins file ([c657725](https://github.com/kilianpaquier/craft/commit/c657725f935b48728ef8298212dcbb2775992fc2))
* **nodejs:** add built dist to semantic-release job and git depth for backmerge ([cc0fee6](https://github.com/kilianpaquier/craft/commit/cc0fee6f0fb9d175d5cbf241f45dd947c47e02b8))
* **semantic-release:** add conventionalcommits preset for commit-analyzer and fix version to 7 ([5138247](https://github.com/kilianpaquier/craft/commit/51382471cf2cf60ccc0ba3a71d4e07f0349f0fb2))


### Chores

* **backmerge:** allow backmerge for all platforms in releaserc ([1ddf381](https://github.com/kilianpaquier/craft/commit/1ddf381c6672dd19663a02203c49b1bedbbc8d5a))
* **deps:** bump golangci-lint action to v5 in templates ([8ead59d](https://github.com/kilianpaquier/craft/commit/8ead59d5700beaef0a8a254a3d58803c3a47ebce))
* **deps:** upgrade dependencies ([b2c74bf](https://github.com/kilianpaquier/craft/commit/b2c74bf6350b354a1d5c1fdb9baca6e08cc51872))
* **deps:** upgrade generated github semantic-release ci to v23 ([d6d592e](https://github.com/kilianpaquier/craft/commit/d6d592eedb6b304fd7a3cadc316ca3cefb06132f))
* **nodejs:** add fields to parsed package.json ([1d73537](https://github.com/kilianpaquier/craft/commit/1d73537586bd45a3709def3fecb9dc1e9cc80c76))
* **releaserc:** disable issue opening on release error ([1552b50](https://github.com/kilianpaquier/craft/commit/1552b50355e4065a60e567ff125e9cc6ce1d1385))


### Code Refactoring

* **auto_release:** move option into ci.options ([072fe6c](https://github.com/kilianpaquier/craft/commit/072fe6c36cf0e07abb443db3487fc9d58f37e5c5))
* **dependabot:** remove default keys and use ci prefix for github actions updates ([5d8b469](https://github.com/kilianpaquier/craft/commit/5d8b469fe1541ca82a1f52da73333e1a75e960f2))

## [1.0.0-beta.3](https://github.com/kilianpaquier/craft/compare/v1.0.0-beta.2...v1.0.0-beta.3) (2024-04-13)


### ⚠ BREAKING CHANGES

* **generate:** openapi generation was removed

### Features

* **github:** add hugo-build needs for release job ([b64f324](https://github.com/kilianpaquier/craft/commit/b64f3241ce376fbaca620fc0016ed280926fcb74))
* **release:** add auto_release option to auto run release job - available for both gitlab and github ([1e3b2c8](https://github.com/kilianpaquier/craft/commit/1e3b2c8e51b7d1e1ab6b6365d7c6084f519833d0))


### Bug Fixes

* **auto_release:** move option to CI configuration part ([d27d379](https://github.com/kilianpaquier/craft/commit/d27d3792c38f435a6ad4acdcdb057b8937322097))
* **releaserc:** changed `failComment` to default value ([513d939](https://github.com/kilianpaquier/craft/commit/513d939461e4d3d75973d14a8acd1aff7c651f8a))


### Documentation

* **auto_release:** add missing doc ([6dcfe6c](https://github.com/kilianpaquier/craft/commit/6dcfe6cb742742f668d7e2f851555ebb5374f675))


### Code Refactoring

* **generate:** remove openapi generation and rework helm generation to take this into account ([aff7b9c](https://github.com/kilianpaquier/craft/commit/aff7b9cff4aa645cba79bcc0f6a0ccc3fb075998))

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
