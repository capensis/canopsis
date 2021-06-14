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

import os

from os import walk
from os.path import join, dirname, exists

from setuptools import setup, find_packages
from setuptools.command.install import install as setup_install

PACKAGE = 'canopsis'
AUTHOR = 'Capensis'
AUTHOR_EMAIL = 'canopsis@capensis.fr'
LICENSE = 'AGPL V3'
ZIP_SAFE = False
URL = 'https://www.capensis.fr/canopsis/'
KEYWORDS = 'Canopsis Hypervision Hypervisor Monitoring'

VERSION = '0.1'

TEST_FOLDERS = ['tests', 'test']


class PostInstall(setup_install):

    def _makedirs(self, unix_path):
        os.environ['CPS_BOOTSTRAP'] = '1'

        from os.path import join as pjoin
        from canopsis.common import root_path

        abspath = '{}/{}'.format(root_path, unix_path)

        path = os.path.sep + pjoin(*abspath.split('/'))

        if os.path.isdir(path):
            return

        if os.path.isfile(path):
            os.unlink(path)

        os.makedirs(path)

    def run(self):
        self._makedirs('var/log/engines')
        self._makedirs('var/cache/canopsis')
        self._makedirs('tmp')
        self._makedirs('etc')

        setup_install.run(self)


def get_cmdclass():
    cmdclass = {
        'install': PostInstall
    }
    return cmdclass


def get_pkgpath():
    return dirname(os.path.realpath(__file__))


def get_scripts(pkgpath):
    """
    Get a list of scripts to install in canopsis env.
    """
    scripts_path = join(pkgpath, 'scripts')
    scripts = []
    for root, _, files in walk(scripts_path):
        for _file in files:
            scripts.append(join(root, _file))

    return scripts


def get_data_files(pkgpath):
    """
    Get a list of data files from the etc directory.
    """

    data_files = []

    """
    Populate data files here.

    DO NOT PUSH CONFIGURATION. EVER.

    You can bundle examples if you want.
    """

    return data_files


def get_install_requires(pkgpath):
    """
    Get a list of requirements from requirements.txt
    """
    reqs = []
    requires_path = join(pkgpath, 'requirements.txt')

    with open(requires_path) as f:
        # remove new lines, extra spaces...
        reqs = f.readlines()

    reqs = filter(None, [r.strip() for r in reqs])
    reqs = [r for r in reqs if not r.startswith('#')]

    return reqs


def get_test_suite(pkgpath):
    test_folder = None

    test_folders = \
        [folder for folder in TEST_FOLDERS if exists(join(pkgpath, folder))]

    if test_folders:
        for test_folder in test_folders:
            return test_folder

    return test_folder


def setup_canopsis(pkgpath):
    """
    :param no_conf:
    :param add_etc: add automatically etc files (default True)
    :type add_etc: bool
    :returns: dict of setup() arguments
    """

    setuptools_args = {}

    # set default parameters if not setted
    setuptools_args['name'] = PACKAGE
    setuptools_args['author'] = AUTHOR
    setuptools_args['author_email'] = AUTHOR_EMAIL
    setuptools_args['license'] = LICENSE
    setuptools_args['zip_safe'] = ZIP_SAFE
    setuptools_args['url'] = URL
    setuptools_args['packages'] = find_packages(exclude=['test.*'])
    setuptools_args['keywords'] = KEYWORDS
    setuptools_args['version'] = VERSION
    setuptools_args['install_requires'] = get_install_requires(pkgpath)
    setuptools_args['test_suite'] = get_test_suite(pkgpath)
    setuptools_args['data_files'] = get_data_files(pkgpath)
    setuptools_args['scripts'] = get_scripts(pkgpath)
    setuptools_args['cmdclass'] = get_cmdclass()

    return setuptools_args


if __name__ == '__main__':
    setuptools_args = setup_canopsis(get_pkgpath())
    setup(**setuptools_args)
