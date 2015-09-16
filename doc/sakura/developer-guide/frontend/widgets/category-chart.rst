Category chart
==============

Overview
--------

The category chart is a Canopsis widget made of a widget core and a display component. The widget core manage underlaying widget data and the `c3 js <http://c3js.org>`_ display component uses data to render the widget depending on user configuration.


widget core
-----------

The category chart core is made of a controller and a template. The template only instanciate the render component.

It also fetch metrics depending on user configuration. The user can select series or metrics to display from the configuration form. The core also prepare widget configuration that will affect the component representation. Data fetch is performed asynchronously and when all identified tasks are complete (configuration, metric and series fetch), then only the widget rendering is performed for performances issues.

The data fetching is mainly based on the metric controller that provider high level methods for metric querying. These method are called with a callback that is one of the widget method doing the job of preparing data for the display component.

display component
-----------------

At the moment, the category chart uses **c3 js** as chart library to compute an interactive chart display.

When data from widget core are ready, then the rendering starts, and prepared data are integrated to the chart instanciation options.

The chart lifecycle depends on the widget lifecycle as c3js chart instance is a property of the canopsis component. It can be destroyed and then instanciated on widget refresh. It can also be recomputed when the user dynamicaly changes the chart display type.

User options make vary the chart generation by only changing the option dictionnary that is used on c3js chart instanciation. Available options are functionnaly described in the `user guide <../../../user-guide/UI/widgets/categorychart.html>`_ section. Most of them are well explained through the c3 js documentation as in the chart option computed values are explicitely labelled with c3 js conventions.

One important custom feature is the series dynamic naming depending on context informations. On metric fetch, the context Id associated to the metric is parsed and then used as a Handlebars template context that the administrator user is able to manipulate thanks to the templating system. Context information availables are **connector, connector_name, component, resource and metric**.

