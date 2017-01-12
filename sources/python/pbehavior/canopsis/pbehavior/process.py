# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task

from canopsis.pbehavior.manager import PBehaviorManager


@register_task
def event_processing(engine, event, pbm=None, logger=None, **kwargs):
    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)

    pbm

    # process pbehavior events...

    logger.error('processing: {}'.format(event))


@register_task
def beat_processing(engine, pbm=None, logger=None, **kwargs):
    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)

    res = pbm.compute_pbehaviors_filters()

    logger.error(res)
