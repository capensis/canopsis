#!/usr/bin/env python
#--------------------------------
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

from caccount import caccount
from cstorage import get_storage

logger = None

##set root account
root = caccount(user="root", group="root")
storage = get_storage(account=root, namespace='object')


def init():
    logger.info(" + Create index of 'perfdata2'")
    storage.get_backend('perfdata2').ensure_index([
        ('co', 1),
        ('re', 1),
        ('me', 1)
    ])
    storage.get_backend('perfdata2').ensure_index([
        ('re', 1),
        ('me', 1)
    ])
    storage.get_backend('perfdata2').ensure_index([
        ('me', 1)
    ])
    storage.get_backend('perfdata2').ensure_index([
        ('tg',    1)
    ])

    logger.info(" + Create index of 'events'")
    storage.get_backend('events').ensure_index([
        ('connector_name', 1),
        ('resource', 1),
        ('component', 1),
        ('state', 1),
        ('state_type', 1),
        ('event_type', 1)
    ])
    storage.get_backend('events').ensure_index([
        ('tags', 1),
        ('source_type', 1)
    ])
    storage.get_backend('events').ensure_index([
        ('component', 1),
        ('resource', 1),
        ('event_type', 1)
    ])
    storage.get_backend('events').ensure_index([
        ('resource', 1),
        ('event_type', 1)
    ])
    storage.get_backend('events').ensure_index([
        ('event_type', 1)
    ])

    logger.info(" + Create index of 'events_log'")
    storage.get_backend('events_log').ensure_index([
        ('connector_name',    1),
        ('resource',        1),
        ('component',        1),
        ('state',            1),
        ('event_type',        1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('component', 1),
        ('resource', 1),
        ('event_type', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('resource', 1),
        ('event_type', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('event_type', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('state_type', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('tags', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('referer', 1)
    ])
    storage.get_backend('events_log').ensure_index([
        ('timestamp', 1)
    ])

def update():
    init()
