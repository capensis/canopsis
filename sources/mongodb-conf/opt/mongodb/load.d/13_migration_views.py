#!/usr/bin/env python
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

from copy import deepcopy
from time import time
from json import dumps

from canopsis.old.storage import get_storage

storage = get_storage(
    namespace='object'
)


def log(message):
    print (message)


def init():
    pass


def update():

    objects = storage.get_backend('object')

    # Get views
    views = objects.find({'crecord_type': 'view'})
    backup = []

    for view in views:
        # Prepare backup
        backup.append(deepcopy(view))

        # Go through the view
        if 'containerwidget' in view:
            if view['loader_id'] == 'view.services':
                if 'items' in view['containerwidget']:
                    find_weathers(view['containerwidget']['items'])

        log('db update {}'.format(view['_id']))

        objects.update(
            {'_id': view['_id']},
            {'$set': {'containerwidget': view['containerwidget']}}
        )

    # Proceed backup
    filepath = '/tmp/canopsis_migration_view_backup_{}'.format(int(time()))

    with open(filepath, 'w') as f:
        f.write(dumps(backup, indent=2))


def find_weathers(items):

    # Recursively find all weathers in the view
    for item in items:
        widget = item['widget']

        # Nested widget search
        if 'items' in widget:
            find_weathers(widget['items'])

        if widget['xtype'] == 'weather' and 'event_selection' in widget:
            log(' processing weather'.format(widget['title']))
            widget['event_selection'] = transform_event_selection(
                widget['event_selection']
            )


def transform_event_selection(event_selection):

    # Business code, update weather content
    selection = []
    while event_selection:
        value = event_selection.pop()
        if isinstance(value, basestring):
            value = {'label': 'label', 'rk': value}
        selection.append(value)
    return selection


#TODO REMOVE
update()
