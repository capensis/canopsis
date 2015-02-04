#!/usr/bin/env python
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

from logging import INFO, getLogger


class ManagerMock(object):
	def __init__(self, logging_level=INFO):

		self.exchange_name_events = 'managerMock'
		self.logger = getLogger(self.exchange_name_events)
		self.data = []

	def push(self, name=None, value=None, meta_data=None):
		self.data.append({'name': name, 'value': value, 'meta_data': 'meta_data'})

	def clean(self):
		self.data = []
