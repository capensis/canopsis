#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from datetime import datetime
import logging
import time
from unittest import TestCase

from canopsis.alerts import AlarmField
from canopsis.alerts.filter import AlarmFilter
from canopsis.alerts.manager import Alerts
from canopsis.context_graph.manager import ContextGraph
from canopsis.check import Check
from canopsis.common.utils import merge_two_dicts
from canopsis.middleware.core import Middleware


class BaseTest(TestCase):

    def setUp(self):
        self.logger = logging.getLogger('alerts')

        self.alarm_storage = Middleware.get_middleware_by_uri(
            'storage-periodical-testalarm://'
        )
        self.config_storage = Middleware.get_middleware_by_uri(
            'storage-default-testconfig://'
        )
        self.context_graph_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.filter_storage = Middleware.get_middleware_by_uri(
            'storage-default-testalarmfilter://'
        )

        self.manager = Alerts()
        self.manager[Alerts.ALARM_STORAGE] = self.alarm_storage
        self.manager[Alerts.CONFIG_STORAGE] = self.config_storage
        self.manager.context_manager[
            ContextGraph.ENTITIES_STORAGE] = self.context_graph_storage
        self.manager[Alerts.FILTER_STORAGE] = self.filter_storage

        self.config_storage.put_element(
            element={
                '_id': 'test_config',
                'crecord_type': 'statusmanagement',
                'bagot_time': 3600,
                'bagot_freq': 10,
                'stealthy_time': 300,
                'stealthy_show': 600,
                'restore_event': True,
                'auto_snooze': False,
                'snooze_default_time': 300,
            },
            _id='test_config'
        )

    def tearDown(self):
        """Teardown"""
        self.alarm_storage.remove_elements()
        self.filter_storage.remove_elements()
        self.context_graph_storage.remove_elements()

    def gen_fake_alarm(self, update={}, moment=None):
        """
        Generate a fake alarm/value.
        """
        if moment is None:
            moment = int(time.mktime(datetime.now().timetuple()))

        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'fake-connector',
                'connector_name': 'fake-connector-name',
                'component': 'c',
                'output': 'out',
                'timestamp': moment
            }
        )

        value = alarm[self.manager[Alerts.ALARM_STORAGE].VALUE]
        value[AlarmField.state.value] = {
            't': moment,
            'val': Check.MINOR
        }
        value[AlarmField.steps.value] = [
            {
                '_t': 'stateinc',
                't': moment,
                'a': 'fake-author',
                'm': 'fake-message',
                'val': Check.MINOR
            }
        ]

        dictio = merge_two_dicts(alarm, update)

        return dictio, value

    def gen_alarm_filter(self, update={}, storage=None):
        """
        Generate a standard alarm filter.
        """
        base = {
            AlarmFilter.LIMIT: 180.0,
            AlarmFilter.FILTER: '',
            AlarmFilter.CONDITION: {},
            AlarmFilter.TASKS: ['alerts.systemaction.state_increase'],
        }

        dictio = merge_two_dicts(base, update)

        return AlarmFilter(dictio, logger=self.logger, storage=storage)
