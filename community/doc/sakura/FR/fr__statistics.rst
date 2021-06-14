.. _FR__Statistics:

==========================
Canopsis Statistics Engine
==========================

This document describes the statistics generation process.

.. contents::
   :depth: 3

References
==========

List of referenced functional requirements:

 - :ref:`FR::Alarm <FR__Alarm>`
 - :ref:`FR::Engine <FR__Engine>`
 - :ref:`FR::Webservice <FR__Webservice>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jean-Baptiste Braun", "2016/22/11", "0.2", "New stats engine, new metrics", ""
   "David Delassus", "2015/10/15", "0.1", "Document creation", ""

Contents
========

Description
-----------

Stats are a feature of canopsis producing metrics about :ref:`alarms
<FR__Alarm>` activity and user sessions durations.

Stats **CAN** displayed on frontend thanks to a proprietary brick :
brick-statstable.

Information details about stats metrics is available in this brick.

.. _FR__Statistics__Desc:

Stats engine
------------

The stats engine has currently only a beat processing. On this beat, it
computes stats from ``*_alarm`` and ``session`` mongo collections, storing
datas in influxdb.

Those condensed data can be retrieved much faster, and quickly filtered via
authors or custom tags.

.. _FR__Statistics__Engine:

Custom tags
~~~~~~~~~~~

:ref:`Alarms <FR__Alarm>` can be tagged via ``extra_fields``. Those extra
fields are exported as tags in influxdb.

.. _FR__Statistics__CustomTags:

Webservice
----------

There is a stats webservice serving data to UI via the brick statstable.

In order to enable it (not enabled by default), add the line ``stats=1`` in
the ``webservices`` section of ``PREFIX/etc/webserver.conf`` and restart
webserver.

.. _FR__Statistics__Webservice:

Functional Tests
================

Functionnal tests have been implemented to check most of produced metrics. But
since there is currently no place in canopsis deployment workflow to insert
it, this work has not been pushed.

.. warning::

   **TODO:** As it is difficult to simulate user authentication and keepalive
   workflow, there are currently no functionnal tests for session duration
   metrics.
