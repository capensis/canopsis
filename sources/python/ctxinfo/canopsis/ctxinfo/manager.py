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

"""The CTXInfoManager is used such as a CTXInfoFunder registry with read and
delete operations.

In such way, all methods execution are delegated to its funders.
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.context.manager import Context


CONF_PATH = 'ctxinfo/ctxinfo.conf'
CATEGORY = 'CTXINFO'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class CTXInfoManager(MiddlewareRegistry):
    """Manage context information.
    """

    DATA_SCOPE = 'ctxinfo'  #: default data scope

    def __init__(self, *args, **kwargs):

        super(CTXInfoManager, self).__init__(*args, **kwargs)

        self.context = Context()

    def _funders_unique(self, funders):
        """Return a couple of (list of funder names, funders is unique).

        :param funders: funder name(s).
        :type funders: str or list
        :return: a couple of (list of funder names, funders is unique) where
            unique is True iif funders is a string and, the list of funder
            names depending on funders type:

        - None: self.configurables.keys().
        - str: [funders].
        - iterable: funders.
        :rtype: tuple
        """

        unique = isinstance(funders, basestring)

        if unique:
            funders = [funders]
        elif funders is None:
            funders = self.configurables.keys()

        result = funders, unique

        return result

    def _add_children(self, pentity_ids, children=True):
        """Add children entity ids to input entity ids if input children is
        True.

        :param pentity_ids: parent pentity_ids.
        :type pentity_ids: str or list
        :param bool children: if True add children entity ids to pentity_ids.
        :return: pentity_ids with children entity ids.
        :rtype: list
        """

        result = pentity_ids  # by default, result is pentity_ids

        if pentity_ids is not None and children:  # if children are requested
            if isinstance(pentity_ids, basestring):
                pentity_ids = [pentity_ids]
            result = pentity_ids[:]  # set result with a copy of pentity ids
            for pentity_id in pentity_ids:
                # get children entity
                entity = self.context.get_entity(pentity_id)
                children = self.context.get_children(entity)
                children_ids = [
                    self.context.get_entity_id(child) for child in children
                ]
                # add children_ids to result
                result += children_ids

        return result

    def get(self, funders=None, entity_ids=None, query=None, children=True):
        """Get information of funders.

        :param funders: funder names in charge of retrieve information. If None
            (default), get information from all funders.
        :type funders: str or list
        :param entity_ids: entity id(s) from where find information. If None
            (default), get all information of all entities.
        :type entity_ids: str or list
        :param dict query: specific query to apply on the method execution.
        :param bool children: if True (default) propagate the method on entity
            children.
        :return: information of funders by funder. If funders is a string, the
            result is specific funders information.
        :rtype: dict
        """

        result = {}

        funders, isunique = self._funders_unique(funders)
        # update entity ids related to children parameter
        entity_ids = self._add_children(
            pentity_ids=entity_ids, children=children
        )
        # deleta the get method to all funders
        for funder in funders:
            fresult = self[funder].get(entity_ids=entity_ids, query=query)
            result[funder] = fresult

        if isunique:
            result = result[funders[0]] if result else None

        return result

    def count(self, funders=None, entity_ids=None, query=None, children=True):
        """Count information of funders.

        :param funders: funder names in charge of counting information. If None
            (default), count information from all funders.
        :type funders: str or list
        :param entity_ids: entity id(s) from where count information. If None
            (default), count all information of all entities.
        :type entity_ids: str or list
        :param dict query: specific query to apply on the method execution.
        :param bool children: if True (default) propagate the method on entity
            children.
        :return: information of funders by funder. If funders is a string, the
            result is specific count of funders information.
        :rtype: dict
        """

        result = {}

        funders, isunique = self._funders_unique(funders)
        # update entity_ids with children if necessary
        entity_ids = self._add_children(
            pentity_ids=entity_ids, children=children
        )
        for funder in funders:
            fresult = self[funder].count(entity_ids=entity_ids, query=query)
            result[funder] = fresult

        if isunique:
            result = result[funders[0]] if result else None

        return result

    def delete(
        self, funders=None, entity_ids=None, query=None, children=True,
        cache=False
    ):
        """Delete information of funders and returns number of information
        deleted per funder and entity id.

        :param funders: funder names in charge of deleting information. If None
            (default), delete information from all funders.
        :type funders: str or list
        :param entity_ids: entity id(s) from where delete information. If None
            (default), delete all information of all entities.
        :type entity_ids: str or list
        :param dict query: specific query to apply on the method execution.
        :param bool children: if True (default) propagate the method on entity
            children.
        :return: information of funders by funder. If funders is a string, the
            result is specific funders information.
        :rtype: dict
        """

        result = {}

        funders, isunique = self._funders_unique(funders)
        # update entity_ids with children if necessary
        entity_ids = self._add_children(
            pentity_ids=entity_ids, children=children
        )
        for funder in funders:
            fresult = self[funder].get(entity_ids=entity_ids, query=query)
            result[funder] = fresult

        if isunique:
            result = result[funders[0]] if result else None

        return result
