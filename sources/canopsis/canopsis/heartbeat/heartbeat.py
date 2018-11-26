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


# __expected_interval_pattern = re.compile(r'^[0-9]*(s|m|h)$')
#
#
# def check_expected_interval(expected_interval):
#     return bool(__expected_interval_pattern.match(expected_interval))


class HeartbeatError(Exception):
    """
    Base Heartbeat error.
    """


class HeartbeatPatternExists(HeartbeatError):
    """
    Heartbeat pattern exists error.
    """


class HeartbeatManager(object):
    """
    Heartbeat service manager abstraction.
    """

    COLLECTION = 'heartbeat'
    __ID_PREFIX = 'heartbeat_'

    def __init__(self, collection, heartbeat_model):
        """

        :param `~.common.collection.MongoCollection` collection: object.
        :param `~.models.heartbeat.HeartBeat` heartbeat_model: object.
        """
        self.__collection = MongoCollection(collection)
        self.__model = heartbeat_model

    def insert_heartbeat_document(self, pattern, expected_interval):
        """
        Add a new Heartbeat.

        :param `dict` pattern: mapping with a string keys
                      and with a string values.
        :param `str` expected_interval: expected Event interval
                     that matches `__expected_interval_pattern` regex pattern.

        :returns: a new Heartbeat ID.
        :rtype: `str`.

        :raises: (`ValueError`, `.HeartbeatPatternExists`,
                  `~.common.collection.CollectionError`, ).
        """
        if not self.__model.check_heartbeat_pattern(pattern):
            raise ValueError('Not valid param: "pattern"')
        if not self.__model.check_expected_interval(expected_interval):
            raise ValueError('Not valid param: "expected_interval"')

        heartbeat_id = self.__ID_PREFIX + get_pattern_hash(pattern)

        if self.find_heartbeat_document(heartbeat_id):
            raise HeartbeatPatternExists()

        return self.__collection.insert({
            "_id": heartbeat_id,
            "pattern": pattern,
            "expected_interval": expected_interval
        })[0]

    def find_heartbeat_document(self, heartbeat_id):
        return self.__collection.find_one({"_id": heartbeat_id})

    def remove_heartbeat_document(self, heartbeat_id):
        """
        Remove Heartbeat by ID.

        :param `str` heartbeat_id: Heartbeat ID.
        :return:
        :raises: (`~.common.collection.CollectionError`, ).
        """
        return self.__collection.remove({"_id": heartbeat_id})

    def list_heartbeat_collection(self):
        return self.__collection.find({})
