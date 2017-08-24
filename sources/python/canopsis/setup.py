#!/usr/bin/env python
# -*- coding: utf-8 -*-

requirements = [
    'potsdb',
    'whisper'
]

from setuptools import setup as _setup, find_packages

from os import walk, getenv
from os.path import join, dirname, expanduser, abspath, basename, exists

import sys

from sys import path, argv
from sys import prefix as sys_prefix


# TODO: set values in a dedicated configuration file
AUTHOR = 'Capensis'
AUTHOR_EMAIL = 'canopsis@capensis.fr'
LICENSE = 'AGPL V3'
ZIP_SAFE = False
URL = 'http://www.canopsis.org'
KEYWORDS = ' Canopsis Hypervision Hypervisor Monitoring'

DEFAULT_VERSION = '0.1'

TEST_FOLDERS = ['tests', 'test']
VERSION = '2.4.6'

def setup(add_etc=True, **kwargs):
    """
    Setup dedicated to canolibs projects.

    :param add_etc: add automatically etc files (default True)
    :type add_etc: bool

    :param kwargs: enrich setuptools.setup method
    """

    # get setup path which corresponds to first python argument
    filename = argv[0]

    _path = dirname(abspath(expanduser(filename)))
    name = basename(_path)

    # add path to python path
    path.append(_path)

    # set default parameters if not setted
    kwargs.setdefault('name', 'canopsis')
    kwargs.setdefault('author', AUTHOR)
    kwargs.setdefault('author_email', AUTHOR_EMAIL)
    kwargs.setdefault('license', LICENSE)
    kwargs.setdefault('zip_safe', ZIP_SAFE)
    kwargs.setdefault('url', URL)
    kwargs.setdefault('package_dir', {'canopsis': 'canopsis'})

    kwargs.setdefault('keywords', kwargs.get('keywords', '') + KEYWORDS)

    # set version
    version = VERSION
    if version is not None:
        kwargs.setdefault('version', version)

    if '--no-conf' not in argv:
        # add etc content if exist and if --no-conf
        if add_etc:
            etc_path = join(_path, 'etc')

            if exists(etc_path):
                data_files = kwargs.get('data_files', [])
                target = getenv('CPS_PREFIX', join(sys_prefix, 'etc'))

                for root, dirs, files in walk(etc_path):
                    files_to_copy = [join(root, _file) for _file in files]
                    final_target = join(target, root[len(etc_path) + 1:])
                    data_files.append((final_target, files_to_copy))
                kwargs['data_files'] = data_files

    else:
        argv.remove('--no-conf')

    # add scripts if exist
    if 'scripts' not in kwargs:
        scripts_path = join(_path, 'scripts')
        if exists(scripts_path):
            scripts = []
            for root, dirs, files in walk(scripts_path):
                for _file in files:
                    scripts.append(join(root, _file))
            kwargs['scripts'] = scripts

    # add description
    if 'long_description' not in kwargs:
        readme_path = join(_path, 'README')
        if exists(readme_path):
            with open(join(_path, 'README')) as f:
                kwargs['long_description'] = f.read()

    # add test
    if 'test_suite' not in kwargs:
        test_folders = \
            [folder for folder in TEST_FOLDERS if exists(join(_path, folder))]
        if test_folders:
            for test_folder in test_folders:
                kwargs['test_suite'] = test_folder
                break

    kwargs['packages'] = find_packages()

    _setup(**kwargs)

if __name__ == '__main__':
    setup()
