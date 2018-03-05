# Canopsis-backend-ui-connector Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Provides adapters to communicate with Canopsis backend

## Content

### adapters

 - profile
 - action
 - baseadapter
 - cancel
 - context
 - crecord
 - cservice
 - entitylink
 - eue
 - event
 - eventlog
 - filter
 - linklist
 - loggedaccount
 - pojo
 - Serie2
 - StatsFilter
 - Storage
 - trap
 - userview
 - userviewsimplemodel

### serializers

 - ctx
 - ctxcomponent
 - ctxmetric
 - ctxresource
 - ctxselector
 - ctxtopology
 - job
 - linklist
 - taskmail
 - ticket

### schemas

 - schema-curve
 - schema-rangecolor
 - schema-serie
 - schema-widgetpreferences



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install canopsis-backend-ui-connector

Then, you need to enable the brick

    $ brickmanager enable canopsis-backend-ui-connector

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'canopsis-backend-ui-connector'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/canopsis-backend-ui-connector/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : [ERROR : The brick is not in a dedicated git repository].

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/serializers/ctx.js |: autodetect xtype |
| src/adapters/cservice.js |: do not use userPreferencesModelName |


### FIXMES



## Additional-info

Minified version : 4 files (size: 52K)
Development version : 36 files (size: 168K)
