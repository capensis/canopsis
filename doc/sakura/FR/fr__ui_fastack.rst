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

Fast ack will have the exact same features of Normal ACK :ref:`ACK <FR__Ack>`.

The puporse grant the ability to ack an event without filling the popup information. 

The Fast ACK feature will can be acces throught a button near the normal ACK button. 
When Fast ACK button is click the will produce and ack with no ticket and a configured message.

Like normal ack, fast ack will have a right on it. To grant acces to thie feature you will need to fist grand acces to normal ack then to fast ack.

An event in Canopsis is the representation of asynchronously incoming data, sent by
a :ref:`connector <FR__Connector>`.

Event Acknowledgment
~~~~~~~~~~~~~~~~~~~~

An ``ack`` event **MUST** contain:

 - an author 
 - a message "auto ACK""
 - a reference to the :ref:`check event <FR__Event__Check>` to acknowledge

.. _FR__Event__Ackremove: