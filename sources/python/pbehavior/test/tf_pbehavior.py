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
import argparse
from pymongo import MongoClient

WEB_HOST = "localhost"
AUTH_KEY = "b71786a6-4c4c-11e7-8da1-0800279471b5"
MONGO_HOST = "localhost"

URL_BASE = "http://{0}:8082/".format(WEB_HOST)
URL_AUTH = "{0}/?authkey={1}".format(URL_BASE, AUTH_KEY)
URL_MONGO = 'mongodb://cpsmongo:canopsis@{0}:27017/canopsis'.format(MONGO_HOST)

ENTITIES_COL = "default_entities"


class BaseTest(unittest.TestCase):

    def setUp(self):
        client = MongoClient(URL_MONGO)
        db = client.canopsis
        self.ent_col = db[ENTITIES_COL]
