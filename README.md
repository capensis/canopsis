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

Tested on commit : 741f946.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/mixins/customsendevent.js | use an adapter for this |
| src/mixins/customsendevent.js | refactor into sub classes |


### FIXMES



## Additional-info

Minified version : 4 files (size: 112K)
Development version : 54 files (size: 396K)
