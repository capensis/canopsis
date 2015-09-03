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

"""The CTXPropManager is used such as a CTXPropRegistry with read and delete
operations.

In such way, all methods execution are delegated to its registries.
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.context.manager import Context


CONF_PATH = 'ctxprop/ctxprop.conf'
CATEGORY = 'CTXPROP'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class CTXPropManager(MiddlewareRegistry):
    """Manage context property.
    """

    class Error(Exception):
        """Handle CTXPropManager errors.
        """

    DATA_SCOPE = 'ctxprop'  #: default data scope

    def __init__(self, *args, **kwargs):

        super(CTXPropManager, self).__init__(*args, **kwargs)

        self.context = Context()

    @property
    def registries(self):
        """Get list of available registries.

        :return: available registries.
        :rtype: list
        """

        return self.configurables.keys()

    def get(self, registries=None, ids=None, query=None, children=True):
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

        result = self._process_providers(
            cmd='get', registries=registries,
            ids=ids, query=query, children=children
        )

        return result

    def count(self, registries=None, ids=None, query=None, children=True):
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

        result = self._process_providers(
            cmd='count', registries=registries,
            ids=ids, query=query, children=children
        )

        return result

    def delete(
            self, registries=None, ids=None, query=None, children=True,
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
            ids in order to delete all existing ctx property.
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

        result = self._process_providers(
            cmd='delete', registries=registries,
            ids=ids, query=query,
            children=children, force=force, cache=cache
        )

        return result

    def ids(self, registries=None, query=None):
        """Get ctx ids from different registries thanks to an input query.

        :param registries: registry names in charge of deleting property. If
            None (default), delete property from all registries.
        :param dict query: specific query to apply on the method execution.
        :return: ctx ids of registries by registry. If registries is a string,
            the result is specific ctx ids registries.
        :rtype: dict
        """

        result = self._process_providers(
            cmd='ids', registries=registries, query=query
        )

        return result

    def _process_providers(self, registries, cmd, **kwargs):
        """Process cmd on all registries with specific kwargs.

        :param list registries: registries on which run cmd.
        :param str cmd: cmd to apply on all registries.
        :return: registry cmd result per by registry. If registries is a
            string, the result is specific registry cmd result.
        :rtype: dict or list
        """

        result = {}

        registries, unique = self._providers_unique(registries)

        # update kwargs ids if necessary
        if 'ids' in kwargs:
            # get ids and children from kwargs
            ids = kwargs['ids']
            children = kwargs.pop('children')  # remove children from kwargs
            kwargs['ids'] = self._add_children(
                pctx_ids=ids, children=children
            )

        for registry in registries:
            provider_cmd = getattr(self[registry], cmd)
            try:
                fresult = provider_cmd(**kwargs)
            except Exception as ex:
                self.logger.error(
                    'Error ({0}) on Funder: {1} with kwargs {2}'.format(
                        ex, registry, kwargs
                    )
                )
            else:
                result[registry] = fresult

        if unique:
            result = result[registries[0]] if result else None

        return result

    def _providers_unique(self, registries):
        """Return a couple of (list of registry names, registries is unique).

        :param str(s) registries: registry name(s).
        :return: a couple of (list of registry names, registries is unique)
            where unique is True iif registries is a string and, the list of
            registry names depending on registries type:

        - None: self.configurables.keys().
        - str: [registries].
        - iterable: registries.
        :rtype: tuple
        """

        unique = isinstance(registries, basestring)

        if unique:
            registries = [registries]
        elif registries is None:
            registries = self.registries

        result = registries, unique

        return result

    def _add_children(self, pctx_ids, children=True):
        """Add children ctx ids to input ctx ids if input children is True.

        :param str(s) pctx_ids: parent pctx_ids.
        :param bool children: if True add children ctx ids to pctx_ids.
        :return: pctx_ids with children ctx ids.
        :rtype: list
        """

        result = pctx_ids  # by default, result is pctx_ids

        if pctx_ids is not None and children:  # if children are requested
            if isinstance(pctx_ids, basestring):
                pctx_ids = [pctx_ids]
            result = pctx_ids[:]  # set result with a copy of pctx ids
            for pctx_id in pctx_ids:
                # get children ctx
                ctx = self.context.get_ctx_by_id(pctx_id)
                children = self.context.get_children(ctx)
                children_ids = [
                    self.context.get_ctx_id(child) for child in children
                ]
                # add children_ids to result
                result += children_ids

        return result
