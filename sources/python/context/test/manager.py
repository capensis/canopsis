#!/usr/bin/env python
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
from canopsis.context.manager import Context


class ContextTest(TestCase):

    def setUp(self):
        self.context = Context(data_scope='test_context')

    def test_ctx_storage(self):

        self.context[Context.CTX_STORAGE].drop()

        context = self.context.context

        count_per_entity_type = 2

        # let's iterate on context items in order to create entities
        for n in range(1, len(context)):
            sub_context = context[:n]
            entity_context = {c: c for c in sub_context[1:]}
            entity = {}
            for i in range(count_per_entity_type):
                entity[Context.NAME] = str(i)
                self.context.put(
                    _type=context[n], context=entity_context, entity=entity
                )

        entities = self.context.find()

        self.assertEqual(
            len(entities), count_per_entity_type * (len(context) - 1)
        )

        for n in range(1, len(context)):
            sub_context = context[:n]
            entity_context = {c: c for c in sub_context[1:]}
            entities = self.context.find(
                _type=context[n], context=entity_context
            )
            self.assertEqual(len(entities), 2)

            _id = self.context.get_entity_id(entities[0])
            self.context.remove(ids=_id)
            entities = self.context.find(
                _type=context[n], context=entity_context
            )
            self.assertEqual(len(entities), 1)

            self.context.remove(_type=context[n], context=entity_context)
            entities = self.context.find(
                _type=context[n], context=entity_context
            )
            self.assertEqual(len(entities), 0)


if __name__ == '__main__':
    main()
