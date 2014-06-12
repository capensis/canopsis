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

"""
Dictionary of configurables by configuration file.
"""
_CONFIGURABLE_BY_CONF_FILES = dict()


def add_configurable(configurable):
    """
    Add input configurable to list of watchers.

    :param configurable: configurable to add
    :type configurable: cconfiguration.Configurable
    """

    for conf_file in configurable.conf_files:

        configurables = _CONFIGURABLE_BY_CONF_FILES.setdefault(
            conf_file, set())

        configurables.add(configurable)


def remove_configurable(configurable):
    """
    Remove configurable to list of watchers.

    :param configurable: configurable to remove
    :type configurable: cconfiguration.Configurable
    """

    for conf_file in configurable.conf_files:

        configurables = _CONFIGURABLE_BY_CONF_FILES.setdefault(
            conf_file, set())

        configurables.remove(configurable)

        if not configurables:

            del _CONFIGURABLE_BY_CONF_FILES[conf_file]


def on_update_conf_file(conf_file):
    """
    Apply configuration on all configurables which watch input
    conf_file.

    :param conf_file: configuration file to reconfigure.
    :type conf_file: str
    """

    configurables = _CONFIGURABLE_BY_CONF_FILES.get(
        conf_file)

    if configurables:

        for configurable in configurables:

            configurable.apply_configuration(
                conf_file=conf_file)
