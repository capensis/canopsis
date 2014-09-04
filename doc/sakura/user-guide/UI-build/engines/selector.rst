.. _selector:

Selector
========

Overview
---------

Selector engine aims to perform aggregation operation on many events,
providing this way a simple aggregated information.

Goals
-----

Allow user to define it's own aggregation information. For now, selector
allows to compute the worst state on an event set.

--------------

Synchronous consume dispatch method
------------------------------------

This method listen to the rabbitmq bus and each time an event becomes
available, a crecord is retrieved from this event. The crecord provide
infomation about witch event implied in the information aggregation.
information are:

-  a list of event id to **include**
-  a list of event id to **exclude**
-  a cfilter describing a mongo style database query which allows
   **custom** event selection

These information make selector compute a set of distinct status values
for these event. From these distinct values, the worst state appears in
a natural way and then, a new event is produced with the selector
information and some other information relative to the selector
description made through the UI.
