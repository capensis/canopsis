Periodic Behavior management
============================

Introduction
------------

Entities can have specific behavior during a period of time:

 * the entity is in downtime (checks and metrics doesn't interfere with SLA)
 * the entity is disabled (no checks or metrics are saved)
 * the entity is in benchmark (no alerts generated when threshold are exceeded)
 * ...

A periodic behavior is:

 * a period (date-time of start, date-time of end, recursion rule, ...) in iCal format
 * a behavior (downtime, disable, benchmark, ...)
 * the entity it acts on

The Periodic Behavior Manager is in charge of storing and retrieving periodic
behavior informations.

Read/Write operations
---------------------

Get all periodic behaviors available for one entity:

.. code-block:: python

   from canopsis.pbehavior.manager import PBehaviorManager, get_query

   manager = PBehaviorManager()

   entity_id = '...'
   pbehaviors = manager.values(
       sources=[entity_id],
       query=get_query('downtime')
   )

Add a new periodic behavior to an entity:

.. code-block:: python

   from canopsis.pbehavior.manager import PBehaviorManager
   from datetime import datetime, timedelta
   from icalendar import Event

   manager = PBehaviorManager()

   entity_id = '...'
   pstart = 0
   pend = 1

   ev = Event()
   ev.add(manager.BEHAVIOR, 'downtime')
   ev.add('summary', 'host update')
   ev.add('dtstart', datetime.fromtimestamp(pstart))
   ev.add('dtend', datetime.fromtimestamp(pend))
   ev.add('dtstamp', datetime.now())
   ev.add('rrule', {'FREQ': 'WEEKLY'})
   ev.add('duration', timedelta(3600))
   ev.add('contact', 'John Smith')

   manager.put(source=entity_id, vevents=[ev])

Determine if an entity is in a periodic behavior:

.. code-block:: python

   from canopsis.pbehavior.manager import PBehaviorManager, get_query

   manager = PBehaviorManager()

   entity_id = '...''
   result = manager.get_after(
       sources=[entity_id],
       query=get_query('benchmark')
   )

   if 'benchmark' in result:
       print('End of benchmark:', result['benchmark'])

Get all entities which are actually in a periodic behavior:

.. code-block:: python

   from canopsis.pbehavior.manager import PBehaviorManager, get_query

   manager = PBehaviorManager()
   entity_ids = manager.whois(query=get_query('downtime'))
