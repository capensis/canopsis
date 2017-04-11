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

        self.template_ent = {'_id': None,
                             'type': 'connector',
                             'name': 'conn-name1',
                             'depends': [],
                             'impact': [],
                             'measurements': [],
                             'infos': {}}
        self.template_ci = {ContextGraphImport.K_ID: None,
                            ContextGraphImport.K_NAME: "Name",
                            ContextGraphImport.K_TYPE: "Type",
                            ContextGraphImport.K_INFOS: {},
                            ContextGraphImport.K_ACTION: None}

        self.template_link = {ContextGraphImport.K_FROM: None,
                              ContextGraphImport.K_TO: None,
                              ContextGraphImport.K_INFOS: {},
                              ContextGraphImport.K_ACTION: None}

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
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"

        self.ctx_import.put_entities(entities.values())

        json = {ContextGraphImport.K_CIS: [{ContextGraphImport.K_ID: "ent1"},
                                         {ContextGraphImport.K_ID: "ent2"},
                                         {ContextGraphImport.K_ID: "ent3"}],
                ContextGraphImport.K_LINKS: [{ContextGraphImport.K_FROM: "ent1",
                                            ContextGraphImport.K_TO: "ent2"},
                                           {ContextGraphImport.K_FROM: "ent1",
                                            ContextGraphImport.K_TO: "ent3"}]}

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(json)

        self._test(ctx, entities)


    def test_no_entities(self):
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"

        self.ctx_import.put_entities(entities.values())

        entities = {}

        json = {ContextGraphImport.K_CIS: [],
                ContextGraphImport.K_LINKS: []}

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(json)

        self._test(ctx, entities)

class AUpdateEntity(BaseTest):

    def setUp(self):
        super(AUpdateEntity, self).setUp()


    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "not_an_entity"

        desc = "The ci of id {0} does not match any existing entity.".format(
                ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(KeyError, desc):
            self.ctx_import._ContextGraphImport__a_update_entity(ci)

    def test_entities(self):
        ent = self.template_ent.copy()
        ent[ContextGraphImport.K_ID] = "ent1"
        self.ctx_import.entities_to_update["ent1"] = ent

        ci = self.template_ci.copy()
        expected =  self.template_ent.copy()

        expected["_id"] = ci[ContextGraphImport.K_ID] = "ent1"
        expected["name"] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected["type"] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected["infos"] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}


        self.ctx_import._ContextGraphImport__a_update_entity(ci)

        self.assertEqualEntities(self.ctx_import.update["ent1"], expected)


class ACreateEntity(BaseTest):

    def test_existing_entity(self):
        ent = self.template_ent.copy()
        ent[ContextGraphImport.K_ID] = "ent1"
        self.ctx_import.entities_to_update["ent1"] = ent

        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "ent1"

        desc = "The ci of id {0} match an existing entity.".format(
                ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(ValueError, desc):
            self.ctx_import._ContextGraphImport__a_create_entity(ci)

    def test_nonexistent_event(self):
        ci = self.template_ci.copy()
        expected =  self.template_ent.copy()

        expected["_id"] = ci[ContextGraphImport.K_ID] = "ent1"
        expected["name"] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected["type"] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected["infos"] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}

        self.ctx_import._ContextGraphImport__a_create_entity(ci)

        self.assertEqualEntities(self.ctx_import.update["ent1"], expected)

class ACreateLink(BaseTest):

    def setUp(self):
        super(ACreateLink, self).setUp()
    
    def test_create_link_e1_e2(self):
        self.ctx_import.update = {'e1':{'impact': []}, 'e2':{'depends': []}}
        self.ctx_import.__a_create_link({
            '_id':'e1-to-e2',
            'from': 'e1',
            'to': 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], ['e2'])
        self.assertEqual(self.ctx_import.update['e2']['depends'], ['e1'])

    def test_create_link_e1_e2(self):
        self.ctx_import.update = {'e2':{'depends': []}}
        self.ctx_import.entities_to_update = {'e1':{'impact': []}}
        self.ctx_import.__a_create_link({
            '_id':'e1-to-e2',
            'from': 'e1',
            'to': 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], ['e2'])
        self.assertEqual(self.ctx_import.update['e2']['depends'], ['e1'])



if __name__ == '__main__':
    main()
