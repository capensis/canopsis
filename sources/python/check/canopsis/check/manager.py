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

__version__ = "0.1"

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.storage import Storage
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.check import Check

#: check manager configuration category
CATEGORY = 'CHECK'
#: check manager conf path
CONF_PATH = 'check/check.conf'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class CheckManager(MiddlewareRegistry):
    """
    Manage entity checking state.

    A state is bound to an entity. Therefore, an entity id is a document state
        id.
    """

    CHECK_STORAGE = 'check_storage'  #: storage name

    ID = Storage.DATA_ID  #: state id field name

    STATE = Check.STATE  #: state field name
    LAST_STATE = 'last'  #: last state field name if criticity != HARD
    COUNT = 'count'  #: last state count if criticity != 0

    HARD = 0  #: hard criticity
    SOFT = 1  #: soft criticity

    #: number of state before updating by criticity levels
    CRITICITY_COUNT = {
        HARD: 1,
        SOFT: 3
    }

    def state(self, ids, state=None, criticity=HARD, cache=False):
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
        # default result is None
        result = {}
        # get state document
        state_documents = self[CheckManager.CHECK_STORAGE].get_elements(
            ids=ids
        )
        # if state document exists
        if state_documents is not None:
            # ensure state_documents is a list
            if isinstance(state_documents, dict):
                state_documents = [state_documents]
            # save id and state field name
            id_field, state_field = CheckManager.ID, CheckManager.STATE
            #result a dictionary of entity id, state value
            result = {}
            for state_document in state_documents:
                entity_id = state_document[id_field]
                entity_state = state_document[state_field]
                result[entity_id] = entity_state

        # if state has to be updated
        if state is not None:
            # save field name for quick access
            id_name = CheckManager.ID
            state_name = CheckManager.STATE
            last_name = CheckManager.LAST_STATE
            count_name = CheckManager.COUNT
            # save storage for quick access
            storage = self[CheckManager.CHECK_STORAGE]
            # get criticity count by criticity level
            criticity_count = CheckManager.CRITICITY_COUNT[criticity]
            # save entity ids
            entity_ids = ids
            # in ensuring it is a set
            if isinstance(entity_ids, basestring):
                entity_ids = {entity_ids}
            if state_documents is not None:
                # for all found documents
                for state_document in state_documents:
                    # get document id
                    _id = state_document[id_name]
                    # remove _id from entity_ids
                    entity_ids.remove(_id)
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
                    new_state_document = {
                        id_name: _id,
                        state_name: entity_state,
                        count_name: count,
                        last_name: last_state
                    }
                    # save new state_document
                    storage.put_element(_id=_id, element=new_state_document)
                    # save state entity in result
                    result[_id] = entity_state
            # for all not found documents
            for entity_id in entity_ids:
                count = 1
                # if criticity is high, entity_state is state
                if criticity_count <= count:
                    entity_state = state
                else:  # else it is unknown
                    entity_state = Check.UNKNOWN
                # create a new document
                new_state_document = {
                    id_name: entity_id,
                    state_name: entity_state,
                    count_name: count,
                    last_name: state
                }
                # save it in storage
                storage.put_element(_id=entity_id, element=new_state_document)
                # and put entity state in the result
                result[entity_id] = entity_state

        # ensure result is a state if ids is a basestring
        if result is not None and isinstance(ids, basestring):
            result = result[ids] if result else None

        return result

    def del_state(self, ids=None):
        """
        Delete states related to input ids. If ids is None, delete all states.

        :param ids: entity ids. Delete all states if ids is None (default).
        :type ids: str or list
        """

        self[CheckManager.CHECK_STORAGE].remove_elements(ids=ids)
