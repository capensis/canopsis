# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

import itertools
import re

from hashlib import md5

from canopsis.common.collection import MongoCollection


def get_pattern_hash(pattern):
    """
    Hash the heartbeat pattern.

    :param `dict` pattern: heartbeat pattern.
    :return: heartbeat pattern hash.
    :rtype: `str`.
    """
    checksum = md5()
    for chunk in itertools.chain(*((k, pattern[k]) for k in sorted(pattern))):
        checksum.update(chunk)
    return checksum.hexdigest()


__expected_interval_pattern = re.compile(r'^[0-9]*(s|m|h)$')


def check_expected_interval(expected_interval):
    return bool(__expected_interval_pattern.match(expected_interval))


class HeartbeatManager(object):
    """
    Heartbeat service manager abstraction.
    """

    COLLECTION = 'heartbeat'
    __ID_PREFIX = 'heartbeat_'

    def __init__(self, collection):
        self.__collection = MongoCollection(collection)

    def __build_heartbeat_id(self, pattern):
        return self.__ID_PREFIX + get_pattern_hash(pattern)

    def create_heartbeat(self, pattern, expected_interval):
        if not check_expected_interval(expected_interval):
            raise ValueError('Not valid param: "expected_interval"')
        heartbeat_id = self.__build_heartbeat_id(pattern)
        return self.__collection.insert({
            "_id": heartbeat_id,
            "pattern": pattern,
            "expected_interval": expected_interval
        })[0]
