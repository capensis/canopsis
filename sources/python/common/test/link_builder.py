#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

import logging
from unittest import TestCase, main

from canopsis.common.associative_table.associative_table import AssociativeTable
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.link_builder.link_builder import HypertextLinkManager
from canopsis.context_graph.process import create_entity
from canopsis.middleware.core import Middleware


class LinkBuilderTest(TestCase):
    """Test the hyperlink table module.
    """
    def setUp(self):
        self.logger = logging.getLogger('webservice')  # TODO: put a real logger

        self.at_storage = Middleware.get_middleware_by_uri(
            'storage-default-testassociativetable://'
        )
        self.at_manager = AssociativeTableManager(storage=self.at_storage)

        self.config = self.at_manager.get('test_hlm')
        self.config.set('basic_link_builder', 'test_blb')
        self.at_manager.save(self.config)

        blb = self.at_manager.get('test_blb')
        blb.set('base_url', 'http://example.com/?id={0}')
        self.at_manager.save(blb)

        self.htl_manager = HypertextLinkManager(self.config, self.logger)

        self.entity = create_entity(
            id='entity-one',
            name='my-entity',
            etype='resource'
        )

    def tearDown(self):
        """Teardown"""
        self.at_storage.remove_elements()

    def test_hypertextlinkbuilder(self):
        res = self.htl_manager.links_for_entity(entity=self.entity)
        self.assertListEqual(res, [''])

        res = self.htl_manager.links_for_entity(entity=self.entity,
                                                options={})
        self.assertListEqual(res, [''])

        options = {
            'base_url': 'http://example.com/{type}'
        }
        res = self.htl_manager.links_for_entity(entity=self.entity,
                                                options=options)
        self.assertListEqual(res, ['http://example.com/resource'])


if __name__ == '__main__':
    main()
