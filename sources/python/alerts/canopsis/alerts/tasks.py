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

from canopsis.task.core import register_task
from canopsis.check import archiver


@register_task('alerts.useraction.ack')
def acknowledge(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.ackremove')
def unacknowledge(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.cancel')
def cancel(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.uncancel')
def restore(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.declareticket')
def declare_ticket(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.assocticket')
def associate_ticket(manager, entity, author, message, event):
    pass


@register_task('alerts.useraction.changestate')
def change_state(manager, entity, author, message, event):
    pass


@register_task('alerts.systemaction.state_increase')
def state_increase(manager, entity, state):
    return archiver.OFF


@register_task('alerts.systemaction.state_decrease')
def state_decrease(manager, entity, state):
    return archiver.OFF


@register_task('alerts.systemaction.status_increase')
def status_increase(manager, entity, status):
    return True


@register_task('alerts.systemaction.status_decrease')
def status_decrease(manager, entity, status):
    return False
