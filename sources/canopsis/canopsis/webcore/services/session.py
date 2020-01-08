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
import json


from canopsis.common.middleware import Middleware
from canopsis.common.utils import singleton_per_scope
from canopsis.common.ws import route
from canopsis.session.manager import Session, SessionError
from .rights import get_manager as get_rights
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


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


def get_username():
    """Returns the username of the logged-in user, or ''."""
    try:
        session = request.environ.get('beaker.session', {})
        user = session.get('user', '')

        # The content of user depends on the authentication method. If the user
        # logged in with HTTP authentication, it contains the username. If they
        # logged in with the loggin form, it contains a dictionnary.
        if isinstance(user, basestring):
            return user
        return user.get('_id', '')
    except AttributeError:
        return ''


def get_creation_time():
    """Returns the _creation_time of the logged-in _creation_time, or ''."""
    try:
        session = request.environ.get('beaker.session', {})
        return session.get('_creation_time', '')
    except AttributeError:
        return 

def get_id_beaker_session():
    """
    Return the id_beaker_session for slecte user's session in Default_session mongoDB 
    """
    creation_time = str(int(get_creation_time()))
    username = str(get_username())
    id_beaker_session_string = username + '_' + creation_time
    id_beaker_session = base64.b64encode(id_beaker_session_string)
    return id_beaker_session


def get_info():
    return get_id_beaker_session(), get_username()


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

    @ws.application.post(
         '/keepalive'
        )
    def keepalive():
        """
        Maintain the current session.
        """
        try :
            data = json.loads(request.body.read())
            visible = data["visible"]
            paths = data["path"]
            id_beaker_session, username = get_info()
            time = session_manager.keep_alive(id_beaker_session,username,visible,paths)
            return gen_json({'description':"Session keepalive","time":time,"visible":visible,"paths":paths})


        except SessionError as e :
            return  gen_json_error({'description':e.value },HTTP_ERROR)

    @ws.application.get(
        '/sessionstart'
    )
    def sessionstart():
        """
        Start a new session.
        """
        try :
            id_beaker_session, username = get_info()
            session_manager.session_start(id_beaker_session,username)
            return gen_json({'description':"Session Start"})
        except SessionError as e :
            return  gen_json_error({'description':e.value},HTTP_ERROR)


    @ws.application.post(
        '/session_hide'
    )
    def  sessionhide():
        try :
            data = json.loads(request.body.read())
            paths = data["path"]
            id_beaker_session, username = get_info()
            session_manager.session_hide(id_beaker_session,username,paths)
        except SessionError as e :
            return  gen_json_error({'description':e.value},HTTP_ERROR)



    @ws.application.get(
            '/sessions'
        )
    def session():
        try :
            params = {}
            params_key = request.query.keys()
            for key in params_key :
                if key == "usernames[]" :
                    params[key] = request.query.getall(key)
                else :
                    params[key] = request.query.get(key)
            id_beaker_session, username = get_info()
            sessions = session_manager.sessions_req(id_beaker_session,params)
            return gen_json({'description':"Sessions", 'sessions':sessions})

        except SessionError as e :
            return  gen_json_error({'description':e.value},HTTP_ERROR)