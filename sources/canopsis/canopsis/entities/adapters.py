#!/usr/bin/env python
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

from .models import Entity


class Adapter(object):

    COLLECTION = 'default_entities'

    def __init__(self, mongo_client):
        self.mongo_client = mongo_client

    def find_all_enabled(self):
        query = {
            'enabled': True
        }
        entities = []
        collection = self.mongo_client[self.COLLECTION]
        for entity in collection.find(query):
            entities.append(make_entity_from_mongo(entity))
        return entities


def make_entity_from_mongo(entity_dict):
    return Entity(
        entity_dict.get('_id'),
        entity_dict.get('name'),
        entity_dict.get('impact'),
        entity_dict.get('depends'),
        entity_dict.get('enable_history'),
        entity_dict.get('measurements'),
        entity_dict.get('enabled'),
        entity_dict.get('infos'),
        entity_dict.get('type')
    )
