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

![View0](https://git.canopsis.net/canopsis-ui-bricks/brick-timeline/raw/master/doc/preview/readme.png)

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



### Lint

Tested on commit : bf96080.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR | <br>/home/illusive/git/canopsis/var/www/src/canopsis/brick-timeline/src/components/timeline/component.js<br>   32:25   error  Expected indentation of 12 space characters but found 24  indent<br>  164:34   error  'v' is not defined                                        no-undef<br>  165:37   error  'v' is not defined                                        no-undef<br>  166:58   error  'v' is not defined                                        no-undef<br>  169:103  error  'v' is not defined                                        no-undef<br><br>✖ 5 problems (5 errors, 0 warnings)<br><br> |


## Code-Notes

### TODOS



### FIXMES



## Additional-info

Minified version : 4 files (size: 28K)
Development version : 3 files (size: 28K)
