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

__version__ = '0.1'

from socket import gethostname

from imp import load_source

from sys import path, argv, exit

from os import listdir
from os.path import expanduser, splitext, join

from canopsis.cli.ccmd import Cmd

modules_path = expanduser('~/opt/ccli/libexec')
path.append(modules_path)


class Application(object):

    class Cli(Cmd):
        def __init__(self, app):
            super(Application.Cli, self).__init__('{0}/'.format(app.prompt))

            self.app = app

    def __init__(self, *args, **kwargs):
        super(Application, self).__init__(*args, **kwargs)

        self.hostname = gethostname()
        self.user = 'root'

        self.prompt = '{0}@{1}:'.format(self.user, self.hostname)

        self.modules = {}

        for module in listdir(modules_path):
            if module[0] != '.':
                modname, ext = splitext(module)

                if ext == '.py':
                    abspath = join(modules_path, module)
                    self.modules[modname] = load_source(modname, abspath)

        for modname in self.modules:
            module = self.modules[modname]

            def caller(s, line):
                module.start_cli(s.app.prompt)

            fname = 'do_{0}'.format(modname)
            caller.__name__ = fname
            Application.Cli.__dict__[fname] = caller

    def do_help(self):
        print("Usage: ccli [module]\n")
        print("Modules:")

        for modname in self.modules:
            print(" - {0}".format(modname))

    def __call__(self):
        print('Welcome to Canopsis CLI')

        try:
            if len(argv) == 2:
                if argv[1] == 'help':
                    self.do_help()

                else:
                    module = argv[1]

                    self.modules[module].start_cli(self.prompt)

            else:
                Application.Cli(self).cmdloop()

        except KeyboardInterrupt:
            print('Received KeyboardInterrupt, exiting...')
            exit(0)
