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

from canopsis.engines import Engine
from canopsis.perfdata.manager import PerfData
from canopsis.context.manager import Context

from copy import deepcopy


class engine(Engine):

    etype = 'perfdata'

    def __init__(self, *args, **kargs):

        super(engine, self).__init__(*args, **kargs)

        self.perfdata = PerfData()

    def work(self, event, *args, **kargs):

        ## Get perfdata
        perf_data = event.get('perf_data')
        perf_data_array = event.get('perf_data_array', [])

        if perf_data_array is None:
            perf_data_array = []

        ### Parse perfdata
        if perf_data:

            self.logger.debug(' + perf_data: {0}'.format(perf_data))

            try:
                perf_data_array = self.perfdata.parse_perfdata(perf_data)

            except Exception as err:
                self.logger.error(
                    "Impossible to parse: {0} ('{1}')".format(perf_data, err))

        self.logger.debug(' + perf_data_array: {0}'.format(perf_data_array))

        ### Add status informations
        event_type = event.get('event_type')

        handled_event_types = ['check', 'selected', 'sla']

        if event_type is not None and event_type in handled_event_types:

            self.logger.debug('Add status informations')

            state = int(event.get('state', 0))
            state_type = int(event.get('state_type', 0))
            state_extra = 0

            # Multiplex state
            cps_state = state * 100 + state_type * 10 + state_extra

            perf_data_array.append(
                {
                    "metric": "cps_state",
                    "value": cps_state
                })

        event['perf_data_array'] = perf_data_array

        # remove perf_data_keys where values are None
        for index, perf_data in enumerate(event['perf_data_array']):

            perf_data_array_with_Nones = event['perf_data_array'][index]

            event['perf_data_array'][index] = {
                name: perf_data_array_with_Nones[name]
                for name in perf_data_array_with_Nones
                if perf_data_array_with_Nones[name] is not None
            }

        self.logger.debug('perf_data_array: {0}'.format(perf_data_array))

        #self.internal_amqp.publish(event, INTERNAL_QUEUE)
        self.on_internal_event(event)

        return event

    def on_internal_event(self, event, msg=None):
        event = deepcopy(event)

        ## Metrology
        timestamp = event.get('timestamp')

        if timestamp is not None:

            perf_data_array = event.get('perf_data_array', [])

            for perf_data in perf_data_array:

                event_with_metric = deepcopy(event)
                event_with_metric['type'] = 'metric'
                event_with_metric[Context.NAME] = perf_data['metric']

                metric_id = self.perfdata.context.get_entity_id(
                    event_with_metric)
                value = perf_data.pop('value', None)

                self.perfdata.put(
                    metric_id=metric_id, points=[(timestamp, value)],
                    meta=perf_data)
