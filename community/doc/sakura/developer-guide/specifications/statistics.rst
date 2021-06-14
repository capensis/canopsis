Statistics
==========


Canopsis computes in background some statistics that allow users to get information about canopsis usage.
Below, the statistics list and **the way they have to be computed** to be available in Canopsis.

The main goal of these statistics is for a canopsis administrator to be able to plot and display them in **charts and widget texts**.
Some users should be able to have a look at the charts and widget texts holding these information.

These metrics have to be computed in an **incremental** way. This means that every event holding an information that have to change any of these metrics will produce a new event with a count metric that will **increment or decrement** the total count of one of these indicators. This way, it is possible **to perform aggregation operation overtime** in the actual Canopsis system.

- Delay between a user connection and deconnection

A backend Session manager manage connection delays from the UI interactions and allow to compute this metric. Every beat interval of the engine stats, ended session are computed and for each of them, the delay between the session start and the session end is computed and produces a metric. This metric is a **GAUGE** as session duration information are absolute times.

    component: #HOSTNAME

    resource: Engine_stats

    metric: cps_session_delay_user_<username>

- Acknowleged alarm by user, perimeter and domain

Produce a metric by single **user, domain, perimeter** if domain and / or perimeter information are available. This indicator is incremented when an alert event is acknowleged and decremented when the alert retrieve it's normal state.
The produced event resource is ack.

    component: #AUTHOR

    resource: ack

    metric: alerts_count_[<domain><perimeter>]


- Delay between an alarm and a user acknowlegement

This metric is a count of ack delay that elapsed between the date when an event passed to alert state until when it is acknowleged.
The produced event resource is ack.

    component: #AUTHOR

    resource: ack

    metric: delay


- New alarms count

This indicator evolve in time with incrementation on an event passing to alert and decrementation when an alert goes back to ok state

    component: __canopsis__

    resource: ack

    metric: alert_count


- Solved alarm with acknowlegement


This metric evolves in time beeing incremented when an alert is solved (back to ok state) when it has previously been acknowleged.

    component: solved_alert

    resource: ack

    metric: count


- Delay between an alarm and it's resolution date for acknowleged alarms


This metric counts the total time between event ack and alert resolution. This is done on change state.

    component: solved_alert

    resource: ack

    metric: delay


- Solved alarm without acknowlegement

This metric evolves in time beeing incremented when an alert is solved (back to ok state) when it has previously NOT been acknowleged.

    component: __canopsis__

    resource: ack

    metric: cps_solved_ack_alarms
