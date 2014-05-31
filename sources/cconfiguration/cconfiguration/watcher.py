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

try:
    from pyinotify import\
        ProcessEvent, ALL_EVENTS, WatchManager, ThreadedNotifier

except ImportError:
    # implement our own ProcessEvent and MASKS
    pass

from os.path import expanduser

_CONFIGURATION_PATH = expanduser('~/etc/')
_CONFIGURABLES_BY_CONFIGURATION_FILE = dict()

_WATCH_MANAGER = WatchManager()


class _EventHandler(ProcessEvent):

    def process_default(self, event):

        configuration_file = event.pathname

        configurables = _CONFIGURABLES_BY_CONFIGURATION_FILE.get(
            configuration_file)

        for configurable in configurables:
            configurable.apply_configuration(
                configuration_files=[expanduser(configuration_file)])


_NOTIFIER = ThreadedNotifier(_WATCH_MANAGER, _EventHandler())

_WATCH_MANAGER.add_watch(_CONFIGURATION_PATH, ALL_EVENTS, rec=True)

_NOTIFIER.start()


def add_watch(configurable):

    for configuration_file in configurable.configuration_files:

        configurables = _CONFIGURABLES_BY_CONFIGURATION_FILE.setdefault(
            configuration_file, set())

        configurables.add(configurable)
