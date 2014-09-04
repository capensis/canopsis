Text cell
=========

"Display a message with a html template."
This widget primary purpose is to display an arbitrary string of text.
It can also display html code and event variables.

.. image:: ../../../_static/images/widgets/text_widget1.png
.. image:: ../../../_static/images/widgets/text_widget2.png

Displaying Variables
--------------------

The text cell widget allows to display event variables.

The list of variables can be seen on [Event specification](https://github.com/capensis/canopsis/wiki/Event-specification) page, Event structure paragraph.
Possible uses of the text cell widget are: add labels on our custom Views, display some html frame (like twitter, icinga etc..) or print variables from the inventory tab as shown on the grabs below. You can use for example {output} to display output of a resource, or directly use {perfdata:%metric_name%:value} ( {perfdata:cps_sel_total:value} ) to display a specific perfdata.

Available Variables
^^^^^^^^^^^^^^^^^^^

Example
^^^^^^^

.. code-block:: text

	{perfdata:cps_sel_total:value}

Mathematical Operations
-----------------------

The widget also allows to parse mathematical operations.

To do so, juste create an html element with a class attribute "mathexpression". When rendering the widget, a math parser will try to evaluate the expression.

This feature can be used to display the addition of two metrics values, to multiply a metric value by a factor, and so on.

Example
^^^^^^^

.. code-block:: html

	<div class="mathexpression">
		{perfdata:metric_1:value} + {perfdata:metric_2:value}
	</div>
