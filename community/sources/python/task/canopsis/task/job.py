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

from canopsis.common.utils import dictproperty
from canopsis.old.record import Record


class Job(Record):
    def __init__(self, *args, **kwargs):
        super(Record, self).__init__(type='job', *args, **kwargs)

    @property
    def task(self):
        return self.data['task']

    @task.setter
    def task(self, value):
        self.data['task'] = value

    @property
    def start(self):
        return self.data['start']

    @start.setter
    def start(self, value):
        self.data['start'] = value

    @property
    def last_execution(self):
        return self.data['last_execution']

    @last_execution.setter
    def last_execution(self, value):
        self.data['last_execution'] = value

    @property
    def rrule(self):
        return self.data['rrule']

    @rrule.setter
    def rrule(self, value):
        self.data['rrule'] = value

    def params_get(self, key):
        return self.data['params'][key]

    def params_set(self, key, value):
        self.data['params'][key] = value

    def params_del(self, key):
        del self.data['params'][key]

    params = dictproperty(params_get, params_set, params_del)
