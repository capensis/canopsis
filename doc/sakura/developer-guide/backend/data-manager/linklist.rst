.. _dev-backend-mgr-linklist:

Link list system
================

You should first have a look at: `How to use link list system <../../../user-guide/engines/linklist.html>`_

Then you learned that link information may come from user association and event specific keys.

Architecture
------------

Link list system is made of two managers:

 - **Linklist** manager's goal is to persist user ``mfilter / link list`` information, with in mind that a link list is a list of dictionary with a label and an url value. That's as simple.
 - **Entitylink** manager's goal is to store current link association for each matched entity.

.. note::

   Database configuration for these manager are written to the **linklist/linklist.conf** file.

Thus, a link list database document have an **arbitrary id**, a **mfilter** and a **link list** information and should looks like:

.. code-block:: javascript

   {
      "_id" : "d4be3c5c-65cc-417c-bbf3-4283539122ab",
      "filterlink" : [
         {
            "url" : "http://canopsis.org",
            "label" : "canopsis"
         }
      ],
      "crecord_type" : "linklist",
      "mfilter" : "{\"$or\":[{\"component\":{\"$eq\":\"mycomponent\"}}]}",
      "name" : "my link list"
   }

This document configuration is generated from the UI and managed by the API python service witch calls the manager methods.
This particular configuration will lead to generate the following entity link information once a scheduled job is processed:

.. code-block:: javascript

   {
      "_id" : "/resource/myconnector/myconnector_name/mycomponent/myresource",
      "event_links" : [
         {
            "url" : "http://perdu.com",
            "label" : "action_url"
         }
      ],
      "computed_links" : [
         {
            "url" : "http://canopsis.org",
            "label" : "canopsis"
         }
      ]
   }

In this document, the ``event_links`` list is computed by the `scheduled jobs <engines/scheduler.html>`_ code, whereas, the ``computed_links`` list is computed by the linklist engines that watch for events key(label)/value(url) described by the **link_field** list property of the engine.

for exemple, when in the engine, **link_field** equals ['action_url'], for each event passing through the work method, ``if 'action_url' in event and event['action_url']`` then information is persisted/overriden by then entity link manager.


When entity link information are filled enough, in the UI, labelled link list are available in the event list widget/view (as shown in the user guide).
