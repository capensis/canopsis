#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from logging import getLogger, DEBUG
from json import loads

from bottle import put, request, HTTPError

## Canopsis
from canopsis.old.storage import get_storage
from canopsis.old.record import Record

#import protection function
from libexec.auth import get_account

logger = getLogger("rights")


@put('/rights/:namespace/:crecord_id')
def change_rights(namespace, crecord_id=None):
    account = get_account()
    storage = get_storage(
        namespace=namespace, account=account, logging_level=DEBUG)

    #get put data
    aaa_owner = request.params.get('aaa_owner', default=None)
    aaa_group = request.params.get('aaa_group', default=None)
    aaa_access_owner = request.params.get('aaa_access_owner', default=None)
    aaa_access_group = request.params.get('aaa_access_group', default=None)
    aaa_access_other = request.params.get('aaa_access_other', default=None)

    if crecord_id is not None:
        record = storage.get(crecord_id, account=account)

    if isinstance(record, Record):
        logger.debug('record found, changing rights/owner')
        #change owner and group
        if aaa_owner is not None:
            record.chown(aaa_owner)
        if aaa_group is not None:
            record.chgrp(aaa_group)

        #change rights
        if aaa_access_owner is not None:
            record.access_owner = loads(aaa_access_owner)
        if aaa_access_group is not None:
            record.access_group = loads(aaa_access_group)
        if aaa_access_other is not None:
            record.access_other = loads(aaa_access_other)

        try:
            storage.put(record, account=account)
        except:
            logger.error('Access denied')
            return HTTPError(403, "Access denied")

    else:
        logger.warning('The record doesn\'t exist')
