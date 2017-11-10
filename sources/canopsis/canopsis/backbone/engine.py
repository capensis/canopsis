from __future__ import unicode_literals

from canopsis.common.pika_amqp import Consumer


class Engine(Consumer):
    def __init__(self, connection, queue, logger):
        """
        :param AmqpConnection amqp_consumer:
        """
        super(Engine, self).__init__(connection, queue)
        self.logger = logger
        self.callback = self.work

    def beat(self):
        raise NotImplementedError('Not implemented yet')

    def work(self):
        raise NotImplementedError('Not implemented yet')
