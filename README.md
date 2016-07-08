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

Tested on commit : cf2988a.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR | <br>/home/tristan/git/canopsis/var/www/src/canopsis/brick-timeline/src/components/timeline/component.js<br>  24:2  error  Mixed spaces and tabs                                   no-mixed-spaces-and-tabs<br>  24:6  error  Expected indentation of 8 space characters but found 4  indent<br>  25:2  error  Mixed spaces and tabs                                   no-mixed-spaces-and-tabs<br>  25:8  error  Missing semicolon                                       semi<br>  26:2  error  Mixed spaces and tabs                                   no-mixed-spaces-and-tabs<br>  27:2  error  Mixed spaces and tabs                                   no-mixed-spaces-and-tabs<br>  27:6  error  Expected indentation of 8 space characters but found 4  indent<br><br>✖ 7 problems (7 errors, 0 warnings)<br><br> |


## Code-Notes

### TODOS



### FIXMES



## Additional-info

Minified version :  files (size: 4,0K)
Development version : 3 files (size: 24K)
