# Merge Request Canopsis

## Description

<message>

## Requirements.

```
Insert the requirements here (canopsis version, docker/packages/, GNU/Linux distribution, …)

```

## How to test

<the minimal and required steps to test and validate the current merge request>

These previous steps are just the minimal and required steps to test the related work, feel free to make more tests.

## Validation criteria

- [ ] The commit messages must respec the `Commit Message Guidelines` of angular.
- [ ] The merge request respect the specification or fix the bug
- [ ] The additions are correctly documented (classes, methods, package, …) and if required, the current documentations are update
- [ ] The documentations about the feature/component is created/updated
- [ ] A new entry in the canopsis/canopsis `CHANGELOG` and `notes de version` is added
- [ ] The unit tests associated must be created/updated and they must pass
- [ ] The associated issue must be linked to the current
- [ ] The process to test the current `Merge Request` is written.
- [ ] The `Merge Request` must be reviewed
- [ ] Python:
  - [ ] The sources must be PEP8 compliant
- [ ] Golang:
  - [ ] The sources must be formated with gofmt
  - [ ] The import must be cleaned with goimport
- [ ] Javascript:
  - [ ] The sources must respect the `Airbnb JavaScript Style Guide`


## `How to test` template

### Python
1. Install canopsis X.Y.Z with the debian/centos package
2. log with canopsis user :`su - canopsis`
3. Install the `Merge Request` with pip install -U `/path/to/canopsis_project/canopsis/sources/canopsis/`
4. (Re)start :`canoctl restart` as root
5. Copy the unit tests in the canopsis environment : `cp /path/to/canopsis_project/canopsis/sources/canopsis/test/* /opt/canopsis/var/lib/canopsis/unittest/canopsis -r`
6. Run `ut_runner` as canopsis
7. Go in :`cd /path/to/canopsis_project/canopsis/sources/functional_testing/` and run `python2 runner.py`

### Golang

### Javascript
