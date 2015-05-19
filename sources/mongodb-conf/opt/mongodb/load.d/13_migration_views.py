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

from canopsis.configuration.dbconfigurationmanager import DBConfiguration
from canopsis.organisation.rights import Rights

from time import time
import json
import sys
import os


dbconf = DBConfiguration()
rights = Rights()


def log(message):
    print(message)


def init():
    pass


def update():
    # Get views
    views = dbconf.find(query={'crecord_type': 'view'})
    backup = []

    for view in views:
        # Go through the view
        if 'containerwidget' in view:
            if 'items' in view['containerwidget']:
                find_weathers(view['containerwidget']['items'])

        log('Updating view: {0}'.format(view['_id']))

        try:
            dbconf.put(view['_id'], view)

        except Exception as err:
            log('Error updating view {0}: {1}'.format(view['_id'], err))
            backup.append(view)

        log('Ensure rights for view: {0}'.format(view['_id']))

        actions = rights.get_action([view['_id']])

        if not actions:
            rights.add(view['_id'], 'Access to view {0}'.format(
                view.get('title', view['_id'])
            ))

    if backup:
        bakdir = os.path.join(
            sys.prefix, 'var', 'cache', 'canopsis', 'migration'
        )

        filepath = os.path.join(bakdir, 'view.{0}.bak'.format(int(time())))

        os.makedirs(bakdir)

        with open(filepath, 'w') as f:
            json.dump(backup, f)

        log('Backup found at: {0}'.format(filepath))


def find_weathers(items):
    # Recursively find all weathers in the view
    for item in items:
        widget = item['widget']

        # Nested widget search
        if 'items' in widget:
            find_weathers(widget['items'])

        if widget['xtype'] == 'weather' and 'event_selection' in widget:
            log('Processing weather in widget: {0}'.format(widget['title']))

            widget['event_selection'] = transform_event_selection(
                widget['event_selection']
            )


def transform_event_selection(event_selection):
    # Business code, update weather content
    selection = []

    while event_selection:
        value = event_selection.pop()

        if isinstance(value, basestring):
            value = {
                'label': 'label',
                'rk': value
            }

        selection.append(value)

    return selection


#TODO REMOVE
update()
