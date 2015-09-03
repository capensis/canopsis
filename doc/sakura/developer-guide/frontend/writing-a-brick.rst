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

.. code:: javascript

   {
     "name": "canopsis-rights",
     "version": "0.1.0",
     "authors": ["Team Canopsis <maintainers@canopsis.org>"],
     "description": "Rights and permission management",
     "main": "init.js",
     "license": "GNU Affero General Public License",
     "homepage": "http://www.canopsis.org/"
   }


Initialization file
-------------------

The initialization file is the only javacript file that is automatically loaded when loading a brick.

Including more files
--------------------

If you want to load more files, you should require them directly from the ```init.js``` file, as the following example show :

.. code:: javascript

   define([
       'canopsis/brick-name/controller/sample',
       'canopsis/brick-name/widgets/sample/controller',
       'link!canopsis/brick-name/widgets/sample/style.css'
   ], function () {});


Javascript files content
------------------------

Javascript code should be wrapped into `Ember Initializers 
<http://guides.emberjs.com/v1.10.0/understanding-ember/dependency-injection-and-service-lookup/#toc_dependency-injection-with-code-register-inject-code>`_.

Here is an example of a file that respects the above guidelines :

.. code:: javascript

   Ember.Application.initializer({
       name: 'UiactionbuttonWidget',
       after: 'WidgetFactory',
       initialize: function(container, application) {
           var WidgetFactory = container.lookupFactory('factory:widget');
           var widget = WidgetFactory('uiactionbutton',{
               tagName: 'span',
               actions: {
                   do: function(action, params) {
                       if(params === undefined || params === null){
                           params = [];
                       }
   
                       this.send(action, params);
                   }
               }
           });
   
           application.register('widget:uiactionbutton', widget);
       }
   });
   

Additionnal content
-------------------

Bricks usually contains inside their repository :

- A readme
- Docstrings (JSDoc format) inside js code
- A user guide

Tooling
-------

To help people respecting convention and producing good quality code, some tooling is available. It is currently a work in progess, and it will become more and more advised to use it, as the tool will grow.

See `Canopsis-ui-toolbelt npm module 
<https://git.canopsis.net/gpluchon/canopsis-ui-toolbelt>`_.

