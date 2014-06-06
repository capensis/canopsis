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

from cmongo import MongoDB


class Store(MongoDB):
	"""
	Manage periodic and timed access to data_type in a mongo collection.
	"""

	DEFAULT_CONFIGURATION_FILE = '~/etc/pyperfstore3.conf'

	def __init__(self, data_type, *args, **kwargs):

		super(Store, self).__init__(
			backend=self.get_collection_name(data_type=data_type),
			_ready_to_conf=False, *args, **kwargs)

		self._data_type = data_type

		if self.auto_conf:
			self.apply_configuration()

	def get_parsers_by_option_by_section(self, *args, **kwargs):

		result = dict()

		mongodb_parsers = super(MongoDB, self).get_parsers_by_option_by_section()

		for section in list(mongodb_parsers.keys()):
			new_section_name = "{0}_{1}_{2}".format(
				section, type(self), self._data_type)
			result[new_section_name] = mongodb_parsers[section]

		return result

	def _get_collection_prefix(self):
		"""
		Protected method to override in order to get collection prefix.
		"""

		raise NotImplementedError()

	def get_collection_name(self, data_type=None):
		"""
		Get collection name managed by this store related to input data_type.
		"""

		if data_type is None:
			data_type = self._data_type

		result = "{0}_{1}".format(
			self._get_collection_prefix(), data_type)

		return result

	def _get_data_id(self, data_id):
		"""
		Get data id if data_id is an entity _id
		"""

		result = data_id

		if data_id.startswith(self._data_type):
			result = data_id[len(self._data_type) + 1:]

		return result
