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

from canopsis.check.manager import CheckManager


def criticity(state_document, state, criticity=CheckManager.HARD):
    """
    Apply criticity on input state_documents with input state.

    :param dict state_document: state document to process.
    :param int state: state to apply on input state document.
    :param int criticity: criticity of state to update.
    :return: state document.
    :rtype: dict
    """

    state_name = CheckManager.STATE
    last_name = CheckManager.LAST_STATE
    count_name = CheckManager.COUNT
    # get criticity count by criticity level
    criticity_count = CheckManager.CRITICITY_COUNT[criticity]
    # by default, result is a copy of state_document
    result = state_document.copy()
    # if state document does not contain state information
    if state_name not in state_document:
        result.update({
            state_name: state,
            last_name: state,
            count_name: 1
        })
    else:
        # get current entity state
        entity_state = state_document[state_name]
        # if state != entity_state
        if state != entity_state:
            # get count and last state
            last_state = state_document[last_name]
            count = state_document[count_name]
            if last_state != state:  # if state != last state
                count = 1  # initialize count
                last_state = state
            else:  # else increment count
                count += 1
            # if state count is equal or greater than crit count
            if count >= criticity_count:
                count = 1  # initialize count
                entity_state = state  # state entity is state
            # construct a new document with state, count and last state
            result.update({
                state_name: entity_state,
                count_name: count,
                last_name: last_state
            })
    return result
