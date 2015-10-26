.. _TR__Alarm:

================
Alarm management
================

This document specifies the alarm management in Canopsis, and its implementation.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Alarm <FR__Alarm>`
 - :ref:`FR::Configuration <FR__Configurable>`
 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Storage <FR__Storage>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/26", "0.1", "Document creation", ""

Contents
========

.. _TR__Alarm__Manager:

Alerts manager
--------------

An ``Alerts`` :ref:`configurable <FR__Configurable>` provides:

 - the ability to archive an :ref:`event <FR__Event>`
 - alarm cycle management operations:
    - create a new one
    - update existing one
    - get last one
    - find old alarms
    - tagging

.. _TR__Alarm__Cycle:

Alarm Cycle
-----------

Alarm cycles are persisted in a :ref:`timed storage <FR__Storage_Type>` with the following informations:

 - data identifier: the entity identifier of the received event
 - value: set of alarm steps (see :ref:`data model <TR__Alarm__DataModel>` bellow)
 - timestamp: date/time of alarm appearance

.. _TR__Alarm__DataModel:

Alarm data model
----------------

The set of alarm steps will compute informations for an easier use:

.. csv-table::
   :header: Field, Description

   state, ``stateinc`` step or ``statedec`` step or ``changestate`` step or ``None``
   status, ``statusinc`` step or ``statusdec`` step or ``None``
   ack, ``ack`` step or ``None`` if unacknowledged
   canceled, ``cancel`` step or ``None`` if uncanceled
   ticket, ``declareticket`` step or ``assocticket`` step or ``None`` if no ticketing informations
   resolved, timestamp of resolution or ``None`` if alarm is still on going
   steps, array of steps (see data models bellow)
   tags, array of tag (as strings) used for filtering

.. _TR__Alarm__DataModel__StateInc:

Alarm step "state increase" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``stateinc``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new state

.. _TR__Alarm__DataModel__StateDec:

Alarm step "state decrease" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statedec``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new state

.. _TR__Alarm__DataModel__StatusInc:

Alarm step "status increase" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statusinc``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new status

.. _TR__Alarm__DataModel__StatusDec:

Alarm step "status decrease" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statusdec``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new status

**NB:** if status decreases to ``OFF``, then the alarm value ``resolved`` is set to this step timestamp.

.. _TR__Alarm__DataModel__Acknowledge:

Alarm step "acknowledge" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``ack``
   t, step timestamp, ``event["timestamp"]``
   a, step author, acknowledgment author
   m, step message, acknowledgment message

.. _TR__Alarm__DataModel__Unacknowledge:

Alarm step "unacknowledge" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``ackremove``
   t, step timestamp, ``event["timestamp"]``
   a, step author, acknowledgment removal author
   m, step message, acknowledgment removal message

**NB:** this step reset the alarm value ``ack`` to ``None``.

.. _TR__Alarm__DataModel__Cancel:

Alarm step "cancel" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``cancel``
   t, step timestamp, ``event["timestamp"]``
   a, step author, alarm canceling author
   m, step message, alarm canceling message

.. _TR__Alarm__DataModel__Restore:

Alarm step "restore" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``uncancel``
   t, step timestamp, ``event["timestamp"]``
   a, step author, alarm restoring author
   m, step message, alarm restoring message

**NB:** this step reset the alarm value ``cancel`` to ``None``.

.. _TR__Alarm__DataModel__DeclareTicket:

Alarm step "declare ticket" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``declareticket``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ticket declaration author
   m, step message, ticket declaration message
   val, ticket number, ``None``

.. _TR__Alarm__DataModel__AssocTicket:

Alarm step "associate ticket" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``assocticket``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ticket association author
   m, step message, ticket association message
   val, ticket number, ``event["ticket"]``

.. _TR__Alarm__DataModel__ChangeState:

Alarm step "change state" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``changestate``
   t, step timestamp, ``event["timestamp"]``
   a, step author, state requalification author
   m, step message, state requalification message
   val, state, new state
