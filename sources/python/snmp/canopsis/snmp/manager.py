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

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry

#: snmp manager configuration category
CATEGORY = 'SNMP'

#: snmp manager configuration path
CONF_PATH = 'snmp/snmp.conf'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class SnmpManager(MiddlewareRegistry):

    # The configuration key to read in the configuration file
    SNMP_STORAGE = 'snmp_storage'

    def __init__(self, *args, **kwargs):
        super(SnmpManager, self).__init__(*args, **kwargs)

    def put(self, oid, rule):
        self[SnmpManager.SNMP_STORAGE].put_element(
            _id=oid, element=rule)

    def get(self, oids=None):
        return self[SnmpManager.SNMP_STORAGE].get_elements(ids=oids)

    def remove(self, oids=None):
        self[SnmpManager.SNMP_STORAGE].remove_elements(ids=oids)
