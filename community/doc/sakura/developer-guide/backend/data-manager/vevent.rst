.. _dev-backend-mgr-vevent:

vEvent management
=================

Introduction
------------

In Canopsis, a **vEvent** is a document containing calendar informations about one or more entities.

Those informations are *iCalendar* expressions, and are described in the `RFC2445 <ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt>`_.

Read/Write operations
---------------------

According to the *iCalendar* documentation, a *vEvent* is identified by a *UID* (**U** nique **Id** entifier).

So, the manager allows you to get a set of *vEvent* if you know their UIDs :

.. code-block:: python

   from canopsis.vevent.manager import VEventManager


   manager = VEventManager()
   vevents = manager.get_by_uids(
       uids,
       # used to limit the number of returned documents
       limit=0,
       skip=0,
       sort=[('dtstart', 1)],
       projection={'id': 0, 'source': 1, 'vevent': 1},
       with_count=False
   )

In order to fetch the *vEvent* related to a set of entities, the manager provide a single method :

.. code-block:: python

   vevents = manager.values(
       sources=[entity_ids],
       # additional filter for storage
       query=None,
       # used to return documents in a period of time
       dtstart=None,
       dtend=None,
       # used to limit the number of returned documents
       limit=0,
       skip=0,
       sort=[('dtstart', 1)],
       projection={'id': 0, 'source': 1, 'vevent': 1},
       with_count=False
   )

