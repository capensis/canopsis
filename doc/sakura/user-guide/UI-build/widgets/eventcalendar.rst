Calendar
========

Overview
--------

"Display a calendar with events (ics-like and regular Canopsis events)"
This widget allows to display events in different types of calendars :

* Monthly view
* Agenda week view
* Agenda day view
* Basic week view
* Basic day view

.. image:: ../../../_static/images/widgets/calendar_widget1.png
.. image:: ../../../_static/images/widgets/calendar_widget2.png

There is two ways to display events on a calendar :

* Stacked : events are shown on a per-day basis on the calendar, in one bubble. When the bubble is clicked, a popup is shown with an extended view of all these events
* By source: Calendar-types events can be shown not stacked, directly on the calendar. Displayed events have to be of calendar type (see calendar event specification)

Customization
-------------

The widget is highly customizable. Thus, it is possible to vary a lot how information is managed by a lot of means :

* By specifying a color for each source of events, and a color for stacked events
* By changing the size of events on the calendar
* By showing or hiding weekends
* By specifying which information about stacked events has to be shown on the popup
* By changing the way calendar header is displayed
* By making the calendar editable, or read-only.

Event reception
---------------

The calendar widget manages incoming events with the websocket technology. From an user point of view, it mostly means that the widget will not need to get refreshed periodically, and will show selected events as they come into Canopsis.

Sending and editing events
--------------------------

When the widget is set as editable, it is possible to send and edit events directly from it.

To add an event, click on the calendar where you want it to be (you also can select a range for the event). A popup will appear, asking you to fill information about the new event, with title and source arguments required.

To edit events, click on the event you want to edit, then modify event information at your convenience.

When events are added or edited, it is sent to canopsis over AMQP. It appears on the calendar only when it has been processed by Canopsis.

Recurrence rule
---------------

Calendar events may also be reccurent : it is possible to repeat any calendar event according to what's known as rules (see icalendar specification). Event rule can be specified in the new/edit event window.
