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
        self.logger.debug("Je passe dans fetch \n\n\n")
        t_serie = serie.copy()
        timewindow = {'start': _from/1000, 'stop': _to/1000, 'timezone':gmtime()}
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

        return finalserie

    def metric_raw(self, results, formula):
        #nmetric = results[1]
        metrics, _ = results
        #  build points dictionnary
        points = {}
        length = False
        for m in metrics:
            cid = m.meta.data_id
            mid = 'metric_' + hashlib.md5(cid)
            mname = self.retreive_metric_name(cid)
            # replace metric name in formula by the unique id
            formula = formula.replace(mname, mid)
            self.logger.debug("Metric {0} - {1}".format(mname, mid))
            points[mid] = m.points
            # make sure we treat the same amount of points by selecting
            # the metric with less points.
            if not length or len(m.points) < length:
                length = len(m.points)
        self.logger.debug('formula:', formula)
        self.logger.debug('points:', points)

        mids = points.keys()
        finalSerie = []

        # now loop over all points to calculate the final serie
        for i in range(length):
            data = {}
            ts = 0
            for j in range(len(mids)):
                mid = mids[j]
                # get point value at timestamp "i" for metric "mid"
                data[mid] = points[mid][i][1]

                # set timestamp
                ts = points[mid][i][0]

            # import data in math context
            math = Formulas(data)
            pointval = math.evaluate(formula)

            # Add computed point in the serie
            finalSerie.append([ts * 1000, pointval])

        self.logger.debug('finalserie:', finalSerie)

        return finalSerie

    def retreive_metric_name(self, name):
        '''
        Impove this method with the Context ID.
        '''
        if name is None:
            return None
        li = name.split('/')
        for i in range(3):
            li.pop(i)
        name = '/'+'/'.join(li)
        return name

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.debug("Start metrics consolidation")
        serie = event
        #manager = PerfData()

        if not serie:
            # Show error message
            self.logger.error('No record found.')

        s_id = serie['metrics']
        self.logger.debug(s_id)
        results, _ = self.perf_data.perfdata(metric_id=s_id)
        #results, count = self.manager.get(metric_id=s_id)
        self.logger.debug(results[0]['points'])
        #self.fetch(serie, 'from', 'to')
        event_id = event['_id']
        # Update crecords informations
        self.crecord_task_complete(event_id)
