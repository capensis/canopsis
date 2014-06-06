import logging

class CamqpMock(object):
	def __init__(self, logging_level=logging.INFO, logging_name="%s-amqp_mock", on_ready=None):

		self.exchange_name_events = 'camqpMock'
		self.logger = logging.getLogger(self.exchange_name_events)
		self.events = []

	def publish(self, event, rk, exchange_name):
		self.events.append(event)

	#TODO some other mock methods