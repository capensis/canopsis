# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

__version__ = "0.1"


from canopsis.middleware.core import Middleware
from canopsis.configuration.parameters import Category, Parameter
from canopsis.common.utils import lookup

from urlparse import urlparse
from functools import reduce

SR_SEPARATOR = ':'


def sender_and_mode(sender):
    """
    Get a tuple of sender name and its mode.
    """

    result = sender, MOM.DIRECT_MODE

    if SR_SEPARATOR in sender:
        splitted_name = sender.split(SR_SEPARATOR)

        result = splitted_name[0], splitted_name[1]

    return result


def receiver_and_callback(receiver):
    """
    Return a tuple of receiver name and callback.

    :return: the tuple (receiver name, callback)
    :rtype: tuple
    """

    result = receiver, None

    if SR_SEPARATOR in receiver:
        splitted_name = receiver.split(SR_SEPARATOR)

        callback = lookup(splitted_name[1])

        result = splitted_name[0], callback

    return result


class MOM(Middleware):
    """
    Dedicated to MOM middleware
    """

    CONF_RESOURCE = 'mom/mom.conf'

    CATEGORY = 'MOM'

    SENDERS = 'senders'
    RECEIVERS = 'receivers'

    DIRECT_MODE = 'direct'  # default mode
    TOPIC_MODE = 'topic'
    FANOUT_MODE = 'fanout'

    def __init__(self, senders=None, receivers=None, *args, **kwargs):
        """
        :param senders: senders to use in order to send data.
        :param receivers: receivers to use in order to receive data.
        """

        super(MOM, self).__init__(*args, **kwargs)

        self.senders = senders

        self.receivers = receivers

    def _get_uri(self):

        result = super(MOM, self)._get_uri()

        # add receivers and senders to the uri
        if result:
            parsed_url = urlparse(result)
            params = parsed_url.params

            if MOM.SENDERS not in params:
                params += '%s=%s' % (
                    MOM.SENDERS,
                    reduce(lambda x, y: x + y, self.senders.keys()))

            if MOM.RECEIVERS not in params:
                params += '%s=%s' % (
                    MOM.RECEIVERS,
                    reduce(lambda x, y: x + y, self.receivers.keys()))

            result = '%s://%s%s?%s' % (
                parsed_url.scheme, parsed_url.netloc, parsed_url.path, params)

        return result

    def _set_uri(self, value):

        super(MOM, self)._set_uri(value)

        parsed_url = urlparse(value)

        # get senders and receivers from params
        if parsed_url.params:
            splitted_params = parsed_url.params.split('&')

            for param in splitted_params:
                # get receivers
                if param.startswith('%s=' % MOM.RECEIVERS):
                    receivers = param[len("%s=" % MOM.RECEIVERS):]
                    self.receivers = receivers
                # get senders
                elif param.startswith('%s=' % MOM.SENDERS):
                    senders = param[len("%s=" % MOM.SENDERS):]
                    self.senders = senders

    @property
    def senders(self):
        return self._senders

    @senders.setter
    def senders(self, value):
        # update a local dictionary of senders by name
        self._senders = {}
        if value is not None:
            for name in value:
                sender, mode = sender_and_mode(name)
                sender = self._get_sender(sender=sender, mode=mode)
                if sender is not None:
                    self._senders[name] = sender

    @property
    def receivers(self):
        return self._receivers

    @receivers.setter
    def receivers(self, value):
        self._receivers = {}
        if value is not None:
            for name in value:
                receiver, callback = receiver_and_callback(name)
                receiver = self._get_receiver(
                    receiver=receiver, callback=callback)
                self._receivers[name] = receiver

    def bind_to(self, receiver, destination):
        """
        Bind input receiver to the input destination.

        :param receiver: receiver name
        :type receiver: str

        :param destination: destination to bind input receiver
        :type destination: str
        """
        raise NotImplementedError()

    def _get_conf_files(self, *args, **kwargs):

        result = super(MOM, self)._get_conf_files(*args, **kwargs)

        result.append(MOM.CONF_RESOURCE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(MOM, self)._conf(*args, **kwargs)

        result += Category(
            Parameter(MOM.SENDERS),
            Parameter(MOM.RECEIVERS))

        return result

    def _init_env(self, conn, *args, **kwargs):

        super(MOM, self)._init_env(conn, *args, **kwargs)

        # update senders and receivers
        self.senders = self.senders.keys()
        self.receivers = self.receivers.keys()

    def _get_sender(self, name):
        """
        Method to implement in order to get sender object related to input
            name.
        """

        raise NotImplementedError()

    def _get_receiver(self, receiver, callback=None):
        """
        Method to implement in order to get sender object related to input
            receiver name and callback.

        :param receiver: receiver name
        :type receiver: str

        :param callback: callback to attach to the receiver if not None.
        :type callback: callable.

        :return: specific receiver object
        """

        raise NotImplementedError()

    def _receive(self, receiver, callback, in_timeout):
        """
        Method to implement in order to get a specific receiver object to this
            mom.

        :param receiver: receiver name to use
        :type name: str

        :param callback: callback to attach to receiver if not None
        :type callback: callable

        :param in_timeout: receive message timeout

        :return: a message if synchronous mode is asked, else None
        """

        raise NotImplementedError()

    def receive(self, receiver, callback=None, in_timeout=None):
        """
        Receive a data (a)synchronously depending on input parameter values.

        :param receiver: receiver name to use.
        :type receiver: str.

        :param callback: callable which take an message in parameter and which
            is attached to input receiver in order to receive messages
            asynchronously if not None.
        :type callback: callable

        :param in_timeout: timout in milliseconds to use in synchronous mode.
            - If positive or 0, use synchronous mode.
            - If None, use self.in_timeout.
        :type in_timeout: int

        :return: message if synchronous is enabled, None in other cases.
        """

        result = None

        if self.connected():

            if in_timeout is None:
                in_timeout = self.in_timeout

            # do asynchronous mode
            if callback is not None:
                self._get_receiver(receiver=receiver, callback=callback)

            # do synchrnous mode
            if in_timeout >= 0:
                try:
                    result = self._receive(
                        receiver=receiver, callback=callback,
                        in_timeout=in_timeout)
                except Exception as e:
                    self.logger.error(
                        'Error while receiving message %s' % e)

        else:
            self.logger.error(
                "Impossible to receive, %s is not connected" % self)

        return result

    def _send(
        self,
        sender, msg, rk, serializer, compression,
        content_type, content_encoding, out_timeout
    ):
        """
        Send a message with input sender resource.

        :param sender: sender to use
        :type sender: object initialized by this middleware
        :param msg: message to send
        :param rk: routing_key
        :param serializer: message serializer into binary digits
        :param compression: message compression mode
        :param content_type: message content type
        :param content_encoding: message content encoding
        :param out_timeout: send timeout to use. If None, use self.out_timeout.
        """

        raise NotImplementedError()

    def send(
        self, msg, rk, sender=None, serializer="json",
        compression=None, content_type=None, content_encoding=None,
        out_timeout=None
    ):
        """
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
        """

        if out_timeout is None:
            out_timeout = self.out_timeout

        if self.connected():

            senders = None

            if isinstance(sender, basestring):
                sender = self._senders_by_name.get(sender)

                if sender is None:
                    self.logger.error("Sender %s is unknown" % sender)

                else:
                    senders = (sender,)

            elif sender is not None:
                senders = (sender,)

            else:
                senders = self._senders_by_name.values()

            for sender in senders:
                self.logger.debug(
                    "Send message to %s in %s" % (rk, sender))
                try:
                    self._send(
                        sender=sender, rk=rk, serializer=serializer,
                        compression=compression, content_type=content_type,
                        content_encoding=content_encoding,
                        out_timeout=out_timeout)

                except Exception as err:
                    self.logger.error("Impossible to send message (%s)" % err)

                else:
                    self.logger.debug(
                        "Message sent message to %s in %s" % (rk, sender))

        else:
            self.logger.error("Impossible to send, %s is not connected" % self)

    def cancel_senders(self, senders=None):
        """
        Cancel senders.

        :param senders: sender names to cancel. If senders is None, cancel all
            senders.
        """

        if senders is None:
            senders = self.senders.keys()

        for sender in senders:
            try:
                del self.senders[senders]
            except KeyError:
                self.logger.debug("No sender %s to delete" % sender)

    def cancel_receivers(self, receivers=None):
        """
        Cancel receivers.

        :param receivers: receiver names to cancel. If receivers is None,
            cancel all receivers.
        """

        if receivers is None:
            receivers = self.receivers.keys()

        for receiver in receivers:
            try:
                del self.receivers[receiver]
            except KeyError:
                self.logger.debug("No receiver %s to delete" % receiver)

    def _disconnect(self, *args, **kwargs):

        self.cancel_receivers()
        self.cancel_senders()
        del self._conn
        self._conn = None
