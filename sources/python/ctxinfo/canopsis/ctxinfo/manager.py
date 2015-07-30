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

    class Error(Exception):
        """Handle CTXInfoManager errors.
        """

    DATA_SCOPE = 'ctxinfo'  #: default data scope

    def __init__(self, *args, **kwargs):

        super(CTXInfoManager, self).__init__(*args, **kwargs)

        self.context = Context()

    @property
    def funders(self):
        """Get list of available funders.

        :return: list of available funders.
        :rtype: list
        """

        return self.configurables.keys()

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

        result = self._process_funders(
            cmd='get', funders=funders,
            entity_ids=entity_ids, query=query, children=children
        )

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

        result = self._process_funders(
            cmd='count',  funders=funders,
            entity_ids=entity_ids, query=query, children=children
        )

        return result

    def delete(
        self, funders=None, entity_ids=None, query=None, children=True,
        force=False, cache=False
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
        :param bool force: if True (False by default), accept to nonify
            entity_ids in order to delete all existing entity information.
        :param bool children: if True (default) propagate the method on entity
            children.
        :return: information of funders by funder. If funders is a string, the
            result is specific funders information.
        :rtype: dict
        """

        # check if force is True if funders is None
        if entity_ids is None and not force:
            raise CTXInfoManager.Error(
                "Impossible to remove all existing ctx information. Use force."
            )

        result = self._process_funders(
            cmd='delete', funders=funders,
            entity_ids=entity_ids, query=query,
            children=children, force=force, cache=cache
        )

        return result

    def entity_ids(self, funders=None, query=None):
        """Get entity ids from different funders thanks to an input query.

        :param funders: funder names in charge of deleting information. If None
            (default), delete information from all funders.
        :param dict query: specific query to apply on the method execution.
        :return: entity ids of funders by funder. If funders is a string, the
            result is specific entity ids funders.
        :rtype: dict
        """

        result = self._process_funders(
            cmd='entity_ids', funders=funders, query=query
        )

        return result

    def _process_funders(self, funders, cmd, **kwargs):
        """Process cmd on all funders with specific kwargs.

        :param list funders: funders on which run cmd.
        :param str cmd: cmd to apply on all funders.
        :return: funder cmd result per by funder. If funders is a string, the
            result is specific funder cmd result.
        :rtype: dict or list
        """

        result = {}

        funders, unique = self._funders_unique(funders)

        # update kwargs entity_ids if necessary
        if 'entity_ids' in kwargs:
            # get entity_ids and children from kwargs
            entity_ids = kwargs['entity_ids']
            children = kwargs.pop('children')  # remove children from kwargs
            kwargs['entity_ids'] = self._add_children(
                pentity_ids=entity_ids, children=children
            )

        for funder in funders:
            funder_cmd = getattr(self[funder], cmd)
            try:
                fresult = funder_cmd(**kwargs)
            except Exception as e:
                self.logger.error(
                    'Error ({0}) on Funder: {1} with kwargs {2}'.format(
                        e, funder, kwargs
                    )
                )
            else:
                result[funder] = fresult

        if unique:
            result = result[funders[0]] if result else None

        return result

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
                entity = self.context.get_entity_by_id(pentity_id)
                children = self.context.get_children(entity)
                children_ids = [
                    self.context.get_entity_id(child) for child in children
                ]
                # add children_ids to result
                result += children_ids

        return result
