#!/usr/bin/env python2
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

import dateutil
import unittest
import re
import json
import requests
from ast import literal_eval
from time import sleep, time
from urllib import urlencode
from pymongo import MongoClient
from canopsis.pbhevior.manager import PBehaviorManager

PBH_RESOURCE = "pbh_resource"
PBH_CONNECTOR = "pbh_connector",
PBH_CON_NAME = PBH_CONNECTOR + "_name"
PBH_COMPONENT= "pbh_component",

RESOURCE_ID = "{0}/{1}".format(PBH_RESOURCE, PBH_COMPONENT)

WEB_HOST = "localhost"
AUTH_KEY = "b71786a6-4c4c-11e7-8da1-0800279471b5"
MONGO_HOST = "localhost"

URL_BASE = "http://{0}:8082".format(WEB_HOST)
URL_AUTH = "{0}/?authkey={1}"
URL_PBEH = "{0}/pbehavior/create?{1}".format(URL_BASE, None)
URL_MONGO = 'mongodb://cpsmongo:canopsis@{0}:27017/canopsis'.format(MONGO_HOST)

ENTITIES_COL = "default_entities"
OBJECT_COL = "object"

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
    "mfilter": "{\"$or\":[{\"connector\":{\"$eq\":\"{0}\"}}]}".format(
        PBH_CONNECTOR),
    "crecord_creation_time": 1497367809,
    "crecord_name": None,
    "description": "I am a nice filter"
    #TODO add in and out pbehaviors
}

DEL_FILTER = {
    "_id": "delete_me",
    "name": "Hi",
    "actions": [{
        "field": "output",
        "type": "override",
        "value": "yep_it's_me"
    }],
    "crecord_type": "filter",
    "mfilter": "{\"$or\":[{\"connector\":{\"$eq\":\"i_am_a_connector\"}}]}",
    "description": "I am a nice filter"
}

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

    def _remove_filter(self):
        # self.obj_col.delete_many(DEL_FILTER)
        pass

    def _authenticate(self):
        session = requests.Session()
        response = session.get(URL_AUTH.format(URL_BASE, AUTH_KEY))
        print("Attempting login on {0}".format(
            URL_AUTH.format(URL_BASE, AUTH_KEY)))

        if re.search("<title>Canopsis | Login</title>", response.text)\
           is not None:
            self.fail("Authentication error.")

        return session

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

    def _create_rrule(self):
        pass

    def setUp(self):
        client = MongoClient(URL_MONGO)
        self.ent_col = client.canopsis[ENTITIES_COL]
        self.obj_col = client.canopsis[OBJECT_COL]

        self.session = self._authenticate()

        self.pbehavior_ids = []

        self.pbm = PBehaviorManager()

    def tearDown(self):
        self._remove_filter()
        # TODO remove entity created in the context


if __name__ == "__main__":
    unittest.main()
