Series and Curves
=================

This document explains how to use the *Performance Data* view.

Introduction
------------

On this view, you can see 3 lists:

 * **Metrics**: lists the known metrics ;
 * **Series**: lists your series and allows you to create/edit/remove them ;
 * **Curves**: lists your curves and allows you to create/edit/remove them.

A **metric** is a named value that can change over time. At the moment, Canopsis knows
4 types of metrics:

 * **GAUGE**: the actual value of the metric is the last incoming value ;
 * **COUNTER**: increment the previous metric's value with this one ;
 * **ABSOLUTE**: always positive value ;
 * **DERIVE**: the current value corresponds to the metric's value derived over time.

A **serie** is an object containing :

 * a name,
 * multiple metrics,
 * a formula, involving the selected metrics,
 * an aggregation, to fetch the metrics aligned on X axis,
 * some metadata (warn, crit, min, max, unit).

Just like a metric, a serie will generate a list of points (where each point is
a pair ``(timestamp,value)``).
Each point is the result of the formula applied to the fetched metrics.

A **curve** describe how to plot a serie or a metric:

 * do we want lines, areas, points, and/or bars ?
 * do we want to plot the values with the curve ?
 * how does look lines, areas, points, and/or bars ?
