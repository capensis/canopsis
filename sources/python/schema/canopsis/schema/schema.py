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


class Schema(object):

    def __init__(self, *args, **kwargs):
        super(Schema, self).__init__(*args, **kwargs)
        pass

    def get_schema(self, schema):
        raise NotImplementedError

    def get_data_type_schemas(self, data_type):
        raise NotImplementedError

    def push_schema(self, data_type, schema, name, schema_type):
        raise NotImplementedError

    def validate_schema(self, schema):
        raise NotImplementedError

    def validate_data(self, data, structural_schema, metamorphic_schema=None,
                      validate_schemas=False):
        raise NotImplementedError

    def transform(self,
                  data,
                  metamorphic_schema,
                  from_structure=None,
                  to_structure=None,
                  validate_data=False,
                  validate_schemas=False,
                  ):
        raise NotImplementedError

