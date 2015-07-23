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

from canopsis.common.init import basestring
from canopsis.common.utils import lookup
from canopsis.configuration.parameters import Parameter
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.check import Check

#: check manager configuration category
CATEGORY = 'CHECK'
#: check manager conf path
CONF_PATH = 'check/check.conf'


class InvalidState(Exception):
    def __init__(self, state, states):
        self.state = state
        self.states = states

    def __str__(self):
        return 'Invalid state: got value {}, expected one of {}'.format(
            self.state,
            self.states
        )


@add_category(CATEGORY, content=Parameter('types', parser=Parameter.array()))
@conf_paths(CONF_PATH)
class CheckManager(MiddlewareRegistry):
    """Manage entity checking state.

    A state is bound to an entity. Therefore, an entity id is a document state
    id.
    """

    CHECK_STORAGE = 'check_storage'  #: storage name

    ID = '_id'  # Storage.DATA_ID  #: state id field name

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

    #: default function to apply when changing of state
    DEFAULT_F = 'canopsis.check.task.criticity'
    valid_states = [0, 1, 2, 3]

    def __init__(self, types=None, *args, **kwargs):

        super(CheckManager, self).__init__(*args, **kwargs)

        self.types = types

    # TODO , is it used, is it usefull to manage state this way
    def state(
        self, ids=None, state=None, criticity=HARD, f=DEFAULT_F, query=None,
        cache=False
    ):
        """Get/update entity state(s).

        :param ids: entity id(s). Default is all entity ids.
        :type ids: str or list
        :param int state: state to update if not None.
        :param int criticity: state criticity level (HARD by default).
        :param f: new state calculation function if state is not None.
        :param dict query: additional query to use in order to find states.
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
            ids=ids, query=query
        )
        # if state document exists
        if state_documents is not None:
            # ensure state_documents is a list
            if isinstance(state_documents, dict):
                state_documents = [state_documents]
            # save id and state field name
            id_field, state_field = CheckManager.ID, CheckManager.STATE
            # result is a dictionary of entity id, state value
            result = {}
            for state_document in state_documents:
                entity_id = state_document[id_field]
                entity_state = state_document[state_field]
                result[entity_id] = entity_state
        # if state has to be updated
        if state is not None:
            # get the right state function
            f = lookup(f) if isinstance(f, basestring) else f
            # save field name for quick access
            id_name = CheckManager.ID
            state_name = CheckManager.STATE
            # save storage for quick access
            storage = self[CheckManager.CHECK_STORAGE]
            # ensure entity_ids is a set
            if isinstance(ids, basestring):
                entity_ids = set([ids])
            elif ids is None:
                if state_documents is None:
                    entity_ids = set()
                else:
                    entity_ids = set([sd[id_name] for sd in state_documents])
            else:
                entity_ids = set(ids)
            # if states exist in DB
            if state_documents is not None:
                # for all found documents
                for state_document in state_documents:
                    # get document id
                    _id = state_document[id_name]
                    # remove _id from entity_ids
                    entity_ids.remove(_id)
                    # get new state with f
                    new_state_document = f(
                        state_document=state_document,
                        state=state,
                        criticity=criticity
                    )
                    # save new state_document if old != new
                    if state_document != new_state_document:
                        storage.put_element(
                            _id=_id, element=new_state_document, cache=cache
                        )
                    # save state entity in result
                    result[_id] = new_state_document[state_name]
            # for all not found documents
            for entity_id in entity_ids:
                # create a new document
                state_document = {
                    id_name: entity_id,
                }
                new_state_document = f(
                    state_document=state_document,
                    state=state,
                    criticity=criticity
                )
                # save it in storage
                storage.put_element(
                    _id=entity_id, element=new_state_document, cache=cache
                )
                # and put entity state in the result
                result[entity_id] = state

        # ensure result is a state if ids is a basestring
        if result is not None and isinstance(ids, basestring):
            result = result[ids] if result else None

        return result

    """
    Simple way to manage states.
    The following methods allow crud operation on a state in database
    with only a data couple of on identifier and an ID
    """

    def del_state(self, ids=None, query=None, cache=False):
        """Delete states related to input ids. If ids is None, delete all
        states.

        :param ids: entity ids. Delete all states if ids is None (default).
        :type ids: str or list
        :param dict query: selection query.
        :param bool cache: storage cache when udpate state.
        """

        return self[CheckManager.CHECK_STORAGE].remove_elements(
            ids=ids, _filter=query, cache=cache
        )

    def put_state(self, entity_id, state, cache=False):
        """
        Allow persistance of a state

        :param entity_id: the identifier for the entity.
        :param state: the state to persist.
        """

        if state not in self.valid_states or not isinstance(state, int):
            raise InvalidState(state, self.valid_states)

        return self[CheckManager.CHECK_STORAGE].put_element(
            _id=entity_id,
            element={'state': state},
            cache=cache
        )

    def get_state(self, ids=None):
        """
        Retrieve state from database depending on an id list

        :param ids: a list of identifier that may have a state in database.
        """
        states = self[CheckManager.CHECK_STORAGE].get_elements(
            ids=ids
        )
        return states
