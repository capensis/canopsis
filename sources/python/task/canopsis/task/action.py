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

from canopsis.utils.init import basestring
from canopsis.task import get_task_with_params, register_task

#: action result
RESULT = 'result'
#: action error
ERROR = 'error'


@register_task
def actions(confs=None, raiseError=False, **kwargs):
    """
    Action which process several input confs and returns a list
        of dict {'result': action result, 'error': action error}.

    :return: a list containing dict of {RESULT: result, ERROR: error}.
    :rtype: list
    """

    result = []

    if confs is not None:
        # ensure confs is a list
        if isinstance(confs, basestring):
            confs = [confs]
        for conf in confs:
            action, params = get_task_with_params(conf)
            params.update(kwargs)
            try:
                action_result = action(**params)
            except Exception as action_error:
                if raiseError:
                    raise

            action_error = None
            result_to_append = {
                RESULT: action_result,
                ERROR: action_error
            }
            result.append(result_to_append)

    return result
