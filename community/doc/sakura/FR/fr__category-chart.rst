.. _FR__Category_chart:

==============
Category chart
==============

The Canopsis category chart is a widget that allow display metrics depending on categories

.. contents::
   :depth: 2

----------
References
----------

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/09/01", "0.3", "Update references", ""
   "David Delassus", "2015/09/01", "0.2", "Rename document", ""
   "Eric Régnier", "2015/08/06", "0.1", "Document creation", ""

--------
Contents
--------

Description
===========

The category chart is a widget that takes works with the actual frontend widget system. It allows to display the last value of a metric in a given period in many representation

The widget features are:

- category representation as gauge, bar, pie, donut, progressbar.
- metric sources are series and raw metrics
- stackable display for bar charts representation
- optional tooltip, labels and legend display on each representation
- customisable metric names

The canopsis system administrator must be able to change the widget display and each of the widget features.

Representation
==============

The chart widget comes with many possible representation for categorised data. they are : gauge, bar, progressbar, donut and pie

Each section of each chart show a specific metric value. These values are show as tooltips, labels on the chart and legend depending on the widget type.

Analyse
=======

The category chart is a widget such as the ones existing in Canopsis frontend. It may work with existing mixins such as periodic refresh that affects it's behavior. This can be done by follwing the Canopsis widget creation conventions.

Data sources for this widget are existing metrics in the canopsis backend system. These data can be reached thanks to existing api services. Shown values are the last available value in a period of ten minutes before the end date that is by default the current date. This widget will also have to work with the live reporting system that will dynamically change the end date of the data fetch.
