#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.common.amqp import AmqpConsumer


class Engine(AmqpConsumer):

    def __init__(self, connection, queue, logger):
        """
        :param AmqpConnection connection: a connection object to the Amqp
        :param str queue: the queue to consume
        :param Logger logger: a logger
        """
        super(Engine, self).__init__(connection, queue)
        self.logger = logger
        self._on_message = self.work

    def beat(self):
        raise NotImplementedError('Not implemented yet')

    def work(self):
        raise NotImplementedError('Not implemented yet')
