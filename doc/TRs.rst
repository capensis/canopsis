TR alarms
---------

.. contents:: Table of contents


Configuration
=============

The widget configuration is organized around menus which look like this:

* General
* Columns
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


Tab 3 : Mixins
^^^^^^^^^^^^^^

This tab is automaticaly added with the schema def.


Data schema
===========

In order for you to understand what is needed to display table, here is the data schema with explanations.  
All needed informations are described in it. If it's not clear, feel free to contact us.  

You have to know that data are provided by an adapter which relies on a webservice.  
In order to use only frontend concepts, we provide 1 dataset (in floder `datasets`) which can be loaded by the adapter.  


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

 {{ linklist category="procedure" }}
 
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
   

Action buttons
^^^^^^^^^^^^^^

In the widget, actions must be shown to the user by a group of buttons that must be callable by a handlebar renderer. 

 {{ actions list="confirm invalidate declareticket pause" }}
 
must show

.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-service-weather/raw/master/doc/screenshots/eventsactions1.png

Actions definitions must be included into a dedicated `ember mixin`.

**To do such a thing, have a look at existing mixins in the project**



Glossary
--------

.. code-block::

    Entity
        An entity is a config item in Canopsis with a type.  
        Type could be `component`, `resource`, `selector`

    Schema
        A schema is a way to represent data.  
        In Canopsis, schemas are in JSON format
 

