========================================
Kombu: MOM implementation over RabbitMQ
========================================

.. module:: canopsis.kombu
    :synopsis: MOM implementation over RabbitMQ.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Objective
=========

This library aims to provide an implementation over RabbitMQ of the MOM paradigm.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

This library is an implementation of the canopsis.mom.MOM class.

The related ``protocol`` is 'kombu'.

The related ``data_types`` are both ``pubsub`` and ``pt2pt``.

URI
---

The URI uses the uri parameters in order to define senders and receivers

example:

- ``kombu://localhost:5672``: default connection with default parameters
- ``kombu-pubsub://canopsis:canopsis@demo.canopsis.org:5672`` connection to canopsis demo

A Canopsis URI respects the convention in taking care only about the scheme and in using other parameters such as user name, password, host, etc. to configure MOM.

A canopsis scheme is as follow ``protocol[-data_type]`` where the data_type is optional but the protocol is required and takes values of embedded technologies in middleware.

Perspectives
------------

This library is further declined into paradigms such as :
- PTP: point to point
- PubSub: publish subscribe

Package contents
================

.. data:: __version__

    Current package version : 0.1

.. data:: RECEIVER_CALLBACK_SEPARATOR = ':'

    Char separation between receiver name and callback.

.. function:: receiver_and_callback(receiver)

   Return a tuple of receiver name and callback.

   :return: the tuple (receiver name, callback)
   :rtype: tuple

.. class:: MOM(canopsis.middleware.Middleware)

   Multi MOM paradigm class.

   .. data:: CATEGORY = 'MOM'

      Configuration category name

   .. data:: CONF_RESOURCE = 'mom/mom.conf'

      MOM conf resource (in addition to ones from the Middleware class).

   .. data:: SENDERS = 'senders'

      Configuration sender ids. List of sender resources able to send data.

   .. data:: RECEIVERS =  'receivers'

      Configuration receiver ids. Lise of receiver resources able to receive data.

   .. data:: senders

      Dictionary of senders by name

   .. data:: receivers

      Dictionary of receivers by name

   .. method:: _connect()

      Protected method to implement in order to connect the MOM.

      :return: True iif self is connected
      :rtype: bool

   .. method:: _get_sender(name)

      Method to implement in order to get sender object related to input name.

   .. method:: _get_receiver(name)

      Method to implement in order to get sender object related to input
            receiver name and callback.

      :param receiver: receiver name
      :type receiver: str

      :param callback: callback to attach to the receiver if not None.
      :type callback: callable.

      :return: specific receiver object.

   .. method:: _receive(receiver, callback, in_timeout)

      Method to override in order to implement message reception in receive\_ method.

      :param receiver: receiver name to use
      :type name: str

      :param callback: callback to attach to receiver if not None
      :type callback: callable

      :param in_timeout: receive message timeout

      :return: a message if synchronous mode is asked, else None

   .. method:: receive(receiver, callback=None, in_timeout=None)

      Receive a data (a)synchronously depending on input parameter values.

        :param receiver: receiver name to use.
        :type receiver: str.

      :param callback: callable which take an message in parameter and which is attached to input receiver in order to receive messages asynchronously if not None.
      :type callback: callable

      :param in_timeout: timout in milliseconds to use in synchronous mode.
         - If positive or 0, use synchronous mode.
         - If None, use self.in_timeout.
      :type in_timeout: int

      :return: message if synchronous is enabled, None in other cases.

   .. method:: _send(sender, msg, rk, serializer, compression, content_type, content_encoding, in_timeout)

      Method to override in order to implement message sending in send method\_.

      :param sender: sender to use
      :type sender: object initialized by this middleware
      :param msg: message to send
      :param rk: routing_key
      :param serializer: message serializer into binary digits
      :param compression: message compression mode
      :param content_type: message content type
      :param content_encoding: message content encoding
      :param out_timeout: send timeout to use. If None, use self.out_timeout.

   .. method: send(msg, rk, sender=None, serializer="json", compression=None, content_type=None, content_encoding=None, out_timeout=None)

      Send a message with input sender resource.

      :param sender: sender name or sender to use. If None, use all senders.
      :type sender: str
      :param msg: message to send
      :param rk: routing_key
      :param serializer: message serializer into binary digits
      :param compression: message compression mode
      :param content_type: message content type
      :param content_encoding: message content encoding
      :param out_timeout: send timeout to use. If None, use self.out_timeout.

   .. method:: cancel_senders(senders=None)

      Cancel senders.

      :param senders: sender names to cancel. If senders is None, cancel all senders.

   .. method:: cancel_receivers(receivers=None)

      Cancel receivers.

      :param receivers: receiver names to cancel. If receivers is None, cancel all receivers.
