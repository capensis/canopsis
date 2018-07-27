# Canopsis-rights Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Rights and permission management

## Content

### components

 - right-checksum
 - right-action
 - rightselector

### functions

 - updateRecord
 - beforeModel
 - beforeModel
 - beforeModel
 - afterModel

### events

 - toggleEditMode
Handle rights management when toggling edit mode.



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install canopsis-rights

Then, you need to enable the brick

    $ brickmanager enable canopsis-rights

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'canopsis-rights'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/canopsis-rights/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : 606530e.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/reopens/adapters/userview.js | Add the correct right to the current user, to allow him to display the view |
| src/reopens/adapters/userview.js | the right already exists, it's an update |
| src/reopens/adapters/userview.js | replace the userview name if it has changed |
| src/components/right-checksum/component.js | not used anymore? check and delete this property if possible |


### FIXMES

| File   | Note   |
|--------|--------|
| src/reopens/routes/authenticated.js | use store#adapterFor |
| src/components/rights-action/component.js | don't use _data, it might lead to unpredictable behaviours! |
| src/components/right-checksum/component.js | don't use "_data"! |


## Additional-info

Minified version : 4 files (size: 56K)
Development version : 28 files (size: 204K)
