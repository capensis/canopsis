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

import argparse
import importlib
from os import listdir
from os.path import isfile, join
import re
import sys
import unittest

pyfile = re.compile('[^_].*\.py$')


def import_module(suite, name):
    """
    Import a module in a suite from his name.

    :param TestSuite suite: suite object
    :param str name: module name (ex. apis.weather)
    """
    try:
        mod = importlib.import_module(name)
    except ImportError as exc:
        print('Cannot import {} ({})'.format(name, exc))
        return

    suite.addTests(loader.loadTestsFromModule(mod))


def add_modules_from_folder(suite, folder):
    """
    Try to import modules from a folder and add them to the test suite.

    :param suite: the suite test
    :param str folder: the folder to analyse
    """
    files = [f for f in listdir(folder) if isfile(join(folder, f)) and pyfile.match(f)]

    for f in files:
        mod_name = 'apis.{}'.format('.'.join(f.split('.')[:-1]))
        import_module(suite, mod_name)


if __name__ == "__main__":
    """
    To run only one test:

    python -m unittest -q apis.<module>[.<class>[.<test]]'
    """

    parser = argparse.ArgumentParser(prog='functionnal testing runner')
    parser.add_argument('module_names', nargs='*',
                        help='Namespace of modules to test (ex. apis.weather)')
    args = parser.parse_args(sys.argv[1:])

    loader = unittest.TestLoader()
    suite = unittest.TestSuite()

    if args.module_names == []:
        add_modules_from_folder(suite=suite, folder='apis')
    else:
        for mod_name in args.module_names:
            import_module(suite, mod_name)

    runner = unittest.TextTestRunner(verbosity=3)
    result = runner.run(suite)
