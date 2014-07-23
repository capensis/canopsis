Definitions
===========

``HOME='~/canopsis'``

``WEBCORE='$HOME/sources/webcore/var/www/canopsis'``

How to create a widget for canopsis
===================================

1. Create a folder into the widget folder.
------------------------------------------

This folder contains a controller and a template, for instance the list widget is initiated as follow:

``mkdir $WEBCORE/components/widgets/weather/``

``touch $WEBCORE/components/widgets/weather/controller.js``

The copntroller is manager thanks to a factory that allow defining the widget and make it benefits of canopsis custom mixins to extend default functionalities.

``touch $WEBCORE/components/widgets/weather/template.js``

2. reference your widget
------------------------

Add a reference in the ``widgetsTemplates`` array in ``$WEBCORE/core/lib/loader/widgets.js`` as follow

.. code-block:: javascript

	//hasJSPart enable js controller file load.
	var widgetsTemplates = [
		... ,
		{ name:'weather', hasJSPart: true }
	];


``touch $HOME/sources/mongodb-conf/opt/mongodb/load.d/views/weather.json``

3. Add a widget schema to mongodb
---------------------------------

The json schema descriptor is usually loaded on Canopsis build-install

Fill it with the widget json basis accordingly:

.. code-block:: javascript

	{
	  "id": "view.weather",
	  "_id": "view.weather",
	  "crecord_type": "view",
	  "crecord_write_time": 1400853201,
	  "enable": true,
	  "containerwidget": {
	    "xtype": "containervbox",
	    "title": "container title vbox",
	    "_id": "test_view_vertical_container",
	    "items": [
	      {
	        "title": "wrapper",
	        "xtype": "widgetwrapper",
	        "widget": {
	          "xtype": "weather",
	          "listed_crecord_type": "weather",
	          "title": "Weather",
	          "toolbar": [
	            "actionbutton-create",
	            "actionbutton-removeselection"
	          ]
	        }
	      }
	    ]
	  },
	  "parent": [],
	  "aaa_access_owner": [
	    "r",
	    "w"
	  ],
	  "aaa_group": "group.CPS_view_admin",
	  "crecord_creation_time": 1400853201,
	  "aaa_access_unauth": [],
	  "aaa_owner": "account.root",
	  "internal": false,
	  "aaa_access_other": [
	    "r"
	  ],
	  "aaa_access_group": [
	    "r"
	  ],
	  "aaa_admin_group": "group.CPS_view_admin",
	  "children": [],
	  "crecord_name": "Accounts view"
	}

Please note that the widget xtype and listed_crecord_type that defines which schema is used to display the widget.

4. Add the widget schema
------------------------

Your newly created widget may need some configuration from UI to work properly. This can be done by defining a schema that describes some properties and how to represent some values for your widget wizard.

for instance, the `weather` widget is configured with the following schema that allow UI to generate the appropriate form to configure the widget.

.. code-block:: javascript

	{
		"type": "object",
		"categories": [{
			"title": "general",
			"keys":["watched_rk"]
		}],
		"properties": {
			"watched_rks": {
				"type": "string"
			},
			"xtype": {
				"type": "string"
			}
		}
	}


5. Publish the widget
---------------------

Run filldb.py as canopsis user:

``python /opt/canopsis/opt/mongodb/filldb.py init``

6. Add your widget to the ember routes
--------------------------------------

In order to make the widget reachable from a canopsis URL, in the **run** children array add

.. code-block:: javascript

	{
		"type": "resource",
		"name": "run",
		"icon": "play",
		"description": "show mmenu",
		"children": [
			...,
			{
			    "type": "resource",
			    "name": "weather",
			    "icon": "filter",
			    "description": "weather"
			}
		]
	}
