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

from canopsis.engines import Engine
from canopsis.old.tools import \
    parse_perfdata, metric_is_internal, internal_metrics
from canopsis.old.storage import get_storage
from canopsis.perfdata.manager import Manager

from copy import deepcopy

from md5 import md5


class engine(Engine):

    etype = 'perfstore3'

    def __init__(self, *args, **kargs):

        super(engine, self).__init__(*args, **kargs)

        storage = get_storage(logging_level=self.logger.level)
        self.entities = storage.get_backend('entities')

        self.manager = Manager(logging_level=self.logger.level)

    def work(self, event, *args, **kargs):

        ## Get perfdata
        perf_data = event.get('perf_data', None)
        perf_data_array = event.get('perf_data_array', [])

        if perf_data_array is None:
            perf_data_array = []

        ### Parse perfdata
        if perf_data:

            self.logger.debug(' + perf_data: {0}'.format(perf_data))

            try:
                perf_data_array = parse_perfdata(perf_data)

            except Exception as err:
                self.logger.error(
                    "Impossible to parse: {0} ('{1}')".format(perf_data, err))

        self.logger.debug(' + perf_data_array: {0}'.format(perf_data_array))

        ### Add status informations
        event_type = event.get('event_type', None)

        if event_type is not None \
                and event_type in ['check', 'selector', 'sla']:

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

        # remove perf_data_akeys where values are None
        for index, perf_data in enumerate(event['perf_data_array']):

            event['perf_data_array'][index] = \
                dict(((key, value) for key, value in perf_data.iteritems()
                    if value is not None))

        self.logger.debug('perf_data_array: {0}'.format(perf_data_array))

        #self.internal_amqp.publish(event, INTERNAL_QUEUE)
        self.on_internal_event(event)

        return event

    def on_internal_event(self, event, msg=None):
        event = deepcopy(event)

        ## Metrology
        timestamp = event.get('timestamp', None)

        if timestamp is not None:

            component = event.get('component', None)

            if component is not None:

                resource = event.get('resource', None)

                perf_data_array = perfstore3_get_perfdata_to_process(event)

                for perf_data in perf_data_array:

                    metric = perf_data.get('metric', None)

                    if metric is not None:

                        if not metric_is_internal(metric):

                            nodeid = md5()

                            nodeid.update(component.encode('ascii', 'ignore'))

                            if resource:
                                nodeid.update(
                                    resource.encode('ascii', 'ignore'))

                            nodeid.update(metric.encode('ascii', 'ignore'))
                            nodeid = nodeid.hexdigest()

                            value = perf_data.pop('value', None)

                            self.manager.put(
                                data_id=nodeid,
                                points_or_point=(timestamp, value),
                                meta=perf_data)

                    else:
                        self.logger.warning(
                            'metric name does not exist: {0}'.format(event))

    def beat(self):
        # Counts metric not in internal metrics for webserver cache purposes
        self.logger.info('Computing cache value for perfdata3 metric count')

        metrics_cursor = self.entities.find(
            {
                'type': 'metric',
                'name': {'$nin': internal_metrics}
            })
        count = len(metrics_cursor)

        self.object.save(
            {
                '_id': 'perfdata2_count_no_internal',
                'crecord_name': 'perfdata2_count_no_internal',
                'count': count
            },
            w=1
        )
        self.logger.info(
            'Cache value for perfdata3 metric count computed > {0}'.format(
                count))
