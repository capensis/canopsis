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

import os
import json

# Set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')

json_path = os.path.expanduser('~/opt/mongodb/load.d')

"""
Upserts json documents files within json_<collection> folder of the load.d
folder Json files are upserted if they contains a 'loader_id' field If json
files contains 'no_update_document' field set to true, upsert is avoided
"""


def init():
    # Iterating over json documents to manage and create paths variables.
    for filename in os.listdir(json_path):
        absolute_path = '{}/{}'.format(json_path, filename)
        if os.path.isdir(absolute_path) and filename.startswith('json_'):
            collection = filename.replace('json_', '')
            for json_filename in os.listdir(absolute_path):
                if json_filename.endswith('.json'):
                    absolute_json_path = '{}/{}'.format(
                        absolute_path,
                        json_filename
                    )
                    try:
                        with open(absolute_json_path) as f:
                            json_data = json.loads(f.read())

                            if type(json_data) == list:
                                for json_document in json_data:
                                    load_document(
                                        json_document,
                                        collection,
                                        json_filename
                                    )
                            elif type(json_data) == dict:
                                load_document(
                                    json_data,
                                    collection,
                                    json_filename
                                )
                    except Exception as e:
                        print ('Unable to load json file {} : {}'.format(
                            absolute_json_path,
                            e
                        ))


def load_document(json_data, collection, json_filename):

    if 'loader_id' not in json_data:

        print (
            ' + Loader_id key not exists in json {} file,\n' +
            ' It must be a uniq document id for your custom documents.\n' +
            ' Cannot process database upsert'.format(json_filename)
        )

    else:
        exists_document = storage.get_backend(collection).find({
            'loader_id': json_data['loader_id']
        }).count()

        if not exists_document:
            do_update(json_data, collection)
            print ('Document not existing, process insert {}'.format(
                json_filename
            ))
        else:
            if ('loader_no_update' not in json_data or
                    not json_data['loader_no_update']):
                print ('Document exists, process update {}'.format(
                    json_filename
                ))
                do_update(json_data, collection)
            else:
                print ('Document is marked as no updatable,' +
                       ' nothing is done for {}'.format(json_filename))


def do_update(json_data, collection):
    record = Record({}).dump()
    for key in json_data:
        record[key] = json_data[key]
    storage.get_backend(collection).update(
        {'loader_id': json_data['loader_id']},
        record,
        upsert=True
    )


def update():
    init()
