.. _snmp:

SNMP
====

It is possible in canopsis to translate SNMP traps to canopsis events by editing snmp translation rules. This can be done by going to the SNMP rule panel in the engine menu. This will take care of possible mib definition that are inserted into Canopsis thanks to the upload button.

Uploading mibs
--------------

..TODO
..`link <../../administrator-guide/>`_

Rule edition
------------

First, create a rule by clicking on the create button on the top right of the screen. Then the snmp rule editor will appear.

|image1|


From this form, you have to type the first letter that a mib module should contain, then a list of mib module will be shown. Then select a module in the list.

|image2|

Selecting a module in the list will perform a search for the related module's mibs availables mibs. Then on the combobox on the right is the mib names.

|image3|

Select a mib name will perform a second search that makes available mib objects for selected module. This search contains object names that can be used in templates. object names templates information can be added by just clicking on them. Module objects in template is given to the snmp engine and will be replaced by the value in the trap event ``snmp_vars`` property that have the same oid than the module's object one. Clicking on a module object name for a template field will append the right template information to the template.

|image4|

In the end, you have to choose the status that the snmp engine's produced event will wear, and then the rule is over and ready to be hanlded by the snmp engine.


.. |image1| image:: ../../_static/images/snmp/ruleedit_1.png
.. |image2| image:: ../../_static/images/snmp/ruleedit_2.png
.. |image3| image:: ../../_static/images/snmp/ruleedit_3.png
.. |image4| image:: ../../_static/images/snmp/ruleedit_4.png

