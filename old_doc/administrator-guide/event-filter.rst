.. include:: ../link.rst

.. _event-filter:

Event filter
============

The engine filter permits to users to choose which events have to be
|processed| by
|canopsis|.

Filter description
------------------

The engine applies filters which are composed of a priority, an action
and a rule.

-  The priority designates the filter application priority order.
-  The action is DROP or PASS designates respectively if an event which
   matches with the rule will be deleted or
   |processed| by
   |canopsis|.
-  The rule parses event field values.

Filter definition
-----------------

Menu
~~~~

Definition of such filters is possible through the view filter
accessible from the filter rules menu.

|image1|

View
~~~~

The view filter permits to add/edit/duplicate/delete and define a
default action for events which match with none filters.

|image2|

Wizard
~~~~~~

In filter editing mode (by pushing the add button or in double clicking
on a filter) a wizard with two tabs permits to define a filter
properties.

Options
_______

The tab "options" defines name, priority and action of a filter.

|image3|

Rules
_____

The tab "rule" defines how an event is chosen by the filter in edition.

.. |image1| image:: ../_static/filter/filter_menu.png
.. |image2| image:: ../_static/filter//filter_view.png
.. |image3| image:: ../_static/filter/filter_options.png
.. |image4| image:: ../_static/filter/filter_rule.png
