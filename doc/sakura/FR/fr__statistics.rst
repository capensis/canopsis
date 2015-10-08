.. _FR__Title:

===========
Statisctics
===========

This document explains Canopsis statistic design.

.. contents::
   :depth: 2

References
==========

- :ref:`FR::Statistics <FR__Statistics>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Eric Régnier", "2015/10/02", "1.0", "Initial functional requirements", ""

Contents
========

.. _FR__Statistics:

Description
-----------

Statistics topic in canopsis consists in internal metric generation about Canopsis activity. This allow Canopsis users to get information about the global platform usage.

Canopsis comes with the idea to be monitorable. This is why it is intersting for canopsis users to be able to have a great understanding about canopsis internal activity. This activity has a wide range of purpose such as observing user acknowlegement workflow management by canopsis users or can just consists in computing some performance indicators.

The way canopsis monitors it's activity is by generating events that contains metric information. These events are published into the canopsis entry AMQP entry point where the event is processed by the perfdata engine. This engine is in charge to extract metric information, store and make them available through the Canopsis API for Canopsis clients such as the Canopsis UI or any program that may query the API.

