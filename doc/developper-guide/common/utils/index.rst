Common.Utils: common utils for all canopsis python libraries
============================================================

.. module:: canopsis.common.utils
    :synopsis: common utils for all canopsis python libraries

.. moduleauthor:: Labejof Jonathan <jlabejof@capensis.fr>
.. sectionauthor:: Labejof Jonathan <jlabejof@capensis.fr>

Common.Utils provides functions and classes useful in the entire canopsis project

Indices and tables
------------------

* :ref:`genindex`
* :ref:`search`

Package contents
----------------

.. function:: resolve_element(path)

    Get a reference to an element known by runtime by the input given path.
    The opposite function is :ref:`get_path`

.. function:: get_path(element)

    Get an absolute runtime path of the input element.
    Opposite to function :ref:`resolve_element`
