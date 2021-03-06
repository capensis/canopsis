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

from __future__ import unicode_literals

from bottle import request

from canopsis.common.converters import id_filter
from canopsis.watcher.manager import Watcher
from canopsis.watcher.links import build_all_links
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    watcher = Watcher()

    @ws.application.get(
        '/api/v2/watchers/<watcher_id:id_filter>'
    )
    def get_watcher(watcher_id):
        """
        Get this particular watcher.

        :param str watcher_id: Entity ID of the watcher
        :returns: <Watcher>
        """
        watcher_obj = watcher.get_watcher(watcher_id)
        if watcher_obj is None:
            return gen_json_error({'description': 'nothing to return'},
                                  HTTP_ERROR)

        return gen_json(watcher_obj)

    @ws.application.post('/api/v2/watchers')
    def create_watcher():
        """
        Create a new watcher.

        Required data:

        JSON as text:

        {
            "mfilter": "JSON-encoded filter that applies on default_entities",
            "display_name": "watcher_name",
            "_id": "uniq_id"
        }

        :returns: nothing
        """
        # element is a full Watcher (dict) to insert
        element = request.json

        if element is None:
            return gen_json_error({'description': 'nothing to insert'},
                                  HTTP_ERROR)

        try:
            watcher_create = watcher.create_watcher(body=element)
        except ValueError as ex:
            return gen_json_error({'description': 'value error: {}'.format(ex)},
                                  HTTP_ERROR)
        except KeyError as ex:
            return gen_json_error({'description': 'key error: {}'.format(ex)},
                                  HTTP_ERROR)
        except TypeError:
            return gen_json_error({'description': 'type error: mfilter '\
                                   'should be a string'},
                                  HTTP_ERROR)
        if watcher_create is None:
            return gen_json_error(
                {'description': 'can\'t decode mfilter'},
                HTTP_ERROR
            )

        return gen_json({})

    @ws.application.put(
        '/api/v2/watchers/<watcher_id:id_filter>'
    )
    def update_watcher(watcher_id):
        """
        Update an existing watcher.

        :param watcher_id: Entity ID of the watcher
        :type watcher_id: str
        :returns: nothing
        """
        dico = request.json

        if dico is None or not isinstance(dico, dict) or len(dico) <= 0:
            return gen_json_error({'description': 'wrong update dict'},
                                  HTTP_ERROR)

        watcher.update_watcher(watcher_id=watcher_id, updated_field=dico)

        return gen_json({})

    @ws.application.delete(
        '/api/v2/watchers/<watcher_id:id_filter>'
    )
    def delete_watcher(watcher_id):
        """
        Delete a watcher, based on his id.

        :param watcher_id: ID of the watcher
        :type watcher_id: str
        :returns: mongo result dict of the deletion
        """
        ws.logger.info('Delete watcher : {}'.format(watcher_id))

        try:
            deletion_dict = watcher.delete_watcher(watcher_id=watcher_id)
        except Exception:
            return gen_json_error({'description': 'cannot delete {}'
                                                  .format(watcher_id)},
                                  HTTP_ERROR)

        return gen_json(deletion_dict)

    @ws.application.get(
        '/api/v2/watchers/compute-watchers-links'
    )
    def compute_watchers_links():
        """
        Force compute of all watchers, once per 10s

        :rtype: bool
        """
        ws.logger.info('Force compute of watcher links')
        build_all_links(watcher.context_graph)

        return gen_json(True)
