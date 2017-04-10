#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.import_ctx import ContextGraphImport
from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware


class BaseTest(TestCase):

    def setUp(self):
        self.ctx_import = ContextGraphImport()
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.organisations_storage = Middleware.get_middleware_by_uri(
            'storage-default-testorganisations://'
        )
        self.users_storage = Middleware.get_middleware_by_uri(
            'storage-default-testusers://'
        )

        self.ctx_import[ContextGraph.ENTITIES_STORAGE] = self.entities_storage
        self.ctx_import[
            ContextGraph.ORGANISATIONS_STORAGE] = self.organisations_storage
        self.ctx_import[ContextGraph.USERS_STORAGE] = self.users_storage

        self.template = {'_id': None,
                         'type': 'connector',
                         'name': 'conn-name1',
                         'depends': [],
                         'impact': [],
                         'measurements': [],
                         'infos': {}}

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.organisations_storage.remove_elements()

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1["depends"])
        sorted(entity1["impact"])
        sorted(entity2["depends"])
        sorted(entity2["impact"])
        self.assertDictEqual(entity1, entity2)


class GetEntitiesToUpdate(BaseTest):

    def _test(self, ctx, entities):
        for id_ in entities.keys():
            try:
                entity = ctx[id_]
            except KeyError:
                self.fail("KeyError : missing key \"ent1\" in ctx")
            else:
                self.assertEqualEntities(entity, entities[id_])

    def test_entities(self):
        entities = {"ent1": self.template.copy(),
                    "ent2": self.template.copy(),
                    "ent3": self.template.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"

        self.ctx_import.put_entities(entities.values())

        json = {ContextGraphImport.CIS: [{"_id": "ent1"},
                                         {"_id": "ent2"},
                                         {"_id": "ent3"}],
                ContextGraphImport.LINKS: [{ContextGraphImport.FROM: "ent1",
                                            ContextGraphImport.TO: "ent2"},
                                           {ContextGraphImport.FROM: "ent1",
                                            ContextGraphImport.TO: "ent3"}]}

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(json)

        self._test(ctx, entities)


    def test_no_entities(self):
        entities = {"ent1": self.template.copy(),
                    "ent2": self.template.copy(),
                    "ent3": self.template.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"

        self.ctx_import.put_entities(entities.values())

        entities = {}

        json = {ContextGraphImport.CIS: [],
                ContextGraphImport.LINKS: []}

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(json)

        self._test(ctx, entities)


if __name__ == '__main__':
    main()
