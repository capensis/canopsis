Serie
=====

This document specifies how series features must be implemented.

References
----------

> -   FR::Serie &lt;FR\_\_Serie&gt;
> -   FR::Metric &lt;FR\_\_Metric&gt;
> -   TR::Metric &lt;TR\_\_Metric&gt;
> -   FR::Configurable &lt;FR\_\_Configurable&gt;
> -   FR::Context &lt;FR\_\_Context&gt;
> -   FR::Schema &lt;FR\_\_Schema&gt;
> -   FR::Engine &lt;FR\_\_Engine&gt;

Updates
-------

Contents
--------

### Frontend Serie

A *SerieAdapter* **MUST** provide a method to:

> -   fetch perfdata &lt;FR\_\_Metric\_\_PerfData&gt; according to the
>     selected metrics &lt;FR\_\_Metric&gt;
> -   aggregate those perfdata using the
>     serie aggregation &lt;FR\_\_Serie\_\_Aggregation&gt;
> -   consolidate the aggregated points using the
>     serie formula &lt;FR\_\_Serie\_\_Formula&gt;

This *SerieAdapter* **SHOULD** be used by a widget which wants to
visualize the series.

In this case, the formula will be a *sand-boxed* JavaScript expression,
where the operators &lt;FR\_\_Serie\_\_Operators&gt; are methods from
the *SerieAdapter* which will returns an array of points according to
the regular expression &lt;FR\_\_Serie\_\_Selection&gt;.

### Backend Serie

A *Serie* configurable registry &lt;FR\_\_Configurable\_\_Registry&gt;
**MUST** provide:

> -   a method to fetch metrics according to a regular expression and,
>     eventually, a set of already selected metrics (to avoid an
>     operation to the Context &lt;FR\_\_Context&gt;)
> -   a method to fetch perfdata using the
>     PerfData configurable &lt;TR\_\_Metric\_\_PerfData&gt;, returning
>     an array of timeseries &lt;TR\_\_Metric\_\_TimeSerie&gt;
> -   a method to consolidate timeseries, using the serie formula

In this case, the formula will be a *sand-boxed* Python expression,
where the operators are functions with the following prototype:

> -   regular expression to select the correct timeseries (using the
>     *Serie* configurable registry) as a parameter
> -   returning a single value (according to the operator)

### Serie engine

A serie engine &lt;FR\_\_Engine&gt; will provide:

> -   a method which periodically check which serie to calculate, and
>     send them on the engine's queue
> -   a method to consume the queue, and calculate the serie

The data schema &lt;FR\_\_Schema\_\_Data&gt; of the
serie configuration &lt;FR\_\_Serie\_\_Configuration&gt; will provide:

> -   `last_computation` as an `integer`, which will be used with the
>     `aggregation_interval` to know if a serie must be calculated or
>     not
