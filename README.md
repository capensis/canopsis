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

Tested on commit : 5ed754a.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR | <br>/home/tristan/git/canopsis/var/www/src/canopsis/brick-timeline/src/components/timeline/component.js<br>  25:13  error  '__' is defined but never used       no-unused-vars<br>  26:13  error  'get' is defined but never used      no-unused-vars<br>  28:13  error  'isArray' is defined but never used  no-unused-vars<br>  28:36  error  Missing semicolon                    semi<br>  38:26  error  Strings must use singlequote         quotes<br>  38:60  error  Missing semicolon                    semi<br>  39:49  error  Missing semicolon                    semi<br>  42:42  error  'moment' is not defined              no-undef<br>  43:42  error  'moment' is not defined              no-undef<br>  43:62  error  Strings must use singlequote         quotes<br>  44:43  error  Strings must use singlequote         quotes<br>  44:48  error  Missing semicolon                    semi<br>  47:30  error  Strings must use singlequote         quotes<br>  48:50  error  Strings must use singlequote         quotes<br>  48:60  error  Missing semicolon                    semi<br>  49:52  error  Strings must use singlequote         quotes<br>  51:30  error  Strings must use singlequote         quotes<br>  52:50  error  Strings must use singlequote         quotes<br>  52:61  error  Missing semicolon                    semi<br>  53:52  error  Strings must use singlequote         quotes<br>  55:30  error  Strings must use singlequote         quotes<br>  56:50  error  Strings must use singlequote         quotes<br>  56:61  error  Missing semicolon                    semi<br>  57:52  error  Strings must use singlequote         quotes<br>  59:30  error  Strings must use singlequote         quotes<br>  60:50  error  Strings must use singlequote         quotes<br>  60:60  error  Missing semicolon                    semi<br>  61:52  error  Strings must use singlequote         quotes<br>  63:30  error  Strings must use singlequote         quotes<br>  64:50  error  Strings must use singlequote         quotes<br>  64:60  error  Missing semicolon                    semi<br>  65:52  error  Strings must use singlequote         quotes<br>  67:30  error  Strings must use singlequote         quotes<br>  68:30  error  Strings must use singlequote         quotes<br>  69:50  error  Strings must use singlequote         quotes<br>  69:63  error  Missing semicolon                    semi<br>  71:30  error  Strings must use singlequote         quotes<br>  72:30  error  Strings must use singlequote         quotes<br>  73:30  error  Strings must use singlequote         quotes<br>  74:50  error  Strings must use singlequote         quotes<br>  74:63  error  Missing semicolon                    semi<br>  83:52  error  Strings must use singlequote         quotes<br>  86:52  error  Strings must use singlequote         quotes<br>  89:52  error  Strings must use singlequote         quotes<br>  92:52  error  Strings must use singlequote         quotes<br>  95:52  error  Missing semicolon                    semi<br>  98:11  error  Missing semicolon                    semi<br><br>✖ 47 problems (47 errors, 0 warnings)<br><br> |


## Code-Notes

### TODOS



### FIXMES



## Additional-info

Minified version :  files (size: 4,0K)
Development version : 4 files (size: 32K)
