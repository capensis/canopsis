# /usr/bin/env python
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

from unittest import TestCase, main
from canopsis.schema.schema import Schema


class TestSchema(TestCase):

    def setUp(self):
        self.s = Schema()

    def test_get_schema(self):
        pass

    def test_get_data_type_schemas(self):
        self.assertRaises(
            NotImplementedError,
            self.s.get_data_type_schemas,
            'data_type')

    def test_push_schema(self):
        pass

    def test_validate_schema(self):
        pass

    def test_validate_data(self):
        pass

    def test_transform(self):
        pass

if __name__ == '__main__':
    main()

