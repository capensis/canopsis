# Listalarm Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Alarms widget for Canopsis

## Content



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install brick-listalarm

Then, you need to enable the brick

    $ brickmanager enable brick-listalarm

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'brick-listalarm'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/brick-listalarm/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : 8721e9c.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/forms/snooze/controller.js | search this value into schema |
| src/forms/snooze/controller.js | refactor the 20 lines below in an utility function "getEditorForAttr" |
| src/forms/snooze/controller.js | use the real schema, not the dict used to create it |


### FIXMES



## Additional-info

Minified version : 4 files (size: 136K)
Development version : 75 files (size: 508K)
