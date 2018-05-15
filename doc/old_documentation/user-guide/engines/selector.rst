.. _user-engines-selector:

Selector
========

This document describes how does the selectors work.

Introduction
------------

A selector is composed of:

 * an event filter (and/or specific events included/excluded)
 * one or more computations

When the engine executes the selector, it will fetch events from the database
according to the configured filter. Then, the computations will be applied to
the set of events, in order to produce a new ``selector`` event.

The available computations are:

 * worst-state : the new event will have the worst event state among the set
 * SLA : a new ``sla`` event will be produced with SLA as metrics

Event Emitted
-------------

The ``selector`` event is very similar to a ``check`` event:

 * it owns a state, computed from the selector's alerts
 * it owns some internal metrics:
    * ``cps_sel_state_off``: number of alerts off
    * ``cps_sel_state_minor``: number of minor alerts
    * ``cps_sel_state_major``: number of major alerts
    * ``cps_sel_state_critical``: number of critical alerts
    * ``cps_sel_ack``: number of acknowledged alerts
    * ``cps_sel_total``: total number of alerts in the selector
 * it is associated to a component only (with the selector's name)
 * its output is the rendered selector's template

Example:

.. code-block:: javascript

   {
     "connector": "canopsis",
     "connector_name": "engine",
     "event_type": "selector",
     "source_type": "component",
     "component": "My Awesome Selector",
     "state": 0,
     "output": "My Selector's template",
     "display_name": "My Awesome Selector",
     "perf_data_array": [(...)]
   }

The ``sla`` event is also similar to a ``check`` event:

 * it owns a state, computed from the SLA metrics
 * it owns some metrics:
    * ``cps_pct_by_0``: percentage of the SLA period in "Off" state
    * ``cps_pct_by_1``: percentage of the SLA period in "Minor" state
    * ``cps_pct_by_2``: percentage of the SLA period in "Major" state
    * ``cps_pct_by_3``: percentage of the SLA period in "Critical" state
    * ``cps_avail``: selector's availability
    * ``cps_avail_duration``: duration when the selector is available
    * ``cps_alerts_duration``: duration when the selector is in alerts
 * it is associated to a resource of the selector's component
 * its output is the rendered SLA's template, configured in the selector

Example:

.. code-block:: javascript

   {
     "connector": "sla",
     "connector_name": "engine",
     "event_type": "sla",
     "source_type": "resource",
     "component": "My Awesome Selector",
     "resource": "sla",
     "state": 0,
     "output": "My SLA's template",
     "display_name": "My Awesome Selector",
     "perf_data_array": [(...)]
   }

See :ref:`dev-spec-sla` for more informations.
