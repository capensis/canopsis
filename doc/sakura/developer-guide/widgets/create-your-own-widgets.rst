Create your own widgets
=======================

Introduction
------------

Making a widget for canopsis is quite simple. First of all you can find a sample in ``/opt/canopsis/var/www/webcore-libs/samples/widget/sample/`` to help you.

This require some basic Javascript and Object Oriented Programming (OOP) knowledge.

Third party widgets goes in ``/opt/canopsis/var/www/widgets/myWidgetFolder``

Files
-----

.. code-block:: text

	└── sample
		├── libs/
		├── sample.css
		├── sample.js
		└── widget.json

* `libs`: Here you put all your externals libs (jquery plugin or whatever else)

* `.css`: It's your css plugin, must have the same name as your widget

* `.js`: The main source code (further explanation below)

Javascript
-----------

Skeleton
^^^^^^^^^

A widget inherit of `cwidget` class that provide auto refresh function and generic base request function.

The skeleton of a widget is the following one:

.. code-block:: javascript

	Ext.define('widgets.thirdparty.sample.sample' , {
	    extend: 'canopsis.lib.view.cwidget',
	    alias: 'widget.sample',
	});

Constructor
^^^^^^^^^^^

.. code-block:: javascript

	initComponent: function() {
	    this.callParent(arguments);
	},

InitComponent is the widget constructor, if you need to execute something before the widget start its building.

.. warning:: Extjs function are not available before the `callParent` function.


Main container
^^^^^^^^^^^^^^^

You can find the main html container in the variable `this.wcontainer`. After this container is rendered the method `afterContainerRender` is called. Exemple with jquery plugin:

.. code-block:: javascript

	afterContainerRender: function() {
	    $('#'+this.wcontainerId).myAwesomePluginJquery({})
	}


Refresh function
^^^^^^^^^^^^^^^^^

.. code-block:: javascript

	doRefresh: function(from, to) {
	}


This function is called at every refresh, you can update your widget and/or plugin in this function.

Widget.json
------------

This file describe your widget:

.. code-block:: javascript

	[{
		"name":  "Sample",
		"version": 1,
		"author": "<AUTHOR>",
		"website": "<WEBSITE>",
		"xtype": "sample",
		"thirdparty": true,
		"description": "Sample Widget",
		"refreshInterval": 300,
		"border" : false,
		"options": [
			{
				"title": "HTML",
				"layout": "anchor",
				"items" : [
					{
						"xtype": "htmleditor",
						"anchor": "100%",
						"name": "innerText",
						"height" : 350,
						"value": ""
					}
				]
			}
		]
	}]

The very important part here is to build your options in the "options" attribute. Those items will be rendered in the wizard when you create a widget. For this one you must know how to create `Extjs form object <http://docs.sencha.com/ext-js/4-2/#!/api/Ext.form.field.Text>`_.

When the user save his widget, all form typed will be available by their name in your widget javascript file. with the widget.json exemple right above, the text typed in "htmleditor" will be avaible by a simple `this.name`.

Be sure that your variable name doesn't collide with other variables (use prefix for var names ?)

.. important:: your widget.json file must be a pure json file, every comma count! you can validate your file with `a json parser <http://json.parser.online.fr/>`_ .


Complete exemple
----------------

.. code-block:: javascript

	Ext.define('widgets.thirdparty.sample.sample' , {
		extend: 'canopsis.lib.view.cwidget',

		alias: 'widget.sample',

		logAuthor: '[sampleWidget]',

		// Setted by Wizard :)
		innerText: undefined,

		refresh_number: 0,

		initComponent: function() {
			log.debug("initComponent", this.logAuthor)
			this.callParent(arguments);
		},

	        afterContainerRender: function() {
	                log.debug('Container just rendered')
	        },

		doRefresh: function(from, to) {
			log.debug("doRefresh", this.logAuthor)
			this.setHtml("refresh_number: " + this.refresh_number +", Html: <br>" + this.innerText);

			this.refresh_number += 1;
		}
	});

Deploying your widget
----------------------

Once your widget is fully functional, you just need to regenerate the minified JavaScript/CSS/... :

.. code-block:: bash

	$ su - canopsis
	$ webcore_minimizer
