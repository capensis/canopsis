#!/usr/bin/env python2.7
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

"""
    Mock for bottle module.
"""


class MockSession(dict):
    def save(self):
        pass


class MockRequest(object):
    def __init__(self, *args, **kwargs):
        super(MockRequest, self).__init__(*args, **kwargs)

        self.params = {}
        self.environ = {}


class MockResponse(object):
    def __init__(self, *args, **kwargs):
        super(MockResponse, self).__init__(*args, **kwargs)

        self.status = 200
        self.headers = {}

    def set_header(self, header, val):
        self.headers[header] = val

    def header(self, name):
        return self.headers.get(name, None)


def mock_redirect(url):
    pass
