## Unreleased (1c1176c..1024fe5)
#### Continuous Integration
- deploy tag v2022.07.19.1257 to adev - (1024fe5) - Vincent Serpoul
#### Features
- inmemory repo (#107) - (58960d0) - claire-tjq
- use requestlogger (#109) - (018c2e6) - Hans Liu
- upgrade handlerwrap v2 to v3 (#108) - (543aa63) - jefferson-crypto
- golangci lint 1.47 (#105) - (1c1176c) - vincentscdc
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.19.1004 - (737c6d9) - Vincent Serpoul

- - -

## Unreleased (d5b187d..104d5aa)
#### Bug Fixes
- add required context in order to have circle oidc token being generated (#103) - (d5b187d) - khongwooilee
#### Continuous Integration
- deploy tag v2022.07.19.1004 to adev - (104d5aa) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.18.1428 - (f5d4d9a) - Vincent Serpoul

- - -

## Unreleased (44b2960..f90435d)
#### Bug Fixes
- add main to push release - (1945e5f) - Vincent Serpoul
- add missing type declaration in circleci (#98) - (44b2960) - khongwooilee
#### Continuous Integration
- deploy tag v2022.07.18.1428 to adev - (f90435d) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.18.1402 - (618e6b1) - Vincent Serpoul

- - -

## Unreleased (28f1377..d4479db)
#### Bug Fixes
- switch to errors.As  (#97) - (28f1377) - BennTayCDC
#### Continuous Integration
- deploy tag v2022.07.18.1402 to adev - (d4479db) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.18.1024 - (cafc4d0) - Vincent Serpoul
#### Refactoring
- use project env var instead of context (#95) - (9072099) - khongwooilee
#### Style
- use the new api standards (#96) - (3e8cf7a) - claire-tjq

- - -

## Unreleased (f08627d..be37901)
#### Continuous Integration
- deploy tag v2022.07.18.1024 to adev - (be37901) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.15.1100 - (ba26c56) - Vincent Serpoul
- update trivy version - (2a059b4) - Vincent Serpoul
- update deps - (f08627d) - Vincent Serpoul

- - -

## Unreleased (184c1a5..4ab75ef)
#### Bug Fixes
- deployment 1.22, remove annotations from yaml - (184c1a5) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.15.1100 to adev - (4ab75ef) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.15.1050 - (51d89a5) - Vincent Serpoul

- - -

## Unreleased (66b8ae0..8492e01)
#### Bug Fixes
- deployment 1.22 - (66b8ae0) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.15.1050 to adev - (8492e01) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.15.1029 - (bc0535e) - Vincent Serpoul

- - -

## Unreleased (5ae1b94..1d3a19a)
#### Bug Fixes
- namespace fix in circleci - (5ae1b94) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.15.1029 to adev - (1d3a19a) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.15.1020 - (6b6f81a) - Vincent Serpoul

- - -

## Unreleased (0da4caa..9038e7c)
#### Bug Fixes
- deployment folder yaml for dev to 1.22 - (0da4caa) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.15.1020 to adev - (9038e7c) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.15.1019 - (2f9155c) - Vincent Serpoul

- - -

## Unreleased (dd9126f..807cd7a)
#### Bug Fixes
- migration to cluster EKS 1.22 dev - (a32f82b) - Vincent Serpoul
- go version - (dd9126f) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.15.1019 to adev - (807cd7a) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.13.2131 - (387bebb) - Vincent Serpoul

- - -

## Unreleased (4c2e2ec..42e87b8)
#### Continuous Integration
- deploy tag v2022.07.13.2131 to adev - (42e87b8) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.12.1312 - (ba762b6) - Vincent Serpoul
#### Refactoring
- rename port to transport (#94) - (35acc97) - BennTayCDC
- general improvements (#91) - (4c2e2ec) - vincentscdc
#### Tests
- **(sqlc)** remove mock, add checks on returned data (#90) - (8628acd) - claire-tjq

- - -

## Unreleased (83b9d7a..c5a42c2)
#### Bug Fixes
- spacing - (fe2334e) - Vincent Serpoul
- remove xz-utils version and make the two dockerfiles closer to eachother - (a0eabe0) - Vincent Serpoul
- use correct deployment env var (#88) - (83b9d7a) - claire-tjq
#### Continuous Integration
- deploy tag v2022.07.12.1312 to adev - (c5a42c2) - Vincent Serpoul
#### Features
- split upx in a different compressor step - (f1d3adb) - Vincent Serpoul
- implement usage of golang-common/database/pginit (#86) - (791af4d) - BennTayCDC
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.07.08.1035 - (1644811) - Vincent Serpoul

- - -

## Unreleased (21c8d26..a0f39c7)
#### Bug Fixes
- update to handlerwrap/v2 (#85) - (dd64f25) - vincentscdc
#### Continuous Integration
- deploy tag v2022.07.08.1035 to adev - (a0f39c7) - Vincent Serpoul
#### Features
- **(dockertest)** add dockertest in ci for unit tests with test db (#56) - (946e0b2) - jonathanyang-cryptocom
- read github username and personal access key from local file (#82) - (123ab59) - kevinchiutw
- use ROLE_ARN in deployment (#67) - (8c6f038) - lttsai
- add otel into grpc (#61) - (21c8d26) - stanley hsieh
#### Miscellaneous Chores
- **(deps)** upgrade deps - (360c887) - Vincent Serpoul
- update CHANGELOG.md for v2022.07.04.1640 - (318d34e) - Vincent Serpoul
#### Refactoring
- **(ci)** parameterize app name (#76) - (9595f7a) - BennTayCDC
- Use useruuidmiddleware in go common (#80) - (c89e219) - gaston-chiu
- use gofrs/uuid in sqlc gen and payment repo (#79) - (63989a0) - jonathanyang-cryptocom
- Clean up unused package (#66) - (fe91f77) - gaston-chiu

- - -

## Unreleased (7c97bbc..d51ec8b)
#### Bug Fixes
- use APP_ENV in local deployment - (6bcc809) - Vincent Serpoul
- remove unnessecary circleci keys - (468fd38) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.07.04.1640 to adev - (d51ec8b) - Vincent Serpoul
#### Documentation
- add circleci and coveralls status badges (#57) - (fa93d76) - Hans Liu
#### Features
- **(repo)** add repositories (#43) - (43c9582) - jonathanyang-cryptocom
- Add controller and service test coverage (#53) - (4c5eee4) - gaston-chiu
- Add buf linter into CI pipeline (#58) - (d5bf644) - gaston-chiu
#### Miscellaneous Chores
- **(localdev)** add port forward in the deploy makefile - (46998ab) - Vincent Serpoul
- update CHANGELOG.md for v2022.06.27.1049 - (4e12aa5) - Vincent Serpoul
#### Refactoring
- **(main)** split into multiple functions (#72) - (25093c5) - BennTayCDC
- Clean up codebase (#47) - (9467ca5) - gaston-chiu
- Refactor error handling (#42) - (7c97bbc) - gaston-chiu

- - -

## Unreleased (9aff9ee..7e5ccfd)
#### Bug Fixes
- linkerd injection - (9aff9ee) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.06.27.1049 to adev - (7e5ccfd) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.27.1033 - (5b9f0d2) - Vincent Serpoul

- - -

## Unreleased (af42429..38c3731)
#### Bug Fixes
- Add payment message test (#44) - (af42429) - gaston-chiu
#### Continuous Integration
- deploy tag v2022.06.27.1033 to adev - (38c3731) - Vincent Serpoul
#### Features
- add linkerd annotation - (cd54beb) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.24.1714 - (3033904) - Vincent Serpoul

- - -

## Unreleased (0b82104..586def8)
#### Bug Fixes
- invert ns and svc - (0b82104) - Vincent Serpoul
#### Continuous Integration
- deploy tag v2022.06.24.1714 to adev - (586def8) - Vincent Serpoul
#### Miscellaneous Chores
- update CHANGELOG.md for v2022.06.24.1640 - (ddb584f) - Vincent Serpoul

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

