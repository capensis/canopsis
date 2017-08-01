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

#from __future__ import unicode_literals

import importlib
from os import listdir
from os.path import isfile, join
import re
import unittest

loader = unittest.TestLoader()
suite = unittest.TestSuite()

path = 'apis'
pyfile = re.compile('[^_].*\.py$')
files = [f for f in listdir(path) if isfile(join(path, f)) and pyfile.match(f)]
for f in files:
    mod_name = 'apis.{}'.format('.'.join(f.split('.')[:-1]))
    try:
        mod = importlib.import_module(mod_name)
    except ImportError:
        print("Cannot import {}.".format(mod_name))
        continue

    suite.addTests(loader.loadTestsFromModule(mod))

runner = unittest.TextTestRunner(verbosity=3)
result = runner.run(suite)
