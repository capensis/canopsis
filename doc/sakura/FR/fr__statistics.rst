.. _FR__Title:

==========
Statistics
==========

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

Indicators
----------

This is a list of the computed indicators that are executed depending in it's nature. This table also describe how to functionnaly produce metrics.

.. csv-table::
   :header: "Component", "Resource", "Metric", "Type", "Production"

   "__canopsis__",  "Engine_stats", "cps_session_delay_user_<username>", "gauge", "Login with a user then logout for at least 5 minutes to produce the metric for the current user."
   "#AUTHOR", "ack", "alerts_count[_<domain><perimeter>]", "counter", "Produce an ack for an alert event depending on whether domain information exists. This metric is incremented when an alert is acknowleged."
   "#AUTHOR", "ack", "delay", "gauge", "Acknowlege an event after an event is on alert. The delta between the alert and the ack define the value of the metric."
   "solved_alarm", "ack", "count", "counter", "A solved alert (state back to 0) that was acknowleged will increment the metric value."
   "solved_alert", "ack", "delay", "gauge", "A solved alert (state back to 0) that was acknowleged will increment the duration metric value."
   "__canopsis__", "ack", "cps_solved_ack_alarms", "counter", "A solved alarm (back to 0 state) that was ack will increment the metric value."
   "__canopsis__", "ack", "cps_solved_not_ack_alarms", "counter", "A solved alarm (back to 0 state) that was NOT ack will increment the metric value."


Performance tests
-----------------

Testing performance for this feature consists in check that canopsis amqp queues does not fall when many ack are sent to the canopsis instance. Considering that Canopsis at the moment can manager about 700 events/s, the ack proportion of these event made by users have to be maximised for performance test issue.
