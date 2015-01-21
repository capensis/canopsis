Text cell
=========

This widget's purpose is to display arbitrary information written by a Canopsis administrator.
Availables features in the widget text are:

 - display arbitrary html content
 - display computed metric values
 - display events information


Arbitrary html content display
------------------------------

When editing the widget text, it is possible to write custom HTML in the text editor. Using the editor buttons would make easier the html edition and it is also possible to write raw html directly.

.. image:: ../../../_static/images/widgets/text_html_edition.png



Metric display system
---------------------

The text widget edition allow selecting performance data series. For more information on how create series from metrics see `series <../../UIv2/serie.html>`_ .
Metrics displayed are the last value for the serie metric computation from the selected date interval witch is by default between **now** and **now - 300 seconds**. If no metric available in this interval, the template value will display `No metric available` as value in the render.

.. image:: ../../../_static/images/widgets/select_series.png

Once serie is selected, it is possible to display it's value in the text area by using the handlebars template system as long as the metric value is embed in the template rendering context under the `serie.serie_name` name.

The code bellow will allow the widget text to display the value of the serie1

.. code-block:: html

	<p>Metric value for serie1 is : {{serie.serie1}}</p>


Event display system
--------------------

.. TODO event value display system
