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

from canopsis.engines import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.event import forger, get_routingkey
from canopsis.perfdata.manager import PerfData
from canopsis.common.math_parser import Formulas
from canopsis.engines.perfdata_utils.perfDataUtils import PerfDataUtils

import hashlib
from time import gmtime


class engine(Engine):
    etype = 'consolidation'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.storage = get_storage(
            namespace='events',
            account=Account(
                user="root",
                group="root"
            )
        )
        self.manager = PerfData()
        self.perf_data = PerfDataUtils()

    def pre_run(self):
        self.storage = get_storage(namespace='object',
            account=Account(user="root", group="root"))
        self.manager = PerfData()

    def publish_aggre_stats(self):

        series_event = forger(
            connector='engine',
            connector_name='engine',
            event_type='perf',
            source_type='resource',
            resource='series_events',
            state=0,
            perf_data_array=self.perf_data_array
        )

        rk = get_routingkey(series_event)

        self.logger.debug('Publishing {0} : {1}'.format(rk, series_event))

        self.amqp.publish(
            series_event,
            rk,
            self.amqp.exchange_name_events
        )

    def fetch(self, serie, _from, _to):
        self.logger.debug("*Start fetch*")
        t_serie = serie.copy()
        timewindow = {'start': _from, 'stop': _to, 'timezone':gmtime()}
        if len(t_serie['metrics']) > 1 and t_serie['aggregate_method'].lower() == 'none':
            self.logger.debug('More than one metric in serie, performing an aggregation')
            self.logger.debug('serie:', t_serie)
            self.logger.debug('aggregation: average - 60s')
            t_serie['aggregate_method'] = 'average'
            t_serie['aggregate_interval'] = 60
        if t_serie['aggregate_method'].lower() == 'none':
            timeserie = {'aggregation':'NONE'}
            results = self.perf_data.perfdata(metric_id=t_serie['metrics'], timewindow=timewindow, timeserie=timeserie)
        else:
            timeserie = {'aggregation':t_serie['aggregate_method'], 'period':{'second':t_serie['aggregate_interval']}}
            results = self.perf_data.perfdata(metric_id=t_serie['metrics'], timewindow=timewindow, timeserie=timeserie)

        formula = t_serie['formula']

        finalserie = self.metric_raw(results, formula)
        self.logger.debug("*End fetch*")

        return finalserie

    def metric_raw(self, results, formula):
        #nmetric = results[1]
        metrics, _ = results
        # Build points dictionnary
        points = {}
        length = False
        for m in metrics:
            cid = m['meta']['data_id']
            mid = 'metric_' + hashlib.md5(cid).hexdigest()
            mname = self.retreive_metric_name(cid)
            # Replace metric name in formula by the unique id
            formula = formula.replace(mname, mid)
            self.logger.debug("Metric {0} - {1}".format(mname, mid))
            points[mid] = m['points']
            # Make sure we treat the same amount of points by selecting
            # The metric with less points.
            if not length or len(m['points']) < length:
                length = len(m['points'])
        self.logger.debug('formula: {}'.format(formula))
        #self.logger.debug('points: {}'.format(points))

        mids = points.keys()
        finalSerie = []

        # Now loop over all points to calculate the final serie
        for i in range(length):
            data = {}
            ts = 0
            for j in range(len(mids)):
                mid = mids[j]
                # Get point value at timestamp "i" for metric "mid"
                data[mid] = points[mid][i][1]

                # Set timestamp
                ts = points[mid][i][0]

            # import data in math context
            math = Formulas(data)
            # describe the mathematic formula here.
            pointval = math.evaluate(formula)

            # Add computed point in the serie
            finalSerie.append([ts * 1000, pointval])
            # Remove data from math context
            math.reset()

        self.logger.debug('finalserie: {}'.format(finalSerie))

        return finalSerie, points[mid]

    def retreive_metric_name(self, name):
        '''
        This method allow to slice data from an existing one.
        TODO: improve this method with the Context ID.
        '''
        if name is None:
            return None
        li = name.split('/')
        for i in range(4):
            li.pop(0)
        name = '/'+'/'.join(li)
        return name

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.debug("Start metrics consolidation")
        serie = event
        if not serie:
            # Show error message
            self.logger.error('No record found.')
        # Test Settings
        t_serie = serie.copy()
        _from = 1425394522
        _to = 1425402296
        t_serie['aggregate_interval'] = 60
        t_serie['aggregate_method'] = 'average'
        timewindow = {'start': _from, 'stop': _to, 'timezone':gmtime()}        

        m_id = serie['metrics']
        timeserie = {'aggregation':t_serie['aggregate_method'], 'period':{'second':t_serie['aggregate_interval']}}
        results, _ = self.perf_data.perfdata(metric_id=t_serie['metrics'], timewindow=timewindow, timeserie=timeserie)
        # End test
        #results, _ = self.perf_data.perfdata(metric_id=m_id)
        self.logger.debug(results[0])
        _, points = self.fetch(serie, _from, _to)
        # Update time series data
        self.manager.put(metric_id=serie['_id'], points=points)

        event_id = event['_id']
        # Update crecords informations
        self.crecord_task_complete(event_id)
