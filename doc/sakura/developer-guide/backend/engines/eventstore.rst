Event store
===========


The event store engine is a legacy system that makes canopsis store/update events into database. This event evolved a lot and recently integrated some features such as bulk event processing. Bulk event processing allow great perfomance increase for database update / store operations.

Good things
-----------

Event insert/update is done thanks to the ``archiver`` python library class and the ``eventstore`` engine itself. Incoming events that have to be stored into database are **queued in memory** until some **conditions** are reached. When conditions are reached, the whole in memory events are processed and depending on their content, are stored or updated.

Processing conditions are the following:

* When queued event count reached a limit then the queue is processed. in the archiver, the queue limit is called ``bulk_limit`` and in the event store it is called ``log_bulk_amount``

* When the last bulk queue processing were at least ``bulk_delay`` seconds ago in archiver and ``log_bulk_delay`` in eventstore, then the whole in memory event queue are processed.

Events to process in event store are **log** event type that have to be stored in ``events log`` collection, whereas in the archiver, event processed are those who have to be inserted if not exists or updated in the ``events`` collection and under certain conditions (state change for instance that produces an alert), into ``events log`` collection.

Bad stuff
---------

At the moment, a lot of business processing is done into the archiver as it is the only engine where the **previous event is fetched** from database and where it is possible to perform comparisons.

However, it would be a bad idea to load previous event from database each time an event is processed in each engine. That's why we have such business code that will in the future be refactored to a better designed location.

In these business processing that needs comparison that we could not put somewhere else at the moment, we get event status processing, acknowlegement computation and some statements cleans depending on the previous statement.
