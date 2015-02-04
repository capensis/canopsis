#!/usr/bin/env python
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
import json
import requests


webserver = 'http://127.0.0.1:8082'

# TODO All


class KnownValues(unittest.TestCase):

    def setUp(self):
        pass

    def login(self):
        pass

    def logout(self):
        pass

    def test_01_init(self):
        pass

    def test_02_test_url(self):
        self.login()
        urls_to_test = [
            '/'
        ]

        for uri in urls_to_test:
            r = requests.get(webserver + uri)
            self.assertEqual(r.status_code, 200)
        self.logout()

if __name__ == "__main__":
    unittest.main(verbosity=2)
