#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.import_ctx import ContextGraphImport, ImportKey
from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware
import copy
import json
from jsonschema.exceptions import ValidationError


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
        self.import_storage = Middleware.get_middleware_by_uri(
            'storage-default-testimport://'
        )

        self.ctx_import[ContextGraph.ENTITIES_STORAGE] = self.entities_storage
        self.ctx_import[
            ContextGraph.ORGANISATIONS_STORAGE] = self.organisations_storage
        self.ctx_import[ContextGraph.USERS_STORAGE] = self.users_storage

        self.uuid = "test"

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
        self.import_storage.remove_elements()

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1["depends"])
        sorted(entity1["impact"])
        sorted(entity2["depends"])
        sorted(entity2["impact"])
        self.assertDictEqual(entity1, entity2)

    @classmethod
    def store_import(self, data, uuid):
        """Store the data in the right directory and with the right ID. The
        return the filename"""
        fname = ImportKey.IMPORT_FILE.format(uuid)
        with open(fname, "w") as fd:
            json.dump(data, fd)

        return fname


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

        self.ctx_import._put_entities(entities.values())

        for entity in entities.values():
            entity[ContextGraphImport.K_DEPENDS] = set(
                entity[ContextGraphImport.K_DEPENDS])
            entity[ContextGraphImport.K_IMPACT] = set(
                entity[ContextGraphImport.K_IMPACT])

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

        fname = self.store_import(json, self.uuid)

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(fname)

        self._test(ctx, entities)

    def test_no_entities(self):
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy()}

        entities["ent1"]["_id"] = "ent1"
        entities["ent2"]["_id"] = "ent2"
        entities["ent3"]["_id"] = "ent3"

        self.ctx_import._put_entities(entities.values())

        for entity in entities.values():
            entity[ContextGraphImport.K_DEPENDS] = set(
                entity[ContextGraphImport.K_DEPENDS])
            entity[ContextGraphImport.K_IMPACT] = set(
                entity[ContextGraphImport.K_IMPACT])

        entities = {}

        json = {ContextGraphImport.K_CIS: [],
                ContextGraphImport.K_LINKS: []}

        fname = self.store_import(json, self.uuid)

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(fname)

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

        self.ctx_import._put_entities(entities.values())

        for entity in entities.values():
            entity[ContextGraphImport.K_DEPENDS] = set(
                entity[ContextGraphImport.K_DEPENDS])
            entity[ContextGraphImport.K_IMPACT] = set(
                entity[ContextGraphImport.K_IMPACT])

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

        fname = self.store_import(json, self.uuid)

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(fname)
        self._test(ctx, entities)


class AUpdateEntity(BaseTest):

    def setUp(self):
        super(AUpdateEntity, self).setUp()

    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "not_an_entity"

        desc = "The ci of id {0} does not match any existing entity.".format(
            ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(ValueError, desc):
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

    def test_nonexistent_entity(self):
        ci = self.template_ci.copy()
        expected = self.template_ent.copy()

        expected["_id"] = ci[ContextGraphImport.K_ID] = "ent1"
        expected["name"] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected["type"] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected["infos"] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}

        expected[ContextGraphImport.K_DEPENDS] = set(
            expected[ContextGraphImport.K_DEPENDS])
        expected[ContextGraphImport.K_IMPACT] = set(
            expected[ContextGraphImport.K_IMPACT])

        self.ctx_import._ContextGraphImport__a_create_entity(ci)

        self.assertEqualEntities(self.ctx_import.update["ent1"], expected)

class ADisableEntity(BaseTest):

    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "ent1"

        desc = "The ci of id {0} does not match any existing entity.".format(
            ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(ValueError, desc):
            self.ctx_import._ContextGraphImport__a_disable_entity(ci)

    def test_entities_single_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = 12345
        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION:
              ContextGraphImport.A_DISABLE,
              ContextGraphImport.K_PROPERTIES: {
                  ContextGraphImport.K_DISABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = set(["1"])
        entities[id1]["impact"] = set(["2"])
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
        timestamp = [67890]

        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION: ContextGraphImport.A_DISABLE,
              ContextGraphImport.K_PROPERTIES:
              {ContextGraphImport.K_DISABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = set(["1"])
        entities[id1]["impact"] = set(["2"])
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {"disable": [12345]}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_disable_entity(ci)

        expected = entities[id1].copy()
        timestamp += [12345]
        expected["infos"] = {ContextGraphImport.K_DISABLE : sorted(timestamp)}

        result = self.ctx_import.update[id1]
        result["infos"]["disable"] = sorted(result["infos"]["disable"])

        self.assertEqualEntities(result, expected)

        with self.assertRaises(KeyError):
            self.ctx_import.update[id2]


class AEnableEntity(BaseTest):

    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "ent1"

        desc = "The ci of id {0} does not match any existing entity.".format(
            ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(ValueError, desc):
            self.ctx_import._ContextGraphImport__a_enable_entity(ci)


    def test_entities_single_timestamp(self):
        id1 = "ent1"
        id2 = "ent2"
        timestamp = 12345
        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION: ContextGraphImport.A_ENABLE,
              ContextGraphImport.K_PROPERTIES: {
                  ContextGraphImport.K_ENABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = set(["1"])
        entities[id1]["impact"] = set(["2"])
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
        timestamp = [67890]

        ci = {ContextGraphImport.K_ID: id1,
              ContextGraphImport.K_ACTION: ContextGraphImport.A_ENABLE,
              ContextGraphImport.K_PROPERTIES:
              {ContextGraphImport.K_ENABLE: timestamp}}

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),}

        entities[id1]["_id"] = id1
        entities[id1]["type"] = "resource"
        entities[id1]["name"] = "entity_1"
        entities[id1]["depends"] = set(["1"])
        entities[id1]["impact"] = set(["2"])
        entities[id1]["measurements"] = ["m1"]
        entities[id1]["infos"] = {"enable": [12345]}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_enable_entity(ci)

        expected = entities[id1].copy()
        timestamp += [12345]
        expected["infos"] = {ContextGraphImport.K_ENABLE : sorted(timestamp)}

        result = self.ctx_import.update[id1]
        result["infos"]["enable"] = sorted(result["infos"]["enable"])

        self.assertEqualEntities(result, expected)

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
        self.ctx_import.update = {'e1':{'impact': set([])},
                                  'e2':{'depends': set([])}}
        self.ctx_import._ContextGraphImport__a_create_link({
            ContextGraphImport.K_ID:'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], set(['e2']))
        self.assertEqual(self.ctx_import.update['e2']['depends'], set(['e1']))

    def test_create_link_e1_e2_2(self):
        self.ctx_import.update = {'e2':{'depends': set([])}}
        self.ctx_import.entities_to_update = {'e1':{'impact': set([])}}
        self.ctx_import._ContextGraphImport__a_create_link({
            ContextGraphImport.K_ID:'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1']['impact'], set(['e2']))
        self.assertEqual(self.ctx_import.update['e2']['depends'], set(['e1']))


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

class NotImplem(BaseTest):
    def update_link(self):
        self.assertRaises(NotImplementedError, self.ctx_import._ContextGraphImport__a_update_link())

    def disable_link(self):
        self.assertRaises(NotImplementedError, self.ctx_import._ContextGraphImport__a_disable_link())

    def enable_link(self):
        self.assertRaises(NotImplementedError, self.ctx_import._ContextGraphImport__a_enable_link())

class ADeleteEntity(BaseTest):

    def test_no_entities(self):
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "not_an_id"

        desc = "No entity found for the following id : {0}".format(
            ci[ContextGraphImport.K_ID])

        with self.assertRaisesRegexp(ValueError, desc):
            self.ctx_import._ContextGraphImport__a_delete_entity(ci)

    def test_entities(self):
        id1 = "ent1"
        id2 = "ent2"
        id3 = "ent3"
        id4 = "ent5"
        id5 = "ent5"
        id6 = "ent6"
        id7 = "ent7"

        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = id1
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DELETE

        entities = {id1: self.template_ent.copy(),
                    id2: self.template_ent.copy(),
                    id3: self.template_ent.copy(),
                    id4: self.template_ent.copy(),
                    id5: self.template_ent.copy(),
                    id6: self.template_ent.copy(),
                    id7: self.template_ent.copy()}

        entities[id1]["_id"] = id1
        entities[id1]["depends"] = [id2, id3, id4]
        entities[id1]["impact"] = [id5, id6, id7]

        entities[id2]["_id"] = id2
        entities[id2]["impact"] = [id1, "dummy0", "dummy1"]

        entities[id3]["_id"] = id3
        entities[id3]["impact"] = [id1, "dummy0"]

        entities[id4]["_id"] = id4
        entities[id4]["impact"] = [id1]

        entities[id5]["_id"] = id5
        entities[id5]["depends"] = [id1, "dummy0", "dummy1"]

        entities[id6]["_id"] = id6
        entities[id6]["depends"] = [id1, "dummy0"]

        entities[id7]["_id"] = id7
        entities[id7]["depends"] = [id1, "dummy0"]

        update_expected = copy.deepcopy(entities)

        del(update_expected[id1])
        update_expected[id2]["impact"].remove(id1)
        update_expected[id3]["impact"].remove(id1)
        update_expected[id4]["impact"].remove(id1)
        update_expected[id5]["depends"].remove(id1)
        update_expected[id6]["depends"].remove(id1)
        update_expected[id7]["depends"].remove(id1)

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_delete_entity(ci)

        for key in update_expected:
            try:
                entity = self.ctx_import.update[key]
            except KeyError:
                self.fail("Missing entity of id {0} in update.".format(key))

            self.assertEqualEntities(entity, update_expected[key])

        delete_expected = [id1]
        self.assertListEqual(self.ctx_import.delete, delete_expected)


    def test_delete_entities_and_related_entities(self):
        """Remove an entity and later delete an entity with a links to the first
        entity."""
        deleted_id = "deleted_ent"
        id_ = "id1"

        entity = self.template_ent.copy()
        entity["_id"] = id_
        entity["impact"] = [deleted_id]

        self.ctx_import.entities_to_update = {id_: entity}

        self.ctx_import.delete = [deleted_id]

        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = id_
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DELETE
        ci[ContextGraphImport.K_TYPE] = "resource"
        ci[ContextGraphImport.K_IMPACT] = [deleted_id]

        try:
            self.ctx_import._ContextGraphImport__a_delete_entity(ci)
        except Exception, e:
            self.fail("The error \"{0}\" was raised".format(e))

        self.assertListEqual(self.ctx_import.delete, [deleted_id, id_])
        self.assertDictEqual(self.ctx_import.update, {})

        self.ctx_import.delete = [deleted_id]
        ci[ContextGraphImport.K_IMPACT] = []
        ci[ContextGraphImport.K_DEPENDS] = [deleted_id]
        try:
            self.ctx_import._ContextGraphImport__a_delete_entity(ci)
        except Exception, e:
            self.fail("The error \"{0}\" was raised".format(e))

        self.ctx_import.update = {}
        self.ctx_import.delete = [deleted_id]
        try:
            self.ctx_import._ContextGraphImport__a_delete_entity(ci)
        except Exception, e:
            self.fail("The error \"{0}\" was raised".format(e))
        self.assertListEqual(self.ctx_import.delete, [deleted_id, id_])

        self.assertDictEqual(self.ctx_import.update, {})


        self.ctx_import.update = {}
        self.ctx_import.delete = [deleted_id]

        ci[ContextGraphImport.K_IMPACT] = [deleted_id]
        ci[ContextGraphImport.K_DEPENDS] = []
        try:
            self.ctx_import._ContextGraphImport__a_delete_entity(ci)
        except Exception, e:
            self.fail("The error \"{0}\" was raised".format(e))
        self.assertListEqual(self.ctx_import.delete, [deleted_id, id_])

        self.assertDictEqual(self.ctx_import.update, {})


class ImportChecker(TestCase):
    """I only check a kind of error on one fields, not every kind of error on
    every fields. I assume that the error will be triggered whatever the fields
    are.
    """

    def setUp(self):

        self.template_ci = {ContextGraphImport.K_ID: "id",
                            ContextGraphImport.K_NAME: "name",
                            ContextGraphImport.K_TYPE: "resource",
                            ContextGraphImport.K_DEPENDS: [],
                            ContextGraphImport.K_IMPACT: [],
                            ContextGraphImport.K_MEASUREMENTS: [],
                            ContextGraphImport.K_INFOS: {},
                            ContextGraphImport.K_ACTION:
                            ContextGraphImport.A_CREATE,
                            ContextGraphImport.K_PROPERTIES: {}}

        self.template_link = {ContextGraphImport.K_ID: "id",
                              ContextGraphImport.K_FROM: "from",
                              ContextGraphImport.K_TO: "to",
                              ContextGraphImport.K_INFOS: {},
                              ContextGraphImport.K_ACTION:
                              ContextGraphImport.A_CREATE,
                              ContextGraphImport.K_PROPERTIES: {}}

        self.template_json = {ContextGraphImport.K_CIS: [],
                              ContextGraphImport.K_LINKS: []}

        self._desc_fail = "import_checker() raise an exception {0}!"

    def test_empty_import(self):
        with self.assertRaises(KeyError):
            import_checker({})

    def test_cis_links(self):
        json = self.template_json
        # check cis with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check cis with wrong type
        json[ContextGraphImport.K_CIS] = {}
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check links with wrong type
        json[ContextGraphImport.K_CIS] = []
        json[ContextGraphImport.K_LINKS] = {}
        with self.assertRaises(ValidationError):
            import_checker(json)

    def test_ci_id(self):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        # check ci.id with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.id with wrong type
        json[ContextGraphImport.K_CIS][0][ContextGraphImport.K_ID] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check missing ci.id
        json[ContextGraphImport.K_CIS][0].pop(ContextGraphImport.K_ID)
        with self.assertRaises(ValidationError):
            import_checker(json)

    def test_ci_name(self):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        # check ci.name with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.name with wrong type
        json[ContextGraphImport.K_CIS][0][ContextGraphImport.K_NAME] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

    def _test_ci_array_string(self, key):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        # check ci.depends with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.depends with wrong type
        json[ContextGraphImport.K_CIS][0][key] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0][key] = [1,2]
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0][key] = [1,"ok"]
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check with out ci.depends
        json[ContextGraphImport.K_CIS][0].pop(key)
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

    def test_ci_depends(self):
        self._test_ci_array_string(ContextGraphImport.K_DEPENDS)

    def test_ci_impact(self):
        self._test_ci_array_string(ContextGraphImport.K_IMPACT)

    def test_ci_measurements(self):
        self._test_ci_array_string(ContextGraphImport.K_MEASUREMENTS)

    def _test_ci_object(self, key, required=False):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        # check ci.{key} with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.depends with wrong type
        json[ContextGraphImport.K_CIS][0][key] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check with out ci.{key}
        if required is True:
            json[ContextGraphImport.K_CIS][0].pop(key)
            try:
                import_checker(json)
            except Exception as e:
                self.fail(self._desc_fail.format(e))

    def test_ci_infos(self):
        self._test_ci_object(ContextGraphImport.K_INFOS, False)

    def _test_action(self, key):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]
        json[ContextGraphImport.K_CIS][0][ContextGraphImport.K_PROPERTIES] \
         = {ContextGraphImport.A_DISABLE: [],
            ContextGraphImport.A_ENABLE: []}

        # check ci.action with good action.
        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = ContextGraphImport.A_CREATE
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = ContextGraphImport.A_DELETE
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = ContextGraphImport.A_DISABLE
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = ContextGraphImport.A_ENABLE
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = ContextGraphImport.A_UPDATE
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.action with an action that did not match the pattern
        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "update_not"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "not_update"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "Update"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "a strange action"
        with self.assertRaises(ValidationError):
            import_checker(json)

    def test_ci_action(self):
        self._test_action(ContextGraphImport.K_CIS)

    def test_ci_type(self):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        # check ci.action with good action.
        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_TYPE] = "resource"
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_TYPE] = "component"
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_TYPE] = "connector"
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.action with an action that did not match the pattern
        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "resource_not"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "not_resource"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "Resource"
        with self.assertRaises(ValidationError):
            import_checker(json)

        json[ContextGraphImport.K_CIS][0]\
            [ContextGraphImport.K_ACTION] = "a strange resource"
        with self.assertRaises(ValidationError):
            import_checker(json)

    def test_ci_properties(self):
        self._test_ci_object(ContextGraphImport.K_PROPERTIES, False)

    def _test_link_string(self, key, required=False):
        json = self.template_json.copy()
        json[ContextGraphImport.K_LINKS] = [self.template_link.copy()]

        # check link.{key} with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check link.{key} with wrong type
        json[ContextGraphImport.K_LINKS][0][key] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check missing link.{key}
        if required is True:
            json[ContextGraphImport.K_LINKS][0].pop(key)
            with self.assertRaises(ValidationError):
                import_checker(json)

    def test_link_id(self):
        self._test_link_string(ContextGraphImport.K_ID, required=True)

    def test_link_from(self):
        self._test_link_string(ContextGraphImport.K_FROM, required=True)

    def test_link_to(self):
        self._test_link_string(ContextGraphImport.K_TO, required=True)

    def _test_link_object(self, key, required=False):
        json = self.template_json.copy()
        json[ContextGraphImport.K_LINKS] = [self.template_link.copy()]

        # check ci.{key} with right type
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # check ci.depends with wrong type
        json[ContextGraphImport.K_LINKS][0][key] = 1
        with self.assertRaises(ValidationError):
            import_checker(json)

        # check with out ci.{key}
        if required is True:
            json[ContextGraphImport.K_LINKS][0].pop(key)
            try:
                import_checker(json)
            except Exception as e:
                self.fail(self._desc_fail.format(e))

    def test_link_infos(self):
        self._test_link_object(ContextGraphImport.K_INFOS, required=False)

    def test_link_action(self):
        self._test_action(ContextGraphImport.K_LINKS)

    def test_link_properties(self):
        self._test_link_object(ContextGraphImport.K_PROPERTIES, required=False)

    def _test_state(self, object_, state):
        json = self.template_json.copy()

        if object_ == ContextGraphImport.K_CIS:
            obj = self.template_ci.copy()
            json[ContextGraphImport.K_CIS] = [obj]
            obj_key = ContextGraphImport.K_LINKS
        elif object_ == ContextGraphImport.K_LINKS:
            obj = self.template_link.copy()
            json[ContextGraphImport.K_LINKS] = [obj]
            obj_key = ContextGraphImport.K_LINKS
        else:
            self.fail("Unrecognized object_ {0}".format(object_))

        if state == ContextGraphImport.K_ENABLE:
            other_state = ContextGraphImport.K_DISABLE
        elif state == ContextGraphImport.K_DISABLE:
            other_state = ContextGraphImport.K_ENABLE
        else:
            self.fail("Unrecognized state {0}".format(object_))

        # {object_}.action : {state} without {object_}.properties.{state}
        obj[ContextGraphImport.K_ACTION] = state
        with self.assertRaises(KeyError):
            import_checker(json)

        # {object_}.action : {state} with {object_}.properties.{state}
        obj[ContextGraphImport.K_PROPERTIES] = {state: []}
        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

        # {object_}.action : {state} with ci.properties.{other_state} and
        # without {object_}.properties.{state}
        obj[ContextGraphImport.K_PROPERTIES] = {other_state: []}
        with self.assertRaises(KeyError):
            import_checker(json)

    def test_ci_disable_with_properties(self):
        self._test_state(ContextGraphImport.K_CIS, ContextGraphImport.A_DISABLE)

    def test_link_disable_with_properties(self):
        self._test_state(ContextGraphImport.K_LINKS,
                         ContextGraphImport.A_DISABLE)

    def test_ci_enable_with_properties(self):
        self._test_state(ContextGraphImport.K_CIS,
                         ContextGraphImport.A_ENABLE)

    def test_link_enable_with_properties(self):
        self._test_state(ContextGraphImport.K_LINKS,
                         ContextGraphImport.A_ENABLE)

    def test_OK_single_ci(self):
        json = self.template_json.copy()
        json[ContextGraphImport.K_CIS] = [self.template_ci.copy()]

        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

    def test_OK_multiple_ci(self):
        json = self.template_json.copy()

        ci0 = self.template_ci.copy()
        ci0[ContextGraphImport.K_ID] = "id0"
        ci1 = self.template_ci.copy()
        ci1[ContextGraphImport.K_ID] = "id1"
        ci2 = self.template_ci.copy()
        ci2[ContextGraphImport.K_ID] = "id2"

        json[ContextGraphImport.K_CIS] = [ci0, ci1, ci2]

        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

    def test_OK_single_link(self):
        json = self.template_json.copy()
        json[ContextGraphImport.K_LINKS] = [self.template_link.copy()]

        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

    def test_OK_multiple_link(self):
        json = self.template_json.copy()

        link0 = self.template_link.copy()
        link0[ContextGraphImport.K_ID] = "id0"
        link1 = self.template_link.copy()
        link1[ContextGraphImport.K_ID] = "id1"
        link2 = self.template_link.copy()
        link2[ContextGraphImport.K_ID] = "id2"

        json[ContextGraphImport.K_LINKS] = [link0, link1, link2]

        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))

    def test_OK_both(self):
        json = self.template_json.copy()

        link0 = self.template_link.copy()
        link0[ContextGraphImport.K_ID] = "id0"
        link1 = self.template_link.copy()
        link1[ContextGraphImport.K_ID] = "id1"
        link2 = self.template_link.copy()
        link2[ContextGraphImport.K_ID] = "id2"

        ci0 = self.template_ci.copy()
        ci0[ContextGraphImport.K_ID] = "id0"
        ci1 = self.template_ci.copy()
        ci1[ContextGraphImport.K_ID] = "id1"
        ci2 = self.template_ci.copy()
        ci2[ContextGraphImport.K_ID] = "id2"

        json[ContextGraphImport.K_LINKS] = [link0, link1, link2]
        json[ContextGraphImport.K_CIS] = [ci0, ci1, ci2]

        try:
            import_checker(json)
        except Exception as e:
            self.fail(self._desc_fail.format(e))


class TestImportFunctions(BaseTest):

    def test_ongoing(self):
        self.assertEqual(self.ctx_import.on_going_in_db(), False)
        self.ctx_import[ContextGraph.IMPORT_STORAGE].put_element({'_id':'id', 'state': 'on going'})
        self.assertEqual(self.ctx_import.on_going_in_db(), True)

    def check_id(self):
        self.assertEqual(self.ctx_import.check_id('id'), False)
        self.ctx_import[ContextGraph.IMPORT_STORAGE].put_element({'_id':'id', 'state': 'on going'})
        self.assertEqual(self.ctx_import.check_id('id'), True)

    def getimporstatus(self):
        self.ctx_import[ContextGraph.IMPORT_STORAGE].put_element({'_id':'id', 'state': 'on going'})
        self.assertEqual(self.ctx_import.get_import_status('id'), 'on going')

if __name__ == '__main__':
    main()
