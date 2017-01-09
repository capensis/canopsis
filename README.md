# Alarms Canopsis Brick

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
    $ brickmanager install brick-alarms

Then, you need to enable the brick

    $ brickmanager enable brick-alarms

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'brick-alarms'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : ee3940e.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS



### FIXMES



## Additional-info

Minified version : 4 files (size: 24K)
Development version : 5 files (size: 40K)
