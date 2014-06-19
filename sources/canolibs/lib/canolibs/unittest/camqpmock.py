import logging

class CamqpMock(object):

	exchange_name_alerts = 'mock_exchange_name_alerts'

	def __init__(self, logging_level=logging.INFO, logging_name="%s-amqp_mock", on_ready=None):

		self.exchange_name_events = 'camqpMock'
		self.logger = logging.getLogger(self.exchange_name_events)
		self.events = []

	def publish(self, event, rk, exchange_name):
		self.events.append(event)

	def clean(self):
		self.events = []

	#TODO some other mock methods