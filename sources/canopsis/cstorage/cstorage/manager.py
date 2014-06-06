#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from cconfiguration import Configurable
from utils import resolve_element


class Manager(Configurable):

	CONFIGURATION_FILE = '~/etc/storage.conf'

	TIMED_STORAGE = 'timed_storage'
	PERIODIC_STORAGE = 'periodic_storage'
	ONE_VALUE_STORAGE = 'one_value_storage'

	SHARED = 'shared'

	STORAGE = 'STORAGE'

	PARSING_RULES = {
		STORAGE: {
			TIMED_STORAGE: str,
			PERIODIC_STORAGE: str,
			ONE_VALUE_STORAGE: str,
			SHARED: bool
		}
	}

	_STORAGE_BY_DATA_TYPE_BY_TYPE = dict()

	def __init__(self, configuration_file=CONFIGURATION_FILE, shared=False,
		*args, **kwargs):

		super(Manager, self).__init__(
			configuration_file=configuration_file, *args, **kwargs)

		self.shared = shared

	def get_parsing_rules(self, *args, **kwargs):
		"""
		Get the right structure to use to parse a configuration file.

		:rtype: list of dict(section, dict(option, parser(option)))
		"""

		result = super(Manager, self).get_parsing_rules(*args, **kwargs)

		parsing_rule = Manager.PARSING_RULES.copy()

		for timed_type in self._get_timed_types():
			parsing_rule[Manager.STORAGE][timed_type] = str

		for periodic_type in self._get_periodic_types():
			parsing_rule[Manager.STORAGE][periodic_type] = str

		for one_value_type in self._get_one_value_types():
			parsing_rule[Manager.STORAGE][one_value_type] = str

		result.append(parsing_rule)

		return result

	def get_storage(self, data_type, storage_type, self_storage_type=None,
		shared=None, *args, **kwargs):

		if shared is None:
			shared = self.shared

		if storage_type is None:
			storage_type = self_storage_type

		elif isinstance(storage_type, str):
			storage_type = resolve_element(storage_type)

		elif callable(storage_type):
			pass

		if shared and \
			storage_type in Manager._STORAGE_BY_DATA_TYPE_BY_TYPE \
			and data_type in Manager._STORAGE_BY_DATA_TYPE_BY_TYPE[storage_type]:
				result = Manager._STORAGE_BY_DATA_TYPE_BY_TYPE[storage_type][data_type]

		else:
			result = storage_type(data_type=data_type, *args, **kwargs)

		if shared:
			Manager._STORAGE_BY_DATA_TYPE_BY_TYPE[storage_type][data_type] = result

		return result

	def get_timed_storage(self, data_type, timed_type=None, shared=None,
		*args, **kwargs):

		result = self.get_storage(data_type=data_type, shared=shared,
			storage_type=timed_type, self_storage_type=self.timed_type,
			*args, **kwargs)

		return result

	def get_periodic_storage(self, data_type, periodic_type=None, shared=None,
		*args, **kwargs):

		result = self.get_storage(data_type=data_type, shared=shared,
			storage_type=periodic_type, self_storage_type=self.periodic_type,
			*args, **kwargs)

		return result

	def get_one_value_storage(self, data_type, one_value_type=None, shared=None,
		*args, **kwargs):

		result = self.get_storage(data_type=data_type, shared=shared,
			storage_type=one_value_type, self_storage_type=self.one_value_type,
			*args, **kwargs)

		return result

	def _set_parameters(self, parameters, error_parameters, *args, **kwargs):

		# set default timed type
		timed_type = parameters.get(Manager.TIMED_STORAGE)
		if timed_type is not None:
			self.timed_type = resolve_element(timed_type)

		# set default periodic type
		periodic_type = parameters.get(Manager.PERIODIC_STORAGE)
		if periodic_type is not None:
			self.periodic_type = resolve_element(periodic_type)

		# set default one value type
		one_value_type = parameters.get(Manager.ONE_VALUE_STORAGE)
		if one_value_type is not None:
			self.one_value_type = resolve_element(one_value_type)

		shared = parameters.get(Manager.SHARED)
		if shared is not None:
			self.shared = shared

		storage_types = {
			Manager.TIMED_STORAGE: self.timed_type,
			Manager.PERIODIC_STORAGE: self.periodic_type,
			Manager.ONE_VALUE_STORAGE: self.one_value_type
		}

		# set attributes with input parameters where name are data_type
		for name, parameter in parameters.iteritems():

			storage = None

			# if parameter is default storage type
			if parameter in storage_types:
				storage = self.get_storage(
					data_type=name, storage_type=storage_types[parameter])

			else:  # get dynamic storage type
				storage = self.get_storage(data_type=name, storage_type=parameter)

			setattr(self, name, storage)

	def _get_timed_types(self, *args, **kwargs):

		return []

	def _get_periodic_types(self, *args, **kwargs):

		return []

	def _get_one_value_types(self, *args, **kwargs):

		return []
