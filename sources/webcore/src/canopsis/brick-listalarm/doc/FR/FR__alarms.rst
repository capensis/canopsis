FR alarms
---------

.. contents:: Table of contents


.. image:: https://git.canopsis.net/canopsis-ui-bricks/brick-alarms/blob/master/doc/screenshots/general_render.png


Synopsis
========

Create a new widget list with the same features than the first one.
Priority is focused on performance : less than 1s to render 50 alarms.


Why ?
=====

Actual list widget is full featured and generic but suffers from a lack of performance. ~4s to render 50 alarms.


Canopsis needs to get a new one with the functional requirements than the existing one plus better performance. ~1s to render 50 alarms.

The new widget style must be the same than the existing one.


Have a look at the view.event view in Canopsis => http://alarms.canopsis.net/en/static/canopsis/index.html#/userview/view.event

The main used theme is still **adminlte** => https://almsaeedstudio.com/themes/AdminLTE/index2.html


Features
========

- Create a Canopsis widget that displays alarms details.

- Give the possibility to execute some actions on alarms

    - Ack
    - Ticket
    - Pbehaviors

- This widget must be responsive
