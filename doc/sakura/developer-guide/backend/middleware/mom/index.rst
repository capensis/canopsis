=====================================================
MOM: Library for managing Message Oriented Middleware
=====================================================

.. module:: canopsis.mom
    :synopsis: middleware library for using MOM.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Objective
=========

This library aims to provide to canopsis developers to focus even more on data to exchange instead of thinking about the way to do the exchange in a MOM paradigm.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

The API is generic in respecting MOM paradigm common properties, and agnostic from technologies.

MOM paradigms are of two types, point to point and pubsub.

As same as any canopsis middleware, its configuration can be done with fine grains or done in a common way to MOM concerns thanks to an uri.

In a data oriented approach, it is useful to take care about the data behavior which is changing related to MOM quality of services (QoS). For example http://www.omg.org/spec/dds4ccm/1.1/PDF/ describes a data in a pubsub middleware which can become an event or a status information depending on QoS values.

Therefore, a MOM configuration depends on a ``data type``.

URI
---

A Canopsis URI respects the convention in taking care only about the scheme and in using other parameters such as user name, password, host, etc. to configure MOM.

A canopsis scheme is as follow ``protocol[-data_type]`` where the ``data_type`` is optional but the protocol is required and takes values of embedded technologies in middleware.

The URI parameters are used in order to define senders and receivers.

Example:
########

The following URI connects to the protocol ``pmom``, at host ``demo.canopsis.org`` with the user ``dude``. The related mom uses two receivers ``alice`` and ``bob`` (where bob is associated to the callback function ``canopsis.pmom.consume``) and the sender named ``chris``:

``pmom://dude@demo.canopsis.org/?receivers=alice,bob:canopsis.pmom.consume&senders=chris``

Perspectives
------------

This library is further declined into paradigms such as :

.. toctree::
   :maxdepth: 2
   :titlesonly:

   pt2pt
   pubsub

And implementations :

.. toctree::
   :maxdepth: 2
   :titlesonly:

   kombu/index

Package contents
================

.. data:: __version__

    Current package version : 0.1

.. data:: SR_SEPARATOR = ':'

    Char separation between receiver/sender name and callback/mode.

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

   .. data:: DIRECT_MODE = 'direct'

      Default sending mode.

   .. data:: TOPIC_MODE = 'topic'

      Topic sending mode.

   .. data:: FANOUT_MODE = 'fanout'

      Fanout sending mode.

   .. data:: senders

      instance attribute which is a dictionary of senders by name

   .. data:: receivers

      instance attribute which is a dictionary of receivers by name

   .. method:: bind_to(receiver, destination)

      Bind input receiver to the input destination.

      :param receiver: receiver name
      :type receiver: str

     :param destination: destination to bind input receiver
      :type destination: str

   .. method:: _get_sender(name)

      Method to implement in order to get sender object related to input name.

   .. method:: _get_receiver(name, mode=DIRECT_MODE)

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

   .. method:: _send(sender, msg, rk, serializer, compression, content_type, content_encoding, out_timeout)

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
