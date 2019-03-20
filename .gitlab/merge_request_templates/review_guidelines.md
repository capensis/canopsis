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

- [ ] The commit messages must respect the `Commit Message Guidelines` of angular.
- [ ] The merge request respects the specification or fixes the bug
- [ ] The additions are correctly documented (classes, methods, package, …) and if required, the current documentation is updated
- [ ] The documentation about the feature/component is created/updated
- [ ] A new entry in the canopsis/canopsis `CHANGELOG` and `notes de version` is added
- [ ] The unit tests associated must be created/updated and they must pass
- [ ] If the current `Merge Request` is a bugfix, a unit test must be written to check if the bug is correctly fixed
- [ ] The associated issue must be linked to the current merge request
- [ ] The process to test the current `Merge Request` is written.
- [ ] The `Merge Request` must be reviewed
- [ ] Python:
  - [ ] The sources must be PEP 8 and PEP 257 compliant
  - [ ] The docstring documentation must follow [sphinx](https://www.sphinx-doc.org/en/master/usage/restructuredtext/domains.html#info-field-lists) format
- [ ] Golang:
  - [ ] The import must be cleaned with goimport
- [ ] Javascript:
  - [ ] The sources must respect the `Airbnb JavaScript Style Guide`


## `How to test` template

### Python
1. Install canopsis X.Y.Z with the debian/centos package
2. Log with canopsis user :`su - canopsis`
3. Install the `Merge Request` with pip install -U `/path/to/canopsis_project/canopsis/sources/canopsis/`
4. (Re)start :`canoctl restart` as root
5. Copy the unit tests in the canopsis environment : `cp /path/to/canopsis_project/canopsis/sources/canopsis/test/* /opt/canopsis/var/lib/canopsis/unittest/canopsis -r`
6. Run `ut_runner` as canopsis
7. Go in :`cd /path/to/canopsis_project/canopsis/sources/functional_testing/` and run `python2 runner.py`

### Golang

### Javascript
