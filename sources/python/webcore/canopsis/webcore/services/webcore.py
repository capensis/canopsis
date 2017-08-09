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

import urllib
from operator import itemgetter

import flask
from canopsis.webcore.apps.flask.helpers import Resource

from canopsis.webcore.utils import gen_json


def inspect_routes(app):
    """
    Inspect all bottle routes.

    Thanks: https://buxty.com/b/2012/05/inspecting-your-routes-in-bottle/

    :param app: the bottle application
    """
    for route_1 in app.routes:
        if 'mountpoint' in route_1.config:
            prefix = route_1.config['mountpoint']['prefix']
            subapp = route_1.config['mountpoint']['target']

            for prefixes, route_2 in inspect_routes(subapp):
                yield [prefix] + prefixes, route_2
        else:
            yield [], route_1


class Methods:

    @staticmethod
    def get_routes(bottle_app, path=None):
        """
        List all routes in the webservice, according to a certain path.

        :param bottle_app: the bottle application
        :param str path: limit the listing to path including this value
        """
        themall = []
        for prefixes, route_3 in inspect_routes(bottle_app):
            if path is None or path in route_3.rule:
                route = {
                    'method': route_3.method,
                    'rule': route_3.rule
                }
                themall.append(route)

        ta = sorted(themall, key=itemgetter('rule'))
        ta = ['{method} -- {rule}'.format(**r) for r in ta]

        return gen_json(ta)

    @staticmethod
    def get_routes_v3(app, path=None):
        # http://flask.pocoo.org/snippets/117/
        routes = []
        for rule in app.url_map.iter_rules():

            options = {}
            for arg in rule.arguments:
                options[arg] = "[{0}]".format(arg)

            methods = urllib.unquote(','.join(rule.methods))
            url = urllib.unquote(flask.url_for(rule.endpoint, **options))

            route = {
                'endpoint': rule.endpoint,
                'methods': methods,
                'url': url
            }

            routes.append(route)

        return sorted(routes)


def exports(ws):

    bottle_app = ws.application  # keep bottle ref before beaker transformation

    @ws.application.get('/api/v2/rule/them/all/<path>')
    @ws.application.get('/api/v2/rule/them/all')
    def get_routes(path=None):
        return Methods.get_routes(bottle_app, path=path)


class APIWebcore(Resource):

    resource_routes = [
        '/api/v3/routes/all',
        '/api/v3/rule/them/all/',
        '/api/v3/rule/them/all/<string:path>'
    ]

    def get(self, path=None):
        return Methods.get_routes_v3(self._app, path=path)


def exports_v3(app, api):
    APIWebcore.init(app, api)
