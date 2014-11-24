#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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

port = 162
interface = '0.0.0.0'

mibs = {
	'1.3.6.1.4.1.8072.4': 'NET-SNMP-AGENT-MIB',
	'1.3.6.1.4.1.674.10892.1': 'MIB-Dell-10892',
	'1.3.6.1.4.1.674.10893.1': 'StorageManagement-MIB',
	'1.3.6.1.4.1.674.10893.1.20.200': 'StorageManagement-MIB'
}

blacklist_enterprise = [
	'1.3.6.1.4.1.8072.3.2.10',
	'1.3.6.1.4.1.311.1.1.3.1.2'
]

blacklist_trap_oid = [
	'1.3.6.1.4.1.674.10892.1.0.1306',
	'1.3.6.1.4.1.674.10892.1.0.1304'
]
