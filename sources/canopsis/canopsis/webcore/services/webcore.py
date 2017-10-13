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
from operator import itemgetter

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


def exports(ws):

    bottle_app = ws.application  # keep bottle ref before beaker transformation

    @ws.application.get('/api/v2/rule/them/all/<path>')
    @ws.application.get('/api/v2/rule/them/all')
    def get_routes(path=None):
        """
        List all routes in the webservice, according to a certain path.

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
