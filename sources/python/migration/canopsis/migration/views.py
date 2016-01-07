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

from canopsis.migration.manager import MigrationModule
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.middleware.core import Middleware
from canopsis.organisation.rights import Rights

from time import time
import json
import sys
import os


CONF_PATH = 'migration/views.conf'
CATEGORY = 'VIEWS'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class ViewsModule(MigrationModule):
    def __init__(self, *args, **kwargs):
        super(ViewsModule, self).__init__(*args, **kwargs)

        self.storage = Middleware.get_middleware_by_uri(
            'storage-default://',
            table='object'
        )
        self.rights = Rights()

    def init(self):
        pass

    def update(self):
        # Get views
        views = self.storage.find_elements(query={'crecord_type': 'view'})
        backup = []

        for view in views:
            # Go through the view
            if 'containerwidget' in view:
                if 'items' in view['containerwidget']:
                    self.find_weathers(view['containerwidget']['items'])

            self.logger.info('Updating view: {0}'.format(view['_id']))

            try:
                self.storage.put_element(element=view, _id=view['_id'])

            except Exception as err:
                self.logger.error(
                    'Unable to update view "{0}": {1}'.format(
                        view['_id'],
                        err
                    )
                )

                backup.append(view)

            self.logger.info('Ensure rights for view: {0}'.format(view['_id']))

            actions = self.rights.get_action([view['_id']])

            if not actions:
                self.rights.add(view['_id'], 'Access to view {0}'.format(
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

            self.logger.info('Backup found at: {0}'.format(filepath))

    def find_weathers(self, items):
        # Recursively find all weathers in the view
        for item in items:
            w = item['widget']

            # Nested widget search
            if 'items' in w:
                self.find_weathers(w['items'])

            if w['xtype'] == 'weather' and 'event_selection' in w:
                self.logger.info(
                    u'Processing weather in widget: {0}'.format(
                        w['title']
                    )
                )

                w['event_selection'] = self.transform_event_selection(
                    w['event_selection']
                )

    def transform_event_selection(self, event_selection):
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
