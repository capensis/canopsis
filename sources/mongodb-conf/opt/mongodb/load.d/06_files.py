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
from canopsis.old.storage import Storage

##set root account
root = Account(user="root", group="root")

logger = None


def init():
    storage = Storage(account=root)

    namespaces = ['files', 'binaries.files', 'binaries.chunks']

    for namespace in namespaces:
        logger.info(" + Drop '%s' collection" % namespace)
        storage.drop_namespace(namespace)


def update():
    update_for_new_rights()


def update_for_new_rights():
    #update briefcase elements
    storage = Storage(namespace='files', account=root)

    dump = storage.find({})

    for record in dump:
        if record.owner.find('account.') == -1:
            record.owner = 'account.%s' % record.owner
        if record.group.find('group.') == -1:
            record.group = 'group.%s' % record.group

    storage.put(dump)
