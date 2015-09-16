.. _FR__Briefcase:

=========
Briefcase
=========

This document describes the briefcase functionality in Canopsis, from backend to
frontend.

.. contents::
   :depth: 2

References
==========

List of referenced functional requirements...

- :ref:`FR::Context <FR__Context>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/09/16", "0.1", "Template creation", ""

Contents
========

.. _FR__Briefcase__Desc:

Description
-----------

In Canopsis, the briefcase allows to associate one or more files to an :ref:`entity <FR__Context__entity>`.

All associated documents are located in a :ref:`file storage <FR__Architecture__DataStore>`.

.. _FR__Briefcase__View:

Frontend view
-------------

The frontend MUST provide a component to explore associated files and folders, with the
following layouts :

 * tree view
 * icon view
 * list view

This component must be filterable by entities.

.. _FR__Briefcase__API:

Backend API
-----------

The backend MUST provide :

 * a briefcase manager
 * a webservice

.. _FR__Briefcase__Manager:

Manager
~~~~~~~

The manager MUST provide the following API :

 * an add method, receiving a file and an entity to associate the file to
 * a fetch method, filterable by entities
 * an update method :
    * to rename/move the file in the entity's *folder*
    * to move a file from an entity to another
 * a remove method, receiving a filter which will remove all matching files

.. _FR__Briefcase__Webservice:

Webservice
~~~~~~~~~~

The webservice will provide a single route with distinct verbs for each method of
the manager.
