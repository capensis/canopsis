# Timeline Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

brick timeline

## Content



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis 
    $ cd var/www
    $ ./bin/brickmanager install brick-timeline

Then, you need to enable the brick

    $ ./bin/brickmanager enable brick-timeline

You can see enabled bricks

    $ su - canopsis
    $ cd var/www
    $ ./bin/brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'brick-timeline'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/brick-timeline/blob/master/doc/index.rst)

## Continuous-Integration

### Tests

The last build was not a full build. Please use the "full-compile" npm script to make test results show up here.

### Lint

Tested on commit : [ERROR : The brick is not in a dedicated git repository].

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :ok: OK |  |


## Code-Notes

### TODOS



### FIXMES



## Additional-info

Minified version :  files (size: 4,0K)
Development version :  files (size: 4,0K)
