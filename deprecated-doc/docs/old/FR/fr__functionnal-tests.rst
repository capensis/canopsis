.. _FR__FunctionnalTesting:

===================
Functionnal Testing
===================

This document describes the requirements for writing functionnal tests in Canopsis.

.. contents::
   :depth: 2

References
==========

List of referenced functional requirements:


Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2016/03/14", "0.2", "Document validation", "Florent Demeulenaere"
   "David Delassus", "2016/03/14", "0.1", "Document creation", ""

Contents
========

.. _FR__FunctionnalTesting__Desc:

Description
-----------

Functionnal tests are based on [lettuce](http://lettuce.it) and **MUST** describe:

 - each use case of a feature, as a scenario
 - each possible data sent to the feature

.. _FR__FunctionnalTesting__Howto:

How-to
------

Functionnal tests description are written in ``<feature-name>.feature`` files in
the folder ``features`` of the Python project.

A file named ``<feature-name>.py`` will be placed with the ``.feature`` file to
provide the steps implentation.

The name ``feature-name`` **MUST** be unique for the whole project, that means:

 - if 2 projects provide an engine, each feature must be named ``engine-<name>`` and not ``engine``
 - if 2 projects provide a task handler, each feature must be named ``taskhandler-<name>`` and not ``taskhandler``
 - if 2 projects provide a webservice, each feature must be named ``webservice-<name>`` and not ``webservice``
 - ...

The content of the ``features`` folder will be installed in ``~/var/lib/canopsis/functionnal-tests``.
The command ``lettuce`` will be called with that folder as argument to execute functionnal tests.
