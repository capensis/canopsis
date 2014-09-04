Time graph
==========

Time graph widget allows to plot metrics representation over time. Y axis represents the metric value whereas the X axis in this graph displays time dimension. It may be used to easily keep track of metric values over time. Information below describes how to use this widget with it's many avaiable options.


.. csv-table:: Timegraph options description
   :header: "Option", "Description"
   :widths: 15, 50

   "General", "Options available for timegraph general purpose"
   "title", "The timegraph title that is displayed on widget board"
   "auto title if available", "Sets a title automaticaly when possible based on metric informations"
   "show border", "Displays a border on the widget when selected"
   "human readable values", "metrics values are shown in a user-friendly way: '10 000' value will be displayed as '10K'"
   "refresh interval", "Defines how often the timegraph will check for new metric information to server. if any, graph is redrawn"
   "time window", "Defines how many information are fetched for server in the past until now. This parameter is made of an integer quantity and a time unit such as minute/hour/day ..."

   "Options", "options for timegraph at first level that are specific to this widget"
   "series type [#f1]_", "Defines how is displayed the graph. It can either be one of the following values: **bars** displays one vertical stick for each metric value. The **line** representation will show lines linking start point to end point in time and will pass by each point between start and stop poiut. **area** is similar to line and will enable area painting between 0 value on Y axis and the current point. Both bars, area and lines color are changeable in advanced mode."
   "line width", "Sets the thinkness of the lines in the graph"
   "marker [#f2]_", "defines the geometrical form used in the graph to replace point. form can be one of `circle, square, diamond, triangle, triangle down`. This property can be overried for any metric individually"
   "marker radius", "Defines the geometrical form's size in pixels used in the graph (markers)."
   "Stacked [#f3]_", "Display all metrics stacked. this means that the first line start from X axis origin (0) to the metric value and the second metric in the list will be represented from the first metric value to the first metric value plus the second metric value. and so on for each metric in the graph "
   "Time window offset", ""
   "History navigation [#f4]_", "Displays a sub graph for a selected period that makes visisble a large time slice for current graph"
   "History time window", "Period to select when displaying history navigation."
   "Display downtimes", "Allows downtimes display on the graph if current graph metric selection are undergoing downtimes on the displayed period."
   "Enable tooltip", "Displays metric values and information on mouseover on a graph metric point."
   "Shared", "Display all metric information for hovered point in graph instead of only point's metric information"


.. [#f1] timegraph_series_typess types
.. image:: ../../../_static/images/widgets/timegraph_series_types.png

.. [#f2] markers
.. image:: ../../../_static/images/widgets/timegraph_markers.png

.. [#f3] stacked option
.. image:: ../../../_static/images/widgets/timegraph_stacked.png

.. [#f4] history navigation
.. image:: ../../../_static/images/widgets/timegraph_history_navigation.png
