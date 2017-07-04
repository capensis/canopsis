#!/usr/bin/env python2
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

import json
import unittest

from canopsis.common.test_base import BaseApiTest, Method
from canopsis.watcher.manager import Watcher
from canopsis.webcore.utils import HTTP_ERROR, HTTP_NOT_ALLOWED, HTTP_OK

# Sample watcher (to insert)
watcher_dict = {
    "alert_level": "minor",
    "crecord_write_time": None,
    "crecord_type": "selector",
    "crecord_creation_time": None,
    "crecord_name": None,
    "description": "a_description",
    "display_name": "a_displayed_name",
    "dosla": True,
    "dostate": True,
    "downtimes_as_ok": True,
    "enable": True,
    "exclude_ids": [],
    "include_ids": [],
    "last_dispatcher_update": None,
    "loaded": False,
    "mfilter": "{\"$or\":[{\"un \":{\"$eq\":\"filtre\"}}]}",
    "output_tpl": "Off: [OFF], Minor: [MINOR], Major: [MAJOR], Critical: [CRITICAL], Ack count [ACK], Total: [TOTAL]",
    "sla_critical": 75,
    "sla_timewindow": {
        "value": 12,
        "durationType": "second",
        "seconds": 12
    },
    "sla_output_tpl": "Available: [P_AVAIL]%, Off: [OFF]%, Minor: [MINOR]%, Major: [MAJOR]%, Critical: [CRITICAL]%, Alerts [ALERTS]%, sla start [TSTART],  time available [T_AVAIL], time alert [T_ALERT]",
    "sla_warning": 90,
    "state": None,
    "state_algorithm": None,
    "state_when_all_ack": "worststate",
    "_id": "watcher_id"
}


class TestWatcherAPI(BaseApiTest):

    def setUp(self):
        self._authenticate()  # default setup

        self.watcher = Watcher()

        self.base = '{}/{}'.format(self.URL_BASE, 'api/v2/watchers')

    def test_CRUD(self):
        watcher_id = watcher_dict['_id']

        # GET
        r = self._send(url=self.base,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP_NOT_ALLOWED)

        r = self._send(url=self.base + '/' + watcher_id,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP_ERROR)

        # POST
        r = self._send(url=self.base,
                       method=Method.post)
        self.assertEqual(r.status_code, HTTP_ERROR)

        r = self._send(url=self.base,
                       method=Method.post,
                       data=json.dumps(watcher_dict))
        self.assertEqual(r.status_code, HTTP_OK)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        # Nothing more to validate on this route

        # PUT
        r = self._send(url=self.base,
                       method=Method.put)
        self.assertEqual(r.status_code, HTTP_NOT_ALLOWED)

        params = {
            "sla_critical": 75,
            "downtimes_as_ok": True,
            "alert_level": "major",
        }
        r = self._send(url=self.base + '/' + watcher_id,
                       method=Method.put,
                       data=json.dumps(params))
        self.assertEqual(r.status_code, HTTP_OK)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        # Nothing more to validate on this route

        # GET (again)
        r = self._send(url=self.base + '/' + watcher_id,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP_OK)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertEqual(data['sla_critical'], 75)
        self.assertEqual(data['alert_level'], "major")

        # DELETE
        r = self._send(url=self.base,
                       method=Method.delete)
        self.assertEqual(r.status_code, HTTP_NOT_ALLOWED)

        r = self._send(url=self.base + '/' + watcher_id,
                       method=Method.delete)
        self.assertEqual(r.status_code, HTTP_OK)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertTrue('ok' in data)
        self.assertEqual(data['ok'], 1.0)

if __name__ == "__main__":
    # Warning ! Can polluate the database
    unittest.main()
