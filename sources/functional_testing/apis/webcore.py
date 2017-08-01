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

from test_base import BaseApiTest, Method, HTTP


class TestWebcoreAPI(BaseApiTest):

    def setUp(self):
        self._authenticate()  # default setup

        self.base = '{}/{}'.format(self.URL_BASE, 'api/v2/rule/them/all')

    def test_RuleThemAll(self):
        r = self._send(url=self.base,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertTrue(len(json) > 0)

        # With a filter
        r = self._send(url=self.base + '/auth',
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
