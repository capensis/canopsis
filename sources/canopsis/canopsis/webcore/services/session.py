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
import base64
from bottle import request, abort
from canopsis.auth.check import check

from canopsis.common.middleware import Middleware
from canopsis.common.utils import singleton_per_scope
from canopsis.common.ws import route
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


def get():
    """
    get the beaker session. If a beaker session dis not exist, check if a basic
    auth is present in the request. If the authentication succeed, return a
    newly created beaker.session with all the info of the suer. If the
    authentication header is invalid return a beaker session with no user.
    """

    beaker_sess = request.environ.get('beaker.session', None)

    if "user" not in beaker_sess:
        # Authorization: Basic
        try:
            auth_header = request.headers["Authorization"]
        except KeyError:
            return beaker_sess
        auth_header = auth_header.replace("Basic ", "")
        try:
            auth_header = base64.b64decode(auth_header)
        except TypeError as exc:
            abort(400, "Authorization headers " + exc.message)
        credential = auth_header.split(":", 1)

        if len(credential) != 2:
            return beaker_sess

        username = credential[0]
        password = credential[1]

        user = get_user(username)

        user = check(mode="plain", user=user, password=password)

        if not user:
            abort(403, 'Forbidden')

        if not user.get('enable', False):
            abort(403, 'Account is disabled')

        beaker_sess["user"] = credential[0]
        beaker_sess["auth_on"] = True
        beaker_sess.save()

    return beaker_sess


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

    @ws.application.get('/sessionstart')
    def sessionstart():
        """
        Start a new session.
        """
        username = request.get('username', None)

        session_manager.session_start(username)
        return {}
