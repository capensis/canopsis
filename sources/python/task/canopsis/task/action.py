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

"""
Rule action functions
"""

from canopsis.common.utils import ensure_iterable
from canopsis.task import get_task_with_params, ActionError, register_task


@register_task
def actions(event, ctx, actions=None, raiseError=False, **kwargs):
    """
    Action which process several input actions and returns a list of results
    """

    result = []

    if actions is not None:

        action_confs = ensure_iterable(actions)
        for action_conf in action_confs:
            action_task, params = get_task_with_params(action_conf)
            try:
                action_result = action_task(event=event, ctx=ctx, **params)
            except Exception as e:
                error = ActionError(e)
                if raiseError:
                    raise error
                result.append(error)
            else:
                result.append(action_result)

    return result
