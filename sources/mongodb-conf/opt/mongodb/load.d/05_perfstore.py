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

import logging
import pyperfstore2

logger = None

manager = pyperfstore2.manager(logging_level=logging.DEBUG)


def init():
    logger.info(" + Drop 'pyperfstore2' collections")
    manager.store.drop()
    add_rotate_task()


def update():
    add_rotate_task()

    # tweak: rename all 'stat' metrics
    metrics = manager.find(mfilter={"co": "stat", "me": {"$regex": "^cps_.*"}})

    for metric in metrics:
        logger.info(
            " + Move '%s.%s.%s' to '%s.%s.%s'" % (
                metric["co"], metric["re"], metric["me"], metric["re"], "selector",
                metric["me"]))

        old_id = metric["_id"]
        del metric["_id"]

        metric["co"] = metric["re"]
        metric["re"] = "selector"
        name = "%s%s%s" % (metric["co"], metric["re"], metric["me"])
        _id = manager.gen_id(name)

        manager.store.create(_id, metric)
        manager.store.remove(_id=old_id)


def add_rotate_task():
    #### TODO: Remove this !!

    from canopsis.old.account import Account
    from canopsis.old.storage import get_storage

    account = Account(user="root", group="root")
    storage = get_storage(account=account, namespace='object')

    _id = 'schedule.pyperfstore_rotate'
    try:
        storage.remove(_id)
    except:
        pass
