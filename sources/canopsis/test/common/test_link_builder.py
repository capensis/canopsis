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
import unittest
from canopsis.common import root_path
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.link_builder.link_builder import HypertextLinkManager
from canopsis.middleware.core import Middleware
import xmlrunner


class LinkBuilderTest(unittest.TestCase):
    """Test the hyperlink table module.
    """
    def setUp(self):
        self.logger = logging.getLogger()
        self.logger.setLevel(logging.DEBUG)

        self.at_storage = Middleware.get_middleware_by_uri(
            'storage-default-testassociativetable://'
        )
        self.at_manager = AssociativeTableManager(collection=self.at_storage._backend,
                                                  logger=self.logger)

        self.config = self.at_manager.create('test_hlm')
        self.config.set('basic_link_builder', {})
        self.at_manager.save(self.config)

        self.htl_manager = HypertextLinkManager(config=self.config.get_all(),
                                                logger=self.logger)

        self.entity = {
            '_id': 'entity-one',
            'type': 'resource',
            'name': 'my-entity',
            'depends': [],
            'impact': [],
            'measurements': {},
            'infos': {}
        }

    def tearDown(self):
        """Teardown"""
        self.at_storage.remove_elements()

    def test_links_for_entity(self):
        res = self.htl_manager.links_for_entity(entity=self.entity)
        self.assertDictEqual(res, {})

        res = self.htl_manager.links_for_entity(entity=self.entity,
                                                options={})
        self.assertDictEqual(res, {})

        options = {
            'obiwan': 'jedi',
            'darthvader': 'sith'
        }
        res = self.htl_manager.links_for_entity(entity=self.entity,
                                                options=options)
        self.assertDictEqual(res, {})

        options = {
            'base_url': 'http://example.com/{type}'
        }
        res = self.htl_manager.links_for_entity(entity=self.entity,
                                                options=options)
        self.assertDictEqual(res, {'links': ['http://example.com/resource']})

        category = 'alexandrin'
        config = {
            'base_url': 'http://example.com/{type}',
            'category': category
        }
        htl_manager2 = HypertextLinkManager(config={'basic_link_builder': config},
                                            logger=self.logger)
        res = htl_manager2.links_for_entity(entity=self.entity)
        self.assertDictEqual(res, {category: ['http://example.com/resource']})


if __name__ == '__main__':
    output = root_path + "tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
