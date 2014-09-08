.. _derogation:

Derogation
==========

Purpose
-------

A derogation is used to perform some actions on incoming events :

-  overriding a field's value ;
-  requalificate the event's state.

Howto to create a derogation
----------------------------

Access to the derogation manager view
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Go to the Derogation Manager, available in the menu "Run" ...

|image1|

... to access to the derogation manager view.

|image2|

Then, click on the add button, you'll see the following form to fill.

Derogation definition form
~~~~~~~~~~~~~~~~~~~~~~~~~~

1. General
^^^^^^^^^^

Type the name of your derogation, and a short description. Then you will
have to define the time interval where your derogation is active. Once
this interval over, the derogation will not be applied anymore.

|image3|

2. Override
^^^^^^^^^^^

This tab allows you to specify default values for some fields of the
derogated event. Click the add button to override a new field :

-  choose the field to override with the first column ;
-  specify the value to set in the second column ;
-  the third column is just a button to delete the overriding.

|image4|

3. Requalificate
^^^^^^^^^^^^^^^^

This tab allows you to select the *Statemap* to use for the event
requalification. Just select the *Statemap* you want, or deselect with
``Ctrl + Left Click`` to disable the event requalification.

|image5|

4. Filter
^^^^^^^^^

This tab allows you to configure the event filter, only events which
match the filter will be derogated. |image6|

.. |image1| image:: ../../../_static/images/derogation/select_derogation_manager.png
	     :height: 100 px
	     :width: 1000 px

.. |image2| image:: ../../../_static/images/derogation/derogation_manager.png
.. |image3| image:: ../../../_static/images/derogation/add_derog_general.png
.. |image4| image:: ../../../_static/images/derogation/add_derog_override.png
.. |image5| image:: ../../../_static/images/derogation/add_derog_requal.png
.. |image6| image:: ../../../_static/images/derogation/add_derog_filter.png
