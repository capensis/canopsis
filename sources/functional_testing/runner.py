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

import importlib
from os import listdir
from os.path import isfile, join
import re
import unittest

pyfile = re.compile('[^_].*\.py$')


def add_modules_from_folder(suite, folder):
    """
    Try to import modules from a folder and add them to the test suite.

    :param suite: the suite test
    :param str folder: the folder to analyse
    """
    files = [f for f in listdir(folder) if isfile(join(folder, f)) and pyfile.match(f)]

    folder = folder.replace("/", ".")

    for f in files:
        mod_name = '{}.{}'.format(folder, '.'.join(f.split('.')[:-1]))
        try:
            mod = importlib.import_module(mod_name)
        except ImportError as exc:
            print('Cannot import {} ({})'.format(mod_name, exc))
            continue

        suite.addTests(loader.loadTestsFromModule(mod))

loader = unittest.TestLoader()
suite = unittest.TestSuite()

add_modules_from_folder(suite=suite, folder='apis')
add_modules_from_folder(suite=suite, folder='apis/alerts')

runner = unittest.TextTestRunner(verbosity=3)
result = runner.run(suite)
