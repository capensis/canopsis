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

import unittest
import re
import json
import requests
from ast import literal_eval
from time import sleep
from urllib import urlencode
from pymongo import MongoClient

# TODO specify the entity use for the test and add it in the PBEHAVIOR
# and FILTER

WEB_HOST = "localhost"
AUTH_KEY = "b71786a6-4c4c-11e7-8da1-0800279471b5"
MONGO_HOST = "localhost"

URL_BASE = "http://{0}:8082".format(WEB_HOST)
URL_AUTH = "{0}/?authkey={1}".format(URL_BASE, AUTH_KEY)
URL_PBEH = "{0}/pbehavior/create?{1}".format(URL_BASE)
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
    "mfilter": "{\"$or\":[{\"connector\":{\"$eq\":\"i_am_a_connector\"}}]}",
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
    "filter": None,
    "author": "Functionnal_test",
    "tstart": None,
    "tstop": None,
    "rrule": None,
    "enabled": "true",
    "connector": None,
    "connector_name": None
}

BEAT = 60


class BaseTest(unittest.TestCase):
    def _insert_filter(self):
        self.obj_col.insert(FILTER)
        print("Waiting {0}s for the beat".format(BEAT))
        sleep(BEAT)

    def _remove_filter(self):
        self.obj_col.delete_one()

    def _authenticate(self):
        session = requests.Session()
        response = session.get(URL_AUTH)

        if re.search("<title>Canopsis | Login</title>", response.text)\
           is not None:
            self.fail("Authentication error.")

        return session

    def _create_pbehavior(self, pbehavior):
        response = self.session.post(URL_PBEH.format(urlencode(pbehavior)))

        response = literal_eval(response.text)

        if response["total"] == 1 and response["sucess"] is True:
            self.pbehavior_ids += response["data"]
        else:
            self.fail("Impossible to insert the pbehavior.")

    def setUp(self):
        client = MongoClient(URL_MONGO)
        self.ent_col = client.canopsis[ENTITIES_COL]
        self.obj_col = client.canopsis[OBJECT_COL]

        self._insert_filter()
        self.session = self._authenticate()

        self.pbehavior_ids = []

    def tearDown(self):
        self._remove_filter()
        # TODO remove entity created in the context


class InFullOutFull(BaseTest):

    pass
