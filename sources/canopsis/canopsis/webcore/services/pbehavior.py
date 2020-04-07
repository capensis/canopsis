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
Webservice for pbehaviors.
"""

from __future__ import unicode_literals

from bottle import request
from json import loads
from six import string_types

from canopsis.common.converters import id_filter
from canopsis.common.ws import route
from canopsis.confng.helpers import cfg_to_bool
from canopsis.pbehavior.manager import PBehaviorManager, PBehavior
from canopsis.pbehavior.utils import check_valid_rrule
from canopsis.watcher.manager import Watcher as WatcherManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
import time


VALID_PBEHAVIOR_PARAMS = [
    '_id', 'name', 'filter_', 'author', 'tstart', 'tstop', 'rrule',
    'enabled', 'comments', 'connector', 'connector_name', 'type_', 'reason',
    'timezone', 'exdate'
]


def check(data, key, type_):
    """
    If key exists in data and data[key] != None, check for type_

    :param dict data: data to check
    :param str key: key in data dict
    :param type type_: data[key] type to check
    :raises ValueError: type_ doesn't match value's type
    """
    if key in data and data[key] is not None:
        if not isinstance(data[key], type_):
            raise ValueError("The {0} must be a {1} not {2}".format(
                key, type_, type(data[key])))


def check_values(data):
    """
    Check if the values present in data respect the specification. If
    the values are correct do nothing. If not, raises an error.

    :param dict data: the data.
    :raises ValueError: a value is invalid.
    """

    # check str values
    for k in ["_id", "name", "author", "rrule", "component", "connector",
              "connector_name", 'type_', 'reason']:
        check(data, k, string_types)

    # check int values
    for k in ["tstart", "tstop"]:
        check(data, k, int)

        if 'tstart' in data and 'tstop' in data:
            if data['tstart'] >= data['tstop'] and \
               data['tstart'] is not None and data['tstop'] is not None:
                raise ValueError('tstop cannot be inferior or equal to tstart')

    # check dict values
    for k in ["comments"]:

        if "comments" not in data:
            continue

        if data["comments"] is None:
            continue

        check(data, k, list)

        for elt in data["comments"]:
            if not isinstance(elt, dict):
                raise ValueError("The list {0} store only {1} not {2}"
                                 .format(k, dict, type(elt)))

    if "filter" in data and isinstance(data["filter"], string_types):
        try:
            data["filter"] = loads(data["filter"])
        except ValueError:
            raise ValueError("Cant decode mfilter parameter: {}"
                             .format(data["filter"]))

    if 'rrule' in data:
        check_valid_rrule(data['rrule'])

    if PBehavior.EXDATE in data:
        if isinstance(data[PBehavior.EXDATE], list):
            for date in data[PBehavior.EXDATE]:
                if not isinstance(date, int):
                    raise ValueError("The date inside exdate must be an int.")
        else:
            raise ValueError("Exdate must be a list.")
    # useful when enabled doesn't exist in document
    if ("enabled" not in data
            or data["enabled"] is None
            or isinstance(data['enabled'], bool)):
        return

    data["enabled"] = cfg_to_bool(data["enabled"])


def create_params(_id, name=None, filter_=None, tstart=None, tstop=None,
                  rrule=None, enabled=None, comments=None, connector=None,
                  connector_name=None, author=None, type_=None, reason=None,
                  timezone=None, exdate=None):
    if exdate is None:
        exdate = []

    params = {
        PBehavior.NAME: name,
        PBehavior.FILTER: filter_,
        PBehavior.AUTHOR: author,
        PBehavior.TSTART: tstart,
        PBehavior.TSTOP: tstop,
        PBehavior.RRULE: rrule,
        PBehavior.ENABLED: enabled,
        PBehavior.COMMENTS: comments,
        PBehavior.CONNECTOR: connector,
        PBehavior.CONNECTOR_NAME: connector_name,
        PBehavior.TYPE: type_,
        PBehavior.REASON: reason,
        PBehavior.TIMEZONE: timezone,
        PBehavior.EXDATE: exdate
    }
    return params


class RouteHandlerPBehavior(object):
    """
    Passthrough class with few checks from the route to the pbehavior
    manager.
    """

    def __init__(self, pb_manager, watcher_manager):
        """
        :param PBehaviorManager pb_manager: pbehavior manager
        :param WatcherManager watcher_manager: watcher manager
        """
        self.watcher_manager = watcher_manager
        self.pb_manager = pb_manager

    def create(self, name, filter_, author,
               tstart, tstop, rrule=None,
               enabled=True, comments=None,
               connector='canopsis', connector_name='canopsis',
               type_=PBehavior.DEFAULT_TYPE, reason='', timezone=None,
               exdate=None, _id=None, replace_expired=False):
        """
        Create a pbehavior.

        :param str name: pb name
        :param str filter_: mongo filter applying on context graph, as string
        :param str author: pb author
        :param int tstart: start timestamp
        :param int tstop: end timestamp
        :param str rrule: RRULE
        :param bool enabled: enable/disable this pb
        :param list comments: list of comments {'author': 'author', 'message': 'msg'}
        :param str connector: connector
        :param str connector_name: connector name
        :param str type_: an associated type_
        :param str reason: a reason to apply this behavior
        :param str _id: the pb id (optional)
        :param bool replace_expired: rename current _id to EXP={_id} if exists or not
        """
        if exdate is None:
            exdate = []

        if comments is None:
            comments = []

        data = {
            PBehavior.NAME: name,
            PBehavior.FILTER: filter_,
            PBehavior.AUTHOR: author,
            PBehavior.TSTART: tstart,
            PBehavior.TSTOP: tstop,
            PBehavior.RRULE: rrule,
            PBehavior.ENABLED: enabled,
            PBehavior.COMMENTS: comments,
            PBehavior.CONNECTOR: connector,
            PBehavior.CONNECTOR_NAME: connector_name,
            PBehavior.TYPE: type_,
            PBehavior.REASON: reason,
            PBehavior.TIMEZONE: timezone,
            PBehavior.EXDATE: exdate
        }

        if _id is not None:
            data[PBehavior.ID] = _id

        check_values(data)

        result = self.pb_manager.create(
            pbh_id=_id,
            name=name,
            filter=filter_,
            author=author,
            tstart=tstart,
            tstop=tstop,
            rrule=rrule,
            enabled=enabled,
            comments=comments,
            connector=connector,
            connector_name=connector_name,
            type_=type_,
            reason=reason,
            timezone=timezone,
            exdate=exdate,
            replace_expired=replace_expired
        )

        return result

    def get_by_eid(self, eid):
        return self.pb_manager.get_pbehaviors_by_eid(eid)

    def read(self, _id, search=None, limit=None, skip=None, current_active_pbh=False):
        """
        Read a pbehavior.

        :param str _id: pb id
        :return: pbehavior
        :rtype: dict
        """
        is_ok = _id is None
        if isinstance(_id, string_types):
            is_ok = True
        elif isinstance(_id, list):
            is_ok = True
            for elt in _id:
                if not isinstance(elt, string_types):
                    is_ok = False

        if not is_ok:
            raise ValueError("_id should be str, a list, None (null) not {}"
                             .format(type(_id)))
        pbehaviors = self.pb_manager.read(_id, search, limit, skip)
        if current_active_pbh is True:
            return self._get_active_only(pbehaviors)
        return pbehaviors

    def _get_active_only(self, pbehaviors_data):
        active_ones = []
        now = int(time.time())
        for pb in pbehaviors_data.get("data", []):
            if self.pb_manager.check_active_pbehavior(now, pb):
                active_ones.append(pb)
        pbehaviors_data["data"] = active_ones
        pbehaviors_data["total_count"] = len(active_ones)
        pbehaviors_data["count"] = len(active_ones)
        return pbehaviors_data

    def update(self, _id, **kwargs):
        """
        Update pbehavior fields. Fields to None will **not** be updated.

        :param str _id: pbehavior id
        """
        params = create_params(_id, **kwargs)
        check_values(params)
        return self.pb_manager.update(_id, **params)

    def update_v2(self, _id, **kwargs):
        """
        Update pbehavior fields. Fields to None will **not** be updated.

        :param str _id: pbehavior id
        """
        params = create_params(_id, **kwargs)
        check_values(params)
        return self.pb_manager.update_v2(_id, **params)

    def delete(self, _id):
        """
        Delete pbehavior.

        :param str _id: pbehavior id
        """
        return self.pb_manager.delete(_id)

    def create_comment(self, pb_id, author, message):
        """
        Create a new comment on a pbehavior.

        :param str pb_id: pbehavior id
        :param str author: author
        :param str message: message
        :returns: comment id
        :rtype: str
        """
        author = str(author)
        message = str(message)

        return self.pb_manager.create_pbehavior_comment(pb_id, author, message)

    def update_comment(self, pb_id, _id, author, message):
        """
        Update an existing comment.

        :param str pb_id: pbehavior id
        :param str _id: comment id
        :param str author: author
        :param str message: message
        """
        author = str(author)
        message = str(message)

        return self.pb_manager.update_pbehavior_comment(
            pb_id, _id, author=author, message=message
        )

    def delete_comment(self, pb_id, _id):
        """
        Delete an existing comment.

        :param str pb_id: pbehavior id
        :param str _id: comment id
        """
        return self.pb_manager.delete_pbehavior_comment(pb_id, _id)


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    pbm = PBehaviorManager(*PBehaviorManager.provide_default_basics())
    watcher_manager = WatcherManager()
    rhpb = RouteHandlerPBehavior(
        pb_manager=pbm, watcher_manager=watcher_manager
    )

    @route(
        ws.application.post,
        name='pbehavior/create',
        payload=[
            'name', 'filter', 'author',
            'tstart', 'tstop', 'rrule',
            'enabled', 'comments',
            'connector', 'connector_name',
            'type_', 'reason', 'timezone', 'exdate'
        ]
    )
    def create(
            name, filter, author,
            tstart, tstop, rrule=None,
            enabled=True, comments=None,
            connector='canopsis', connector_name='canopsis',
            type_=PBehavior.DEFAULT_TYPE, reason='', timezone=None,
            exdate=None
    ):
        """
        Create a pbehavior.
        """
        return rhpb.create(
            name, filter, author, tstart, tstop, rrule,
            enabled, comments, connector, connector_name, type_, reason,
            timezone, exdate
        )

    @ws.application.post('/api/v2/pbehavior')
    def create_v2():
        """
        Create a pbehavior.

        required keys: name str, filter dict, comments list of dict with
        author message, tstart int, tstop int, author str

        optionnal keys: rrule str, enabled bool, _id str

        :raises ValueError: invalid keys sent.
        """
        try:
            elements = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        if elements is None:
            return gen_json_error(
                {'description': 'nothing to insert'},
                HTTP_ERROR
            )

        invalid_keys = []

        # keep compatibility with APIv1
        if 'filter' in elements:
            elements['filter_'] = elements.pop('filter')

        for key in elements.keys():
            if key not in VALID_PBEHAVIOR_PARAMS:
                invalid_keys.append(key)
                elements.pop(key)
        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        replace_expired = False
        try:
            replace_expired = int(request.params['replace_expired']) == 1
        except:
            pass

        try:
            elements['replace_expired'] = replace_expired
            return rhpb.create(**elements)
        except TypeError:
            return gen_json_error(
                {'description': 'The fields name, filter, author, tstart, tstop are required.'},
                HTTP_ERROR
            )
        except ValueError as exc:
            return gen_json_error(
                {'description': '{}'.format(exc.message)},
                HTTP_ERROR
            )

    @ws.application.put('/api/v2/pbehavior/<pbehavior_id:id_filter>')
    def update_v2(pbehavior_id):
        """
        Update a pbehavior.

        :raises ValueError: invalid keys sent.
        """
        try:
            elements = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        if elements is None:
            return gen_json_error(
                {'description': 'nothing to update'},
                HTTP_ERROR
            )

        invalid_keys = []

        # keep compatibility with APIv1
        if 'filter' in elements:
            elements['filter_'] = elements.pop('filter')

        for key in elements.keys():
            if key not in VALID_PBEHAVIOR_PARAMS:
                invalid_keys.append(key)
                elements.pop(key)
        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        try:
            return rhpb.update_v2(pbehavior_id, **elements)
        except TypeError as te:
            return gen_json_error(
                {'description': str(
                    'The fields name, filter, author, tstart, tstop are required.')},
                HTTP_ERROR
            )
        except ValueError as exc:
            return gen_json_error(
                {'description': '{}'.format(exc.message)},
                HTTP_ERROR
            )

    @route(
        ws.application.get,
        name='pbehavior/read',
        payload=['_id', 'search', 'limit', 'skip', 'current_active_pbh']
    )
    def read(_id=None, search=None, limit=None, skip=None, current_active_pbh=False):
        """
        Get a pbehavior.
        """
        return rhpb.read(_id, search=search, limit=limit, skip=skip, current_active_pbh=current_active_pbh)

    @route(
        ws.application.put,
        name='pbehavior/update',
        payload=[
            '_id',
            'name', 'filter',
            'tstart', 'tstop', 'rrule',
            'enabled',
            'timezone', 'exdate'
        ]
    )
    def update(
            _id,
            name=None, filter=None,
            tstart=None, tstop=None, rrule=None,
            enabled=None, comments=None,
            connector=None, connector_name=None,
            author=None, type_=None, reason=None, timezone=None, exdate=None
    ):
        """
        Update a pbehavior.
        """
        return rhpb.update_v2(
            _id=_id,
            name=name,
            filter_=filter,
            tstart=tstart,
            tstop=tstop,
            rrule=rrule,
            enabled=enabled,
            comments=comments,
            connector=connector,
            connector_name=connector_name,
            author=author,
            type_=type_,
            reason=reason,
            timezone=timezone,
            exdate=exdate
        )

    @route(
        ws.application.delete,
        name='pbehavior/delete',
        payload=['_id']
    )
    def delete(_id):
        """/pbehavior/delete : delete the pbehaviour that match the _id

        :param _id: the pbehaviour id
        :returns: a dict with two field. "acknowledged" that True if the
        delete is a sucess. False, otherwise.
        :rtype: dict
        """
        return rhpb.delete(_id)

    @ws.application.delete('/api/v2/pbehavior/<pbehavior_id:id_filter>')
    def delete_v2(pbehavior_id):
        """Delete the pbehaviour that match the _id

        :param pbehavior_id: the pbehaviour id
        :return: a dict with two field. "acknowledged" that True if the
        delete is a sucess. False, otherwise.
        :rtype: dict
        """
        ws.logger.info('Delete pbehavior : {}'.format(pbehavior_id))

        return gen_json(rhpb.delete(pbehavior_id))

    @ws.application.get('/api/v2/pbehavior_byeid/<entity_id:id_filter>')
    def get_by_eid(entity_id):
        """
        Return pbehaviors that apply on entity entity_id.
        """
        return gen_json(rhpb.get_by_eid(entity_id))

    @route(
        ws.application.post,
        name='pbehavior/comment/create',
        payload=['pbehavior_id', 'author', 'message']
    )
    def create_comment(pbehavior_id, author, message):
        """/pbehavior/comment/create : create a comment on the given pbehavior.

        :param _id: the pbehavior id
        :param author: author name
        :param message: the message to store in the comment.
        :returns: In case of success, return the comment id. None otherwise.
        """
        return rhpb.create_comment(pbehavior_id, author, message)

    @route(
        ws.application.put,
        name='pbehavior/comment/update',
        payload=['pbehavior_id', '_id', 'author', 'message']
    )
    def update_comment(pbehavior_id, _id, author=None, message=None):
        """/pbehavior/comment/update : create a comment on the given pbehavior.

        :param pbehavior_id: the pbehavior id
        :param _id: the comment id
        :param author: author name
        :param message: the message to store in the comment.
        :returns: In case of success, return the updated comment. None otherwise.
        """
        return rhpb.update_comment(pbehavior_id, _id, author, message)

    @route(
        ws.application.delete,
        name='pbehavior/comment/delete',
        payload=['pbehavior_id', '_id']
    )
    def delete_comment(pbehavior_id, _id):
        """/pbehavior/comment/delete : delete a comment on the given pbehavior.

        :param pbehavior_id: the pbehavior id
        :param _id: the comment id
        :returns: a dict with two field. "acknowledged" that contains True if
        delete has successed. False, otherwise.
        :rtype: dict
        """
        return rhpb.delete_comment(pbehavior_id, _id)

    @ws.application.get(
        '/api/v2/compute-pbehaviors'
    )
    def compute_pbehaviors():
        """
        Force compute of all pbehaviors, once per 10s

        :rtype: bool
        """
        ws.logger.info('Force compute on all pbehaviors')
        pbm.compute_pbehaviors_filters()
        pbm.launch_update_watcher(watcher_manager)

        return gen_json(True)
