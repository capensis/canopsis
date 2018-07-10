#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

import re

from canopsis.common.link_builder.link_builder import HypertextLinkBuilder
from canopsis.common.mongo_store import MongoStore
from canopsis.common.utils import merge_two_dicts, get_sub_key
from canopsis.confng import Configuration, Ini
from canopsis.logger import Logger

ALERTS_COLLECTION = 'periodical_alarm'
LOG_PATH = 'var/log/engines/context-graph.log'
PARAM_REG = re.compile('\{([a-zA-Z0-9_\.\*]+)\}')
SEPARATOR = '*'


class BasicAlarmLinkBuilder(HypertextLinkBuilder):

    """
    Basic builder which read a base_url parameter, and enrich it with entity
    values.
    """

    def __init__(self, options={}):
        super(BasicAlarmLinkBuilder, self).__init__(options=options)
        self.logger = Logger.get('context-graph', LOG_PATH)

        conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)
        mongo = MongoStore(config=conf_store)
        self.alerts_collection = mongo.get_collection(name=ALERTS_COLLECTION)

    def build(self, entity, options={}):
        opt = merge_two_dicts(self.options, options)
        alarm = self.alerts_collection.find_one({'d': entity['_id']})

        if 'base_url' in opt:
            url = opt['base_url']
            hay = {}
            for m in re.finditer(PARAM_REG, opt['base_url']):
                needles = m.group(0).strip('{').strip('}').split('.')
                needle = SEPARATOR.join(needles)
                value = ''
                if needles[0] == 'alarm':
                    needle = SEPARATOR.join(needles[1:])
                    value = get_sub_key(alarm, '.'.join(needles[1:]))
                else:
                    value = get_sub_key(entity, '.'.join(needles))

                url = url[:m.start()] + '{' + needle + '}' + url[m.end():]
                hay[needle] = value

            self.logger.debug(url.format(**hay))
            return {self.category: [url.format(**hay)]}

        return {}
