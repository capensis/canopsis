#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

from celerylibs import listing

BROKER_HOST             = "localhost"
BROKER_PORT             = 5672
BROKER_USER             = "guest"
BROKER_PASSWORD         = "guest"
BROKER_VHOST            = "canopsis"
CELERY_RESULT_BACKEND       = "amqp"
CELERY_IMPORTS          = listing.tasks('~/etc/tasks.d')

# informations here http://celery.github.com/celery/configuration.html#id1
CELERY_TASK_RESULT_EXPIRES  = 1800

CELERYD_LOG_LEVEL       = 'INFO'

CELERYD_TASK_TIME_LIMIT     = 1800
CELERYD_CONCURRENCY     = 5
