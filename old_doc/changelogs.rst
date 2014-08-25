Changelogs
==========

` <#version-201311>`__\ Version ``201311``
------------------------------------------

` <#new-features->`__\ New features :
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

-  Debian 7 (Wheezy) compatibility

-  Dependencies :

   -  Add ChromaJS
   -  Update libevent-devel on RHEL6 (1.4.13-1 -> 1.4.13-4)
   -  libevent-doc 1.4.13-4 and libevent-headers 1.4.13-4 for RHEL6
   -  Redis integration
   -  Python module psutil==1.0.1
   -  Python module python-ldap
   -  openldap-devel for Centos 6, RHEL5 and RHEL6

-  Installation improvements :

   -  Do not erase views permissions on update

-  Engines improvements :

   -  Better use of AMQP
   -  Better MongoDB integration
   -  Add cleaner engine
   -  Refactoring of consolidation
   -  Configure amqp2engines service via ~/etc/amqp2engines.conf
   -  pyperfstore2 and selector engines are now asynchronous

-  New engines :

   -  perfstore2\_rotate: performing rotation of perfdata from Redis to
      MongoDB
   -  event\_filter: Filter incoming events according to list of rules
      (allows the event, or drop it)

-  Aggregation improvements :

   -  SUM operator
   -  Better time interval calculation (option "Round time interval" in
      Wizard)
   -  Possibility to perform aggregation on perfdata displayed in Text
      Cell Widget

-  General UI improvements :

   -  Search bar
   -  No Read/Write access to user group on view creation
   -  Import/Export objects (selectors, topologies, views, ...) as JSON

-  Accounts improvements :

   -  Avatar support
   -  Change password support
   -  LDAP based authentication
   -  Possibility to duplicate an account

-  Widgets improvements :

   -  New wizard UI
   -  Customize metrics
   -  Time window offset
   -  Enable/Disable human readable values
   -  Categorized Bar Graph (based on Diagram widget)
   -  Multiple Y axis on Line Graph/Bar Graph
   -  Colorize Y axis with associated curve color
   -  Possibility to have area and line curves in Line Graph
   -  Remove deprecated Pie and Categorized Graph widget (use Diagram
      widget now)
   -  Custom Warning/Critical threshold in wizard
   -  Text Cell working correctly when Live Reporting is active
   -  Custom name for selectors/topologies in Weather Widget
   -  Show/Hide selector's or topology's name in Weather Widget

-  New widgets :

   -  Pogressbar
   -  Trends
   -  Mini chart

-  New views:

   -  Filter rules : configure behavior of event\_filter engine

-  Briefcase improvements :

   -  Upload button
   -  Filters
   -  File's extension check
   -  Video/Images/PDF/Unknown files handling

-  More unittests

More than 140 bugfixes.

` <#version-201307>`__\ Version ``201307``
------------------------------------------

-  Widget wizard refactoring with default and advanced mode
-  Add thumbs for select widget
-  Add perfdata cache (``Redis``) for store plain data and relieve
   ``MongoDB``
-  Compress plain perfdata on all the day
-  Rebuild ``consolidation`` engine
-  Tweak ``tag`` engine performance
-  Add widget ``Trends``
-  Add widget ``Progressbar``
-  Use ``MongoDB`` aggregator instead map/reduce for ``selector`` engine
-  Display ``ETA`` (Estimated Time of Arrival) in trend tooltip of
   ``line_graph``
-  Add avatar personalisation
-  Many bugfixes

` <#howto-update>`__\ Howto Update:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Before process, it's good thing to make a full backup or snapshot.

::

    # with "root"
    pip install --upgrade git+http://github.com/socketubs/ubik.git@0.1

    # with "canopsis"
    hypcontrol stop
    ubik update
    ubik install redis-conf
    ubik upgrade
    hypcontrol start
    pyperfstore2 update

` <#version-201303>`__\ Version ``201303``
------------------------------------------

-  All JavaScript core is minimized for better load time
-  Rebuild reporting process and time format in webservice exchanges
-  New engine ``Consolidation``: Aggregate vertically (mean, sum ...)
   many metrics
-  New engine ``Topology`` (beta): Add dependencies for state computing
-  New main bar disposition
-  Deal with webservices with your ``auth_key``
-  Create your own widget
   (`wiki <https://github.com/capensis/canopsis/wiki/Create-your-own-widgets>`__)
-  Add time navigation on widgets ``Line graph``
-  Add second Y-Axis on widgets ``Line graph`` when unit of metrics are
   different
-  Add flags on line graph to symbolize events
-  New wiki index page:
   `Index <https://github.com/capensis/canopsis/wiki>`__
-  Now, we use `Transifex <https://www.transifex.com/>`__ for
   collaborative translations, you can contribute
   `here <https://www.transifex.com/projects/p/canopsis/>`__
-  Add debuging tools for javascript ui (capserjs)
-  Many bugfixes
-  Updated packages: ``webcore``, ``pyperfstore2``,
   ``wkhtmltopdf-libs``, ``canolibs``, ``amqp2engines``, ``canotools``,
   ``celery-libs``, ``webcore-libs``, ``python-libs``, ``mongodb-conf``

Last edited by William Pain, August 05, 2013

