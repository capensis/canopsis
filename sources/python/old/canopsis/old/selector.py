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

#import logging
from canopsis.old.record import Record
from canopsis.old.downtime import Downtime
from canopsis.old.event import get_routingkey, forger

from time import time
from json import dumps, loads
from logging import getLogger


class Selector(Record):
    def __init__(
        self, storage,
        _id=None,
        name=None,
        namespace='events',
        record=None,
        logging_level=None
        ):
        self.type = 'selector'
        self.storage = storage

        if name and not _id:
            self._id = self.type + "." + storage.account._id + "." + name

        ## Default vars
        self.namespace = namespace

        self.dostate = True

        self.data = {}
        self.mfilter = {}
        self.include_ids = []
        self.exclude_ids = []
        self.changed = False
        self.rk = None

        self.last_event = None

        self.output_tpl = "{cps_sel_state_0} Ok, {cps_sel_state_1} \
        Warning, {cps_sel_state_2} Critical"

        self.sel_metric_prefix = "cps_sel_"
        self.sel_metric_name = self.sel_metric_prefix + "state_%s"

        self.logger = getLogger('Selector')
        self.cdowntime = Downtime(self.logger)
        if logging_level:
            self.logger.setLevel(logging_level)


        self.load(record.dump())

    def load(self, dump):

        self.logger.debug('Loading selector from record')

        mfilter = dump.get('mfilter', '{}')


        if type(mfilter) == dict:
            self.mfilter = mfilter
        elif mfilter == None:
            self.mfilter = {}
        else:
            try:
                self.mfilter = loads(mfilter)
            except Exception as e:
                self.logger.warning('invalid mfilter for selector {} : {} '.format(
                    dump.get('display_name', 'noname'), e))


        self._id = dump.get('_id')
        self.namespace = dump.get('namespace', self.namespace)
        self.dostate = dump.get('dostate')
        self.display_name = dump.get('display_name', 'noname')
        self.rk = dump.get('rk', self.rk)
        self.include_ids = dump.get('include_ids', self.include_ids)
        self.exclude_ids = dump.get('exclude_ids', self.exclude_ids)
        self.output_tpl = dump.get('output_tpl', None)
        self.sla_rk = dump.get('sla_rk', None)

        self.data = dump

    ## Build MongoDB query to find every id matching event
    def makeMfilter(self):
        self.logger.debug("Make filter:")
        (ifilter, efilter, mfilter) = ({}, {}, {})

        if self.include_ids:
            if len(self.include_ids) == 1:
                ifilter = {'_id': self.include_ids[0]}
            else:
                ifilter = {'_id': {'$in': self.include_ids}}

        if self.exclude_ids:
            if len(self.exclude_ids) == 1:
                efilter = {'_id': {'$ne': self.exclude_ids[0]}}
            else:
                efilter = {'_id': {'$nin': self.include_ids}}

        if self.mfilter:
            mfilter = self.mfilter

        self.logger.debug(" + ifilter: %s" % ifilter)
        self.logger.debug(" + efilter: %s" % efilter)
        self.logger.debug(" + mfilter: %s" % mfilter)

        #Adds downtime elements to ignore in query
        downtime = self.cdowntime.get_filter()
        if downtime:
            self.logger.debug(
                ' + Selector downtime exclusion %s' % downtime)

        ## Tweaks
        if not mfilter and not ifilter and not efilter:
            self.logger.warning("%s: Invalid filter" % self.name)
            return None

        if mfilter and not ifilter and not efilter:
            if downtime:
                return {'$and': [downtime, mfilter]}
            else:
                return mfilter

        if not mfilter and ifilter and not efilter:
            if downtime:
                return {'$and': [downtime, ifilter]}
            else:
                return ifilter

        if not mfilter and not ifilter and efilter:
            return None

        if mfilter and ifilter and not efilter:
            filters = {"$or": [mfilter, ifilter]}
            if downtime:
                return {'$and': [downtime, filters]}
            else:
                return filters

        if mfilter and not ifilter and efilter:
            filters = [mfilter, efilter]

            if downtime:
                filters = [downtime] + filters
            return {'$and': filters}

        and_clause = [{"$or": [mfilter, ifilter]}, efilter]

        if downtime:
            and_clause = downtime + and_clause

        return {"$and": and_clause}


    def getState(self):
        self.logger.debug("getStates:")

        # Build MongoDB filter
        mfilter = self.makeMfilter()
        if not mfilter:
            self.logger.debug("  + Invalid filter")
            return ({}, 3, 1)

        # Check filter
        self.logger.debug(" + filter: %s" % mfilter)

        self.logger.debug(" + selector statment agregation")
        result = self.storage.get_backend(namespace=self.namespace).aggregate([
                {'$match': mfilter},
                {'$project': {
                    '_id': True,
                    'state': True,
                    'state_type': True,
                    'previous_state': True,
                }},
                {'$group': {
                        '_id': {
                            'state': '$state',
                            'state_type': "$state_type",
                            'previous_state': "$previous_state",
                        },
                        'count': {'$sum': 1}}}])

        self.logger.debug(" + result: %s" % result)

        states = {}
        total = 0
        for state in result['result']:
            key = state['_id']['state']

            if state['_id'].get('state_type', 1) == 0:
                key = state['_id'].get('previous_state', key)

            states[key] = states.get(key, 0) + state['count']
            total += state['count']

        self.logger.debug(" + namespace: {}".format(self.namespace))
        self.logger.debug(" + states: {}".format(states))
        self.logger.debug(" + total: {}".format(total))


        state, state_type = 0, 1
        # Compute worst state
        for s in [0, 1, 2, 3]:
            if s in states:
                state = s


        result = self.storage.get_backend(namespace=self.namespace).aggregate([
                {'$match': mfilter},
                {
                    "$group": {
                        "_id": { "isAck": "$ack.isAck" },
                        "count": {
                            "$sum": {
                                "$cond": [ "$ack.isAck", 1, 0 ]
                            }
                        }
                    }
                },
                {
                    "$project": {
                        "_id": 0,
                        "ack": "$_id",
                        "count": 1
                    }
                }
        ])

        ack_count = result['result'][0]['count']

        self.logger.debug(" + result for ack : %s" % result)

        return (states, state, state_type, ack_count)

    def event(self):

        ### Transform Selector to Canopsis Event
        self.logger.debug("To Event:")

        # Get state
        (states, state, state_type, ack_count) = self.getState()

        # Build output
        total = 0
        for s in states:
            states[s] = int(states[s])
            total += states[s]

        send_ack = total == ack_count


        self.logger.debug(" + state: %s" % state)
        self.logger.debug(" + state_type: %s" % state_type)

        perf_data_array = []
        long_output = ""
        output = ""

        self.logger.debug(" + total: %s" % total)

        # Create perfdata array
        output_data = {}
        for i in [0, 1, 2, 3]:
            value = 0
            try:
                value = states[i]
            except:
                pass

            metric = self.sel_metric_name % i
            output_data[metric] = value
            perf_data_array.append(
                {"metric": metric, "value": value, "max": total})
            self.logger.info('metric %s : %s' % (metric, value))
        perf_data_array.append(
            {"metric": self.sel_metric_prefix + "total", "value": total})

        output_data['total'] = total

        # Fill Output template
        self.logger.debug(" + output TPL: %s" % self.output_tpl)
        output = self.output_tpl
        if output and output_data:
            for key in output_data:
                output = output.replace("{%s}" % key, str(output_data[key]))


        # Debug
        self.logger.debug(" + Display Name: %s" % self.display_name)
        self.logger.debug(" + output: %s" % output)
        self.logger.debug(" + long_output: %s" % long_output)
        self.logger.debug(" + perf_data_array: %s" % perf_data_array)

        # Build Event
        event = forger(
            connector="selector",
            connector_name="engine",
            event_type="selector",
            source_type="component",
            component='selector',
            resource=self.display_name,
            state=state,
            state_type=state_type,
            output=output,
            long_output=long_output,
            perf_data=None,
            perf_data_array=perf_data_array,
            display_name=self.display_name
        )

        # Extra field
        event["selector_id"] = str(self._id)

        # Build RK
        rk = get_routingkey(event)

        # Save RK
        if not self.rk:
            self.logger.debug("Set RK to '%s'" % rk)
            self.storage.update(self._id, {'rk': rk})
            self.rk = rk

        # Cache event
        self.last_event = event

        return (rk, event, send_ack)
