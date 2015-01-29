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

from canopsis.cli.ccmd import Cmd

from pprint import PrettyPrinter

from pyperfstore2 import manager

from os import system


class Cli(Cmd):
    def __init__(self, prompt):
        super(Cli, self).__init__('{0}:redis'.format(prompt))

        self.printer = PrettyPrinter(indent=2)
        self.manager = manager()

        self.check_length = 250
        self.key_info = 'gte_values_in_redis'

    def _get_data(self, _id):
        data = self.manager.store.redis.lrange(_id, 0, -1)

        def clean_point(p):
            p[0] = int(p[0])
            p[1] = float(p[1])

            return p

        return [clean_point(p.split('|')) for p in data]

    def _get_redis_data(self):
        keys = self.manager.store.redis.keys('*')
        keys.remove('perfstore2:rotate:plan')

        for key in keys:
            self.manager.store.redis_pipe.llen(key)

        result = self.manager.store.redis_pipe.execute()
        data = {}

        for index, key in enumerate(keys):
            data[key] = {
                'data_in_redis': self._get_data(key),
                self.key_info: {
                    'test_value': self.check_length,
                    'is_gte': (result[index] >= self.check_length)
                }
            }

        return data

    def _get_meta(self, data):
        def db_query_chunk(data):
            keys = []

            for key in data:
                keys.append(key)

                if len(keys) == 1000:
                    yield keys
                    keys = []

            yield keys

        for keys_chunk in db_query_chunk(data):
            if keys_chunk:
                query = self.manager.store.collection.find({
                    '_id': {'$in': keys_chunk}
                })

                for result in query:
                    data[result['_id']]['meta_in_db'] = result

    def do_stats(self, line):
        """ Show statistics about perfdata stored in Redis """

        data = self._get_redis_data()
        self._get_meta(data)

        lte = gte = 0

        for key in data:
            if data[key][self.key_info]['is_gte']:
                gte += 1

            else:
                lte += 1

        print('Redis key count where:')
        print(' - Less than {0}: {1}'.format(self.check_length, lte))
        print(' - More than {0}: {1}'.format(self.check_length, gte))
        print(' - Total: {0}'.format(lte + gte))

    def do_show(self, line):
        """ Show Redis content """

        data = self._get_redis_data()
        self._get_meta(data)

        self.printer.pprint(data)

    def do_rotate(self, line):
        """ Run 'pyperfstore2 rotate' """

        system('pyperfstore2 rotate')

    def do_rotateall(self, line):
        """ Rotate absolutely every perfdata stored in Redis """

        data = self._get_redis_data()
        self._get_meta(data)

        for key in data:
            self.manager.rotate(key)

    def help_stats(self):
        print(self.do_stats.__doc__)

    def help_show(self):
        print(self.do_show.__doc__)

    def help_rotate(self):
        print(self.do_rotate.__doc__)

    def help_rotateall(self):
        print(self.do_rotateall.__doc__)


def start_cli(prompt):
    try:
        mycli = Cli(prompt)
        mycli.cmdloop()

    except Exception as err:
        print("Impossible to start module: {0}".format(err))
