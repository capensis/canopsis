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
LOG_PATH = 'var/log/linkbuilder.log'
PARAM_REG = re.compile('\{([a-zA-Z0-9_\.]+)\}')


class BasicAlarmLinkBuilder(HypertextLinkBuilder):

    """
    Basic builder which read a base_url parameter, and enrich it with entity
    values.
    """

    def __init__(self, options={}):
        super(BasicAlarmLinkBuilder, self).__init__(options=options)
        self.logger = Logger.get('linkbuilder', LOG_PATH)

        conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)
        mongo = MongoStore(config=conf_store)
        self.alerts_collection = mongo.get_collection(name=ALERTS_COLLECTION)

    def custom_format(self, phrase, entity, alarm):
        """
        Translate a custom format string, with entity and alarm informations.
        Alarm informations must be prefixed with "alarm."

        Ex: http://{infos.url.value}/{alarm.v.resource}

        :param str phrase: a template string
        :param dict entity: the entity
        :param dict alarm: the alarm
        :returns: a string with replaced parameters
        :rtype: str
        """
        hay = {}
        for m in re.finditer(PARAM_REG, phrase):
            needles = m.group(0).strip('{').strip('}').split('.')
            needle = '_'.join(needles)  # used as a parameter name in format !
            value = ''
            if alarm is not None and needles[0] == 'alarm':
                value = get_sub_key(alarm, '.'.join(needles[1:]))
            else:
                value = get_sub_key(entity, '.'.join(needles))

            if value is None:
                raise ValueError(
                    "Value {} is missing or None".format('.'.join(needles)))

            phrase = phrase[:m.start()] + '{' + needle + '}' + phrase[m.end():]
            hay[needle] = value

        return phrase.format(**hay)

    def build(self, entity, options={}):
        alarm = options.pop('alarm', None)
        opt = merge_two_dicts(self.options, options)
        if alarm is None:
            alarm = self.alerts_collection.find_one({'d': entity['_id']})
        links = {}

        if 'base_url' in opt:
            try:
                link = self.custom_format(opt['base_url'], entity, alarm)
                self.logger.debug(link)

                category = 'Liens'
                label = 'URL'

                if 'category' in opt:
                    category = opt['category']

                if 'label' in opt:
                    label = opt['label']

                links[category] = [{'label': label,
                                    'link': link}]

                return links

            except ValueError:
                return {}

            except Exception as err:
                self.logger.exception(
                    'Unhandled error {} : {}'.format(type(err), err))
                return {}

        return {}
