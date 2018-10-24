.. _FR__Engine:

===============
Canopsis Engine
===============

This document describes the functionality provided by a Canopsis engine.

.. contents::
   :depth: 2

References
==========

 - :ref:`FR::Event <FR__Event>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/06", "0.1", "Document creation", ""

Contents
========

.. _FR__Engine__Desc:

Description
-----------

An engine is a **daemon**, processing :ref:`events <FR__Event>`.

It **MUST** provide:

 - an :ref:`event <FR__Event>` processing algorithm
 - an algorithm to execute every ``X`` seconds (this **MUST** be configurable)
 - an algorithm to consume data from Canopsis
