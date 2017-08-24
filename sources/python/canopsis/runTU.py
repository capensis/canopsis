#!/usr/bin/env python2.7
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

from os.path import isdir, join
from sys import executable
from subprocess import Popen, PIPE
from argparse import ArgumentParser
from glob import glob


def test(project):
    """Run a test"""
    print('**********************************')
    print('running tests for {0}'.format(project))
    print('**********************************')

    p = Popen([executable, join(project, 'setup.py'), 'test'], stdout=PIPE)
    return p.communicate()[0]


if __name__ == '__main__':
    parser = ArgumentParser(description='Run tests on canolibs')
    parser.add_argument(
        'projects',
        metavar='projects',
        type=str,
        nargs='*',
        help='projects to run tests on'
    )

    args = parser.parse_args()


    if not args.projects: # we run tests for all projects if none is supplied
        for element in glob('*'):
            if isdir(element):
                test(element)
                raw_input('Press a key to continue...')

    else:
        for project in args.projects:
            test(project)
