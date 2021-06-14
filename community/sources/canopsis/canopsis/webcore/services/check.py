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

from canopsis.common.ws import route
from canopsis.check.manager import CheckManager

manager = CheckManager()


def exports(ws):

    @route(ws.application.get, name='check')
    @route(
        ws.application.post,
        payload=['ids', 'state', 'criticity', 'cache'],
        name='check'
    )
    def state(ids=None, state=None, criticity=CheckManager.HARD, cache=False):
        """
        Get entity states.

        :param ids: entity ids.
        :type ids: str or list
        :param int state: state to update if not None.
        :param int criticity: state criticity level (HARD by default).
        :param bool cache: storage cache when udpate state.

        :return: entity states by entity id or one state value if ids is a str.
            None if ids is a str, related entity does not exists and no update
            is required.
        :rtype: int or dict
        """

        result = manager.state(
            ids=ids, state=state, criticity=criticity, cache=cache
        )

        return result

    @route(ws.application.delete, payload=['ids', 'cache'], name='check')
    def del_state(ids=None, cache=False):
        """
        Delete states related to input ids. If ids is None, delete all states.

        :param ids: entity ids. Delete all states if ids is None (default).
        :type ids: str or list
        :param bool cache: storage cache when udpate state.
        """

        result = manager.del_state(ids=ids, cache=cache)

        return result
