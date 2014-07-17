#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record

import hashlib

##set root account
root = Account(user="root", group="root")

logger = None


def init():
    storage = get_storage(account=root, namespace='object')

    curves = [
        {'line_color': 'B7CA79', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 75, 'zIndex': 0, 'area_color': None, 'label': 'Free', 'metric': 'free'},
        {'line_color': 'B1221C', 'dashStyle': 'Solid', 'invert': True, 'area_opacity': 50, 'zIndex': 0, 'area_color': None, 'label': 'Upload', 'metric': 'if_octets-tx'},
        {'line_color': 'ABC8E2', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 50, 'zIndex': 0, 'area_color': None, 'label': 'Download', 'metric': 'if_octets-rx'},
        {'line_color': 'f11f0d', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 0, 'area_color': None, 'label': 'Load longterm', 'metric': 'load-longterm'},
        {'line_color': 'e97b15', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 1, 'area_color': None, 'label': 'Load midterm', 'metric': 'load-midgterm'},
        {'line_color': 'f3d30b', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 2, 'area_color': None, 'label': 'Load shortterm', 'metric': 'load-shortterm'},
        {'line_color': 'e97b15', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 1, 'area_color': None, 'label': 'Load midterm', 'metric': 'load-midterm'},
        {'line_color': '795344', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 75, 'zIndex': 0, 'area_color': None, 'label': 'Used', 'metric': 'used'},
        {'line_color': 'f11f0d', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 0, 'area_color': None, 'label': 'Load longterm', 'metric': 'load15'},
        {'line_color': 'e97b15', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 1, 'area_color': None, 'label': 'Load midterm', 'metric': 'load5'},
        {'line_color': 'f3d30b', 'dashStyle': 'Solid', 'invert': False, 'area_opacity': 30, 'zIndex': 2, 'area_color': None, 'label': 'Load shortterm', 'metric': 'load1'},
        {'line_color': 'FF9300', 'dashStyle': 'Dash', 'invert': False, 'area_opacity': 75, 'zIndex': 10, 'area_color': None, 'label': 'Warning', 'metric': 'pl_warning'},
        {'line_color': 'FF0000', 'dashStyle': 'Dash', 'invert': False, 'area_opacity': 75, 'zIndex': 10, 'area_color': None, 'label': 'Critical', 'metric': 'pl_critical'},
        {'line_color': 'BDBDBD', 'dashStyle': 'Solid', 'metric': 'cps_state_3', 'label': 'Unknown', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'FF0000', 'dashStyle': 'Solid', 'metric': 'cps_state_2', 'label': 'Critical', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'FF9300', 'dashStyle': 'Solid', 'metric': 'cps_state_1', 'label': 'Warning', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B7CA79', 'dashStyle': 'Solid', 'metric': 'cps_state_0', 'label': 'Ok', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'BDBDBD', 'dashStyle': 'Solid', 'metric': 'cps_sel_state_3', 'label': 'Unknown', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B9121B', 'dashStyle': 'Solid', 'metric': 'cps_sel_state_2', 'label': 'Critical', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'FF9300', 'dashStyle': 'Solid', 'metric': 'cps_sel_state_1', 'label': 'Warning', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B7CA79', 'dashStyle': 'Solid', 'metric': 'cps_sel_state_0', 'label': 'Ok', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'BDBDBD', 'dashStyle': 'Solid', 'metric': 'cps_pct_by_state_3', 'label': 'Unknown', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B9121B', 'dashStyle': 'Solid', 'metric': 'cps_pct_by_state_2', 'label': 'Critical', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'FF9300', 'dashStyle': 'Solid', 'metric': 'cps_pct_by_state_1', 'label': 'Warning', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B7CA79', 'dashStyle': 'Solid', 'metric': 'cps_pct_by_state_0', 'label': 'Ok', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'BDBDBD', 'dashStyle': 'Solid', 'metric': 'cps_statechange_3', 'label': 'Unknown', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B9121B', 'dashStyle': 'Solid', 'metric': 'cps_statechange_2', 'label': 'Critical', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'FF9300', 'dashStyle': 'Solid', 'metric': 'cps_statechange_1', 'label': 'Warning', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B7CA79', 'dashStyle': 'Solid', 'metric': 'cps_statechange_0', 'label': 'Ok', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False},
        {'line_color': 'B9121B', 'dashStyle': 'Solid', 'metric': 'cps_statechange_nok', 'label': 'Not ok', 'zIndex': -30, 'area_opacity': 20, 'area_color': None, 'invert': False}
    ]

    for curve in curves:
        _id = hashlib.sha1(curve['metric']).hexdigest().upper()
        try:
            storage.get(_id)
        except:
            logger.info(" + Create curve '%s'" % curve['metric'])
            record = Record(data=curve, _id=_id, name=curve['metric'], _type='curve')
            record.chmod('g+w')
            record.chmod('o+r')
            record.chgrp('group.CPS_curve_admin')
            storage.put(record)


def update():
    init()
