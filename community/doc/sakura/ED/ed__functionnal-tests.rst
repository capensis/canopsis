.. _ED__FunctionnalTests:

=================
Functionnal Tests
=================

This document describes how to run the functionnal tests

.. contents::
   :depth: 2

References
==========

List of referenced functional requirements:

- :ref:`FR::Functionnal Tests <FR__FunctionnalTests>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2016/03/14", "1.0", "Document validation", "Florent Demeulenaere"
   "David Delassus", "2016/03/14", "1.0", "Document creation", ""

Contents
========

.. _ED__FunctionnalTests__Location:

Location
--------

All tests are located in ``~/var/lib/canopsis/functionnal-tests``, you **SHOULD**
be in that folder when running tests.

.. _ED__FunctionnalTests__RunAll:

Running all tests
-----------------

The command ``lettuce .`` will run all the tests located in the current folder.

The following commands are also available:

 - ``lettuce -r .``: run the scenarios in a random order to avoid interference
 - ``lettuce --with-xunit --xunit-file=report.xml .``: to create a jUnit XML report
 - ``lettuce --failfast .``: to stop testing on the first failure
 - ``lettuce --pdb .``: used to debug, will run an interactive debugger upon error

For more informations: ``lettuce --help``.

.. _ED__FunctionnalTests__RunOne:

Run specific test
-----------------

The command ``lettuce <feature-name>.feature`` will run all the scenario from a
specific test.
