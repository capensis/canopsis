# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from canopsis.common.utils import singleton_per_scope
from canopsis.router.manager import RouterManager
from canopsis.engines.core import publish


def event_processing(engine, event, manager=None, logger=None, **_):
    if manager is None:
        manager = singleton_per_scope(RouterManager)

    if manager.match_filters(event):
        event = manager.apply_patchs(event)
        rk = manager.get_routing_key(event)

        publish(
            publisher=engine.amqp,
            event=event,
            rk=rk,
            exchange=manager.exchange,
            logger=logger
        )


def beat_processing(engine, manager=None, logger=None, **_):
    if manager is None:
        manager = singleton_per_scope(RouterManager)

    manager.reload()
