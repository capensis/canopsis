Statistics
==========


Canopsis computes in background some statistics that allow users to get information about canopsis usage.
Below, the statistics list and **the way they have to be computed** to be available in Canopsis.

As metrics have to be displayed with arbitrary time periods, agregation information have to be computed on the fly. It may be interesting to make evolve the metric system to be able to retrieve information from produced metrics instead of producing more metrics that would use disk space. This point of view have to be moderated by the fact that performance issue that may result with some too high level on the fly operations.

For instance, it may be possible to solve the cps_user_ack_delay metric on global (sum of all users metrics) thanks to an appropriate api that allow such operation

The main goal of these statistics is for a canopsis administrator to be able to plot and display them in **charts and widget texts**.
Some users should be able to have a look at the charts and widget texts holding these information.

The following metrics expects to be computed over time, accurately and with no double computation.


 - Delay between a user connection and deconnection

   ``[cps_session_delay: seconds, integer]`` A backend Session manager allow compute required to be able to produce these metrics

 - Acknowleged alarm by user, perimeter and domain

   ``[cps_ack_alarm: integer]`` Should produce a metric by single **[user, perimeter, domain]**

 - Delay between an alarm and a user acknowlegement. Also compute global value (for all users)

   ``[cps_user_ack_delay: seconds, integer]`` The produced metrics are for individual users. The global one is a sum of all users metrics

 - New alarms count

   ``[cps_new_alarms: integer]``

 - Acknowleged alarms

   ``[cps_ack_alarm: integer]``

 - Solved alarm with acknowlegement

   ``[cps_solved_ack_alarms: integer]``

 - Solved alarm without acknowlegement

   ``[cps_solved_not_ack_alarms: integer]``

 - Delay between an alarm and it's resolution date for acknowleged alarms

   ``[cps_delay_ack_alarm_solve: integer]``


 It must also be possible to perform aggregation operation on a custom period with a start and stop date such as:

 - minimum
 - maximum
 - mean
 - sum

This should be acheived thanks to the serie's system.

