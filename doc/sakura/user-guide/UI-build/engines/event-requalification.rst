.. include:: ../../../includes/links.rst

Event Requalification
=====================

Purpose
-------

Some supervisors use differents status code for their check results. If
we want to connect one of those supervisor to **Canopsis**, we need to
convert its status codes to the standard used in *Canopsis* :

-  0 for OK
-  1 for WARNING
-  2 for CRITICAL
-  3 for UNKNOWN

*Canopsis* use a **state map** to resolve the standardized status code
from any original status code. A *state map* contains a list of codes
associated to a *Canopsis* status, for example :

+------+-------------+
| CODE |   STATUT    |
+======+=============+
| 0    |     OK      |
+------+-------------+
| 1    |   WARNING   |
+------+-------------+
| 2    |   WARNING   |
+------+-------------+
| 3    |   CRITICAL  |
+------+-------------+
| 4    |   CRITICAL  |
+------+-------------+
| 5    |   CRITICAL  |
+------+-------------+
| 6    |   UNKNOWN   |
+------+-------------+


When you create a derogation with a state map, every new event, matching
the derogation's filter, will have a new field ``real_state`` set to the
original state of the event, and the field ``state`` will be updated
accordingly to the *state map* defined (if the code is not in the *state
map*, then UNKNOWN is used).

Howto
-----

Create a state map
~~~~~~~~~~~~~~~~~~

Go to the Statemap Manager available in the menu "Run" ...

|image1|

... in otder to open the state map manager view ...

|image2|

... and click on the add button, you'll see the following form :

|image3|

Click on the Add button to add a new state association :

-  in the first column, set the state code to associate ;
-  in the second column, set the corresponding *Canopsis* state ;
-  the third column is used to remove the association.

.. NOTE:: If you set two associations like this :

  -  0 associated with OK
  -  2 associated with CRITICAL

  Then, the state code 1 will be associated with UNKNOWN.

Integration with derogation
~~~~~~~~~~~~~~~~~~~~~~~~~~~

See |derogation| .

In the tab "Requalificate", you will see the list of all your *state
maps*, just select the one you want.

.. |image1| image:: ../../../_static/images/requalification/menu_statemap_manager.png
.. |image2| image:: ../../../_static/images/requalification/statemap_manager.png
.. |image3| image:: ../../../_static/images/requalification/add_statemap.png
