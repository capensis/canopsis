.. _TR__Monitoring:

==========
Monitoring
==========

This document describes the Monitoring requirements of Canopsis.

.. contents::
   :depth: 2

References
==========

List of referenced functional requirements:

 - :ref:`FR::Monitoring <FR__Monitoring>`
 - :ref:`FR::Event <FR__Event>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2016/03/14", "0.2", "Document validation", "Florent Demeulenaere"
   "David Delassus", "2016/03/08", "0.1", "Document creation", ""

Contents
========

.. _TR__Monitoring__Check:

Monitoring check
----------------

A monitoring check is composed of:

 - a *Nagios* check
 - a schema named ``monitoringplugin.<CHECKNAME>.json`` in ``~/etc/schema.d``:

.. _TR__Monitoring__Schema:

Command schema
--------------

The schema **MUST** contains the following meta informations:

 - ``command``:
    - ``binpath``: path to the *Nagios* check executable
    - ``args``: mapping of *Nagios* check command line arguments to schema fields

.. _TR__Monitoring__Command:

Command definition
------------------

A command is a JSON file matching the :ref:`schema <TR__Monitoring__Schema>`, and
**MUST** be located in ``~/etc/monitoring/commands/<COMMAND-NAME>.json``.

.. _TR__Monitoring__Runner:

Check runner
------------

The check runner **MUST** receive a command name, and:

 - using the :ref:`command definition <TR__Monitoring__Command>`, run the *Nagios* check
 - parse the *Nagios* check output to produce a :ref:`check event <FR__Event__Check>` with:
    - ``state``: the exit code of the check
    - ``output``: the message produced by the check
    - ``perfdata``: data produced after the message by the check
    - ``long_output``: the rest of the output
