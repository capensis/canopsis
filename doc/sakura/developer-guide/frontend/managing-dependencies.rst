.. _dev-frontend-deps:

Managing externals dependencies and JS libs
*******************************************

Synopsis
========

This guide describes the new workflow for managing frontend external dependencies.


Introduction
============

External frontend dependencies are the more often JS and css frameworks and tools. They can be used all across the frontend when they are loaded.

Before Sakura, libs where managed directly by adding and committing files into externals repo (under "webcore-libs" folder). The new workflow uses Bower, a package manager for frontend js applications and websites. It helps by providing a way to install packages the way they were designed to, to describes those packages, and to manage dependencies and versions of provided libraries.

It also helps by allowing to refetch automatically every listed frontend dependency needed by the Canopsis frontend.



Webcore-libs as a package
=========================

The webcore-libs folder is now registered as a private Bower package.

It contains a ``bower.json`` file that lists metadata and dependencies of the package.


How to install libraries
========================

Let's assume we want to add a packages that contains a XYZ package into Canopsis.

Please be sure that a ``.bowerrc`` file is present into the root of the ``canopsis-externals`` repo of the Canopsis sources.

Then go into the ``webcore-libs`` folder, and add the bower package :

.. code-block:: bash

   bower install XYZ --save


the ``--save`` option registers the library into the webcore-libs' ``bower.json``.

Now you can publish the changes to the ``bower.json`` file and add the library folder to the ``canopsis-externals`` repository.

.. NOTE::

   TODO: library wrappers
