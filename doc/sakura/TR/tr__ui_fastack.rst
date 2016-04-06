.. _TR__UI_FastACK:

====================
UI Fast ACK Function
====================

This document describes the Fast Ack UI Button

References
==========

 - :ref:`FR::Architecture <FR__Architecture>`
 - :ref:`FR::Engine <FR__Engine>`
 - :ref:`TR::Package <TR__Package>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Vincent CANDEAU", "2016/04/04", "0.1", "Document creation", ""

Contents
========

Description
-----------

This Button will produce an ack event configured in its :ref:`FastACk <FR__FastACK>`.

ACK Configuration:

It **MUST** contains:

 - a sender, composed of:
 - an message that will be **auto ACK**
 - an event ref
 
List of file modified to add this feature :
 - ACK Template ( uibase => template => ack )
 - ACK Selection template ( uibase => template => ackselection )
 - Mixin send event ( monitoring => mixins => sendevent )
 - JSON Right list
 - Init Schema file of Monitoring UI Brick
 - Add mixin.sendevent.json in UI Brick
 
 
Case: OK
~~~~~~~~~~~~~~~~~~~~~~~~~~~

- An new label state on the UI will be show instead of the ack button.

Case: KO
~~~~~~~~~~~~~~~~~~~~~~~~~~~

 - a message explaining the error
