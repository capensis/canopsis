# Canopsis-ui Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Provides objects and templates that are not technically required by the UI, but that are always provided with Canopsis to provide some functionnal aspects of a Canopsis application

## Content



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install canopsis-ui

Then, you need to enable the brick

    $ brickmanager enable canopsis-ui

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'canopsis-ui'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/canopsis-ui/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : b49ef43.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/reopens/views/application.js | uncomment while ready |
| src/components/rruleeditor/component.js | move this in utils, somewhere ... |


### FIXMES



## Additional-info

Minified version : 4 files (size: 68K)
Development version : 15 files (size: 156K)
