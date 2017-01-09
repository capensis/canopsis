FR alarms
---------

.. contents:: Table of contents


.. image:: ../_static/


Synopsis
========

Create a new widget list with the same features than the first one.  
Priority is focused on performance : less than 1s to render 50 alarms.


Why ?
=====

Actual list widget is full featured and generic but suffers from a lack of performance. ~4s to render 50 alarms.  
Canopsis needs to get a new one with the functional requirements than the existing one plus better performance. ~1s to render 50 alarms


Features
========

- Create a Canopsis widget that displays alarms details.

- Give the possibility to execute some actions on alarms

    - Ack
    - Ticket
    - Pbehaviors
    - Linklist
    
- This widget must be responsive

