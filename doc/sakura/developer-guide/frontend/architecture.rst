Architecture
************


1. Presentation
---------------

Canopsis UI is built on top on the `Emberjs <http://www.emberjs.com>`_ framework. This framework is use in Canopsis as a skeletton that let developper plug functionnalities all around the Canopsis core. Canopsis UI also uses some other frameworks like `requirejs <http://requirejs.org>`_ or the famous css framework `bootstrap <http://getbootstrap.com>`_ In this documentation will be explained how canopsis core is designed and what are it's mechanics. It is also explained how to integrate new features in order they best fit existing content.

Canopsis UI once installed on a server is served by the `Gunicorn <http://gunicorn.org>`_ webserver with some static **javascript, html, css** files and some **dynamic content** computed and server throught the `Canopsis API </developer-guide/API/index.html>`_

Please note that the folder referential is based on the **webcore** folder which in Canopsis sources is located in the subfolder ``canopsis/sources/webcore/var/www/canopsis`` . This path once Canopsis built is by default ``/opt/canopsis/var/www/canopsis``. Some resources may however only be avaiable from other location. This is the case for *locales* translation *.po* files that are located by default in ``/opt/canopsis/locale`` or schemas that are documented in the `model </developer-guide/uiv2/model_layer.html>`_ part of the Canopsis documentation.

1. The Entry point
------------------

When the *index.html* file is queried by a client to the server, index.html loads a few files (mainly css ones) directly. However, it includes the **canopsis.js** file that is the entry point of the canopsis UI system. File loading is handled by requirejs and once loaded, the canopsis.js file instanciate the require.js framework recursively and this leads to the call of each canopsis requirement for the initial load.

2. Dependencies overview
------------------------

Canopsis UI is made of many dependencies that together render the User Interface once loaded. These dependencies are made of:

 - `editors <#>`_
 - `factories <#>`_
 - `loaders </developer-guide/uiv2/architecture.html>`_
 - `renderers <#>`_
 - `schemas <#>`_
 - `templates <#>`_
.. - `widgets <#>`_

3. Canopsis dependencies Loaders
--------------------------------

The initial application load though canopsis.js calls some Canopsis javascript files which role is to load Canopsis UI resources such as Ember js components, asynchonously. These loaders manage resources such as **templates, ember components, Canopsis forms, some tools or widgets**. All those resources are part of the canopsis UI and are various combination of HTML templates and/or javascript view/controllers. The dependencies for each of these loaders are hardcoded into each file but their design let further integration as easy as copying a line and fill it with a new value. Below, some loader example may light about how to integrate some features from low application level.

The template loader is based on a list of template that ember will require from the template folder located in the webcore/template/<name of the template>.html

.. code-block:: javascript

    { name: 'application' }

This requirement object will tell the Canopsis template loader to fetch the application.html file then it's content will be compiled and stored in the ``Ember.TEMPLATE`` object.

Another more complex example is the editor component loaders that loads both javascript and html template depending on parameters into Canopsis UI. One of these editor call for <editor_name> == right is:

.. code-block:: javascript

   { name: 'rights', js: 'v' },

the name information in this object describes the component folder to load and the js attribute define what should be loaded. Available parameters to load js are:

   - none: does not load extra file
   - **c** : loads ``app/editors/<editor_name>/controller.js``
   - **v** : loads ``app/editors/<editor_name>/view.js``
   - **w** : loads ``app/editors/<editor_name>/component.html`` and ``app/editors/<editor_name>/component.js``

These options can be mixed together to use for exemple in `cw` case a controller in the web component.

Each of these options will load ``app/editors/<editor_name>/template.html`` too.

The third example is the factory loader which will load a factory js file:

.. code-block:: javascript

   { name:'editor', url: 'app/lib/factories/editor' }

Here is how to load an editor.


Loaders are mostly based on the same model and adding a new entry is as simple as adding a well formatted object in a list depending on the type of loaded resource.


Canopsis UI Model
-----------------

The Canopsis UI Model system is based upon json schemas that describes datatypes for each document type managed into canopsis. Those schemas are used in both front office an back office in order to keep redundancy in the project. see more `model </developer-guide/uiv2/model_layer.html>`_

Widgets
-------

Widgets are components used in Canopsis UI. They are made of a controller and a template and they can be parametrized in order to best fit users need. see more `widgets </developer-guide/uiv2/widgets.html>`_


