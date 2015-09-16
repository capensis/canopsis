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

   "David Delassus", "2015/09/16", "0.3", "Add widget and view", ""
   "David Delassus", "2015/09/16", "0.2", "Update frontend description", ""
   "David Delassus", "2015/09/16", "0.1", "Document creation", ""

Contents
========

.. _FR__Briefcase__Desc:

Description
-----------

In Canopsis, the briefcase allows to associate one or more files to an :ref:`entity <FR__Context__entity>`.

All associated documents are located in a :ref:`file storage <FR__Architecture__DataStore>`.

Frontend
--------

.. _FR__Briefcase__Component:

Briefcase component
~~~~~~~~~~~~~~~~~~~

The frontend MUST provide a component to explore associated files and folders, with the
following layouts :

 * list view : each file is listed with its full path
 * tree view : each folder is expandable/collapsible
 * icon view : each file and folder is represented as an icon with the basename bellow

The component is configured with a filter (on entities and file(s) metadata), and
also provide a search-box to filter at runtime.

.. _FR__Briefcase__Widget:

Briefcase widget
~~~~~~~~~~~~~~~~

A widget, using the Briefcase component, SHOULD be available to let the user customize
the filter used by the component.

.. _FR__Briefcase__View:

Briefcase view
~~~~~~~~~~~~~~

A view containing the widget with no particular filter SHOULD be available in order
to offer a global view on all associated file(s). 

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
