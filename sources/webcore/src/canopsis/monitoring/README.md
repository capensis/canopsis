# Monitoring Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Monitoring-related features for Canopsis

## Content



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install monitoring

Then, you need to enable the brick

    $ brickmanager enable monitoring

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'monitoring'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/monitoring/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : 14d6578.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/mixins/sendevent.js | use an adapter for this |
| src/mixins/sendevent.js | refactor into sub classes |
| src/widgets/weather/controller.js | avoid using 0 as limit. A better practivce should be used, like limiting to 1000 and display a warning if payload.length > 1000 |


### FIXMES

| File   | Note   |
|--------|--------|
| src/components/cfiltereditor/component.js | Canopsis object is not accessible anymore |


## Additional-info

Minified version : 4 files (size: 120K)
Development version : 48 files (size: 360K)
