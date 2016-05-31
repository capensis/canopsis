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

from canopsis.engines.core import Engine

from bson import BSON
from json import loads
from time import time

from canopsis.common.init import basestring, PYVER
from canopsis.common.utils import ensure_unicode, forceUTF8
from canopsis.monitoring.parser import PerfDataParser


class engine(Engine):
    etype = "cleaner"

    def work(self, body, msg, *args, **kargs):
        ## Sanity Checks
        rk = msg.delivery_info['routing_key']
        if not rk:
            raise Exception("Invalid routing-key '%s' (%s)" % (rk, body))

        #self.logger.info( body )
        ## Try to decode event
        if isinstance(body, dict):
            event = body
            # force utf8 only if python version is 2
            if PYVER < '3':
                event = forceUTF8(event)
        else:
            self.logger.debug(" + Decode JSON")
            try:
                if isinstance(body, basestring):
                    try:
                        event = loads(body)
                        self.logger.debug("   + Ok")
                    except Exception as err:
                        try:
                            self.logger.debug(" + Try hack for windows string")
                            # Hack for windows FS -_-
                            event = loads(body.replace('\\', '\\\\'))
                            self.logger.debug("   + Ok")
                        except Exception as err:
                            try:
                                self.logger.debug(" + Decode BSON")
                                bson = BSON(body)
                                event = bson.decode()
                                self.logger.debug("   + Ok")
                            except Exception as err:
                                raise Exception(err)

            except Exception as err:
                self.logger.error("   + Failed (%s)" % err)
                self.logger.debug("RK: '%s', Body:" % rk)
                self.logger.debug(body)
                raise Exception("Impossible to parse event '%s'" % rk)

        event['rk'] = ensure_unicode(rk)

        if "resource" in event:
            if not isinstance(event['resource'], basestring):
                event['resource'] = ''
            else:
                event['resource'] = ensure_unicode(event['resource'])
            if not event['resource']:
                del event['resource']

        # Clean tags field
        event['tags'] = event.get('tags', [])

        tags = event['tags']

        if isinstance(tags, basestring) and tags != "":
            event['tags'] = [event['tags']]

        elif not isinstance(tags, list):
            event['tags'] = []

        event["timestamp"] = int(event.get("timestamp", time()))

        event["state"] = event.get("state", 0)
        event["state_type"] = event.get("state_type", 1)
        event["event_type"] = event.get("event_type", "check")

        default_status = 0 if not event["state"] else 1
        event["status"] = event.get("status", default_status)

        event['output'] = event.get('output', '')

        # Get perfdata
        perf_data = event.get('perf_data')
        perf_data_array = event.get('perf_data_array')

        if perf_data_array is None:
            perf_data_array = []

        # Parse perfdata
        if perf_data:
            self.logger.debug(' + perf_data: {0}'.format(perf_data))

            try:
                perf_data_array += PerfDataParser(perf_data).perf_data_array

            except Exception as err:
                self.logger.error(
                    "Impossible to parse perfdata from: {0} ({1})".format(
                        event, err
                    )
                )

            event['perf_data_array'] = perf_data_array

        return event
