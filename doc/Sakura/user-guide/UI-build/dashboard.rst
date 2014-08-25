.. include:: /Sakura/includes/links.rst

.. _Sakura_dashboard:

Dashboard
=========

Canopsis WebUI is available through your navigator.
Fire up your prefered browser, enter the IP of the server Canopsis
has been installed on, port 8082 and you should be able to access it.
The default view, called Dashboard, is used to display the various
components and resources reports. Although the preferred method to
customise content shown by Canopsis is to create additional
|Sakura_Views|, it's also possible to modify
Dashboard. Below I explain how to add a new
|Sakura_widget| on it. ( :ref:`Figure 1 <figure1>` shows
default dashboard, while :ref:`Figure 2 <figure1>` in contrast shows a customised one).


.. |link| replace:: `figure 1 </Sakura/images/dashboard/default_dashboard.png>`__
.. |link2| replace:: `figure 2 </Sakura/images/dashboard/dashboard.png>`__

.. _figure1:

+-----------------------+-------------+
| |default_dashboard|   | |dashboard| |
+=======================+=============+
| |link|                | |link2|     |
+-----------------------+-------------+


To start customizing Canopsis dashboard, click on the Build menu (see
:ref:`Figure 3 <figure3>` ) then choose "Edit active view". Once in Edit mode you will be
able to add and remove or customize any existing
|Sakura_widgets| To add a new
|Sakura_widget|, with the mouse click and
select an empty area over the dashboard then release the button. (see
:ref:`Figure 4 <figure3>` )


.. |link3| replace:: `figure 3 </Sakura/images/dashboard/build_menu.png>`__
.. |link4| replace:: `figure 4 </Sakura/images/dashboard/add_widget.png>`__

.. _figure3:

+----------------+----------------+
| |build_menu|   | |add_widget|   |
+================+================+
| |link3|        | |link4|        |
+----------------+----------------+



A new |Sakura_widget| dialog should popup, enter the necessary information click save on the bottom right corner of the dialog (see :ref:`Figure 5 and 6 <figure5>` ).



.. |link5| replace:: `figure 5 </Sakura/images/dashboard/add_widget1.png>`__
.. |link6| replace:: `figure 6 </Sakura/images/dashboard/add_widget2.png>`__

.. _figure5:

+----------------+----------------+
| |add_widget1|  | |add_widget2|  |
+================+================+
| |link5|        | |link6|        |
+----------------+----------------+



Here we've added a new |Sakura_widget| of type |Sakura_bar_graph| that's
setup to display cpu charge history of system processes. Let's save the
current view ( :ref:`Figure 7 <figure7>` ) and see the result ( :ref:`Figure 8 <figure7>` ). Please refer to
dedicated wiki pages for detailed info over the different types of
|Sakura_widgets| and existing
|Sakura_metrics| . 




.. |link7| replace:: `figure 7 </Sakura/images/dashboard/add_widget3.png>`__
.. |link8| replace:: `figure 8 </Sakura/images/dashboard/add_widget4.png>`__

.. _figure7:

+----------------+-------------------+
| |add_widget3|  | |add_widget4|     |
+================+===================+
| |link7|        | |link8|           |
+----------------+-------------------+

.. |default_dashboard| image:: /Sakura/images/dashboard/default_dashboard.png  
                :height: 65 px
                :width: 325 px
.. |dashboard| image:: /Sakura/images/dashboard/dashboard.png  
                :height: 65 px
                :width: 325 px
.. |build_menu| image:: /Sakura/images/dashboard/build_menu.png
                :height: 65 px
                :width: 325 px
.. |add_widget| image:: /Sakura/images/dashboard/add_widget.png
                :height: 65 px
                :width: 325 px
.. |add_widget1| image:: /Sakura/images/dashboard/add_widget1.png
                :height: 65 px
                :width: 325 px
.. |add_widget2| image:: /Sakura/images/dashboard/add_widget2.png
                :height: 65 px
                :width: 325 px
.. |add_widget3| image:: /Sakura/images/dashboard/add_widget3.png
                :height: 65 px
                :width: 325 px
.. |add_widget4| image:: /Sakura/images/dashboard/add_widget4.png
                :height: 65 px
                :width: 325 px
