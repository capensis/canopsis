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

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record
from socket import getfqdn
import os
import json
import pprint

pp = pprint.PrettyPrinter(indent=2)

# Set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')

json_path = os.path.expanduser('~/opt/mongodb/load.d')

# Can enable debug trace
DEBUG = False

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


# Replace system should be externalized
def replace_value(replace_string):

    # Macro to replace in documents with the value as tuple (MACRO, REPLACE)
    replacements = [
        ('[[HOSTNAME]]', getfqdn()),
    ]

    for replacement in replacements:
        replace_string = replace_string.replace(
            replacement[0],
            replacement[1]
        )

    return replace_string


def set_list_value(replace_list):

    length = len(replace_list)

    for x in xrange(length):

        if isinstance(replace_list[x], basestring):
            replace_list[x] = replace_value(replace_list[x])

        if isinstance(replace_list[x], dict):
            set_dict_value(replace_list[x])

        if isinstance(replace_list[x], list):
            set_list_value(replace_list[x])


def set_dict_value(replace_dict):

    if isinstance(replace_dict, dict):

        for key in replace_dict:

            if isinstance(replace_dict[key], basestring):
                replace_dict[key] = replace_value(replace_dict[key])

            if isinstance(replace_dict[key], dict):
                set_dict_value(replace_dict[key])

            if isinstance(replace_dict[key], list):
                set_list_value(replace_dict[key])


def hooks(json_data):
    """
    This function is a hook on json document that are being updated
    into database and allow call transformation method that may change
    records
    """
    # replace [[HOSTNAME]] macro by the current host value in record
    set_dict_value(json_data)


def do_update(json_data, collection):

    record = Record({}).dump()

    for key in json_data:
        record[key] = json_data[key]

    compare_record = record.copy()

    hooks(record)

    if DEBUG and record != compare_record:
        print 'Differences found\n # before \n{}\n\n # after\n {}'.format(
            pp.pformat(compare_record),
            pp.pformat(record)
        )

    storage.get_backend(collection).update(
        {'loader_id': json_data['loader_id']},
        record,
        upsert=True
    )


def update():
    init()
