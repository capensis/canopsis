#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

from canopsis.configuration import Parameter, Configuration
from canopsis.common.utils import resolve_element
from canopsis.storage.manager import Manager


class Benchmark(Manager):
    """
    Aims to run benchmark scenarios
    """

    CONF_FILE = '~/etc/bench.conf'

    CATEGORY = 'BENCHMARK'

    STORAGES = 'storages'
    SIZE = 'size'
    COUNT = 'count'
    TIME = 'time'
    SCENARIOS = 'scenarios'

    def __init__(
        self, storages=None, size=0, count=0, time=0, scenarios=None,
        *args, **kwargs
    ):

        super(Benchmark, self).__init__(*args, **kwargs)

        self.storages = storages
        self.size = size
        self.count = count
        self.time = time
        self.scenarios = scenarios

    def run(self, scenarios=None, *args, **kwargs):
        """
        Run input scenario with self in parameter

        :param scenarios: scenarios to run
        :type scenarios: list of {str, callable}

        :return: a dictionary of scenarios result by scenario entry
        :rtype: dict(scenario, result)
        """

        result = dict()

        if scenarios is None:
            scenarios = self.scenarios.split(',')

        # convert strings to callable objects if required
        _scenarios = [resolve_element(scenario) if isinstance(scenario, str)
        else scenario for scenario in scenarios]

        for index, scenario in enumerate(_scenarios):
            result[scenarios[index]] = scenario(self, *args, **kwargs)

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(Benchmark, self)._get_conf_files(*args, **kwargs)

        result.append(Benchmark.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Benchmark, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Benchmark.CATEGORY,
            new_content=(
                Parameter(Benchmark.STORAGES, self.storages),
                Parameter(Benchmark.COUNT, self.count, int),
                Parameter(Benchmark.SIZE, self.size, int),
                Parameter(Benchmark.TIME, self.time, int),
                Parameter(Benchmark.SCENARIOS, self.scenarios, str)))

        return result

    def _configure(self, conf, *args, **kwargs):

        super(Benchmark, self)._configure(conf=conf, *args, **kwargs)

        values = conf[Configuration.VALUES]

        # set shared
        self._update_parameter(values, Benchmark.STORAGES)
        self._update_parameter(values, Benchmark.COUNT)
        self._update_parameter(values, Benchmark.SIZE)
        self._update_parameter(values, Benchmark.TIME)
        self._update_parameter(values, Benchmark.SCENARIOS)
