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

#from time import sleep

from test_base import BaseApiTest, Method, HTTP

PBH_RESOURCE = "pbh_resource"
PBH_CONNECTOR = "pbh_connector"
PBH_CON_NAME = PBH_CONNECTOR + "_name"
PBH_COMPONENT = "pbh_component"

RESOURCE_ID = "{0}/{1}".format(PBH_RESOURCE, PBH_COMPONENT)

FILTER = {
    "_id": "delete_me",
    "break": False,
    "crecord_write_time": 1497428303,
    "enable": True,
    "name": "Hi",
    "actions": [{
        "field": "output",
        "type": "override",
        "value": "yep_it's_me"
    }],
    "priority": 1,
    "run_once": True,
    "crecord_type": "filter",
    "mfilter": "{{\"$or\":[{{\"connector\":{{\"$eq\":\"{conn}\"}}}}]}}"
               .format(conn=PBH_CONNECTOR),
    "crecord_creation_time": 1497367809,
    "crecord_name": None,
    "description": "I am a nice filter"
    #TODO add in and out pbehaviors
}

EVENT = {
    "component": PBH_COMPONENT,
    "resource": PBH_RESOURCE,
    "source_type": "resource",
    "event_type": "check",
    "connector": PBH_CONNECTOR,
    "connector_name": PBH_CON_NAME,
    "output": "oh my god ! an awesome output",
    "state": 0
}


"""
TODO: cleanup this (we cannot have managers here and requests are made by
BaseApiTest class)

WEB_HOST = "localhost"
AUTH_KEY = "b71786a6-4c4c-11e7-8da1-0800279471b5"
MONGO_HOST = "localhost"

URL_BASE = "http://{0}:8082".format(WEB_HOST)
URL_AUTH = "{0}/?authkey={1}"
URL_PBEH = "{0}/pbehavior/create?{1}".format(URL_BASE, None)
URL_SEND = "{0}/event".format(WEB_HOST)
URL_MONGO = 'mongodb://cpsmongo:canopsis@{0}:27017/canopsis'.format(MONGO_HOST)

ENTITIES_COL = "default_entities"
OBJECT_COL = "object"

DEL_FILTER = {
    "_id": "delete_me",
    "name": "Hi",
    "actions": [{
        "field": "output",
        "type": "override",
        "value": "yep_it's_me"
    }],
    "crecord_type": "filter",
    "mfilter": "{{\"$or\":[{{\"connector\":{{\"$eq\":\"{conn}\"}}}}]}}"
               .format(conn=PBH_CONNECTOR),
    "description": "I am a nice filter"
}

to_delete = ["break", "crecord_write_time", "crecord_write_time", "priority",
             "run_once", "crecord_type", "crecord_creation_time",
             "crecord_name"]

DEL_FILTER = FILTER.copy()
for key in to_delete:
    if key in DEL_FILTER:
        DEL_FILTER.pop(key)

PBEHAVIOR = {
    "name": "A name",
    "filter": {"_id": RESOURCE_ID},
    "author": "Functionnal_test",
    "tstart": None,
    "tstop": None,
    "rrule": None,
    "enabled": True,
    "connector": PBH_CONNECTOR,
    "connector_name": PBH_CON_NAME
}

BEAT = 2

class BaseTest(unittest.TestCase):

    def _insert_filter(self, **kwargs):
        filter_ = FILTER.copy()

        for key in kwargs:
            filter_[key] = kwargs[key]

        self.obj_col.insert(filter_)
        print("Waiting {0}s for the beat".format(BEAT))
        sleep(BEAT)

    def _create_pbehavior(self, in_=True, **kwargs):
        pb = PBEHAVIOR.copy()
        middle = time()
        pb["tstart"] = middle - 5 * 60
        pb["tstop"] = middle + 5 * 60

        if in_ is False:
            pb["tstart"] += 3600
            pb["tstop"] += 3600
        elif in_ is not True:
            self.fail("Nope")

        key_to_add = PBEHAVIOR.keys()
        key_to_add.remove("tstart")
        key_to_add.remove("tstop")

        for key in key_to_add:
            if key in kwargs:
                pb[key] = kwargs[key]

        return pb

    def _push_pbehavior(self, pbehavior):
        response = self.session.post(URL_PBEH.format(urlencode(pbehavior)))

        response = literal_eval(response.text)

        if response["total"] == 1 and response["sucess"] is True:
            self.pbehavior_ids += response["data"]
        else:
            self.fail("Impossible to insert the pbehavior.")

    def setUp(self):
        #client = MongoClient(URL_MONGO)
        #self.ent_col = client.canopsis[ENTITIES_COL]
        #self.obj_col = client.canopsis[OBJECT_COL]

        self.pbehavior_ids = []

        #self.pbm = PBehaviorManager()
        #self.cm = ContextGraph()
"""


class TestPbehaviorAPI(BaseApiTest):

    def setUp(self):
        self._authenticate()  # default setup

        self.base = '{}/{}'.format(self.URL_BASE, 'pbehavior')

    def test_pbehavior(self):
        pbehavior_id = 'fake_pbehavior_id'

        # GET
        r = self._send(url=self.base + '/read')
        self.assertEqual(r.status_code, HTTP.OK.value)
        #json = r.json()
        #self.assertEqual(json['total'], 0)

        r = self._send(url=self.base + '/read',
                       params={'_id': pbehavior_id})
        json = r.json()
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertEqual(json['total'], 0)
        self.assertIsNone(json['data'])

        # TODO: test create/update/delete too...
        print('!!! Incomplete crud tests !!!')

    #def test_in_OK_out_OK(self):
    def in_OK_out_OK(self):
        # TODO: finalize this test (whatever it does)
        kwargs = {}
        #kwargs["rrule"] = self._create_rrule()
        pb_in = self._create_pbehavior(**kwargs)
        pb_out = self._create_pbehavior(in_=False, **kwargs)
        #self.pbm.compute_pbehaviors_filters()

        kwargs = {}
        kwargs["in"] = [pb_in]
        kwargs["out"] = [pb_out]
        self._insert_filter(**kwargs)

        r = self._send(url=self.URL_BASE,
                       data=EVENT,
                       method=Method.post)
        self.assertEqual(r.status_code, HTTP.OK.value)
        #print(r)

        #print("Waiting for the event to be handle by the engines")
        #sleep(5)

        #res = self.cm.get_entities_by_id(RESOURCE_ID)[0]
        #print(res)
