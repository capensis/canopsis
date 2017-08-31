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

"""
Webservice for session managment.
"""

from __future__ import unicode_literals
from bottle import request

from canopsis.common.utils import singleton_per_scope
from canopsis.common.ws import route
from canopsis.middleware.core import Middleware
from canopsis.session.manager import Session
from .rights import get_manager as get_rights


def get_user(_id=None):
    """
    Find current user, or return None
    """
    session = request.environ.get('beaker.session')
    user = session.get('user', {})

    if not _id:
        _id = user.get('_id', None)

    if not _id:
        return None

    rights = get_rights()
    user = rights.get_user(_id)

    if user:
        user['rights'] = rights.get_user_rights(_id)

    return user


def create(user):
    """
    Create a user session.
    """
    session = request.environ.get('beaker.session')
    session['user'] = user
    session['auth_on'] = True
    session.save()

    return session


def delete():
    """
    Delete user user session.
    """
    session = request.environ.get('beaker.session')
    session.delete()


def exports(ws):
    """
    Expose session routes.
    """

    kwargs = {
        'collection': Middleware.get_middleware_by_uri(
            Session.SESSION_STORAGE_URI
        )._backend
    }
    session_manager = singleton_per_scope(Session, kwargs=kwargs)

    @route(ws.application.get, name='account/me', adapt=False)
    def get_me():
        """
        Return the user account.
        """
        user = get_user()
        user.pop('id', None)
        user.pop('eid', None)

        return user

    @route(ws.application.get, payload=['username'])
    def keepalive(username):
        """
        Maintain the current session.
        """
        session_manager.keep_alive(username)

    @route(ws.application.get, payload=['username'])
    def sessionstart(username):
        """
        Start a new session.
        """
        session_manager.session_start(username)
