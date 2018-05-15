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
from canopsis.ctxprop.manager import CTXPropManager


def exports(ws):

    manager = CTXPropManager()

    @route(ws.application.get)
    def registries():
        """Get list of available registries.

        :return: available registries.
        :rtype: list
        """

        return manager.registries

    @route(ws.application.get)
    def get(registries=None, ids=None, query=None, children=True):
        """Get property of registries.

        :param str(s) registries: registry names in charge of retrieving
            property. If None (default), get property from all registries.
        :param str(s) ids: ctx id(s) from where find property. If None
            (default), get all property of all entities.
        :param dict query: specific query to apply on the method execution.
        :param bool children: if True (default) propagate the method on ctx
            children.
        :return: property of registries by registry. If registries is a string,
            the result is specific registries property.
        :rtype: dict
        """

        result = manager.get(
            registries=registries, ids=ids, query=query, children=children
        )

        return result

    @route(ws.application.get)
    def count(registries=None, ids=None, query=None, children=True):
        """Count property of registries.

        :param str(s) registries: registry names in charge of counting
            property. If None (default), count property from all registries.
        :param str(s) ids: ctx id(s) from where count property. If None
            (default), count all property of all entities.
        :param dict query: specific query to apply on the method execution.
        :param bool children: if True (default) propagate the method on ctx
            children.
        :return: property of registries by registry. If registries is a string,
            the result is specific count of registries property.
        :rtype: dict
        """

        result = manager.count(
            registries=registries, ids=ids, query=query, children=children
        )

        return result

    @route(ws.application.delete)
    def delete(
            registries=None, ids=None, query=None, children=True,
            force=False, cache=False
    ):
        """Delete property of registries and returns number of property
        deleted per registry and ctx id.

        :param str(s) registries: registry names in charge of deleting
            property. If None (default), delete property from all registries.
        :param str(s) ids: ctx id(s) from where delete property. If None
            (default), delete all property of all entities.
        :param dict query: specific query to apply on the method execution.
        :param bool force: if True (False by default), accept to nonify
            ctx_ids in order to delete all existing ctx property.
        :param bool children: if True (default) propagate the method on ctx
            children.
        :return: property of registries by registry. If registries is a string,
            the result is specific registries property.
        :rtype: dict
        """

        # check if force is True if registries is None
        if ids is None and not force:
            raise CTXPropManager.Error(
                "Impossible to remove all existing ctx property. Use force."
            )

        result = manager.delete(
            registries=registries, ids=ids, query=query,
            children=children, force=force, cache=cache
        )

        return result

    def ids(registries=None, query=None):
        """Get ctx ids from different registries thanks to an input query.

        :param registries: registry names in charge of deleting property. If
            None (default), delete property from all registries.
        :param dict query: specific query to apply on the method execution.
        :return: ctx ids of registries by registry. If registries is a string,
            the result is specific ctx ids registries.
        :rtype: dict
        """

        result = manager.ids(registries=registries, query=query)

        return result
