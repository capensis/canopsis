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
from __future__ import unicode_literals

from canopsis.task.core import register_task

from canopsis.perfdata.manager import PerfData
from canopsis.context.manager import Context

from copy import deepcopy


perfdatamgr = PerfData()


@register_task
def event_processing(engine, event, manager=None, logger=None, **kwargs):
    """Perfdata engine synchronous processing.
    """

    if manager is None:
        manager = perfdatamgr

    # Get perfdata
    perf_data = event.get('perf_data')
    perf_data_array = event.get('perf_data_array', [])

    if perf_data_array is None:
        perf_data_array = []

    # Parse perfdata
    if perf_data:
        logger.debug(' + perf_data: {0}'.format(perf_data))

        try:
            perf_data_array += manager.parse_perfdata(perf_data)

        except Exception as err:
            logger.error(
                "Impossible to parse perfdata from: {0} ({1})".format(
                    event, err
                )
            )

    logger.debug(' + perf_data_array: {0}'.format(perf_data_array))

    # Add status informations
    event_type = event.get('event_type')

    handled_event_types = ['check', 'selector', 'sla']

    if event_type is not None and event_type in handled_event_types:

        logger.debug('Add status informations')

        state = int(event.get('state', 0))
        state_type = int(event.get('state_type', 0))
        state_extra = 0

        # Multiplex state
        cps_state = state * 100 + state_type * 10 + state_extra

        perf_data_array.append(
            {
                "metric": "cps_state",
                "value": cps_state
            }
        )

    event['perf_data_array'] = perf_data_array

    # remove perf_data_keys where values are None
    for index, perf_data in enumerate(event['perf_data_array']):

        perf_data_array_with_Nones = event['perf_data_array'][index]

        event['perf_data_array'][index] = {
            name: perf_data_array_with_Nones[name]
            for name in perf_data_array_with_Nones
            if perf_data_array_with_Nones[name] is not None
        }

    logger.debug('perf_data_array: {0}'.format(perf_data_array))

    event = deepcopy(event)

    # Metrology
    timestamp = event.get('timestamp', None)

    if timestamp is not None:

        perf_data_array = event.get('perf_data_array', [])

        for perf_data in perf_data_array:

            perf_data = perf_data.copy()
            event_with_metric = deepcopy(event)
            event_with_metric['type'] = 'metric'
            event_with_metric[Context.NAME] = perf_data.pop('metric')

            encoded_event_with_metric = {}
            for k, v in event_with_metric.items():
                try:
                    k = k.encode('utf-8')
                except:
                    pass
                try:
                    v = v.encode('utf-8')
                except:
                    pass
                encoded_event_with_metric[k] = v

            metric_id = manager.context.get_entity_id(
                encoded_event_with_metric
            )
            value = perf_data.pop('value', None)

            manager.put(
                metric_id=metric_id, points=[(timestamp, value)],
                meta=perf_data, cache=True
            )

    return event
