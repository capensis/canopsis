Canopsis Python Project
=======================

Every Python component of Canopsis is distributed via a Python package, using
the namespace ``canopsis``.

For example, if I have to develop something specific to nagios, I will create the
following Python package : ``canopsis.nagios``.

A project exists within the ``sources/python`` folder and **must** contain a
``setup.py``. We provide some useful tools for this, in the module ``canopsis.common.setup``
(first package to be installed).

