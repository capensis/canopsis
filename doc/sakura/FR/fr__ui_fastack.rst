.. _FR__UI_FastACK:

==============
Canopsis Event
==============

This document describes the Fast Ack UI Button

.. contents::
   :depth: 2

----------
References
----------


-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Vincent CANDEAU", "2016/04/04", "0.1", "Document creation", ""

--------
Contents
--------

Description
-----------

The event list get an UI ACK Button that popup a form in order to send ack on an event. 
This new features will do the same thing faster. The popup will not be show and the ACK event will be sent directly.  

An event in Canopsis is the representation of asynchronously incoming data, sent by
a :ref:`connector <FR__Connector>`.

Event Acknowledgment
~~~~~~~~~~~~~~~~~~~~

An ``ack`` event **MUST** contain:

 - an author 
 - a message "auto ACK""
 - a reference to the :ref:`check event <FR__Event__Check>` to acknowledge

.. _FR__Event__Ackremove: