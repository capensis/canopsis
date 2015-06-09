.. _dev-backend-storage-periodic:

Periodic Storage
================

This type of storage is used to store data that comes in periodically, for example :

 * the state of a component/resource : scheduled every 5 minutes ;
 * a performance data : transported by the check, scheduled every 5 minutes ;
 * the output of a component/resource...

Searches and deletions are made using a time-window (start/end timestamps).
But insertions and update are made using the data id.
