#!/usr/bin/env python
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

"""
Session manager definition.
"""

from __future__ import unicode_literals
from uuid import uuid4
from time import time

from canopsis.confng import Configuration, Ini

DEFAULT_ALIVE_SESSION_DURATION = 300


class SessionError(Exception):
    """
    Exception for Session object
    """

    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)


class Session(object):
    """
    Manage session informations.
    """

    CONF_PATH = 'etc/session/session.conf'

    SESSION_STORAGE_URI = 'mongodb-default-session://'

    SESSION_COLLECTION = 'session_collection'

    def __init__(
            self,
            collection,
    ):
        """
        :param MongoCursor collection: the collection where user sessoins are located
        """

        self.session_collection = collection

        self.config = Configuration.load(self.CONF_PATH, Ini)
        session = self.config.get('SESSION', {})

        self.alive_session_duration = int(
            session.get(
                'alive_session_duration',
                DEFAULT_ALIVE_SESSION_DURATION))

    def keep_alive(self, id_beaker_session, username, visible, path):
        """
        Keep session alive by setting the ``last_check`` field
        to current timestamp.

        :type visible:boolean
        :type visible: string
        :returns: last check timestamp
        :rtype: timestamp
        """

        if not self.is_session_active(id_beaker_session):
            raise SessionError("Session Not Valid")

        now = int(time())
        session_current = self.session_collection.find_one(
            {'id_beaker_session': id_beaker_session, 'last_ping': {"$gt": now - self.alive_session_duration}})
        if not visible:
            self.session_collection.update({'_id': session_current['_id']}, {
                                           '$set': {'last_ping': now}})
            return now

        else:
            view_id = path[0]
            if len(path) == 2:
                tab_id = path[1]
            else:
                tab_id = None
            tab_duration = session_current['tab_duration']
            if path is session_current['last_visible_path']:
                tab_duration[view_id][tab_id] += now - \
                    session_current['last_visible_ping']
            else:
                if isinstance(tab_duration, dict) and view_id in tab_duration:
                    if isinstance(
                            tab_duration[view_id],
                            dict) and tab_id in tab_duration[view_id]:
                        tab_duration[view_id][tab_id] += now - \
                            session_current['last_ping']
                    elif view_id and isinstance(tab_duration[view_id], int):
                        tab_duration[view_id] += now - \
                            session_current['last_ping']
                    else:
                        tab_duration[view_id][tab_id] = now - \
                            session_current['last_ping']

                elif len(path) != 2:
                    tab_duration[view_id] = now - session_current['last_ping']
                else:
                    tab_duration[view_id] = {
                        tab_id: now - session_current['last_ping']}
            session_current['visible_duration'] += now - \
                session_current['last_ping']
            session_current['last_visible_ping'] = now
            session_current['last_visible_path'] = path
            session_current['last_ping'] = now

        self.session_collection.update({'_id': session_current['_id']}, {
                                       '$set': session_current})

        return now

    def is_session_active(self, id_beaker_session):
        """
        Check if session is active.
        If the session isn't found, then it is considered inactive.

        :param username: user identifier
        :type username: string
        :returns: check if session is active or not
        :rtype: bool
        """
        now = int(time())
        session = self.session_collection.find_one({'id_beaker_session': id_beaker_session, 'last_ping': {
                                                   "$gt": now - self.alive_session_duration}})

        if session is None:
            return False

        return True

    def session_start(self, id_beaker_session, username):
        """
        Make session active for a user.

        :param username: user identifier
        :type username: string
        :returns: when the has started
        :rtype: timestamp or None
        """

        if self.is_session_active(id_beaker_session) is False:
            now = int(time())
            element = {
                'id_beaker_session': id_beaker_session,
                'username': username,
                'start': now,
                'last_ping': now,
                'last_visible_ping': now,
                'last_visible_path': 'None',
                'visible_duration': 0,
                'tab_duration': {}

            }

            self.session_collection.update({'_id': str(uuid4())},
                                           element,
                                           upsert=True)

            return now

        else:
            return None

    def session_hide(self, id_beaker_session, username, path):

        if not self.is_session_active(id_beaker_session):
            raise SessionError("Session Not Valid")

        now = int(time())
        session_current = self.session_collection.find_one(
            {'id_beaker_session': id_beaker_session, 'last_ping': {"$gt": now - self.alive_session_duration}})

        view_id = path[0]
        if len(path) == 2:
            tab_id = path[1]
        else:
            tab_id = None
        tab_duration = session_current['tab_duration']
        if path is session_current['last_visible_path']:
            tab_duration[view_id][tab_id] += now - \
                session_current['last_visible_ping']
        else:
            if isinstance(tab_duration, dict) and view_id in tab_duration:
                if isinstance(
                        tab_duration[view_id],
                        dict) and tab_id in tab_duration[view_id]:
                    tab_duration[view_id][tab_id] += now - \
                        session_current['last_ping']
                elif view_id and isinstance(tab_duration[view_id], int):
                    tab_duration[view_id] += now - \
                        session_current['last_ping']
                else:
                    tab_duration[view_id][tab_id] = now - \
                        session_current['last_ping']

            elif len(path) != 2:
                tab_duration[view_id] = now - session_current['last_ping']
            else:
                tab_duration[view_id] = {
                    tab_id: now - session_current['last_ping']}
        session_current['visible_duration'] += now - \
            session_current['last_ping']
        session_current['last_visible_ping'] = now
        session_current['last_visible_path'] = path
        session_current['last_ping'] = now

        self.session_collection.update({'_id': session_current['_id']}, {
                                       '$set': session_current})
        return now

    def sessions_req(self, id_beaker_session, params):
        if not self.is_session_active(id_beaker_session):
            raise SessionError("Session Not Valid")

        now = int(time())
        req = {}
        if "active" in params:
            if params["active"] == "true":

                req['last_ping'] = {"$gt": now - self.alive_session_duration}
            else:
                req['last_ping'] = {"$lt": now - self.alive_session_duration}

        if "usernames[]" in params:
            names = []
            for name in params["usernames[]"]:
                names.append(name)
            req["username"] = {"$in": names}

        if "started_after" in params:
            req["start"] = {"$gt": params["started_after"]}

        if "stopped_before" in params:
            req["last_ping"] = {"$lt": params["stopped_before"]}

        sessions = list(self.session_collection.find(req))
        return sessions
