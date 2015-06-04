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

from fabric.api import run
from os.path import dirname, expanduser

projects = (
    'common',
    'configuration',
    'timeserie',
    'event',
    'check',
    'middleware',
    'rpc',
    'mom',
    'storage',
    'schema',
    'mongo',
    'kombu',
    'context',
    'perfdata',
    'old',
    'task',
    'graph',
    'topology',
    'engines',
    'connectors',
    'tools',
    'cli',
    'topology',
    'organisation',
    'auth',
    'snmp')


def setup(cmd="install", projects=projects):
    """
    Run setup cmd on all projects.
    """
    # find __file__ directory
    path = dirname(expanduser(__file__))

    cmd_path = "python {0}/{{0}}/setup.py {1}".format(path, cmd)

    for project in projects:
        # run setup command
        run(cmd_path.format(project))
