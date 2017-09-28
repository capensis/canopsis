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

from os import walk, getenv
from os.path import join, dirname, expanduser, abspath, exists

from sys import path, argv
from sys import prefix as sys_prefix

from setuptools import setup as _setup

PACKAGE = 'canopsis'
AUTHOR = 'Capensis'
AUTHOR_EMAIL = 'canopsis@capensis.fr'
LICENSE = 'AGPL V3'
ZIP_SAFE = False
URL = 'http://www.canopsis.org'
KEYWORDS = 'Canopsis Hypervision Hypervisor Monitoring'

VERSION = '0.1'

TEST_FOLDERS = ['tests', 'test']


def get_pkgpath():
    # get setup path which corresponds to first python argument
    filename = argv[0]

    _path = dirname(abspath(expanduser(filename)))

    return _path

def find_scripts(pkgpath):
    scripts_path = join(pkgpath, 'scripts')
    scripts = []
    for root, _, files in walk(scripts_path):
        for _file in files:
            scripts.append(join(root, _file))

    return scripts

def get_data_files(pkgpath):
    etc_path = join(pkgpath, 'etc')
    data_files = []

    target = getenv('CPS_PREFIX', join(sys_prefix, 'etc'))

    for root, _, files in walk(etc_path):
        files_to_copy = [join(root, _file) for _file in files]
        final_target = join(target, root[len(etc_path) + 1:])
        data_files.append((final_target, files_to_copy))

    return data_files

def get_install_requires(pkgpath):
    requirements = []
    requires_path = join(pkgpath, 'requirements.txt')

    with open(requires_path) as f:
        # remove new lines, extra spaces...
        requirements = [r.strip() for r in f.readlines()]

    return requirements

def get_description(pkgpath):
    readme_path = join(pkgpath, 'README.md')
    description = None

    with open(readme_path) as f:
        description = f.read()

    return description

def get_test_suite(pkgpath):
    test_folder = None

    test_folders = \
        [folder for folder in TEST_FOLDERS if exists(join(pkgpath, folder))]

    if test_folders:
        for test_folder in test_folders:
            return test_folder

    return test_folder

def setup(add_etc=True):
    """
    :param add_etc: add automatically etc files (default True)
    :type add_etc: bool
    :param setuptools_args: enrich setuptools.setup method
    """

    pkgpath = get_pkgpath()
    data_files = []
    setuptools_args = {}

    # add path to python path
    path.append(pkgpath)

    if '--no-conf' not in argv:
        if add_etc:
            data_files = get_data_files(pkgpath)
    else:
        argv.remove('--no-conf')

    # set default parameters if not setted
    setuptools_args['name'] = PACKAGE
    setuptools_args['author'] = AUTHOR
    setuptools_args['author_email'] = AUTHOR_EMAIL
    setuptools_args['license'] = LICENSE
    setuptools_args['zip_safe'] = ZIP_SAFE
    setuptools_args['url'] = URL
    setuptools_args['packages'] = ['canopsis']
    setuptools_args['package_dir'] = {'': pkgpath}
    setuptools_args['keywords'] = KEYWORDS
    setuptools_args['version'] = VERSION
    setuptools_args['install_requires'] = get_install_requires(pkgpath)
    setuptools_args['long_description'] = get_description(pkgpath)
    setuptools_args['test_suite'] = get_test_suite(pkgpath)
    setuptools_args['data_files'] = data_files

    _setup(**setuptools_args)


setup()
