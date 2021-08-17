# Developers: How to Make a Release

- [ ] Decide what the new version should be. In this example, __`v1.11.16[-stable]`__ will be used.
- [ ] `git checkout master`
- [ ] `make lint` and `make test` are passing on master. :white_check_mark:
  > This is important because the artifacts to be included with the release will be generated
  by the CI workflows. If linting or tests fail, the workflows will be interrupted
  and artifacts will not be generated.
- [ ] `git checkout release/v1.11.16`
- [ ] Edit `params/version.go` making the necessary changes to version information. (To `-stable` version.) _Gotcha:_ make sure this passes linting, too.
- [ ] `git commit -m "bump version from v1.11.16-unstable to v1.11.16-stable"`
- [ ] `git tag -a v1.11.16`
- [ ] `git push origin v1.11.16`
  > Push the tag to the remote. I like to do it this way because it triggers the tagged version on CI before the branch/PR version,
  expediting artifact delivery.
- [ ] Edit `params/version.go` making the necessary changes to version information. (To `-unstable` version.)
- [ ] `git commit -m "bump version from v1.11.16-stable to v1.11.17-unstable"`
- [ ] `git push origin`
  > Push the branch. This will get PR'd, eg. https://github.com/etclabscore/core-geth/pull/197
- [ ] Draft a new release, following the existing patterns for naming and notes. https://github.com/ubiq/go-ubiq/releases/new
    - Define the tag the release should be associated with (eg `v1.11.16`).
