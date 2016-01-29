#!/usr/bin/env python
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

from __future__ import print_function

from unittest import TestCase, main

from canopsis.storage.core import Storage
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import conf_paths, add_category
from canopsis.configuration.model import Parameter
from canopsis.common.utils import lookup


@conf_paths('storage/test.conf')
@add_category('TEST',
    content=[
        Parameter('protocols', Parameter.array),
        Parameter('data_type', Parameter.array),
        Parameter('data_scopes', Parameter.array),
        Parameter('storages', Parameter.array),
        Parameter('params', eval)
    ]
)
class BaseTestConfiguration(Configurable):
    """Default test configuration."""


class BaseStorageTest(TestCase):
    """Class to override in order to load a set of storages depending on
    specific protocols, data_scopes and data_types.
    """

    def _testconfcls(self):

        return BaseTestConfiguration

    def test(self):
        """To override in order to run tests related to one storage."""

        map(lambda x: self._test(x), self.storages)

    def tearDown(self):
        """Drop content of storages."""

        map(lambda x: x.drop(), self.storages)

    def setUp(self):
        """initialize storages"""

        self.storages = []

        testconf = self._testconfcls()

        if testconf.storages:

            for storage in testconf.storages:

                storagecls = lookup(storage)

                storage = storagecls(
                    data_scope=data_scope, conf_paths=testconf.conf_paths,
                    **testconf.params
                )

                self.storages.append(storage)

        else:

            for protocol in testconf.protocols:

                for data_type in testconf.data_types:

                    for data_scope in testconf.data_scopes:

                        storage = Storage.get_middleware(
                            protocol=protocol, data_type=data_type,
                            data_scope=data_scope,
                            conf_paths=testconf.conf_paths,
                            **testconf.params
                        )

                        self.storages.append(storage)


if __name__ == '__main__':
    main()
