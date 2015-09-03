Writing a brick
===============

Writing new features as Canopsis frontend addons is supported and encouraged. Developpers might want to enhance their Canopsis with widgets, adapters, templates, and so on.

To reach this goal, Canopsis frontend is divided in bricks. Bricks are a way to split features in deliverable packages.

Requirements
------------

A brick is a folder whose purpose is to stay in the ~/var/www/canopsis directory. This folders contains :

 * A ```bower.json``` manifest file at the root folder of the brick (mandatory)
 * a ```init.js``` file. (mandatory)
 * Javascript files.
 * Javascript libraries
 * HTMLBars templates
 * Css files
 * Images
 * Basically whatever you need to fetch from the frontend code, or need as build tools, ...

Thus, bricks are compartmentalised bits of the UI.
The goal of a brick maintainer is to make the frontend work with or without its brick, whatever the sibling enabled bricks are.

The developper should still take into account that there are some core bricks always activated (mostly ```core``` and ```uibase``` bricks).

Manifest file
-------------

The manifest file contains a JSON structure that shows up some brick meta-information, such as:

- The license of the brick
- Informations about the maintainer
- The brick dependencies
- The brick description
- The brick version number
- The brick Homepage

To initialize a manifest and run a wizard, you can run the command ```bower init``` inside your brick.

Sample brick manifest file :

```
{
  "name": "canopsis-rights",
  "version": "0.1.0",
  "authors": [
    "Team Canopsis <maintainers@canopsis.org>"
  ],
  "description": "Rights and permission management",
  "main": "init.js",
  "license": "GNU Affero General Public License",
  "homepage": "http://www.canopsis.org/"
}
```

Initialization file
-------------------

The initialization file is the only javacript file that is automatically loaded when loading a brick.

Including more files
--------------------

If you want to load more files, you should require them directly from the ```init.js``` file.
