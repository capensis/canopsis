# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from b3j0f.conf import Configurable, category, Parameter, Array
from b3j0f.requester.driver.composite import DriverComposite
from b3j0f.requester.driver.custom import obj2driver
from b3j0f.requester.driver.ctx import Context
from b3j0f.utils.path import lookup

from canopsis.middleware.core import Middleware
from six import string_types


CONF_PATH = 'sysreq/middleware.conf'
CATEGORY = 'SYSREQ'
CONTENT = [
    Parameter(name='middlewares', type=Array(string_types)),
    Parameter(name='managers', type=Array(string_types))
]


@Configurable(paths=CONF_PATH, conf=category(CATEGORY, *CONTENT))
class SysReq(Middleware):

    __protocol__ = 'sysreq'

    def __init__(self, *args, **kwargs):
        super(SysReq, self).__init__(*args, **kwargs)

        self.dirty = True

    @property
    def middlewares(self):
        if not hasattr(self, '_middlewares'):
            self.middlewares = None

        return self._middlewares

    @middlewares.setter
    def middlewares(self, value):
        if value is None:
            value = []

        self._middlewares = []

        for middleware in value:
            if isinstance(middleware, string_types):
                middleware = Middleware.get_middleware_by_uri(middleware)
                self._middlewares.append(middleware)

            elif isinstance(middleware, Middleware):
                self._middlewares.append(middleware)

            else:
                raise TypeError(
                    'Expected string or Middleware, got: {0}'.format(
                        middleware.__class__.__name__
                    )
                )

        self.dirty = True

    @property
    def managers(self):
        if not hasattr(self, '_managers'):
            self.managers = None

        return self._managers

    @managers.setter
    def managers(self, value):
        if value is None:
            value = []

        self._managers = []

        for manager in value:
            if isinstance(manager, string_types):
                cls = lookup(manager)
                manager = cls()

            self._managers.append(manager)

        self.dirty = True

    @property
    def driver(self):
        if not hasattr(self, '_driver') or self.dirty:
            drivers = [
                obj2driver(obj)
                for obj in self.middlewares + self.managers
            ]

            self._driver = DriverComposite(drivers)
            self.dirty = False

        return self._driver

    def __call__(self, cruds, ctx=None):
        if ctx is not None and not isinstance(ctx, Context):
            ctx = Context(ctx)

        t = self.driver.open(ctx=ctx, cruds=cruds)
        return t.commit()
