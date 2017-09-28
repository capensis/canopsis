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

from cmd import Cmd

import sys

from datetime import datetime

from logging import INFO

from canopsis.old.storage import Storage


# Main object
class Cmd(Cmd):
    def __init__(self, prompt):
        super(Cmd, self).__init__()
        self.prompt = prompt + '> '

    def do_quit(self, line):
        return True

    def help_quit(self):
        print("Exit CLI")

    def help_help(self):
        print("Show help message (type help <topic>)")

    # shortcuts
    do_exit = do_quit
    help_exit = help_quit


class Browser(Cmd):
    def __init__(self, prompt, account, namespace='object', crecord_type=None):
        super(Browser, self).__init__(prompt)
        self.account = account
        self.namespace = namespace
        self.crecord_type = crecord_type
        self.storage = Storage(
            account, namespace=namespace, logging_level=INFO)

    def do_ls(self, crecord_type=None):
        if self.crecord_type:
            records = self.storage.find({'crecord_type': self.crecord_type})

        elif crecord_type:
            records = self.storage.find({'crecord_type': crecord_type})

        else:
            records = self.storage.find()

        self.print_records(records)

    def do_cat(self, _id):
        try:
            if _id != '*':
                record = self.storage.get(_id)
                record.cat()

        except Exception as err:
            print("Impossible to cat {0}: {1}".format(_id, err))

    def do_dump(self, _id):
        try:
            if _id != '*':
                record = self.storage.get(_id)
                record.cat(dump=True)

        except Exception as err:
            print("Impossible to dump {0}: {1}".format(_id, err))

    def do_rm(self, _id):
        try:
            self.storage.remove(_id)

        except Exception as err:
            print("Impossible to remove {0}: {1}".format(_id, err))

    def do_cd(self, path):
        if path == "..":
            return True

    def print_records(self, records):
        print("Total: {0}".format(len(records)))

        lines = []

        for record in records:
            line = []

            line.append(record.owner)
            line.append(record.group)

            line.append(str(sys.getsizeof(record)))

            date = datetime.fromtimestamp(record.write_time)
            line.append(str(date))

            line.append(record.type)

            line.append(str(record._id))

            line.append(str(record.name))

            # self.columnize(line, displaywidth=200)
            lines.append(line)

        # Quick and dirty ...

        max_ln = {}
        for line in lines:
            i = 0
            for word in line:
                try:
                    if len(word) > max_ln[i]:
                        max_ln[i] = len(word)
                except:
                    max_ln[i] = len(word)

                i += 1

        # new_lines = []
        for line in lines:
            i = 0
            new_line = ""
            for word in line:
                empty = ""
                nb = max_ln[i] - len(word)

                for s in range(nb + 2):
                    empty += " "

                new_line += word + empty
                i += 1

            print(new_line)
