Description
-----------

Annotation library like Java's annotation with reflective concerns for Python.

.. image:: https://img.shields.io/pypi/l/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: License

.. image:: https://img.shields.io/pypi/status/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: Development Status

.. image:: https://img.shields.io/pypi/v/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: Latest release

.. image:: https://img.shields.io/pypi/pyversions/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: Supported Python versions

.. image:: https://img.shields.io/pypi/implementation/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: Supported Python implementations

.. image:: https://img.shields.io/pypi/wheel/b3j0f.annotation.svg
   :target: https://travis-ci.org/b3j0f/annotation
   :alt: Download format

.. image:: https://travis-ci.org/b3j0f/annotation.svg?branch=master
   :target: https://travis-ci.org/b3j0f/annotation
   :alt: Build status

.. image:: https://coveralls.io/repos/b3j0f/annotation/badge.png
   :target: https://coveralls.io/r/b3j0f/annotation
   :alt: Code test coverage

.. image:: https://img.shields.io/pypi/dm/b3j0f.annotation.svg
   :target: https://pypi.python.org/pypi/b3j0f.annotation/
   :alt: Downloads

.. image:: https://readthedocs.org/projects/b3j0fannotation/badge/?version=master
   :target: https://readthedocs.org/projects/b3j0fannotation/?badge=master
   :alt: Documentation Status

.. image:: https://landscape.io/github/b3j0f/annotation/master/landscape.svg?style=flat
   :target: https://landscape.io/github/b3j0f/annotation/master
   :alt: Code Health

Links
-----

- `Homepage`_
- `PyPI`_
- `Documentation`_

Installation
------------

pip install b3j0f.annotation

Features
--------

What does mean annotations in a reflective way:

- one annotation can annotate several objects at a time (modules, classes, functions, instances, builtins, annotation like themselves, etc.).
- such as a reflective object, they could have their own behavior and lifecycle independently to annotated elements.

This library provides the base Annotation class in order to specialize your own annotations, and several examples of useful annotation given in different modules such as:

- async: dedicated to asynchronous programming.
- interception: annotations able to intercept callable object calls.
- call: inherits from interception module and provides annotations which allow to do checking on callable objects.
- check: annotations which check some conditions such as type of annotated targets, max number of annotated elements, etc.
- oop: useful in object oriented programming like allowing to weave mixins.

Examples
--------

Perspectives
------------

- Cython implementation.

Donation
--------

.. image:: https://liberapay.com/assets/widgets/donate.svg
   :target: https://liberapay.com/b3j0f/donate
   :alt: I'm grateful for gifts, but don't have a specific funding goal.

.. _Homepage: https://github.com/b3j0f/annotation
.. _Documentation: http://b3j0fannotation.readthedocs.org/en/master/
.. _PyPI: https://pypi.python.org/pypi/b3j0f.annotation/


