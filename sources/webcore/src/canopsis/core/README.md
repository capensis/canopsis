# Core Canopsis Brick

## Index

- [Description](#description)
- [Content](#content)
- [Installation](#installation)
- [Usage](#usage)
- [Continuous-integration](#continuous-integration)
- [Code-notes](#code-notes)
- [Additional-info](#additional-info)

## Description

Core UI brick. Provides the base application layer for the Canopsis UI

## Content



## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ su - canopsis
    $ brickmanager install core

Then, you need to enable the brick

    $ brickmanager enable core

You can see enabled bricks

    $ su - canopsis
    $ brickmanager list
    [u'core', u'uibase', u'monitoring', ..., **u'core'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/core/blob/master/doc/index.rst)

## Continuous-Integration

### Tests



### Lint

Tested on commit : 4e09104.

| Target | Status | Log |
| ------ | ------ | --- |
| Lint   | :negative_squared_cross_mark: ERROR |  |


## Code-Notes

### TODOS

| File   | Note   |
|--------|--------|
| src/routes/application.js | check if this is still used |
| src/lib/abstractclassregistry.js | manage element with add and remove methods |
| src/lib/factories/widget.js | check if this is still needed, as mixins are in configuration now |
| src/lib/utils/notification.js | doing it clean |
| src/lib/utils/notification.js | doing it clean |
| src/lib/utils/notification.js | doing it clean |
| src/lib/utils/dom.js | remove this file and requirements from elsewhere |
| src/lib/utils/data.js | change parentElement term to something more descriptive |
| src/lib/utils/event.js | delete this, as it looks more like a registry than an util |
| src/lib/utils/widgetSelectors.js | implement Key exist feature |
| src/view/mixineditdropdown.js | @gwen check if it's possible to remove this class |
| src/view/validationtextarea.js | move this to components dir |
| src/view/validationtextfield.js | move this to components dir |
| src/view/listline.js | @gwen check if it's possible to remove this class |
| src/view/formwrapper.js | watch out ! garbage collector might not work here! Possible memory leak. |
| src/view/formwrapper.js | "on" without "off" |
| src/view/tabledraggableth.js | @gwen check if it's possible to remove this class, or move it to uibase |
| src/mixins/embeddedrecordserializer.js | dynamize |
| src/mixins/inspectableitem.js | refactor the 20 lines below in an utility function "getEditorForAttr" |
| src/mixins/consolemanager.js | move this to development brick |
| src/components/renderer/component.js | check why there is a property dependant on "shown_columns" in here. As it is a List Widget property, it does not seems relevant at all. |
| src/components/editor/component.js | check if still used |
| src/controller/login.js | delete store in this#destroy |
| src/controller/form.js | refactor this |
| src/controller/partialslotable.js | put this in arrayutils |
| src/controller/widget.js | manage this with utils.problems |
| src/forms/modelform/controller.js | search this value into schema |
| src/forms/modelform/controller.js | use the real schema, not the dict used to create it |


### FIXMES

| File   | Note   |
|--------|--------|
| src/routes/userview.js | don't use jquery in here, it's for views ! |
| src/components/editor/component.js | auto-detect if we need standalone mode or not, stop using a variable, for a better comprehension |
| src/controller/userview.js | wrapper does not seems to have a widget |
| src/controller/recordinfopopup.js | do not use jquery for that kind of things on a controller |
| src/forms/widgetform/controller.js | this works when "xtype" is "widget" |


## Additional-info

Minified version : 4 files (size: 240K)
Development version : 114 files (size: 700K)
