#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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


##set root account
root = Account(user="root", group="root")

logger = None


def init():
    storage = get_storage(account=root, namespace='object')

    curves = [
        {
            'crecord_name': 'simplepoints',
            'lines': False,
            'areas': False,
            'points': True,
            'bars': False,
            'line_width': 1,
            'bar_width': 10,
            'line_style': 'line',
            'point_shape': 'circle',
            'area_opacity': 1
        },
        {
            'crecord_name': 'simpleline',
            'lines': True,
            'areas': False,
            'points': False,
            'bars': False,
            'line_width': 1,
            'bar_width': 10,
            'line_style': 'line',
            'point_shape': 'circle',
            'area_opacity': 1
        },
        {
            'crecord_name': 'simplebar',
            'lines': False,
            'areas': False,
            'points': False,
            'bars': True,
            'line_width': 1,
            'bar_width': 10,
            'line_style': 'line',
            'point_shape': 'circle',
            'area_opacity': 1
        },
        {
            'crecord_name': 'simplearea',
            'lines': True,
            'areas': True,
            'points': False,
            'bars': False,
            'line_width': 1,
            'bar_width': 10,
            'line_style': 'line',
            'point_shape': 'circle',
            'area_opacity': 1
        }
    ]

    for curve in curves:
        _id = 'curve.{0}'.format(curve['crecord_name'])

        try:
            storage.get(_id)

        except KeyError:
            logger.info(" + Create curve '%s'" % curve['crecord_name'])
            record = Record(
                data=curve,
                _id=_id,
                name=curve['crecord_name'],
                _type='curve'
            )

            record.chmod('g+w')
            record.chmod('o+r')
            record.chgrp('group.CPS_curve_admin')
            storage.put(record)


def update():
    init()
