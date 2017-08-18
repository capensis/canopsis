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

from canopsis.engines.core import Engine, publish
from canopsis.event import forger

from time import time

from canopsis.engines.collectd_utils import types

import re

class engine(Engine):
    etype = 'collectdgw'

    def new_amqp_queue(self, *args, **kwargs):
        """
        Override AMQP queue creation (ignore possible parameters,
        they aren't needed here)
        """

        self.amqp.add_queue(
            self.amqp_queue, ['collectd'], self.on_collectd_event, "amq.topic",
            auto_delete=False)

    def on_collectd_event(self, body, msg):
        start = time()
        error = False

        prog = re.compile('^(PUTVAL) ("(.+)"|([^\s]+)) (interval=.+) ([^\s]+)$')
        match = re.match(prog, body)
        if match:
            if match.group(3):
                cnode = match.group(3)
            else:
                cnode = match.group(4)
            collectd_info = [match.group(1),cnode,match.group(5),match.group(6)]
            self.logger.debug(body)
            action = collectd_info[0]
            self.logger.debug(" + Action: %s" % action)

            if len(collectd_info) == 4 and action == "PUTVAL":
                cnode = collectd_info[1].split("/")
                component = cnode[0]
                resource = cnode[1]
                metric = cnode[2]
                options = collectd_info[2]
                values = collectd_info[3]

                self.logger.debug(" + Options: %s" % options)
                self.logger.debug(" + Component: %s" % component)
                self.logger.debug(" + Resource: %s" % resource)
                self.logger.debug(" + Metric: %s" % metric)
                self.logger.debug(" + Raw Values: %s" % values)

                values = values.split(":")

                perf_data_array = []

                ctype = None
                try:
                    ## Know metric
                    ctype = types[metric]
                except:
                    try:
                        ctype = types[metric.split('-')[0]]
                        metric = metric.split('-')[1]
                    except Exception as err:
                        self.logger.error(
                            "Invalid format '%s' (%s)" % (body, err))
                        return None

                try:
                    timestamp = int(float(values[0]))
                    values = values[1:]
                    self.logger.debug("   + Timestamp: %s" % timestamp)
                    self.logger.debug("   + Values: %s" % values)

                except Exception as err:
                    self.logger.error(
                        "Impossible to get timestamp or values (%s)" % err)
                    return None

                self.logger.debug(" + metric: %s" % metric)
                self.logger.debug(" + ctype: %s" % ctype)

                if ctype:
                    try:
                        i = 0
                        for value in values:
                            name = ctype[i]['name']
                            unit = ctype[i]['unit']
                            vmin = ctype[i]['min']
                            vmax = ctype[i]['max']

                            if vmin == 'U':
                                vmin = None

                            if vmax == 'U':
                                vmax = None

                            if name == "value":
                                name = metric

                            if metric != name:
                                name = "%s-%s" % (metric, name)

                            data_type = ctype[i]['type']

                            value = float(value)

                            self.logger.debug("     + %s" % name)
                            self.logger.debug(
                                "       -> %s (%s)" % (value, data_type))
                            i += 1

                            perf_data_array.append(
                                {
                                    'metric': name, 'value': value,
                                    'type': data_type, 'unit': unit,
                                    'min': vmin, 'max': vmax
                                }
                            )

                    except Exception as err:
                        self.logger.error(
                            "Impossible to parse values '%s' (%s)" % (
                                values, err))

                if perf_data_array:
                    self.logger.debug(
                        ' + perf_data_array: %s', perf_data_array)

                    event = forger(
                        connector='collectd',
                        connector_name='collectd2event',
                        component=component,
                        resource=resource,
                        timestamp=timestamp,
                        source_type='resource',
                        event_type='perf',
                        perf_data_array=perf_data_array
                    )

                    self.logger.debug("Send Event: %s" % event)

                    ## send event on amqp
                    publish(publisher=self.amqp, event=event)

            else:
                error = True
                self.logger.error("Invalid collectd Action (%s)" % body)

        else:
            self.logger.error("Invalid collectd Message (%s)" % body)

        if error:
            self.counter_error += 1

        self.counter_event += 1
        self.counter_worktime += time() - start
