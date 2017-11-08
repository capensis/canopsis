from __future__ import unicodeliterals


class Engine(object):
    def __init__(self, logger, amqp_connection):
        """
        :param AmqpConnection amqp_connection:
        """
        self.logger = logger
        self.amqp_connection = amqp_connection

    def beat(self):
        raise NotImplementedError('Not implemented yet')

    def work(self):
        raise NotImplementedError('Not implemented yet')
