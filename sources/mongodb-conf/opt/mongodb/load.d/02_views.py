#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

import os
import sys
import json

logger = None

views_path = os.path.expanduser('~/opt/mongodb/load.d/views')

# Set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')


def init():

    for path, folders, files in os.walk(views_path):
        for filename in files:
            filepath = os.path.join(path, filename)

            with open(filepath) as f:
                try:
                    data = json.loads(f.read())
                except Exception, e:
                    print (
                        '\n + Error while loading JSON file {} :' +
                        ' {}, aborting...\n'.format(filepath, e))
                    sys.exit(1)

                try:
                    _id = data.pop('id')
                except KeyError as err:
                    print >>sys.stderr, "Can't parse view, missing key:", err
                    sys.exit(1)

                create_view(_id, filename, data)


def update():
    pass


def create_view(_id, name, data, mod='o+r', autorm=True, internal=False):
    # Delete old view

    try:
        record = storage.get('view.%s' % _id)
        if autorm:
            storage.remove(record)
        else:
            return record
    except:
        pass

    logger.info(" + Create view '%s'" % name)
    record = Record(data, _type='view', name=name, group='group.CPS_view')
    record.chmod(mod)
    storage.put(record)
    return record


def update_view_for_new_metric_format():
    records = storage.find(
        {'crecord_type': 'view'},
        namespace='object',
        account=root)
    for view in records:
        container = view.data['container']

        for item in container['items']:
            widget = item['widget']
            nodesObject = {}

            # check with flotchart migration from highchart
            if 'xtype' in item['data']:
                xtype = item['data']['xtype']

                if xtype == 'line_graph':
                    xtype = 'timegraph'
                elif xtype == 'bar_graph':
                    xtype = 'timegraph'
                    item['data']['SeriesType'] = 'bars'
                elif xtype == 'diagram':
                    xtype = 'category_graph'

                item['data']['xtype'] = xtype

            # check if old format
            if 'nodes' in item['data']:
                itemNodes = item['data']['nodes']

                if isinstance(itemNodes, list):
                    itemXtype = item['data']['xtype']

                    if itemXtype == 'weather':
                        print('Ignore for weather widget')
                        break

                    # update for text widget
                    if itemXtype == 'text' or itemXtype == 'topology_viewer':
                        print('Update widget text/topology_viewer format')
                        item['data']['inventory'] = item['data']['nodes']
                        del item['data']['nodes']
                        break

                    for node in itemNodes:
                        try:
                            nodesObject[node['id']] = node

                            # write extra_fields in node root
                            if 'extra_field' in node:
                                nodesObject[node['id']].update(
                                    node['extra_field']
                                )

                                # build ccustom in view
                                del node['extra_field']
                        except Exception as error:
                            print(
                                'An error occured for the ' +
                                'following widget: {}'.format(error))
                            print(item)

                    item['data']['nodes'] = nodesObject
                    print(item['data']['nodes'])

                # check between commits
                if 'ccustom' in item['data']:
                    if isinstance(item['data']['ccustom'], dict):
                        ccustoms = item['data']['ccustom'].iteritems()
                        for nodeId, customValue in ccustoms:
                            if nodeId in itemNodes:
                                itemNodes[nodeId].update(customValue)
                        del item['data']['ccustom']

    storage.put(records)
