.. _ED__Title:

========================
Writing a frontend brick
========================

This guide will cover the process of brick creation, and show how it is possible to add features to the frontend, while preserving the initial frontend codebase, as the brick is assimilable to a self-contained plugin.


.. contents::
   :depth: 2

References
==========

List of referenced functional requirements...

- :ref:`FR::Title <FR__Title>`
- :ref:`TR::Title <TR__Title>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Gwenael Pluchon", "2016/02/16", "1.0", "First draft", ""

Contents
========

.. _ED__Title__Desc:


Brick creation
--------------

Creating the brick
^^^^^^^^^^^^^^^^^^

A boilerplate brick has been created in a git repository to help developers to kickstart this step.

Currently, the repository is available at https://github.com/gwenaelp/brick-boilerplate.

Bricks usually have to be contained in a dedicated repository, to provide a better packaging and delivery process. Thus, you will have to fork this brick into a repository.

Then, clone your brick into the frontend source directory.

Initializing the project
^^^^^^^^^^^^^^^^^^^^^^^^

To help maintaining bricks, the use of a tool (canopsis-ui-toolbelt) is advised. It has to be installed only on development environments, inside the brick.

To install canopsis-ui-toolbelt, go into the brick folder, and execute the following command :

.. code-block:: bash

   npm install

npm will take care to install the tooling and its dependencies at the right place.

To complete the brick initialization step, it is require to reconfigure the brick manifest (located into the `manifest.json` file).

The manifest have several options, but few are really required to start a project. To see how to configure a manifest in depth, please refer the  `canopsis-ui-toolbelt documentation <https://git.canopsis.net/canopsis/canopsis-ui-toolbelt#manifest>`_.

The only options that are relevant at this steps are :

- name : The name of the brick. It is the same as the repository and the brick folder names
- description : A short description of the brick purpose
- repository : repository information, in the format as npm's package.json
- author : Your name
- license : The license under which you want to release the brick

Once you have changed these options to fit your needs, you have to regenerate some files using the tooling :

.. code-block:: bash

   npm run compile

This command is supposed to take care to regenerate files that are automatically handled by the tooling. These files are those located in the root directory, except the `manifest.json` file.

This compilation step will take care to propagate the modifications done into the manifest file, into :

- README.md
- package.json
- bower.json
- and so on, and so forth...


Folder hierarchy
^^^^^^^^^^^^^^^^

Now that the brick is set up, let's have a look on the folder hierarchy.

- the src folder contains most of the source code.
- the requirejs-modules folder can contain code that uses requirejs. It is mostly used to manage javascript libraries.
- the dist folder contains a compiled version of the source code.
- the doc folder contains the documentation (API, user guides, ...).
- the tests folder contains tests scenarios, whether they are functionnal tests or unit tests.
- the externals folder can contain external libraries, such as jquery plugins for instance.
- the schemas folder can contain Json schemas, that will be used as models in the frontend application.

Environments
^^^^^^^^^^^^

Bricks manifests can contain an "envMode" property. This property can take two values : "production" or "development". If this property is not assigned into the manifest, it is considered that the brick is in production mode.

- Production mode have to be used when the brick is ready to be deployed, and it is assumed that is is stable. The minified version of the code will be executed.

- Development mode is used while editing brick code. The non-minified source code will be executed. This mode is more ressource-consuming, but it allows to help developers to debug their brick.

Starting to code
^^^^^^^^^^^^^^^^

Once that everything is set up, you can start implementing features into the brick.

To add some code, you just have to create code files into the src folder.

The brick system supports at the moment :

- Javascript files (.js)
- Css files (.css)
- Handlebars files (.hbs). Note that for handlebars files, the file name corresponds to the template name, except for components.

Except for components (see the dedicated part), you can manage your source folder hierarchy however you want. However, it is advised to keep this folder tied to this hierarchy :

- adapters
- components
- controllers
- forms
- mixins
- serializers
- templates
- views
- widgets
- ...

Note that for components, the following directory structure must be respected :

- src
   - components
      - <componentName>
         - component.js
         - template.hbs

Note that when you add or remove files to the source folder, the tests folder, or the requirejs-modules folder, it is mandatory to re-run the `npm run compile` command, to re-generate application entry points (`init.js`, `init.test.js`, ...).

Writing javascript code
^^^^^^^^^^^^^^^^^^^^^^^

Javascript code must be contained into Ember initializers, to be able to manage dependency injection properly :

.. code-block:: javascript

   Ember.Application.initializer({
     name: 'MyObjectInitializer',
     after: 'InheritedObjectInitializer'
     initialize: function(container, application) {
       var InheritedObject = container.lookupFactory('object:inherited-object');

       var MyObject = InheritedObject.extend({});

       application.register('object:my-object', MyObject);
     }
   });

.. WARNING::
   As of now, we will assume that our code is always wrapped into an initializer. Examples will not show the initializer wrapping to preserve code clarity and conciseness.

Writing components
^^^^^^^^^^^^^^^^^^

Components are elements that can be called directly from handlebars template, and data is directly binded to them through templates :

.. code-block:: handlebars

   {{component-checkbox checked=checkedValue class="toggle"}}

For a complete guide about writing components, please refer to Ember guides :
https://guides.emberjs.com/v1.11.0/components/

Writing editors and renderers
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Editors and renderers elements that shows up on the interface when it's needed to display or edit data in a nifty way. They usually calls components from the template when it's needed to provide some javascript function calls.

.. NOTE::
   Editors usually appears on forms, and renderers on list cells.

The data that is supposed to be handled by editors and renderers is schema's properties.

To find the suitable editor for a data, The frontend internal logic is the following :

- If the property contains a role, find (by looking for the template "editor-<role>.hbs") if an editor is available.
- Then, if no editor is found, it tries to find the editor for the data type (with the same logic as above)
- Finally, if nothing is found, it falls back to a default input.

The logic is the same for renderers.

Naming conventions
""""""""""""""""""

- Editors templates must start with "editor-" ("editor-boolean.hbs" for instance)
- Renderers templates must start with "renderer-" ("renderer-boolean.hbs" for instance)

Writing widgets
^^^^^^^^^^^^^^^

Widgets are the elements that are displayed on userviews. Basically, an userview is a set of widgets that are saved into it, with their configuration.

Widgets are composed by :

- A controller
- A schema (that generates a model)
- A template
- Eventually, a view mixin (that is applied to the view, to allow injecting code on the widget view)

So here, the MVC pattern is in use. It is advised to try to stick to it and :

- Write the data handling logic into the controller.
- Describe the widget editable properties into the schema.
- Use the template to do the layout.
- Interact with the template within the view mixin.

Writing adapters
^^^^^^^^^^^^^^^^

Adapters are objects dedicated to contains all the logic to interact with a dedicated backend API. There should be no backend-specific code elsewhere on the code base, making the application compatible with different backends by enabling some other adapters.

To see more information about writing adapters, please refer to the `Ember Adapter guide <https://guides.emberjs.com/v1.11.0/models/customizing-adapters/>`_.


Writing mixins
^^^^^^^^^^^^^^

Mixins are initially present to encourage code reuse : They can be added at the runtime to classes to enrich the classes algorithms.

In the frontend, mixin internals have been pushed further. They can be applied to widgets directly from the view edition mode, and be delivered with a schema to provide persistant settings that are configurable through a form.

If a mixin is provided with a schema, they must have the same name.

Importing and using external JS libraries
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^


Testing
-------

Tests have to be written in the "tests" brick folder. 

Here is a test file example :

.. code-block:: javascript

   module('core'); //This must be the brick name

   test('Creating a view with an empty text widget', function() {
       visit('/userview/view.event');

       expect(2); //We have 2 assertions in our test

       click('.nav-tabs-custom a.dropdown-toggle');
       click('.nav-tabs-custom .fa.fa-plus');

       waitForElement('input[name=crecord_name]').then(function(){
           fillIn('input[name=crecord_name]', 'test');
           click('.modal-dialog .btn-primary');
       });
       click('.nav-tabs-custom a.dropdown-toggle');
       click('.nav-tabs-custom .fa.fa-pencil');
       click('.btn-add-widget');

       waitForElement('.modal-dialog .ember-text-field').then(function(){
           equal(find('.box-title').length, 0, 'No widget on the view');
           fillIn('.modal-dialog .ember-text-field', 'text');
           click('.modal-dialog .panel-default:first a');
           click('.modal-dialog .list-group-item a');
           click('.modal-dialog .btn-primary');
           click('.modal-dialog .btn-primary');
           waitForElement('.box-title').then(function(){
               equal(find('.box-title').text(), "< Untitled text widget >", 'an untitled text widget is present');
           });
       });
   });

When adding a test file, the brick file generation step have to be run to modify the testing entry point of the brick ("init.test.js").

More information about available test helpers is available at :

- `Ember Test guide <https://guides.emberjs.com/v1.11.0/testing/test-helpers/>`_
- "Tests" brick documentation

Brick delivery
--------------

Once the brick has been written, it is important to perform several checks and manipulations to ensure the brick is ready to be shipped to client installations.

- Set the environment to "production" ("envMode": "production")
- Run the minification toolchain on the brick (npm run minify)
- Recompile all the automatically generated brick files (npm run compile)
- Ensure test suites are all green
- Run the doc compilation (npm run doc)
- Ensure the brick lint is OK (npm run lint)

Once these steps has been completed, commit the brick to it's upstream repository.
