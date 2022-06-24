## Unreleased (0b82104..0b82104)
#### Bug Fixes
- invert ns and svc - (0b82104) - Vincent Serpoul

- - -

## Unreleased (bfe7d19..0c59276)
#### Bug Fixes
- change APP_ENVIRONMENT to APP_ENV to match monaco-k8s - (a9cb2c3) - Vincent Serpoul
- update liveness and readiness initialDelaySeconds to decrease deploy time - (bfe7d19) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.06.24.1640 to adev - (0c59276) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.24.1632 - (3100868) - Vincent Serpoul

- - -

## Unreleased (61ef976..7e2a0f1)
#### Bug Fixes
- switch to otel collector - (61ef976) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.06.24.1632 to adev - (7e2a0f1) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.24.1624 - (ee927a9) - Vincent Serpoul

- - -

## Unreleased (6aac388..5fa283e)
#### Bug Fixes
- wait for image building before deployment (#34) - (07d58cd) - Yepeng Liang
#### Continuous Integration
- deploy tag v2022.06.24.1624 to adev - (5fa283e) - Vincent Serpoul
#### Features
- Add payment APIs (#19) - (84a06d8) - gaston-chiu
- Add payment plan repo and implementation (#27) - (6aac388) - jonathanyang-cryptocom
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.22.1647 - (16773ed) - Vincent Serpoul
- upgrade deps - (9a617e0) - Vincent Serpoul
- default collector set to the adev collector - (6cf5148) - Vincent Serpoul

- - -

## Unreleased (220c061..2ed0789)
#### Bug Fixes
- endpoint change to 204 (#31) - (220c061) - vincentscdc
#### Continuous Integration
- deploy tag v2022.06.22.1647 to adev - (2ed0789) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.22.1551 - (b6abe58) - Vincent Serpoul

- - -

## Unreleased (7d891eb..e415250)
#### Bug Fixes
- change health endpoint according to monaco-k8s spec (#30) - (1184f61) - vincentscdc
- update dev deployment and tagging in one go (#29) - (5c5a5a9) - Yepeng Liang
#### Continuous Integration
- deploy tag v2022.06.22.1551 to adev - (e415250) - Vincent Serpoul
- deploy tag v2022.06.22.1412 to adev - (2155685) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.22.1412 - (0f3042a) - Vincent Serpoul
- update CHANGELOG.md for v2022.06.22.1412 - (7d891eb) - Vincent Serpoul

- - -

## Unreleased (4c4d20d..aef1720)
#### Bug Fixes
- base local image name - (aef1720) - Vincent Serpoul
- improve localdev environment - (89c2249) - Vincent Serpoul
- rename module name and fix compile errors (#9) - (c04dee3) - vincentscdc
#### Features
- add production deployment pipeline (#24) - (00012e1) - khongwooilee
- Integrate Spectral into CI for API spec linting when spec is changed (#13) - (1b1b000) - jonathanyang-cryptocom
- CircleCI and deploy files (#7) - (c93ce29) - Yepeng Liang
- Use buf to manage code gen (#8) - (0f638ab) - gaston-chiu
- add grpc (#6) - (9e4e8c7) - gaston-chiu
- improve message when files not gofumpt-ed - (a8a88b6) - Vincent Serpoul
- initialize code based on bnpl (#2) - (4c4d20d) - Crypto-Kyle-Yu
#### Miscellaneous Chores
- change APP_NAME to golang-reference-api - (882dee8) - Vincent Serpoul
- update go version - (c934f65) - Vincent Serpoul

