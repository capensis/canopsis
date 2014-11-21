#!/usr/bin/env python
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

import unittest
import json

from canopsis.old.webservices import Webservices

WS = None

get_auth_uri = [
    # account
    '/account/me',
    '/account/',
    '/account/account.root',
    # rest
    '/rest/object',
    '/rest/object/account',
    '/rest/object/account/account.root',
]


class KnownValues(unittest.TestCase):
    def setUp(self):
        self.view_id = 'root.view.mytest'
        self.view_name = 'mytest'

        self.directory_id = 'root.view.mydir'
        self.directory_name = 'mydir'

        self.parent_directory = 'directory.root.root'
        pass

    def test_01_Init(self):
        global WS
        WS = Webservices()

    def test_02_TestURL(self):
        for uri in get_auth_uri:
            success = False
            try:
                WS.get(uri)
                success = True
            except:
                pass

            if success:
                raise Exception("%s is open !" % uri)

    def test_03_Login(self):
        WS.login('root', 'root')

    def test_04_TestURLAuthed(self):
        WS.login('root', 'root')

        for uri in get_auth_uri:
            WS.get(uri)

    def test_05_Putting_View(self):
        WS.login('root', 'root')
        output = WS.put_view(
            self.view_id, self.view_name, self.parent_directory)
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_06_Putting_Directory(self):
        WS.login('root', 'root')
        output = WS.put_view(
            self.directory_id, self.directory_name, self.parent_directory,
            leaf=False)
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_07_renaming_view(self):
        WS.login('root', 'root')
        output = WS.rename_view(self.view_id, 'myNewName')
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_08_Changing_directory(self):
        WS.login('root', 'root')
        output = WS.change_view_parent(self.view_id, self.directory_id)
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_09_Moving_back_view(self):
        WS.login('root', 'root')
        output = WS.change_view_parent(self.view_id, self.parent_directory)
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_10_delete_view_or_dir(self):
        WS.login('root', 'root')
        output = WS.delete_view_or_dir([self.view_id, self.directory_id])
        try:
            output = json.loads(output)
        except:
            raise Exception('Invalid webserver output')
        WS.valid_server_response(output)

    def test_99_Logout(self):
        WS.logout()


if __name__ == "__main__":
    unittest.main(verbosity=2)
