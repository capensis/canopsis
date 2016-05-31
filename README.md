# Querybuilder Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Query builder editor

## Content



## Screenshots

![View0](https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder/raw/master/doc/preview/querybuilder1.png)![View1](https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder/raw/master/doc/preview/querybuilder2.png)

## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis 
    $ cd var/www
    $ ./bin/brickmanager install brick-querybuilder

Then, you need to enable the brick

    $ ./bin/brickmanager enable brick-querybuilder

You can see enabled bricks

    $ su - canopsis
    $ cd var/www
    $ ./bin/brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'brick-querybuilder'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder/blob/master/doc/index.rst)

## Continuous-Integration

### Tests

The last build was not a full build. Please use the "full-compile" npm script to make test results show up here.

### Lint

Tested on commit : 96ac821.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/components/querybuilder/component.js | activate this when events will be triggered at rule drop, probably on version 2.3.1 of querybuilder |


### FIXMES



## Additional-info

Minified version : 5 files (size: 40K)
Development version : 4 files (size: 44K)
