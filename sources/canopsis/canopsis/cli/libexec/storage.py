#!/usr/bin/env python
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

from canopsis.cli.ccmd import Cmd, Browser
from canopsis.old.account import Account

from os import system


class Cli(Cmd):
    def __init__(self, prompt):
        super(Cli, self).__init__("%sstorage" % prompt)
        self.myprompt = "%sstorage" % prompt

    def do_cd(self, namespace):
        Browser(
            "%s/%s" % (self.myprompt, namespace),
            Account(user="root", group="root"), namespace).cmdloop()

    def do_mongo(self, line):
        system('mongo canopsis')


def start_cli(prompt):
    try:
        mycli = Cli(prompt)
        mycli.cmdloop()
    except Exception as err:
        print("Impossible to start module: %s" % err)
