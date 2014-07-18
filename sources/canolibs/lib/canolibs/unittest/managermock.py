import logging

class ManagerMock(object):
	def __init__(self, logging_level=logging.INFO):

		self.exchange_name_events = 'managerMock'
		self.logger = logging.getLogger(self.exchange_name_events)
		self.data = []

	def push(self, name=None, value=None, meta_data=None):
		self.data.append({'name': name, 'value': value, 'meta_data': 'meta_data'})

	def clean(self):
		self.data = []

	#TODO some other mock methods