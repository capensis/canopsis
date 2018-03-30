.. _TR__Statistics:

=============================
Canopsis Statistics Generator
=============================

This document specifies how the Statistics Generator is implemented.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Statistics <FR__Statistics>`
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

Stats engine
------------

This engine does detect, process and tag 3 different objects :
  - expired sessions (via session manager)
  - opened alarms
  - closed alarms

.. _TR__Statistics__Engine:

Expired sessions
~~~~~~~~~~~~~~~~

Expired sessions are sessions for which now() - ``last_check`` is greater than
``alive_session_duration`` (property in ``PREFIX/etc/session/session.conf``).

A session is in this state if ``/keepalive`` route has not been recently
called.

Those sessions are closed by the engine by setting ``active`` to ``False`` and
``session_stop`` to now().

.. _TR__Statistics__Sessions:

Opened and closed alarms
~~~~~~~~~~~~~~~~~~~~~~~~

On each beat, the engine retrieves alarms that are not tagged with
``stats-opened`` on one hand and alarms not tagged with ``stats-resolved``
while being resolved on the other hand.

The first category allows to produce ``alarm_opened_count`` metric and the
second one ``alarm_ack_delay``, ``alarm_ack_solved_delay`` and
``alarm_solved_delay`` metrics. The engine tags corresponding alarms
afterwards.

.. info:: Metrics that are computed when an alarm is closed are stored with
  now() as timestamp. It implies if an alarm is closed 3 days after its
  acknowledgement, the ack will be recorded 3 days a posteriori.

.. info:: The four ``alarm_*`` metrics have an attached value. This value is
  worth a duration, except for alarm_opened_count, for which a duration is not
  relevent (only timestamp matters). For homogeneity purpose, a value is
  stored, but it's always 1.

.. _TR__Statistics__Alarms:
