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

from pkgutil import extend_path

from setuptools import setup as _setup, find_packages

from os.path import join, dirname, expanduser, abspath, basename

from sys import path, argv

from canopsis.common.utils import resolve_element

import canopsis

# TODO: set values in a dedicated configuration file
AUTHOR = 'Capensis'
AUTHOR_EMAIL = 'canopsis@capensis.fr'
LICENSE = 'AGPL V3'
ZIP_SAFE = False
URL = 'http://www.canopsis.org'
KEYWORDS = ' Canopsis Hypervision Hypervisor Monitoring'


def setup(description, keywords, **kwargs):

    # get setup path which corresponds to first python argument
    filename = argv[0]

    _path = dirname(abspath(expanduser(filename)))
    name = basename(_path)

    # add path to python path
    path.append(_path)

    # extend canopsis path with new sub modules and packages
    canopsis.__path__ = extend_path(canopsis.__path__, canopsis.__name__)

    # get package
    package = resolve_element("canopsis.{0}".format(name))

    # set default parameters if not setted
    kwargs.setdefault('name', package.__name__)
    kwargs.setdefault('version', package.__version__)
    kwargs.setdefault('author', AUTHOR)
    kwargs.setdefault('author_email', AUTHOR_EMAIL)
    kwargs.setdefault('license', LICENSE)
    kwargs.setdefault('zip_safe', ZIP_SAFE)
    kwargs.setdefault('url', URL)
    kwargs.setdefault('package_dir', {'': _path})

    kwargs.setdefault('keywords', kwargs.get('keywords', '') + KEYWORDS)

    if 'packages' not in kwargs:
        packages = find_packages(where=_path, exclude=['test'])
        kwargs['packages'] = packages

    if 'long_description' not in kwargs:
        with open(join(_path, 'README')) as f:
            kwargs['long_description'] = f.read()

    if 'test_suite' not in kwargs:
        kwargs['test_suite'] = 'test'

    _setup(**kwargs)
