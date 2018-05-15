#!/usr/bin/env python2.7
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
