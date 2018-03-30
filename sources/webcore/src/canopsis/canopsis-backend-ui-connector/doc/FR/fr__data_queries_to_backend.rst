.. _FR__CalendarWidget:

===========================
Data queries to the backend
===========================

This document describes the way data queries should be handled in the frontend.

.. contents::
   :depth: 3


References
==========

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Gwenael Pluchon", "2016/02/10", "0.1", "Document creation", ""

Contents
========

.. _FR__Title__Desc:

Description
-----------

Data managed from the frontend must respect some design rules :

- the data usage and the query algorithms must have the less coupling possible. Thus, the `adapter <https://en.wikipedia.org/wiki/Adapter_pattern>`_ pattern should be used, throughout `Ember Data Adapters <http://emberjs.com/api/data/classes/DS.Adapter.html>`_.
- Adapters should be implemented by source, not by data. It is not useful at all to implement an adapter by data type. An adapter must implement a way to communicate with a backend technology or version.
- Keep the choice of the adapter to a `strategy pattern <https://en.wikipedia.org/wiki/Strategy_pattern>`_. Do not pick manually an adapter to retreive a data. Thus, the adapter choice logic could be changed if necessary.
