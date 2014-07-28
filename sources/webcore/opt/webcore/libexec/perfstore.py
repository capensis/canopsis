#!/usr/bin/env python
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

import logging
import json
import time
from datetime import datetime

from bottle import route, get, post, put, delete, request, HTTPError, response

#import protection function
from libexec.auth import get_account

# Modules
from canopsis.old.storage import get_storage

from canopsis.old.tools import clean_mfilter

from canopsis.perfdata.manager import PerfData
from canopsis.timeserie import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period
from canopsis.storage.periodic import PeriodicStorage

from canopsis.old.account import Account

storage = get_storage(
    namespace='object', account=Account(user="root", group="root"))

manager = None

SERVICE_NAME = 'perfstore'

logger = logging.getLogger(SERVICE_NAME)


def load():
    global logger
    global manager

    manager = PerfData(log_lvl=logger.level)

def unload():
    global manager
    del manager

group_managing_access = ['group.CPS_perfdata_admin']
#########################################################################


#### POST@
@post('/{0}/values'.format(SERVICE_NAME))
@post('/{0}/values/:start/:stop'.format(SERVICE_NAME))
def perfstore_values_route(start=None, stop=None):
    """subset selection param allow filter metrics with exclusion periods and component,resource,hostgroup exclusion"""
    return perfstore_nodes_get_values(
        start=start,
        stop=stop,
        metas=request.params.get('nodes', default=None),
        aggregate_method=request.params.get('aggregate_method', default=None),
        aggregate_interval=request.params.get('aggregate_interval', default=None),
        aggregate_max_points=request.params.get('aggregate_max_points', default=None),
        aggregate_round_time=request.params.get('aggregate_round_time', default=None),
        consolidation_method=request.params.get('consolidation_method', default=None),
        timezone=request.params.get('timezone', default=0),
        subset_selection=request.params.get('subset_selection', default={}))


@get('/{0}'.format(SERVICE_NAME))
@get('/{0}/get_all_metrics'.format(SERVICE_NAME))
def perfstore_get_all_metrics():
    return perfstore_get_all_metrics(
        limit=int(request.params.get('limit', default=20)),
        start=int(request.params.get('start', default=0)),
        search=request.params.get('search', default=None),
        filter=request.params.get('filter', default=None),
        sort=request.params.get('sort', default=None),
        show_internals=request.params.get('show_internals', default=False))


def perfstore_nodes_get_values(
    start=None,
    stop=None,
    metas=None,
    aggregate_method=None,
    aggregate_interval=None,
    aggregate_max_points=None,
    aggregate_round_time=None,
    consolidation_method=None,
    timezone=0,
    subset_selection=dict()):

    # set subset_selection
    if subset_selection:
        try:
            subset_selection = json.loads(subset_selection)
            logger.debug('subset selection found: {0}'.format(subset_selection))
        except:
            subset_selection = {}
            logger.warning('Unable to load subset_selection filters from params')

    # set timewindow
    timewindow = TimeWindow(start=None, stop=None, timezone=timezone)

    # set timeserie if required
    if aggregate_method is not None:
        if aggregate_max_points:
            aggregate_max_points = int(aggregate_max_points)
        if aggregate_round_time:
            aggregate_round_time = aggregate_round_time == 'true'
        if aggregate_interval:
            period = Period(second=int(aggregate_interval))
        else:
            period = None
        timeserie_kwargs = {
            'max_points': aggregate_max_points,
            'period': period,
            'round_time': aggregate_round_time,
            'aggregation': aggregate_method,
        }

        timeserie = TimeSerie(**timeserie_kwargs)
    else:
        timeserie = None

    if manager is None:
        load()

    output = []

    if not metas:
        logger.warning("Invalid arguments")
        return HTTPError(400, "Invalid arguments. metas are necassaries")

    metas = json.loads(metas)

    logger.debug("POST:")
    logger.debug(" + metas: %s" % metas)
    logger.debug(" + aggregate_method: %s" % aggregate_method)
    logger.debug(" + aggregate_interval: %s" % aggregate_interval)
    logger.debug(" + aggregate_max_points: %s" % aggregate_max_points)
    logger.debug(" + aggregate_round_time: %s" % aggregate_round_time)
    logger.debug(" + consolidation_method: %s" % consolidation_method)
    logger.debug(" + timezone: %s" % timezone)

    for meta in metas:
        _id = meta.get('id', None)

        # update specific timewindow
        mstart = meta.get('from', timewindow.start())
        mstop = meta.get('to', timewindow.stop())
        meta_timewindow = TimeWindow(start=mstart, stop=mstop, timezone=timezone)

        # get perfdata values for each meta
        if _id is not None:
            output += perfstore_get_values(
                meta=meta,
                timewindow=meta_timewindow,
                timeserie=timeserie,
                subset_selection=subset_selection)

    if aggregate_method and consolidation_method and timeserie and len(output):
        consolidation_serie = TimeSerie(
            aggregation=consolidation_method,
            fill = timeserie.fill,
            round_time=timeserie.round_time,
            period=timeserie.period)

        points = [point for point in [result['values'] for result in output]]

        points.sort(key=itemgetter(0))

        consolidation_serie.calculate(points=points, timewindow=timewindow)

        output = [{
            'node': output[0]['node'],
            'metric': consolidation_method,
            'bunit': None,
            'type': 'GAUGE',
            'values': points
        }]

    output = {'total': len(output), 'success': True, 'data': output}

    return output


def perfstore_get_all_metrics(limit=20, start=0, search=None, filter=None, sort=None, show_internals=False):

    logger.debug("perfstore_get_all_metrics:")

    if manager is None:
        load()

    if filter:
        try:
            filter = json.loads(filter)
        except ValueError as ve:
            logger.error("Impossible to decode filter: {0}: {1}".format(filter, ve))
            filter = None

    show_internals = show_internals == "true"

    msort = []
    if sort:
        sort = json.loads(sort)
        for item in sort:
            direction = 1
            if str(item['direction']) == "DESC":
                direction = -1
            _property = str(item['property'])
            if _property == 'co':
                _property = 'component'
            elif _property == 're':
                _property = 'resource'
            elif _property == 'me':
                _property = 'name'
            msort.append((_property, direction))
    else:
        msort.append(('component', 1))

    logger.debug(" + limit: {0}".format(limit))
    logger.debug(" + start: {0}".format(start))
    logger.debug(" + search: {0}".format(search))
    logger.debug(" + sort: {0}".format(sort))
    logger.debug(" + msort: {0}".format(msort))
    logger.debug(" + filter: {0}".format(filter))
    logger.debug(" + show_internals: {0}".format(show_internals))

    mfilter = None

    if isinstance(filter, list):
        if len(filter) > 0:
            mfilter = filter[0]
        else:
            logger.error(" + Invalid filter format")

    elif isinstance(filter, dict):
        mfilter = filter

    if search:
        # Todo: Tweak this ...
        fields = ['component', 'resource', 'name']
        mor = []
        search = search.split(' ')
        if len(search) == 1:
            for field in fields:
                mor.append({field: {'$regex': '.*{0}.*'.format(search[0]), '$options': 'i'}})

            mfilter = {'$or': mor}
        else:
            mfilter = {'$and': []}
            for word in search:
                mor = []
                for field in fields:
                    mor.append({field: {'$regex': '.*{0}.*'.format(word), '$options': 'i'}})
                mfilter['$and'].append({'$or': mor})

    use_hint = False

    if not show_internals:
        if mfilter:
            mfilter['internal'] = False
            use_hint = True
        else:
            mfilter = {'internal': False}

    logger.debug(" + mfilter: {0}".format(mfilter))

    mfilter = clean_mfilter(mfilter)
    mfilter['type'] = 'metric'
    data = manager.entities.find(mfilter, limit=limit, skip=start, data=False, sort=msort)

    if use_hint:
        data.hint([('type', 1), ('component', 1), ('resource', 1), ('id', 1)])

    total = data.count()

    if isinstance(data, dict):
        data = [data]

    elif data is not None:
        l = list()
        for document in data:
            # TODO : remove this when uiv2 will be available
            document['component'], document['co'] = document.get('co'), document.get('component')
            document['resource'], document['re'] = document.get('re'), document.get('resource')
            document['name'], document['me'] = document.get('me'), document.get('name')
            l.append(document)
        data = l

    else:
        data = list()

    logger.debug(" + data: {0}".format(data))

    result = {'success': True, 'data' : data, 'total' : total}

    logger.debug(" + result: {0}".format(result))

    return result


### manipulating meta
@delete('/perfstore', checkAuthPlugin={'authorized_grp': group_managing_access})
@delete('/perfstore/:_id', checkAuthPlugin={'authorized_grp': group_managing_access})
def remove_meta(_id=None):

    if not _id:
        _id =  json.loads(request.body.readline())
    if not _id:
        return HTTPError(500, "No ids provided, bad request")

    if not isinstance(_id, list):
        _id = [_id]

    logger.debug('delete {0}: '.format(_id))

    for item in _id:
        if isinstance(item, dict):
            manager.entities.remove(_id=item['_id'])
        else:
            manager.entities.remove(_id=item)


@put('/perfstore', checkAuthPlugin={'authorized_grp': group_managing_access})
def update_meta(_id=None):
    data = json.loads(request.body.readline())

    if not isinstance(data, list):
        data = [data]

    for item in data:
        try:
            if not _id:
                _id = item['_id']

            if '_id' in item:
                del item['_id']

            manager.entities.update({'_id': _id}, {'$set': item})

        except Exception as err:
            logger.warning('Error while updating meta_id: {0}'.format(err))
            return HTTPError(500, "Error while updating meta_id: {0}".format(err))

from canopsis.engines.pyperfstore3.store import Store


#### POST@
@route('/perfstore/perftop')
@route('/perfstore/perftop/:start/:stop')
def perfstore_perftop(start=None, stop=None):

    data = []

    limit                   = int(request.params.get('limit', default=10))
    sort                    = int(request.params.get('sort', default=1))
    mfilter                 = request.params.get('mfilter', default={})
    get_output              = request.params.get('output', default=False)
    time_window             = int(request.params.get('time_window', default=86400))
    threshold               = request.params.get('threshold', default=None)
    threshold_direction     = int(request.params.get('threshold_direction', default=-1))
    expand                  = request.params.get('expand', default=False)
    percent                 = request.params.get('percent', default=False)
    threshold_on_pct        = request.params.get('threshold_on_pct', default=False)
    report                  = request.params.get('report', default=False)

    export_csv              = request.params.get('csv', default=False)
    export_fields           = request.params.get('fields', default="[]")

    if percent == 'true':
        percent = True
    elif percent == 'false':
        percent = False

    if report == 'true':
        report = True
    elif report == 'false':
        report = False

    if threshold_on_pct == 'true':
        threshold_on_pct = True
    elif threshold_on_pct == 'false':
        threshold_on_pct = False

    sort_on_percent = False
    if percent is True:
        sort_on_percent = True

    if mfilter:
        try:
            mfilter = json.loads(mfilter)
        except ValueError as err:
            logger.error("Impossible to decode mfilter: {0}: {1}".format(mfilter, err))
            mfilter = None

    try:
        export_fields = json.loads(export_fields)

    except ValueError as err:
        logger.error("Impossible to decode export_fields: {0}: {1}".format(export_fields, err))
        export_fields = []


    if threshold:
        threshold = float(threshold)

    if expand == 'true':
        expand = True
    else:
        expand = False

    if stop:
        stop = int(stop)
    else:
        stop = int(time.time())

    if start:
        start = int(start)
    else:
        start = stop - time_window

    logger.debug("PerfTop:")
    logger.debug(" + mfilter:     %s" % mfilter)
    logger.debug(" + get_output:  %s" % get_output)
    logger.debug(" + limit:       %s" % limit)
    logger.debug(" + threshold:   %s" % threshold)
    logger.debug(" + threshold_direction:   %s" % threshold_direction)
    logger.debug(" + sort:        %s" % sort)
    logger.debug(" + expand:       %s" % expand)
    logger.debug(" + report:       %s" % report)
    logger.debug(" + percent:       %s" % percent)
    logger.debug(" + threshold_on_pct:       %s" % threshold_on_pct)
    logger.debug(" + time_window: %s" % time_window)
    logger.debug(" + start:       %s (%s)" % (start, datetime.utcfromtimestamp(start)))
    logger.debug(" + stop:        %s (%s)" % (stop, datetime.utcfromtimestamp(stop)))
    logger.debug(" + export csv:  %s" % export_csv)
    logger.debug(" + export fields: %s" % str(export_fields))

    mfilter =  clean_mfilter(mfilter)

    # find the right type entity
    entities = manager.entities.find(mfilter, limit=1, projection={'nodeid'})
    entities.hint([('type', 1), ('component', 1), ('resource', 1), ('name', 1)])

    entity = None

    try:
        entity = entities[0]
    except KeyError:
        logger.error('No entity found with filter: {0}'.format(mfilter))
        return HTTPError(500, 'No entity found with filter: {0}'.format(mfilter))

    meta_data = manager.get_meta(data_id=entity['nodeid'], limit=1)
    mtype = None

    try:
        mtype = meta[0]
    except KeyError:
        pass

    def check_threshold(value):
        result = True
        if threshold:
            result = (threshold_direction == -1 and value >= threshold) or \
                (threshold_direction == 1 and value <= threshold)
        return result

    if mtype:
        mtype = mtype.get('type', 'GAUGE')

        logger.debug(" + mtype:    %s" % mtype)

        if mtype != 'COUNTER' and not expand and not report:

            # Quick method, use last value
            metrics = manager.get_meta(data_id=entity['nodeid'], sort=[('last_value', Store.ASC)], limit=limit)

            if isinstance(metrics, dict):
                metrics = [metrics]

            for metric in metrics:
                if (percent or threshold_on_pct) and 'max' in metric and 'last_value' in metric:
                    metric['pct'] = round(((metric['last_value'] * 100) / metric['max']) * 100) / 100

                if threshold_on_pct:
                    val = metric['pct']
                else:
                    val = metric['last_value']

                if check_threshold(val):
                    data.append(metric)

        else:

            metric_limit = 0

            if expand:
                metric_limit = 1

            #clean mfilter
            mfilter =  clean_mfilter(mfilter)

            metrics = manager.get_meta(data_id=entity['nodeid'], limit=limit)

            metrics.sort('last_value', sort)
            if isinstance(metrics, dict):
                metrics = [metrics]

            for metric in metrics:
                # Recheck type
                mtype = metric.get('type', 'GAUGE')
                if mtype != 'COUNTER' and not expand and not report:
                    logger.debug(" + Metric '%s' (%s) is not a COUNTER" % (metric['name'], metric['_id']))

                    if (percent or threshold_on_pct) and 'max' in metric and 'last_value' in metric:
                            metric['pct'] = round(((metric[PerfData.LAST_VALUE] * 100)/ metric['max']) * 100) / 100

                    if threshold_on_pct:
                        val = metric['pct']

                    else:
                        val = metric[PerfData.LAST_VALUE]

                    if check_threshold(val):
                        data.append(metric)

                else:

                    points = []
                    if mtype != 'COUNTER' and not expand:
                        # Get only one point
                        timewindow = TimeWindow(start=stop, stop=stop)
                        point = manager.get_point(data_id=metric[TimedStore.DATA_ID], timestamp=stop)
                        if point:
                            points = [point]
                    else:
                        # grt points from 'start' to 'stop'
                        timewindow = TimeWindow(start=start, stop=stop)
                        points = manager.get_points(data_id=metric[TimedStore.DATA_ID], timewindow=timewindow)

                    if len(points):
                        if expand:

                            del metric[TimedStore.DATA_ID]
                            for point in points:
                                if check_threshold(point[1]):
                                    nmetric = metric.copy()
                                    nmetric['last_timestamp'] = point[0]
                                    nmetric['last_value'] = point[1]
                                    if (percent or threshold_on_pct) and 'max' in nmetric and 'last_value' in nmetric:
                                        nmetric['pct'] = round(((nmetric['last_value'] * 100)/ nmetric['max']) * 100) / 100
                                    data.append(nmetric)
                        else:
                            # keep last point
                            metric['last_timestamp'] = points[len(points)-1][0]
                            metric['last_value'] = points[len(points)-1][1]

                            if (percent or threshold_on_pct) and 'max' in metric and 'last_value' in metric:
                                    metric['pct'] = round(((metric['last_value'] * 100)/ metric['max']) * 100) / 100

                            if threshold_on_pct:
                                val = metric['pct']
                            else:
                                val = metric['last_value']

                            if check_threshold(val):
                                data.append(metric)

        # Calculate most recurrent output
        if get_output:
            logs = get_storage(namespace='events_log', account=get_account())

            for item in data:
                evfilter = {'$and': [
                    {
                        'component': item['co'],
                        'resource': item.get('re', {'$exists': False}),
                        'state': {'$ne': 0}
                    },{
                        'timestamp': {'$gt': start}
                    },{
                        'timestamp': {'$lt': stop}
                    }
                ]}

                records = logs.find(evfilter)

                outputs = {}

                for record in records:
                    output = record.data['output']

                    if output not in outputs:
                        outputs[output] = 1

                    else:
                        outputs[output] += 1

                last_max = 0

                for output in outputs:
                    if outputs[output] > last_max:
                        item['output'] = output
                        last_max = outputs[output]

        reverse = True
        if sort == 1:
            reverse = False

        if sort_on_percent:
            for item in data:
                if not 'pct' in item:
                    item['pct'] = -1
            data = sorted(data, key=lambda k: k['pct'] , reverse=reverse)[:limit]
        else:
            data = sorted(data, key=lambda k: k['last_value'], reverse=reverse)[:limit]
    else:
        logger.debug("No records found")

    if not export_csv:
        return {'success': True, 'data' : data, 'total' : len(data)}

    else:
        response.headers['Content-Disposition'] = 'attachment; filename="perftop.csv"'
        response.headers['Content-Type'] = 'text/csv'

        exported = None

        logger.debug(' + Data: %s' % str(data))

        for entry in data:
            row = []

            for field in export_fields:
                value = entry.get(field, '')

                if isinstance(value, basestring):
                    value = value.replace('"', '""')
                    value = u'"{0}"'.format(value)

                else:
                    value = str(value)

                row.append(value)

            if exported:
                exported = u"{0}\n{1}".format(exported, u','.join(row))

            else:
                exported = u','.join(row)

        logger.debug(' + Exported: %s' % exported)

        return exported

########################################################################
# Functions
########################################################################

def perfstore_get_values(meta, timewindow, timeserie, subset_selection={}):

    logger.debug("Perfstore get values:")
    logger.debug(" + meta: {0}".format(meta))
    logger.debug(" + timewindow: {0}".format(timewindow))
    logger.debug(" + timeserie: {0}".format(timeserie))

    result = list()

    period = meta.get(PeriodicStorage.PERIOD)
    aggregation = meta.get(PeriodicStorage.AGGREGATION)

    points = []

    _id = meta.get('id')

    points, meta_data = manager.get(
        data_id=_id,
        with_meta=True,
        period=period,
        aggregation=aggregation,
        timewindow=timewindow)

    points = exclude_points(points, subset_selection)

    logger.debug(" + points: {0}, meta_data: {1}".format(points, meta_data))

    # For UI display
    if len(points) == 0 and meta_data and meta_data[0][1]['type'] == 'COUNTER':
        points = [(start, 0), (stop, 0)]

    if len(points) and meta_data and meta_data[0][1]['type'] == 'COUNTER':
        # Insert null point for aggregation
        points.insert(0, [points[0][0], 0])

    if timeserie is not None:
        points = timeserie.calculate(points=points, timewindow=timewindow)

    if meta_data:
        meta_data = meta_data[0][1]

    else:
        meta_data = dict()

    if points and meta:
        result.append(
            {'node': _id,
            'metric': meta.get('metric'),
            'values': points,
            'bunit': meta_data.get('unit'),
            'min': meta_data.get('min'),
            'max': meta_data.get('max'),
            'thld_warn': meta_data.get('thd_warn'),
            'thld_crit': meta_data.get('thd_crit'),
            'type': meta_data.get('type')})

    return result


def exclude_points(points, subset_selection={}):
    """unit test
    assert(exclude_points([[0,1],[0.5,2],[1,1],[2,3],[4,5],[3,1],[5,2]],{'intervals':[{'from':1,'to':3}]})\
     == [[0, 1], [0.5, 2], [1, None], [2, None], [4, 5], [3, None], [5, 2]], True)
    """

    # Compute exclusion periods and set a point to None value (for UI purposes) if point is in any exclusion period.
    exclusion_points = []
    if subset_selection and 'exclusions' in subset_selection:
        logger.debug('Interval exclusions detected, will apply it to output data')
        # Iterate over database point list for current metric.
        for value in points:
            is_excluded = False
            # Takes care of exclusion intervals given in parameters.
            for interval in subset_selection['exclusions']:
                if value[0] >= interval['from'] and value[0] <= interval['to']:
                    is_excluded = True
                    break
            if is_excluded:
                # Add a point that UI won t dispay.
                exclusion_points.append([value[0], None])
            else:
                # Nothing to do, just keep the original point
                exclusion_points.append(value)
        # returns the new computed point list for given exclusion interval
        return exclusion_points
    else:
        return points
