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

from canopsis.task.core import register_task
from canopsis.engines.core import publish


@register_task('router.actions.drop')
def action_drop(manager, event, **_):
    return None


@register_task('router.actions.pass')
def action_pass(manager, event, **_):
    return event


@register_task('router.actions.override')
def action_override(manager, event, field=None, value=None, **_):
    if field is not None and value is not None:
        if field in event and isinstance(event[field], list):
            event[field].append(value)

        else:
            event[field] = value

    return event


@register_task('router.actions.remove')
def action_remove(manager, event, key=None, element=None, met=0, **_):
    if key:
        if element:
            if met:
                for i, to_del in enumerate(event, key):
                    if to_del['name'] == element:
                        del event[key][i]
                        break

            elif isinstance(event[key], dict):
                del event[key][element]

            elif isinstance(event[key], list):
                event[key].remove(element)

        else:
            del event[key]

    return event


@register_task('router.actions.execjob')
def action_execjob(manager, event, job=None, publisher=None, **_):
    if job:
        record = manager[manager.CONFIG_STORAGE].get_elements(
            ids=job
        )

        if record is not None:
            record['context'] = event
            publish(
                publisher=publisher,
                event=record,
                rk='Engine_scheduler',
                exchange='amq.direct'
            )

    return event
