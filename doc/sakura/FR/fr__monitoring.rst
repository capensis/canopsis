.. _FR__Monitoring:

==========
Monitoring
==========

This document describes the monitoring functionality in Canopsis.

.. contents::
   :depth: 2

References
==========

List of referenced functional requirements:

 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Schema <FR__Schema>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2016/03/14", "0.2", "Document validation", "Florent Demeulenaere"
   "David Delassus", "2016/03/08", "0.1", "Document creation", ""

Contents
========

.. _FR__Monitoring__Desc:

Description
-----------

Canopsis provides the ability to run *Nagios* checks and generate
:ref:`check events <FR__Event__Check>`.

.. _FR__Monitoring__Check:

Monitoring check
----------------

Each check **MUST** be provided and described by a
:ref:`data schema <FR__Schema__Data>`.

.. _FR__Monitoring__Runner:

Check runner
------------

The check runner will parse a command configuration file matching the associated
data schema, and execute the corresponding *Nagios* check.

It **MUST** produce a check event using the informations parsed from the check output.