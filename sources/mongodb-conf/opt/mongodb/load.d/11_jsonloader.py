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

import os
import json

##set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')

json_path = os.path.expanduser('~/opt/mongodb/load.d')

"""
Upserts json documents files within json_<collection> folder of the load.d folder
Json files are upserted if they contains a 'loader_id' field
If json files contains 'no_update_document' field set to true, upsert is avoided
"""

def init():

    for filename in os.listdir(json_path):
        absolute_path = '{}/{}'.format(json_path, filename)
        if os.path.isdir(absolute_path) and filename.startswith('json_'):
            collection = filename.replace('json_', '')
            for json_filename in os.listdir(absolute_path):
                if json_filename.endswith('.json'):
                    absolute_json_path = '{}/{}'.format(absolute_path, json_filename)
                    try:
                        with open(absolute_json_path) as f:
                            json_data = json.loads(f.read())
                            if '_id' in json_data : 
                                print ('Malformated insert document. A json loaded document must not contain _id key')
                            else:
                                if 'loader_id' not in json_data:
                                    print (' + Loader_id key not exists in json {} file,\n' \
                                    ' It must be a uniq document id for your custom json documents.\n' \
                                    ' Cannot process database upsert'.format(json_filename))
                                else:
                                    if 'no_update_document' in json_data and json_data['no_update_document']:
                                        print ('Document is marked as no updatable, nothing is done for {}'.format(json_filename))
                                    else:
                                        storage.get_backend(collection).update({'loader_id': json_data['loader_id']}, json_data, upsert=True)
                                        print ('{} information upserted'.format(json_filename))

                    except Exception as e:
                        print ('Unable to load json file {} : {}'.format(absolute_json_path, e))


def update():
    init()
init()
