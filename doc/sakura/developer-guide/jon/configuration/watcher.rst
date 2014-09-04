===========================================================
Watcher : Module aiming to apply reconfiguration at runtime
===========================================================

.. contents:
    maxdepth: 2

.. module:: canopsis.configuration.watcher
    :synopsis: configuration watcher library

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

Permits to (un)bind a Configurable object to configuration files modification/creation. Such bound Configurable have their apply_configuration method which is called when watched configuration files are modified.

Technical description
=====================

.. function:: add_configurable(configurable)

    Bind input configurable to all its configuration files

.. function:: remove_configurable(configurable)

    Unbind input configurable to all its configuration files

.. function:: on_update_conf_file(conf_file)

    call the method Configurable.apply_configuration on all bound Configurable to input conf_file

.. class:: Watcher(Configurable)

    In charge of watch configuration files modification

.. data:: sleeping_time

    sleeping time between two reading

.. method:: start(self)

    start self watcher

.. method:: stop(self)

    stop self watcher

.. data:: DEFAULT_SLEEPING_TIME = 5

    default sleeping time

.. data:: _WATCHER

    private singleton Watcher object

.. function:: start_watch()

    start global WATCHER

.. function:: stop_watch()

    stop global WATCHER

.. function:: change_sleeping_time(sleeping_time)

    change global WATCHER sleeping_time