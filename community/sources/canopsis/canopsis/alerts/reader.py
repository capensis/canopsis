# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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
Alarm reader manager.

TODO: replace the storage class parameter with a collection (=> rewriting count())
"""

from __future__ import unicode_literals

from canopsis.logger import Logger
from canopsis.confng import Configuration, Ini
from canopsis.common.middleware import Middleware


from canopsis.common.collection import MongoCollection


class AlertsReader(object):
    """
    Alarm cycle managment.

    Used to retrieve events related to alarms in a TimedStorage.
    """

    LOG_PATH = 'var/log/alertsreader.log'
    ALERTS_CONF_PATH = 'etc/alerts/manager.conf'
    ALERTS_STORAGE_URI = 'mongodb-periodical-alarm://'

    def __init__(self, logger, config, storage):
        """
        :param logger: a logger object
        :param config: a confng instance
        :param storage: a storage instance
        :param pbehavior_manager: a pbehavior manager instance
        """
        self.logger = logger
        self.config = config
        self.alarm_storage = storage
        self.alarm_collection = MongoCollection(self.alarm_storage._backend)

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canospis.confng.simpleconf.Configuration,
                      canopsis.storage.core.Storage,
                      canopsis.pbehavior.manager.PBehaviorManager]
        """
        logger = Logger.get('alertsreader', cls.LOG_PATH)
        conf = Configuration.load(cls.ALERTS_CONF_PATH, Ini)
        alerts_storage = Middleware.get_middleware_by_uri(
            cls.ALERTS_STORAGE_URI)
        return (logger, conf, alerts_storage)
