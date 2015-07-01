Statistics
==========


Canopsis computes in background some statistics that allow users to get information about canopsis usage.
Below, the statistics list and **the way they have to be computed** to be available in Canopsis.

The main goal of these statistics is for a canopsis administrator to be able to plot and display them in **charts and widget texts**.
Some users should be able to have a look at the charts and widget texts holding these information.

The following metrics expects to be computed over time, accurately and with no double computation.


 - Delay between a user connection and deconnection

   A backend Session manager allow compute required to be able to produce these metrics

 - Acknowleged alarm count by user, perimeter and domain

   Should produce a metric by single **[user, perimeter, domain]**

 - Delay between an alarm and a user acknowlegement. Also compute global value (for all users)

   The produced metrics are for individual users. The global one is a sum of all users metrics

 - New alarms count

 - Acknowleged alarms

 - Solved alarm with acknowlegement

 - Solved alarm without acknowlegement

 - Delay between an alarm and it's resolution date for acknowleged alarms


 It must also be possible to perform aggregation operation on a custom period with a start and stop date such as:

 - minimum
 - maximum
 - mean
 - sum

This should be acheived thanks to the serie's system.

