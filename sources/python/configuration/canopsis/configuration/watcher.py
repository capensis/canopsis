#!/usr/bin/env python
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

try:
    from threading import Timer
except ImportError:
    from dummythreading import Timer

from stat import ST_MTIME

from os import stat
from os.path import exists, expanduser

from .parameters import Configuration, Parameter, Category
from .configurable import Configurable

"""
Dictionary of (mtime, configurables) by configuration file.
"""
_CONFIGURABLES_BY_CONF_FILES = {}
_MTIME_BY_CONF_FILES = {}


def add_configurable(configurable):
    """
    Add input configurable to list of watchers.

    :param configurable: configurable to add
    :type configurable: canopsis.configuration.Configurable
    """

    for conf_file in configurable.conf_files:

        conf_file = expanduser(conf_file)

        configurables = _CONFIGURABLES_BY_CONF_FILES.setdefault(
            conf_file, set())

        configurables.add(configurable)


def remove_configurable(configurable):
    """
    Remove configurable to list of watchers.

    :param configurable: configurable to remove
    :type configurable: canopsis.configuration.Configurable
    """

    for conf_file in configurable.conf_files:

        conf_file = expanduser(conf_file)

        configurables = _CONFIGURABLES_BY_CONF_FILES.get(conf_file)

        if configurables:
            configurables.remove(configurable)

            if not configurables:
                del _CONFIGURABLES_BY_CONF_FILES[conf_file]


def on_update_conf_file(conf_file):
    """
    Apply configuration on all configurables which watch input
    conf_file.

    :param conf_file: configuration file to reconfigure.
    :type conf_file: str
    """

    conf_file = expanduser(conf_file)

    configurables = _CONFIGURABLES_BY_CONF_FILES.get(conf_file)

    if configurables:

        for configurable in configurables:

            configurable.apply_configuration(conf_file=conf_file)

# default value for sleeping_time
DEFAULT_SLEEPING_TIME = 5


class Watcher(Configurable):
    """
    Watches all sleeping_time
    """

    CONF_FILE = 'configuration/watcher.conf'

    CATEGORY = 'WATCHER'
    SLEEPING_TIME = 'sleeping_time'

    DEFAULT_CONFIGURATION = Configuration(
        Category(CATEGORY,
            Parameter(SLEEPING_TIME, value=DEFAULT_SLEEPING_TIME, parser=int)))

    def __init__(self, sleeping_time=DEFAULT_SLEEPING_TIME, *args, **kwargs):

        super(Watcher, self).__init__(*args, **kwargs)

        self._sleeping_time = sleeping_time
        self._timer = None

    @property
    def sleeping_time(self):
        return self._sleeping_time

    @sleeping_time.setter
    def sleeping_time(self, value):
        """
        Change value of sleeping_time
        """
        self._sleeping_time = value
        # restart the timer
        self.stop()
        self.start()

    def _get_conf_files(self, *args, **kwargs):

        result = super(Watcher, self)._get_conf_files(*args, **kwargs)

        result.append(Watcher.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Watcher, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Watcher.CATEGORY,
            new_content=(
                Parameter(Watcher.SLEEPING_TIME, parser=int)))

        return result

    def run(self):
        global _CONFIGURABLES_BY_CONF_FILES

        for conf_file in _CONFIGURABLES_BY_CONF_FILES:

            # check file exists
            if exists(conf_file):
                # get mtime
                mtime = stat(conf_file)[ST_MTIME]
                # compare with old mtime
                old_mtime = _MTIME_BY_CONF_FILES.setdefault(conf_file, mtime)
                if old_mtime != mtime:
                    on_update_conf_file(conf_file)

    def _run(self):
        self._timer = Timer(self.sleeping_time, self._run)
        self._timer.start()
        self.run()

    def start(self):
        self._timer = Timer(self.sleeping_time, self._run)
        self._timer.start()

    def stop(self):
        if self._timer is not None:
            self._timer.cancel()


# To execute once the watcher will be ready
#_WATCHER = Watcher()
#_WATCHER.apply_configuration()


def start_watch():
    """
    Start default watcher
    """

    _WATCHER.start()


def stop_watch():
    """
    Stop default watcher
    """

    _WATCHER.Stop()


def change_sleeping_time(sleeping_time):
    """
    Change of sleeping_time for default watcher
    """

    _WATCHER.sleeping_time = sleeping_time
