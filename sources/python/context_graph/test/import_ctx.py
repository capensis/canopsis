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
                self.fail("KeyError : missing key {0} in ctx".format(id_))
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

        json = {ContextGraphImport.K_CIS: [{ContextGraphImport.K_ID: "ent1",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_DELETE},
                                           {ContextGraphImport.K_ID: "ent2",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_DELETE},
                                           {ContextGraphImport.K_ID: "ent3",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_DELETE}],
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

    def test_action_delete(self):
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy(),
                    "ent4": self.template_ent.copy(),
                    "ent5": self.template_ent.copy(),
                    "ent6": self.template_ent.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"
        entities["ent3"]["depends"] = ["ent4"]
        entities["ent3"]["impact"] = ["ent5", "ent6"]
        entities["ent4"]["_id"] = "ent4"
        entities["ent5"]["_id"] = "ent5"
        entities["ent6"]["_id"] = "ent6"

        self.ctx_import.put_entities(entities.values())

        json = {ContextGraphImport.K_CIS: [{ContextGraphImport.K_ID: "ent1",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_UPDATE},
                                           {ContextGraphImport.K_ID: "ent2",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_UPDATE},
                                           {ContextGraphImport.K_ID: "ent3",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_DELETE}],
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
        expected = self.template_ent.copy()

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
        expected = self.template_ent.copy()

        expected["_id"] = ci[ContextGraphImport.K_ID] = "ent1"
        expected["name"] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected["type"] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected["infos"] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}

        self.ctx_import._ContextGraphImport__a_create_entity(ci)

        self.assertEqualEntities(self.ctx_import.update["ent1"], expected)

class ADisableEntity(BaseTest):

    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "ent1"

        desc = "The ci of id {0} does not match any existing entity.".format(
            ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(KeyError, desc):
            self.ctx_import._ContextGraphImport__a_disable_entity(ci)

    def test_entities_single_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = 12345
        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION:
              ContextGraphImport.A_DISABLE,
              ContextGraphImport.K_INFOS: {
                  ContextGraphImport.K_DISABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = ["1"]
        entities[id1]["impact"] = ["2"]
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_disable_entity(ci)

        expected = entities[id1].copy()
        expected["infos"] = {ContextGraphImport.K_DISABLE : [timestamp]}

        self.assertEqualEntities(self.ctx_import.update[id1], expected)

        with self.assertRaises(KeyError):
            self.ctx_import.update[id2]

    def test_entities_multiple_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = [12345, 67890]

        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION:
              ContextGraphImport.A_DISABLE,
              ContextGraphImport.K_INFOS: {
                  ContextGraphImport.K_DISABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = ["1"]
        entities[id1]["impact"] = ["2"]
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_disable_entity(ci)

        expected = entities[id1].copy()
        expected["infos"] = {ContextGraphImport.K_DISABLE : timestamp}

        self.assertEqualEntities(self.ctx_import.update[id1], expected)

        with self.assertRaises(KeyError):
            self.ctx_import.update[id2]


def AEnableEntity(BaseTest):

    def test_entities_single_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = 12345
        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION: ContextGraphImport.A_ENABLE,
              ContextGraphImport.K_INFOS: {
                  ContextGraphImport.K_ENABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = ["1"]
        entities[id1]["impact"] = ["2"]
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_enable_entity(ci)

        expected = entities[id1].copy()
        expected["infos"] = {ContextGraphImport.K_ENABLE : [timestamp]}

        self.assertEqualEntities(self.ctx_import.update[id1], expected)

        with self.assertRaises(KeyError):
            self.ctx_import.update[id2]

    def test_entities_multiple_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = [12345, 67890]

        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION: ContextGraphImport.A_ENABLE,
              ContextGraphImport.K_INFOS: {
                  ContextGraphImport.K_ENABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = ["1"]
        entities[id1]["impact"] = ["2"]
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_enable_entity(ci)

        expected = entities[id1].copy()
        expected["infos"] = {ContextGraphImport.K_ENABLE : timestamp}

        self.assertEqualEntities(self.ctx_import.update[id1], expected)

        with self.assertRaises(KeyError):
            self.ctx_import.update[id2]

class ChangeStateEntity(BaseTest):

    def test(self):
        state = "not_a_state"
        desc = "{0} is not a valid state.".format(state)
        with self.assertRaisesRegexp(ValueError, desc):
            self.ctx_import._ContextGraphImport__change_state_entity(None,
                                                                     state)
class ACreateLink(BaseTest):
    def test_create_link_e1_e2(self):
        self.ctx_import.update = {'e1':{'impact': []}, 'e2':{'depends': []}}
        self.ctx_import._ContextGraphImport__a_create_link({
            '_id':'e1-to-e2',
            'from': 'e1',
            'to': 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], ['e2'])
        self.assertEqual(self.ctx_import.update['e2']['depends'], ['e1'])

    def test_create_link_e1_e2_2(self):
        self.ctx_import.update = {'e2':{'depends': []}}
        self.ctx_import.entities_to_update = {'e1':{'impact': []}}
        self.ctx_import._ContextGraphImport__a_create_link({
            '_id':'e1-to-e2',
            'from': 'e1',
            'to': 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], ['e2'])
        self.assertEqual(self.ctx_import.update['e2']['depends'], ['e1'])


class ADeleteLink(BaseTest):
    def test_delete__link_e1_e2(self):
        self.ctx_import.update = {'e1':{'impact': ['e2']}, 'e2':{'depends': ['e1']}}
        self.ctx_import._ContextGraphImport__a_delete_link(
            {'_id': 'e1-to-e2', 'from': 'e1', 'to': 'e2'}
        )
        self.assertEqual(self.ctx_import.update['e1']['impact'], [])
        self.assertEqual(self.ctx_import.update['e2']['depends'], [])

    def test_delete_link_e1_e2_2(self):
        self.ctx_import.update = {'e2':{'depends': ['e1']}}
        self.ctx_import.entities_to_update = {'e1':{'impact': ['e2']}}
        self.ctx_import._ContextGraphImport__a_delete_link(
            {'_id': 'e1-to-e2', 'from': 'e1', 'to': 'e2'}
        )
        self.assertEqual(self.ctx_import.update['e1']['impact'], [])
        self.assertEqual(self.ctx_import.update['e2']['depends'], [])

class AUpdateLink(BaseTest):

    def test(self):
        with self.assertRaises(NotImplementedError):
            self.ctx_import._ContextGraphImport__a_update_link(None)

class ADisableLink(BaseTest):

    def test(self):
        with self.assertRaises(NotImplementedError):
            self.ctx_import._ContextGraphImport__a_disable_link(None)

class AEnableLink(BaseTest):

    def test(self):
        with self.assertRaises(NotImplementedError):
            self.ctx_import._ContextGraphImport__a_enable_link(None)


if __name__ == '__main__':
    main()
