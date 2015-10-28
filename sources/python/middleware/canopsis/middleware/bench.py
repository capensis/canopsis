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

from canopsis.configuration.model import Parameter, Configuration
from canopsis.common.utils import lookup
from canopsis.middleware.registry import MiddlewareRegistry


class Benchmark(MiddlewareRegistry):
    """
    Aims to run benchmark scenarios
    """

    CONF_PATH = 'middleware/bench.conf'  #: configuration path

    CATEGORY = 'BENCH'  #: category name

    MIDDLEWARES = 'middlewares'  #: middleware to compare
    DATA_SIZE = 'data_size'  #: data size to bench
    DATA_COUNT = 'data_count'  #: data count to bench
    ITERATION = 'iteration'  #: iteration
    TIME = 'time'  #: maximal time between iterations
    SCENARIOS = 'scenarios'  #: bench scenarios

    def __init__(
        self, middlewares=None, data_size=0, data_count=0, iteration=1, time=0,
        scenarios=None,
        *args, **kwargs
    ):

        super(Benchmark, self).__init__(*args, **kwargs)

        self.middlewares = middlewares
        self.data_size = data_size
        self.data_count = data_count
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
        _scenarios = [lookup(scenario) if isinstance(scenario, basestring)
        else scenario for scenario in scenarios]

        for index, scenario in enumerate(_scenarios):
            result[scenarios[index]] = scenario(self, *args, **kwargs)

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(Benchmark, self)._get_conf_files(*args, **kwargs)

        result.append(Benchmark.CONF_PATH)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Benchmark, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Benchmark.CATEGORY,
            new_content=(
                Parameter(Benchmark.MIDDLEWARES, self.middlewares),
                Parameter(Benchmark.DATA_COUNT, self.data_count, int),
                Parameter(Benchmark.DATA_SIZE, self.data_size, int),
                Parameter(Benchmark.TIME, self.time, int),
                Parameter(Benchmark.ITERATION, self.iteration, int),
                Parameter(Benchmark.SCENARIOS, self.scenarios)))

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(Benchmark, self)._configure(conf=unified_conf, *args, **kwargs)

        values = unified_conf[Configuration.VALUES]

        # set shared
        self._update_parameter(values, Benchmark.MIDDLEWARES)
        self._update_parameter(values, Benchmark.DATA_COUNT)
        self._update_parameter(values, Benchmark.DATA_SIZE)
        self._update_parameter(values, Benchmark.TIME)
        self._update_parameter(values, Benchmark.SCENARIOS)
