.. _FR__Event:

==============
Canopsis Event
==============

This document describes the structure of events in Canopsis.

.. contents::
   :depth: 2

----------
References
----------

List of referenced function requirements:

 - :ref:`FR::Connector <FR__Connector>`
 - :ref:`FR::Schema <FR__Schema>`
 - :ref:`FR::Context <FR__Context>`
 - :ref:`FR::Metric <FR__Metric>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/06", "0.1", "Document creation", "Jonathan Labéjof"

--------
Contents
--------

Description
-----------

An event in Canopsis is the representation of asynchronously incoming data, sent by
a :ref:`connector <FR__Connector>`.

They are described by a **CEvent** :ref:`data schema <FR__Schema__Data>` per type of event.

Each event contains typed informations (listed bellow), associated to one or more :ref:`entities <FR__Context__Entity>`.

Logging events
--------------

Those events are used to log informations to Canopsis.

.. _FR__Event__Log:

Event Log
~~~~~~~~~

A ``log`` event **MUST** contain:

 - a message
 - a severity level

And it **MAY** contain:

 - a detailed message

.. _FR__Event__User:

Event User
~~~~~~~~~~

A ``user`` event **MUST** contain:

 - an author
 - a message
 - a criticality level

And it **MAY** contain:

 - a detailed message

.. _FR__Event__Comment:

Event Comment
~~~~~~~~~~~~~

A ``comment`` event **MUST** contain:

 - an author
 - a message
 - a criticality level
 - a reference to the commented :ref:`user event <FR__Event__User>`

And it **MAY** contain:

 - a detailed message

Supervising events
------------------

Those events are used to store changes in a supervision environment to Canopsis.

.. _FR__Event__Check:

Event Check
~~~~~~~~~~~

A ``check`` event **MUST** contain:

 - a state
 - a message

And it **MAY** contain:

 - a state specification
 - a detailed message

.. _FR__Event__Selector:

Event Selector
~~~~~~~~~~~~~~

A ``selector`` event **MUST** contain:

 - a state
 - a message

And it **MAY** contain:

 - a displayed name
 - a state specification
 - a detailed message

.. _FR__Event__Trap:

Event Trap
~~~~~~~~~~

A ``trap`` event **MUST** contain:

 - a state
 - a severity level
 - an :ref:`OID <FR__SNMP__OID>`

.. _FR__Event__Changestate:

Event ChangeState
~~~~~~~~~~~~~~~~~

A ``changestate`` event **MUST** contain:

 - an author
 - a state
 - a message
 - a reference to the :ref:`check event <FR__Event__Check>` to modify

.. _FR__Event__Downtime:

Event Downtime
~~~~~~~~~~~~~~

A ``downtime`` event **MUST** contain:

 - an author
 - a message
 - a period

.. _FR__Event__Cancel:

Event Cancel
~~~~~~~~~~~~

A ``cancel`` event **MUST** contain:

 - an author
 - a message
 - a reference to the :ref:`check event <FR__Event__Check>` to cancel

.. _FR__Event__Uncancel:

Event Uncancel
~~~~~~~~~~~~~~

An ``uncancel`` event **MUST** contain:

 - an author
 - a message
 - a reference to the canceled :ref:`check event <FR__Event__Check>`

Ticketing events
----------------

Those events are used to represent interactions with a CMDB.

.. _FR__Event__Declareticket:

Event Declareticket
~~~~~~~~~~~~~~~~~~~

A ``declareticket`` event **MUST** contain:

 - an author
 - a message
 - a reference to the :ref:`check event <FR__Event__Check>` to create a ticket for

.. _FR__Event__Assocticket:

Event Assocticket
~~~~~~~~~~~~~~~~~

A ``declareticket`` event **MUST** contain:

 - an author
 - a message
 - a ticket ID
 - a reference to the :ref:`check event <FR__Event__Check>` to assign the ticket to

Acknowledging events
--------------------

Those events are used to manage supervising events.

.. _FR__Event__Ack:

Event Acknowledgment
~~~~~~~~~~~~~~~~~~~~

An ``ack`` event **MUST** contain:

 - an author
 - a message
 - a reference to the :ref:`check event <FR__Event__Check>` to acknowledge

.. _FR__Event__Ackremove:

Event Ackremove
~~~~~~~~~~~~~~~

An ``ackremove`` event **MUST** contain:

 - an author
 - a message
 - a reference to the :ref:`check event <FR__Event__Check>` to *unacknowledge*

Performance events
------------------

Those events are used to store metrics in Canopsis.

.. _FR__Event__Perf:

Event perf
~~~~~~~~~~

A ``perf`` event **MUST** contain:

 - one or more :ref:`metric <FR__Metric>`

**NB:** ``perf`` events are not stored and can be included in every other events.

Connector events
----------------

Those events are used to control remote :ref:`connectors <FR__Connector>`.

Event enable
~~~~~~~~~~~~

An ``enableconnector`` has no supplementary informations.

Event disable
~~~~~~~~~~~~~

A ``disableconnector`` has no supplementary informations.

Event getconf
~~~~~~~~~~~~~

A ``getconfconnector`` has no supplementary informations.

Event setconf
~~~~~~~~~~~~~

A ``setconfconnector`` has no supplementary informations.

Event getstate
~~~~~~~~~~~~~~

A ``getstateconnector`` has no supplementary informations.
