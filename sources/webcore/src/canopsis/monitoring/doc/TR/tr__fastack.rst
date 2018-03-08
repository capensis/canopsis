.. _TR__UI_FastACK:

====================
UI Fast ACK Function
====================

This document describes the Fast Ack UI Button

References
==========

 - :ref:`FR::FastACk <FastACk>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Vincent CANDEAU", "2016/04/04", "0.1", "Document creation", ""

Contents
========

Description
-----------

This feature procude a configured ack event. 


This ack event will have :
 - No ticket number
 - Preconfigured message (Mixins Options)
 - Event ID ref

Files modification
------------------

ACK Template ( uibase => template => ack ):
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

    Add an small button with the icon glyphicon-saved. It has sendevent fastack action assign
    This button need the right fastack to be displayed

ACK Selection template ( uibase => template => ackselection ):
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

    Add an small button with the icon glyphicon-saved. It has sendevent fastack action assign
    This button need the right fastack to be displayed

Mixin send event ( monitoring => mixins => sendevent ):
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    fastack event_processors is the exact copy of the ack event_processors.
    The handle code of the fastack event_processors have been modified:
    
       - Remove the popup form function
       - Add "fastackmsg" mixin option retrieving
       - Set output with the value of fastackmsg
       - Submit select events for processing

JSON Right list:
~~~~~~~~~~~~~~~~

    Add the right actionbutton_fastack inside the json

Add mixin.sendevent.json inside Monitoring Brick:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

    Mixin option of the sendevent with the param fastackmsg.