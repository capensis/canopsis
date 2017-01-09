TR alarms
---------

.. contents:: Table of contents


Configuration
=============

The widget configuration is organized around menus which look like this:

* General
* Columns
* Info Popup
* Mixins

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/conf_tab_alarms.png


Tab 1 : General
^^^^^^^^^^^^^^^

* The widget needs a title.
* If there are no user preferences, default_sort_key and default_sort_dir may be defined


Tab 2 : Columns
^^^^^^^^^^^^^^^

Like the existing "list widget", this one should allow users to select columns to be shown on the table

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/conf_columns.png


Tab 3 : Info Popup
^^^^^^^^^^^^^^^^^^

The user must be able to define a template (with handlebars) which could be compiled when he clicks on a cell.
For example, if I click on a connector cell, a popup will be displayed with informations from the template.


Tab 4 : Mixins
^^^^^^^^^^^^^^

This tab is automaticaly added with the schema def.


Data schema
===========

In order for you to understand what is needed to display table, here is the data schema with explanations.
All needed informations are described in it. If it's not clear, feel free to contact us.

You have to know that data are provided by an adapter which relies on a webservice.
In order to use only frontend concepts, we provide 1 dataset (in folder `datasets`) which can be loaded by the adapter.


.. code-block::

  return this.ajax('/static/canopsis/brick-alarms/datasets/ds1.json', 'GET', {data: ""});

These data are already available in the alamrs controller. the "fetchData" method of the controller get these data.


**TO BE COMPLETED BY FLO**

.. code-block::

    fetchData: function() {
        var controller = this;

        var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:serviceweatherdata');
        adapter.findQuery('serviceweatherdata', {}).then(function (result) {
            // onfullfillment
            console.log('Raw data', result);
        }, function (reason) {
            // onrejection
            console.error('ERROR in the adapter: ', reason);
        });
    },

the "result" variable give you the data.


Raw schema
^^^^^^^^^^

Here is the raw schema => https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/schemas/crecord.alarms.json
**TO BE COMPLETED BY FLO**

You can find some datasets compliant with schema here : https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/tree/master/datasets


WebService
^^^^^^^^^^

**TO BE COMPLETED**

.. code-block::

  def get_alarms(
             tstart,
             tstop,
             opened=True,
             resolved=False,
             consolidations=[],
             filter={},
             search='',
             sort_key='opened',
             sort_dir='DESC',
             skip=0,
             limit=50
     ):
         """
         Return filtered, sorted and paginated alarms.
         :param int tstart: Beginning timestamp of requested period
         :param int tstop: End timestamp of requested period
         :param bool opened: If True, consider alarms that are currently opened
         :param bool resolved: If True, consider alarms that have been resolved
         :param list consolidations: List of extra columns to compute for each
           returned alarm. Extra columns are "pbehaviors" and/or "linklist".
         :param dict filter: Mongo filter. Keys are UI column names.
         :param str search: Search expression in custom DSL
         :param str sort_key: Name of the column to sort
         :param str sort_dir: Either "ASC" or "DESC"
         :param int skip: Number of alarms to skip (pagination)
         :param int limit: Maximum number of alarms to return
         :returns: List of sorted alarms + pagination informations
         :rtype: dict
         """


Main description
^^^^^^^^^^^^^^^^

.. csv-table:: Alarm main description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "Main", "_id", "string", "04f2372b-8410-40b6-a5ce-7dc3a3f0ade1", "Unique ID of an alarm"
   "", "", "string", "/component/bra/iva/eqw", "Uinique ID of the entity concerned by the alarm"
   "", "t", "timestamp", "1462399200", "Date of the alarm creation"
   "", "v", "list", "", "Contents of an alarm"


V description
^^^^^^^^^^^^^

.. csv-table:: V description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "V", "connector", "string", "nagios", "Type of the connector source"
   "", "connector_name", "string", "prodnagios1", "Name of the connector"
   "", "component", "string", "a_component", "Name of the component"
   "", "resource", "string", "a_resource", "Name of the resource"
   "", "output", "string", "a_output", "Current output of the alarm"
   "", "solved", "timestamp", "1462399200", "Date of the end of the alarm. If null, alarm is still alive"

Extra description
^^^^^^^^^^^^^^^^^

.. csv-table:: Extra description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "Extra", "e1", "string", "Extra1", "Extra fields that come with the alarm"
   "", "e2", "string", "Extra2", "Extra fields that come with the alarm"


State description
^^^^^^^^^^^^^^^^^

.. csv-table:: State description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "state", "a", "string", "John Doe", "Author which has generated this state"
   "", "_t", "string", "stateinc statedec changestate", "Type of the step"
   "", "m", "string", "Resource 9 in state 0", "Message that comes with the state"
   "", "t", "number/timestamp", "1476673252", "Timestamp of the state"
   "", "val", "number [0-3]", "0", "Value of state"


Status description
^^^^^^^^^^^^^^^^^^

.. csv-table:: Status description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "status", "a", "string", "John Doe", "Author which has generated this status"
   "", "_t", "string", "statusinc statusdec changestatus", "Type of the step"
   "", "m", "string", "Component 10 in status 3", "Message that comes with the status"
   "", "t", "number/timestamp", "1476673252", "Timestamp of the status"
   "", "val", "number [0-3]", "0", "Value of status"


Snooze description
^^^^^^^^^^^^^^^^^^

.. csv-table:: Snooze description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "snooze", "a", "string", "John Doe", "Author which has generated this snooze"
   "", "_t", "string", "snooze", "Type of the step"
   "", "m", "string", "Resource 9 is snoozed for 600s", "Message that comes with the snooze"
   "", "t", "number/timestamp", "1476654503", "Timestamp of the snooze (begining)"
   "", "val", "number/timestamp", "1476655103", "Timestamp of the end of snooze"

ACK description
^^^^^^^^^^^^^^^

.. csv-table:: ACK description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "ack", "a", "string", "John Doe", "Author which has generated this ack"
   "", "_t", "string", "ack", "Type of the step"
   "", "m", "string", "ack from MMA", "Message that comes with the ack"
   "", "t", "number/timestamp", "1476654503", "Timestamp of the ack"
   "", "val", "string", "null", "N/A"

Ticket description
^^^^^^^^^^^^^^^^^^

.. csv-table:: Tikcet description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "ticket", "a", "string", "John Doe", "Author which has generated this ticket"
   "", "_t", "string", "declareticket", "Type of the step"
   "", "m", "string", "ticket from MMA", "Message that comes with the ticket"
   "", "t", "number/timestamp", "1476654503", "Timestamp of the ticket"
   "", "val", "string", "null", "N/A"


Cancel description
^^^^^^^^^^^^^^^^^^

.. csv-table:: Cancel description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "cancel", "a", "string", "John Doe", "Author which has cancelled the alarm"
   "", "_t", "string", "cancel", "Type of the step"
   "", "m", "string", "alarm was cancelled from MMA", "Message that comes with the cancel action"
   "", "t", "number/timestamp", "1476654503", "Timestamp of the cancel"
   "", "val", "string", "null", "N/A"


Linklist description
^^^^^^^^^^^^^^^^^^^^

.. csv-table:: Linklist description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "linklist", "url", "string", "http://urltoticket.local/?id=entity_id", "Url associated to a label"
   "", "label", "string", "Ticket", "Label associated to an url"


Linklist is a list of urls associated to the entity.
Links must appear in the modal like potentialy any other variables but with special helper.

The goal is to let the user access a handlebar renderer

 {{ linklist category="procedure" }}

Must return something like

 foreach linklist with category = "procedure"

 <a href="http://urloflinklist">labeloflinklist</a><br>


Pbehavior description
^^^^^^^^^^^^^^^^^^^^^

.. csv-table:: Pbehavior description
   :header: "Structure", "Attribute", "Type", "Example", "Description"
   :widths: 5, 10, 5, 10, 30

   "pbehavior", "behavior", "string", "maintenance pause", "Name of the behavior"
   "", "isActive", "boolean", "True False", "Is the pbehavior active ?"
   "", "dtstart", "number/timestamp", "1476705600", "Timestamp of the begin of pbehavior"
   "", "dtstop", "number/timestamp", "1476706600", "Timestamp of the end of pbehavior"
   "", "rrule", "structure attr1 : string, attr2 : string", "text=Every Week, rule='FREQ=WEEKLY'", "Reccurent rule of the behavior"



Actions / Events
================

In the widget, users may be able to execute actions.
In the Canopsis world, actions are performed via sending messages to a AMQP bus.

Listing
^^^^^^^

Here is a list of actions that need to be handled by the widget :

.. csv-table:: Actions description
   :header: "Action", "Type", "Goal", "Attrs description"
   :widths: 5, 5, 15, 30

   "confirm", "changestate", "Change criticity of an alarm", "See `Schema <https://git.canopsis.net/canopsis/canopsis/blob/develop/sources/python/alerts/etc/schema.d/cevent.changestate.json>`_. "
   "invalidate", "changestate", "Change criticity of an alarm", "See `Schema <https://git.canopsis.net/canopsis/canopsis/blob/develop/sources/python/alerts/etc/schema.d/cevent.changestate.json>`_. "
   "pause", "pbehavior", "Change criticity of an alarm", "{}"
   "declareticket", "declareticket", "Call a API/email of an external tool to create a ticket", "See `Schema <https://git.canopsis.net/canopsis/canopsis/blob/develop/sources/python/alerts/etc/schema.d/cevent.declareticket.json>`_. "
   "assocticket", "assocticket", "Add a ticket number into Canopsis", "See `Schema <https://git.canopsis.net/canopsis/canopsis/blob/develop/sources/python/alerts/etc/schema.d/cevent.assocticket.json>`_. "


Rendering
=========

Main table
^^^^^^^^^^

* The main table must repect adminlte standards  https://almsaeedstudio.com/themes/AdminLTE/pages/tables/simple.html

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/general_render.png


* It must be responsive (big screen, desktop, mobile)
* 50 tr must be shown in 1 second, not more.
* Pagination (done by the backend)
* Sort (done by the backend)


Responsive
^^^^^^^^^^

As the widget is a table, the responsive feature can take args to perform.  
The user must be able to spécify columns that can be not printed if display does not permit it.  
In the widget conf, user must be able to select these columns.


.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/responsive_list.png

Column Renderering
^^^^^^^^^^^^^^^^^^

The user must be able to select columns and order he wants to show on the main table within the widget conf form.

Some data have to be shown with a renderer.
For example, a timestamp must use a special timestamp renderer.
The mapping between data and renderer is done via the schema.


.. code-block::

 "opened": {
       "stored_name": "t",
       "role" : "timestamp"
     },

With these informations, you know that you have to call the renderer below

.. code-block::

 $ cat uibase/src/renderers/renderer-timestamp.hbs
 {{!*
  * @renderer timestamp
 }}
 {{#unless attr.options.hideDate}}
     <div>{{timestamp value attr}}</div>
 {{/unless}}
 {{#if attr.options.canDisplayAgo}}
     <small class="text-muted">
         <span class="glyphicon glyphicon-time"></span>
         {{timeSince value}}
     </small>
 {{/if}}

If there is no role associated with the attribute, you have to render value as string.

Customfilterlist
^^^^^^^^^^^^^^^^

In the widget, users must be able to set data filters.  

This is done by using a lib called **querybuilder**.  
The library is already included in Canopsis.  
Filters are formatted as mongodb filters.

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/customfilterlist.png

Generated filters must be forwarded as webservice parameters.

**TO BE COMPLETED BY FLO**


Array Search
^^^^^^^^^^^^

The widget must show a input to make searches

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/search.png

A dsl is provided by the backend to perform searches.
You can find it here : https://git.canopsis.net/canopsis/canopsis/blob/develop/sources/python/alerts/etc/alerts/search/grammar.bnf

Finaly, you can find some general informations about searches here : https://git.canopsis.net/canopsis/canopsis/blob/develop/doc/sakura/FR/fr__alarms_tray.rst#search-dsl

Before sending a query to the default route, you need to validate the expression provided by users.  
Once it is validated, you can perform search by using the default route.  
If it's not validated, you must inform user of that. A message telling about the wrong expression

**TO BE COMPLETED BY FLO** => Donner les infos de la route à appeler avec ses paramètres
alerts/search/validate?expression=<EXPRESSION>



Action buttons
^^^^^^^^^^^^^^

In the widget, a column must be dedicated to user actions.

In the widget conf form, there must be a checkbox to do such a thing.
Actions are shown only if the user is authorized to. Don't forget to include this contraint.

Here are available actions :

* ACK / FastACK / UnACK  (glyphicon-saved / glyphicon-ok / glyphicon-ban-circle)
* Declareticket (fa-ticket)
* Assocticket (fa-thumb-tack)
* Cancel alarm (glyphicon-trash)
* Change criticity (fa-exclamation-triangle)
* Restore Alarm (glyphicon-share-alt)
* Snooze alarm (fa-clock-o)

Each action is associated with a font

Executing an action is the same thing as sending an event.

Action forms must be picked from the actual "list widget".  
For example, ACK form look like thie :

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/ackform.png


Massive actions can be performed too by seclecting multiple alarms 


.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/massiveactions.png

**TO BE COMPLETED BY FLO**


**Rules that apply to actions**

* Except **snooze action**, all actions apply to ack'ed alarms
* **Restore Alarm** apply to Cancelled alarms


Info Popup
^^^^^^^^^^

When set, a popup can be displayed by clicking in a cell.
Popup results from a template compilation which can be defined by the user.

The user must be able to set multiple infopopup on multiple columns.

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/recordinfopopup.png


Linklist
^^^^^^^^

As said before, linklists are links with categories that are attached to an entity. 
The widget has to display it like on screenshots


.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/linklistrender.png


Pbehavior
^^^^^^^^^

The widget must be able to display pbehaviors if there is some.  
Pick an icon from library and make a renderer for that.  
Pbehaviors must be displayed like **ack** or **ticket**


Timeline
^^^^^^^^

The TR you have to show in the main table describe an alarm.
There are many other informations available by calling another webservice, **steps**.

In the main table, each tr must show a "+" that will call a component that represent steps.

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/timeline.png


**TO BE COMPLETED BY FLO** Comment instancie t-on le composant timeline ?

Live Reporting
^^^^^^^^^^^^^^


In Canopsis, users are able to select data that fit timeperiod.  

User first clicks on 

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/livereporting1.png


And then, he selects period

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/raw/master/doc/screenshots/livereporting2.png


**From** and **to** are then provided to the widget as timestamps

**TO BE COMPLETED BY FLO**

Glossary
--------

.. code-block::

    Entity
        An entity is a config item in Canopsis with a type.
        Type could be `component`, `resource`, `selector`

    Schema
        A schema is a way to represent data.
        In Canopsis, schemas are in JSON format
