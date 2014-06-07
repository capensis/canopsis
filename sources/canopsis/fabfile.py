#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

from fabric.api import run
from fabric.context_managers import cd
from os.path import dirname, abspath, join


def run_cmd(cmd):
    """
    Run setup cmd on all projects.
    """

    # find __file__ directory
    path = dirname(abspath(__file__))

    cmd_path = "python {0}/{1}/setup.py {2}".format(path, '{0}', cmd)

    def run_cmd_path(sub_directory=''):
        """
        Call fabric run function on sub_directory setup.py.
        """

        # get absolute sub-path
        sub_path = join(path, sub_directory)
        # change directory
        with cd(sub_path):
            run('pwd')
            # run setup command
            run(cmd_path.format(sub_directory))

    # setup canopsis
    run_cmd_path()

    # setup cconfiguration
    run_cmd_path("cconfiguration")

    # setup timeserie
    run_cmd_path("ctimeserie")

    # setup cstorage
    run_cmd_path("cstorage")

    # setup ccontext
    run_cmd_path("ccontext")

    # setup cperfdata
    run_cmd_path("cperfdata")

    # setup ctopology
    run_cmd_path("ctopology")

    # setup cmongo
    run_cmd_path("cmongo")


if __name__ == '__main__':

    from sys import argv

    cmd = 'install'

    if len(argv) == 2:
        cmd = argv[1]

    run_cmd(cmd)
