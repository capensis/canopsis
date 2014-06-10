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
import os
import sys
import getpass

projects = (
    '.',
    'cconfiguration',
    'ctimeserie',
    'cstorage',
    'ccontext',
    'cperfdata',
    'ctopology',
    'cmongo')

def run_cmd(huser="canopsis", cmd="install"):
    """
    Run setup cmd on all projects.
    """
    # find __file__ directory
    path = os.path.dirname(os.path.abspath(__file__))

    for project in projects:
        command = "python {0}/{1}/setup.py {2}".format(path, project, cmd)
        # run setup command
        run(command)

