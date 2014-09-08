======================================
RPC: Library for remote procedure call
======================================

.. module:: canopsis.rpc
    :synopsis: middleware library for using RPC.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Objective
=========

This library aims to provide to canopsis developers to focus even more on data to exchange instead of thinking about the way to do the exchange in a RPC paradigm.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

The API is generic in respecting RPC paradigm common properties, and agnostic from technologies.

As same as any canopsis middleware, its configuration can be done with fine grains or done in a common way to MOM concerns thanks to an uri.

URI
---

A Canopsis URI respects the convention in taking care only about the scheme and in using other parameters such as user name, password, host, etc. to configure RPC.

A canopsis scheme is as follow ``protocol[-data_type]`` where the ``data_type`` is optional but the protocol is required and takes values of embedded technologies in middleware.

The URI parameters are used in order to define senders and receivers.

Perspectives
------------

This library is further declined into protocols such as :

- WS: Web services
- RMI
- XMLRPC: XML-RPC
- SOAP
- REST

Package contents
================

.. data:: __version__

    Current package version : 0.1

.. class:: RPC(canopsis.middleware.Middleware)

   RPC paradigm class.

   .. data:: CATEGORY = 'RPC'

      Configuration category name

   .. data:: CONF_RESOURCE = 'rpc/rpc.conf'

      RPC conf resource (in addition to ones from the Middleware class).
