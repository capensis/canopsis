#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.import_ctx import ContextGraphImport, ImportKey,\
    Manager
from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware
import copy
import json
import os
import time
from jsonschema.exceptions import ValidationError

class Keys:

    ID = "_id"
    TYPE = "type"
    NAME = "name"
    DEPENDS = "depends"
    IMPACT = "impact"
    MEASUREMENTS = "measurements"
    INFOS = "infos"
    DISABLE_HISTORY = "disable_history"
    ENABLE_HISTORY = "enable_history"
    ENABLED = "enabled"


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

        self.uuid = "test"

        self.template_ent = {Keys.ID: None,
                             Keys.TYPE: 'connector',
                             Keys.NAME: 'conn-name1',
                             Keys.DEPENDS: [],
                             Keys.IMPACT: [],
                             Keys.MEASUREMENTS: [],
                             Keys.INFOS: {}}
        self.template_ci = {ContextGraphImport.K_ID: "id",
                            ContextGraphImport.K_NAME: "Name",
                            ContextGraphImport.K_TYPE: "resource",
                            ContextGraphImport.K_INFOS: {},
                            ContextGraphImport.K_ACTION:
                            ContextGraphImport.A_CREATE}

        self.template_link = {ContextGraphImport.K_FROM: None,
                              ContextGraphImport.K_TO: None,
                              ContextGraphImport.K_INFOS: {},
                              ContextGraphImport.K_ACTION:
                              ContextGraphImport.A_CREATE}

        self.template_json = {ContextGraphImport.K_CIS: [],
                              ContextGraphImport.K_LINKS: []}

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.users_storage.remove_elements()
        try:
            os.remove(ImportKey.IMPORT_FILE.format(self.uuid))
        except:
            pass

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1[Keys.DEPENDS])
        sorted(entity1[Keys.IMPACT])
        sorted(entity2[Keys.DEPENDS])
        sorted(entity2[Keys.IMPACT])
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

        entities["ent1"][Keys.ID] = "ent1"
        entities["ent2"][Keys.ID] = "ent2"
        entities["ent3"][Keys.ID] = "ent3"

        self.entities_storage.put_elements(entities.values())

        for entity in entities.values():
            entity[ContextGraphImport.K_DEPENDS] = set(
                entity[ContextGraphImport.K_DEPENDS])
            entity[ContextGraphImport.K_IMPACT] = set(
                entity[ContextGraphImport.K_IMPACT])

        data = self.template_json.copy()

        cis = [self.template_ci.copy(),
               self.template_ci.copy(),
               self.template_ci.copy()]

        links = [{ContextGraphImport.K_FROM: ["ent1"],
                  ContextGraphImport.K_TO: "ent2",
                  ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE},
                 {ContextGraphImport.K_FROM: ["ent1"],
                  ContextGraphImport.K_TO: "ent3",
                  ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE}]

        for i in range(len(cis)):
            cis[i][ContextGraphImport.K_ACTION] = ContextGraphImport.A_DELETE
            cis[i][ContextGraphImport.K_TYPE] = ContextGraphImport.RESOURCE
            cis[i][ContextGraphImport.K_ID] = "ent{0}".format(i + 1)

        data[ContextGraphImport.K_CIS] = cis
        data[ContextGraphImport.K_LINKS] = links
        fname = self.store_import(data, self.uuid)

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(
            fname)

        self._test(ctx, entities)

    def test_no_entities(self):
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy()}

        entities["ent1"][Keys.ID] = "ent1"
        entities["ent2"][Keys.ID] = "ent2"
        entities["ent3"][Keys.ID] = "ent3"

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

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(
            fname)

        self._test(ctx, entities)

    def test_action_delete(self):
        entities = {"ent1": self.template_ent.copy(),
                    "ent2": self.template_ent.copy(),
                    "ent3": self.template_ent.copy(),
                    "ent4": self.template_ent.copy(),
                    "ent5": self.template_ent.copy(),
                    "ent6": self.template_ent.copy()}

        entities["ent1"][Keys.ID] = "ent1"
        entities["ent2"][Keys.ID] = "ent2"
        entities["ent3"][Keys.ID] = "ent3"
        entities["ent3"][Keys.DEPENDS] = ["ent4"]
        entities["ent3"][Keys.IMPACT] = ["ent5", "ent6"]
        entities["ent4"][Keys.ID] = "ent4"
        entities["ent5"][Keys.ID] = "ent5"
        entities["ent6"][Keys.ID] = "ent6"

        self.ctx_import._put_entities(entities.values())

        for entity in entities.values():
            entity[ContextGraphImport.K_DEPENDS] = set(
                entity[ContextGraphImport.K_DEPENDS])
            entity[ContextGraphImport.K_IMPACT] = set(
                entity[ContextGraphImport.K_IMPACT])

        json = {ContextGraphImport.K_CIS: [{ContextGraphImport.K_ID: "ent1",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_UPDATE,
                                            ContextGraphImport.K_TYPE:
                                            ContextGraphImport.RESOURCE},
                                           {ContextGraphImport.K_ID: "ent2",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_UPDATE,
                                            ContextGraphImport.K_TYPE:
                                            ContextGraphImport.RESOURCE},
                                           {ContextGraphImport.K_ID: "ent3",
                                            ContextGraphImport.K_ACTION:
                                            ContextGraphImport.A_DELETE,
                                            ContextGraphImport.K_TYPE:
                                            ContextGraphImport.RESOURCE}],
                ContextGraphImport.K_LINKS: []}

        fname = self.store_import(json, self.uuid)

        ctx = self.ctx_import._ContextGraphImport__get_entities_to_update(
            fname)
        self._test(ctx, entities)


class AUpdateEntity(BaseTest):

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

        expected[Keys.ID] = ci[ContextGraphImport.K_ID] = "ent1"
        expected[Keys.NAME] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected[Keys.TYPE] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected[Keys.INFOS] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}

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

        expected[Keys.ID] = ci[ContextGraphImport.K_ID] = "ent1"
        expected[Keys.NAME] = ci[ContextGraphImport.K_NAME] = "other_name"
        expected[Keys.TYPE] = ci[ContextGraphImport.K_TYPE] = "other_type"
        expected[Keys.INFOS] = ci[ContextGraphImport.K_INFOS] = {"infos": "infos"}

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
                    id2: self.template_ent.copy(), }

        entities[id1][Keys.ID] = id1
        entities[id1][Keys.TYPE] = "resource"
        entities[id1][Keys.NAME] = "entity_1"
        entities[id1][Keys.DEPENDS] = set(["1"])
        entities[id1][Keys.IMPACT] = set(["2"])
        entities[id1][Keys.MEASUREMENTS] = ["m1"]
        entities[id1][Keys.INFOS] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_disable_entity(ci)

        expected = entities[id1].copy()
        expected[Keys.INFOS] = {Keys.DISABLE_HISTORY: [timestamp],
                                Keys.ENABLED: False}

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
                    id2: self.template_ent.copy(), }

        entities[id1][Keys.ID] = id1
        entities[id1][Keys.TYPE] = "resource"
        entities[id1][Keys.NAME] = "entity_1"
        entities[id1][Keys.DEPENDS] = set(["1"])
        entities[id1][Keys.IMPACT] = set(["2"])
        entities[id1][Keys.MEASUREMENTS] = ["m1"]
        entities[id1][Keys.INFOS] = {Keys.DISABLE_HISTORY: [12345],
                                Keys.ENABLED: False}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_disable_entity(ci)

        expected = entities[id1].copy()
        timestamp += [12345]
        expected[Keys.INFOS] = {Keys.DISABLE_HISTORY: sorted(timestamp),
                                Keys.ENABLED: False}

        result = self.ctx_import.update[id1]
        result[Keys.INFOS][Keys.DISABLE_HISTORY] = sorted(result[Keys.INFOS][Keys.DISABLE_HISTORY])

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
                    id2: self.template_ent.copy(), }

        entities[id1][Keys.ID] = id1
        entities[id1][Keys.TYPE] = "resource"
        entities[id1][Keys.NAME] = "entity_1"
        entities[id1][Keys.DEPENDS] = set(["1"])
        entities[id1][Keys.IMPACT] = set(["2"])
        entities[id1][Keys.MEASUREMENTS] = ["m1"]
        entities[id1][Keys.INFOS] = {}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_enable_entity(ci)

        expected = entities[id1].copy()
        expected[Keys.INFOS] = {Keys.ENABLE_HISTORY: [timestamp],
                                Keys.ENABLED: True}

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
                    id2: self.template_ent.copy(), }

        entities[id1][Keys.ID] = id1
        entities[id1][Keys.TYPE] = "resource"
        entities[id1][Keys.NAME] = "entity_1"
        entities[id1][Keys.DEPENDS] = set(["1"])
        entities[id1][Keys.IMPACT] = set(["2"])
        entities[id1][Keys.MEASUREMENTS] = ["m1"]
        entities[id1][Keys.INFOS] = {Keys.ENABLE_HISTORY: [12345]}

        self.ctx_import.entities_to_update = entities

        self.ctx_import._ContextGraphImport__a_enable_entity(ci)

        expected = entities[id1].copy()
        timestamp += [12345]
        expected[Keys.INFOS] = {Keys.ENABLE_HISTORY: sorted(timestamp),
                                Keys.ENABLED: True}

        result = self.ctx_import.update[id1]
        result[Keys.INFOS][Keys.ENABLE_HISTORY] = sorted(result[Keys.INFOS][Keys.ENABLE_HISTORY])

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
        self.ctx_import.update = {'e1': {Keys.IMPACT: set()},
                                  'e2': {Keys.DEPENDS: set()}}
        self.ctx_import._ContextGraphImport__a_create_link({
            ContextGraphImport.K_ID: 'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1'][Keys.IMPACT], set(['e2']))
        self.assertEqual(self.ctx_import.update['e2'][Keys.DEPENDS], set(['e1']))

    def test_create_link_e1_e2_2(self):
        self.ctx_import.update = {'e2': {Keys.DEPENDS: set()}}
        self.ctx_import.entities_to_update = {'e1': {Keys.IMPACT: set()}}
        self.ctx_import._ContextGraphImport__a_create_link({
            ContextGraphImport.K_ID: 'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'
        })
        self.assertEqual(self.ctx_import.update['e1'][Keys.IMPACT], set(['e2']))
        self.assertEqual(self.ctx_import.update['e2'][Keys.DEPENDS], set(['e1']))


class ADeleteLink(BaseTest):
    def test_delete_link_e1_e2(self):
        self.ctx_import.update = {'e1': {Keys.IMPACT: set(['e2'])},
                                  'e2': {Keys.DEPENDS: set(['e1'])}}
        self.ctx_import._ContextGraphImport__a_delete_link({
            ContextGraphImport.K_ID: 'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'})
        self.assertEqual(self.ctx_import.update['e1'][Keys.IMPACT], set())
        self.assertEqual(self.ctx_import.update['e2'][Keys.DEPENDS], set())

    def test_delete_link_e1_e2_2(self):
        self.ctx_import.update = {'e2': {Keys.DEPENDS: set(['e1'])}}
        self.ctx_import.entities_to_update = {'e1': {Keys.IMPACT: set(['e2'])}}
        self.ctx_import._ContextGraphImport__a_delete_link({
            ContextGraphImport.K_ID: 'e1-to-e2',
            ContextGraphImport.K_FROM: ['e1'],
            ContextGraphImport.K_TO: 'e2'})
        self.assertEqual(self.ctx_import.update['e1'][Keys.IMPACT], set())
        self.assertEqual(self.ctx_import.update['e2'][Keys.DEPENDS], set())


class NotImplem(BaseTest):
    def update_link(self):
        self.assertRaises(NotImplementedError,
                          self.ctx_import._ContextGraphImport__a_update_link())

    def disable_link(self):
        self.assertRaises(NotImplementedError,
                          self.ctx_import._ContextGraphImport__a_disable_link())

    def enable_link(self):
        self.assertRaises(NotImplementedError,
                          self.ctx_import._ContextGraphImport__a_enable_link())


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

        entities[id1][Keys.ID] = id1
        entities[id1][Keys.DEPENDS] = [id2, id3, id4]
        entities[id1][Keys.IMPACT] = [id5, id6, id7]

        entities[id2][Keys.ID] = id2
        entities[id2][Keys.IMPACT] = [id1, "dummy0", "dummy1"]

        entities[id3][Keys.ID] = id3
        entities[id3][Keys.IMPACT] = [id1, "dummy0"]

        entities[id4][Keys.ID] = id4
        entities[id4][Keys.IMPACT] = [id1]

        entities[id5][Keys.ID] = id5
        entities[id5][Keys.DEPENDS] = [id1, "dummy0", "dummy1"]

        entities[id6][Keys.ID] = id6
        entities[id6][Keys.DEPENDS] = [id1, "dummy0"]

        entities[id7][Keys.ID] = id7
        entities[id7][Keys.DEPENDS] = [id1, "dummy0"]

        update_expected = copy.deepcopy(entities)

        del(update_expected[id1])
        update_expected[id2][Keys.IMPACT].remove(id1)
        update_expected[id3][Keys.IMPACT].remove(id1)
        update_expected[id4][Keys.IMPACT].remove(id1)
        update_expected[id5][Keys.DEPENDS].remove(id1)
        update_expected[id6][Keys.DEPENDS].remove(id1)
        update_expected[id7][Keys.DEPENDS].remove(id1)

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
        entity[Keys.ID] = id_
        entity[Keys.IMPACT] = [deleted_id]

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


class ImportChecker(BaseTest):
    """I only check a kind of error on one fields, not every kind of error on
    every fields. I assume that the error will be triggered whatever the fields
    are.
    """

    def put_dummy_entity(self, id_):
        entity = self.template_ent.copy()
        entity[Keys.ID] = id_
        self.ctx_import._put_entities(entity)

    def setUp(self):
        super(ImportChecker, self).setUp()

        self.uuid = "test"
        self._desc_fail = "The check of the import raise an exception {0}!"

    def test_empty_import(self):
        self.store_import({}, self.uuid)
        with self.assertRaisesRegexp(ValidationError,
                                     "CIS and LINKS should be an array."):
            self.ctx_import.import_context(self.uuid)

    def test_cis_links(self):
        data = self.template_json.copy()
        self.store_import(data, self.uuid)
        # check cis with right type
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        # check cis with wrong type
        data[ContextGraphImport.K_CIS] = {}
        self.store_import(data, self.uuid)
        with self.assertRaisesRegexp(ValidationError,
                                     "CIS should be an array."):
            self.ctx_import.import_context(self.uuid)

        # check links with wrong type
        data[ContextGraphImport.K_CIS] = []
        data[ContextGraphImport.K_LINKS] = {}
        self.store_import(data, self.uuid)
        with self.assertRaisesRegexp(ValidationError,
                                     "LINKS should be an array."):
            self.ctx_import.import_context(self.uuid)

        # check links and cis with wrong type
        data[ContextGraphImport.K_CIS] = {}
        data[ContextGraphImport.K_LINKS] = {}
        self.store_import(data, self.uuid)
        with self.assertRaisesRegexp(ValidationError,
                                     "CIS and LINKS should be an array."):
            self.ctx_import.import_context(self.uuid)

    def test_ci_id_ok(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = "id"
        data[ContextGraphImport.K_CIS] = [ci]

        # check ci.id with right type
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_ci_id_wrong_type(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_ID] = 1
        data[ContextGraphImport.K_CIS] = [ci]

        # check ci.id with wrong type
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

    def test_ci_id_no_id(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        ci.pop(ContextGraphImport.K_ID)
        data[ContextGraphImport.K_CIS] = [ci]

        # check missing ci.id
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

    def test_ci_name_ok(self):
        data = self.template_json.copy()
        data[ContextGraphImport.K_CIS] = [self.template_ci.copy()]
        self.store_import(data, self.uuid)

        # check ci.name with right type
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_ci_name_wrong_type(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        ci[ContextGraphImport.K_NAME] = 1
        data[ContextGraphImport.K_CIS] = [ci]
        self.store_import(data, self.uuid)

        # check ci.name with wrong type
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

    def test_ci_measurements_empty(self):
        data = self.template_json.copy()
        data[ContextGraphImport.K_CIS] = [self.template_ci.copy()]
        self.store_import(data, self.uuid)

        # check ci.key with right type
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_ci_measurements(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        ci[ContextGraphImport.K_MEASUREMENTS] = ["measurmement1",
                                                 "measurmement2"]
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_ci_measurements_wrong_type(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        ci[ContextGraphImport.K_MEASUREMENTS] = 1
        self.store_import(data, self.uuid)

        # check ci.key with wrong type
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        ci[ContextGraphImport.K_MEASUREMENTS] = [1, 2]
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        ci[ContextGraphImport.K_MEASUREMENTS] = [1, "ok"]
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

    def test_ci_measurements_no_key(self):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        # check with out ci.key
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def _test_ci_object(self, key, required=False):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        # check ci.depends with wrong type
        ci[key] = 1
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        self.assertListEqual(result, [])

        # check ci without ci.{key}
        if required is True:
            ci.pop(key)
            self.store_import(data, self.uuid)
            try:
                self.ctx_import.import_context(self.uuid)
            except Exception as e:
                self.fail(self._desc_fail.format(repr(e)))

            result = self.ctx_import.get_entities_by_id(ci[
                ContextGraphImport.K_ID])
            self.assertListEqual(result, [])

        # check ci.{key} with right type
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected = self.template_ent.copy()
        expected[Keys.ID] = ci[ContextGraphImport.K_ID]
        expected[Keys.NAME] = ci[ContextGraphImport.K_NAME]
        expected[Keys.TYPE] = ci[ContextGraphImport.K_TYPE]
        self.assertEqualEntities(result[0], expected)

    def test_ci_infos(self):
        self._test_ci_object(ContextGraphImport.K_INFOS, False)

    def _test_action(self, key):
        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        # Create an entity
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_CREATE
        ci[ContextGraphImport.K_ID] = "id"
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected = self.template_ent.copy()
        expected[Keys.ID] = ci[ContextGraphImport.K_ID]
        expected[Keys.NAME] = ci[ContextGraphImport.K_NAME]
        expected[Keys.TYPE] = ci[ContextGraphImport.K_TYPE]

        self.assertEqualEntities(result[0], expected)

        # Disable an entity
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DISABLE
        ci[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_DISABLE:
                                               [12345]}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected = self.template_ent.copy()
        expected[Keys.ID] = ci[ContextGraphImport.K_ID]
        expected[Keys.NAME] = ci[ContextGraphImport.K_NAME]
        expected[Keys.TYPE] = ci[ContextGraphImport.K_TYPE]
        expected[Keys.INFOS] = {
            Keys.DISABLE_HISTORY: ci[ContextGraphImport.K_PROPERTIES][ContextGraphImport.K_DISABLE],
            Keys.ENABLED: False
        }

        self.assertEqualEntities(result[0], expected)

        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DISABLE
        ci[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_DISABLE:
                                               [54321]}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected[Keys.INFOS][Keys.DISABLE_HISTORY].append(54321)

        self.assertEqualEntities(result[0], expected)

        # Enable an entity
        ci.pop(ContextGraphImport.K_PROPERTIES)
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_ENABLE
        ci[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_ENABLE:
                                               [67890]}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected[Keys.INFOS][Keys.ENABLE_HISTORY] = [67890]
        expected[Keys.INFOS][Keys.ENABLED] = True

        self.assertEqualEntities(result[0], expected)

        ci.pop(ContextGraphImport.K_PROPERTIES)
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_ENABLE
        ci[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_ENABLE:
                                               [9876]}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        expected[Keys.INFOS][Keys.ENABLE_HISTORY].append(9876)
        expected[Keys.INFOS][Keys.ENABLED] = True
        self.assertEqualEntities(result[0], expected)

        # Update an entity
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_UPDATE
        ci[ContextGraphImport.K_NAME] = "An other name"
        ci.pop(ContextGraphImport.K_INFOS)
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        expected[ContextGraphImport.K_NAME] = ci[ContextGraphImport.K_NAME]
        result = self.ctx_import.get_entities_by_id(
            ci[ContextGraphImport.K_ID])
        self.assertEqualEntities(result[0], expected)

        # Not a correct action
        ci[ContextGraphImport.K_ACTION] = "update_not"
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        ci[ContextGraphImport.K_ACTION] = "not_update"
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        ci[ContextGraphImport.K_ACTION] = "Update"
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        ci[ContextGraphImport.K_ACTION] = "a strange action"
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        # Delete an entity
        ci[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DELETE
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id("id")
        expected = []
        self.assertListEqual(result, expected)

    def test_ci_action(self):
        self._test_action(ContextGraphImport.K_CIS)

    def _test_ci_type(self, resource):

        if resource not in [ContextGraphImport.COMPONENT,
                            ContextGraphImport.CONNECTOR,
                            ContextGraphImport.RESOURCE]:
            self.fail("Unrecognized type")

        data = self.template_json.copy()
        ci = self.template_ci.copy()
        data[ContextGraphImport.K_CIS] = [ci]

        ci[ContextGraphImport.K_TYPE] = resource
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        result = self.ctx_import.get_entities_by_id("id")
        expected = self.template_ent.copy()
        expected[ContextGraphImport.K_ID] = ci[ContextGraphImport.K_ID]
        expected[ContextGraphImport.K_TYPE] = ci[ContextGraphImport.K_TYPE]
        expected[ContextGraphImport.K_NAME] = ci[ContextGraphImport.K_NAME]

        self.assertEqualEntities(result[0], expected)

    def test_ci_type_resource(self):
        self._test_ci_type(ContextGraphImport.RESOURCE)

    def test_ci_type_component(self):
        self._test_ci_type(ContextGraphImport.COMPONENT)

    def test_ci_type_connector(self):
        self._test_ci_type(ContextGraphImport.CONNECTOR)

    def test_ci_properties(self):
        self._test_ci_object(ContextGraphImport.K_PROPERTIES, False)

    def _test_link_string(self, key, required=False):
        self.put_dummy_entity("id")

        data = self.template_json.copy()
        link = self.template_link.copy()
        link[ContextGraphImport.K_FROM] = []
        link[ContextGraphImport.K_TO] = "id"
        data[ContextGraphImport.K_LINKS] = [link]

        # check link.{key} with right type
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        # check link.{key} with wrong type
        link[key] = 1
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        # check missing link.{key}
        if required is True:
            link.pop(key)
            self.store_import(data, self.uuid)
            with self.assertRaises(ValidationError):
                self.ctx_import.import_context(self.uuid)

    def test_link_id(self):
        self._test_link_string(ContextGraphImport.K_ID, required=False)

    def test_link_to(self):
        self._test_link_string(ContextGraphImport.K_TO, required=True)

    def test_link_from(self):
        self.put_dummy_entity("id")

        data = self.template_json.copy()
        link = self.template_link.copy()
        link[ContextGraphImport.K_FROM] = []
        link[ContextGraphImport.K_TO] = "id"
        data[ContextGraphImport.K_LINKS] = [link]

        # check link.from with empty list
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        # check link.from with list of string
        link[ContextGraphImport.K_FROM] = ["id", "id"]
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

        # check link.from with list of int
        link[ContextGraphImport.K_FROM] = [1, 2]
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        # check link.from with wront type
        link[ContextGraphImport.K_FROM] = 1
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

    def _test_link_object(self, key, required=False):
        data = self.template_json.copy()
        link = self.template_link.copy()
        link[ContextGraphImport.K_FROM] = []
        link[ContextGraphImport.K_TO] = "id"
        data[ContextGraphImport.K_LINKS] = [link]

        # check ci.depends with wrong type
        link[key] = 1
        self.store_import(data, self.uuid)
        with self.assertRaises(ValidationError):
            self.ctx_import.import_context(self.uuid)

        self.put_dummy_entity("id")

        # check with out ci.{key}
        if required is True:
            link.pop(key)
            self.store_import(data, self.uuid)
            try:
                self.ctx_import.import_context(self.uuid)
            except Exception as e:
                self.fail(self._desc_fail.format(repr(e)))

        # check ci.{key} with right type
        link[key] = {}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_link_infos(self):
        self._test_link_object(ContextGraphImport.K_INFOS, required=False)

    def test_link_action(self):
        self._test_action(ContextGraphImport.K_LINKS)

    def test_link_properties(self):
        self._test_link_object(ContextGraphImport.K_PROPERTIES, required=False)

    def _test_state(self, object_, state):
        self.put_dummy_entity("id")

        data = self.template_json.copy()

        if object_ == ContextGraphImport.K_CIS:
            obj = self.template_ci.copy()
            obj[ContextGraphImport.K_ID] = "id"
            data[ContextGraphImport.K_CIS] = [obj]
            obj_key = ContextGraphImport.K_LINKS
        elif object_ == ContextGraphImport.K_LINKS:
            obj = self.template_link.copy()
            data[ContextGraphImport.K_LINKS] = [obj]
            obj_key = ContextGraphImport.K_LINKS
            obj[ContextGraphImport.K_FROM] = []
            obj[ContextGraphImport.K_TO] = "id"
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
        self.store_import(data, self.uuid)
        with self.assertRaises(KeyError):
            self.ctx_import.import_context(self.uuid)

        # {object_}.action : {state} with ci.properties.{other_state} and
        # without {object_}.properties.{state}
        obj[ContextGraphImport.K_PROPERTIES] = {other_state: []}
        self.store_import(data, self.uuid)
        with self.assertRaises(KeyError):
            self.ctx_import.import_context(self.uuid)

        # {object_}.action : {state} with {object_}.properties.{state}
        obj[ContextGraphImport.K_PROPERTIES] = {state: [123456]}
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_ci_disable_with_properties(self):
        self._test_state(ContextGraphImport.K_CIS,
                         ContextGraphImport.A_DISABLE)

    def test_link_disable_with_properties(self):
        # self._test_state(ContextGraphImport.K_LINKS,
        #                  ContextGraphImport.A_DISABLE)
        data = self.template_json.copy()
        link = self.template_link.copy()
        data[ContextGraphImport.K_LINKS] = [link]
        link[ContextGraphImport.K_FROM] = []
        link[ContextGraphImport.K_TO] = "id"
        link[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DISABLE
        link[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_DISABLE:
                                                 [123456]}
        self.store_import(data, self.uuid)

        self.put_dummy_entity("id")
        with self.assertRaises(NotImplementedError):
            self.ctx_import.import_context(self.uuid)

    def test_ci_enable_with_properties(self):
        self._test_state(ContextGraphImport.K_CIS,
                         ContextGraphImport.A_ENABLE)

    def test_link_enable_with_properties(self):
        # self._test_state(ContextGraphImport.K_LINKS,
        #                  ContextGraphImport.A_ENABLE)
        data = self.template_json.copy()
        link = self.template_link.copy()
        data[ContextGraphImport.K_LINKS] = [link]
        link[ContextGraphImport.K_FROM] = []
        link[ContextGraphImport.K_TO] = "id"
        link[ContextGraphImport.K_ACTION] = ContextGraphImport.A_DISABLE
        link[ContextGraphImport.K_PROPERTIES] = {ContextGraphImport.A_DISABLE:
                                                 [123456]}
        self.store_import(data, self.uuid)

        self.put_dummy_entity("id")
        with self.assertRaises(NotImplementedError):
            self.ctx_import.import_context(self.uuid)

    def test_OK_single_ci(self):
        data = self.template_json.copy()
        data[ContextGraphImport.K_CIS] = [self.template_ci.copy()]
        self.store_import(data, self.uuid)

        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_OK_multiple_ci(self):
        data = self.template_json.copy()

        ci0 = self.template_ci.copy()
        ci0[ContextGraphImport.K_ID] = "id0"
        ci1 = self.template_ci.copy()
        ci1[ContextGraphImport.K_ID] = "id1"
        ci2 = self.template_ci.copy()
        ci2[ContextGraphImport.K_ID] = "id2"

        data[ContextGraphImport.K_CIS] = [ci0, ci1, ci2]
        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_OK_single_link(self):
        data = self.template_json.copy()
        link = self.template_link.copy()
        self.put_dummy_entity("id1")
        self.put_dummy_entity("id2")
        link[ContextGraphImport.K_FROM] = ["id2"]
        link[ContextGraphImport.K_TO] = "id1"
        data[ContextGraphImport.K_LINKS] = [link]

        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_OK_multiple_link(self):
        data = self.template_json.copy()

        self.put_dummy_entity("id0")
        self.put_dummy_entity("id1")
        self.put_dummy_entity("id2")

        link0 = self.template_link.copy()
        link0[ContextGraphImport.K_TO] = "id0"
        link0[ContextGraphImport.K_FROM] = ["id1"]
        link1 = self.template_link.copy()
        link1[ContextGraphImport.K_TO] = "id1"
        link1[ContextGraphImport.K_FROM] = ["id2"]
        link2 = self.template_link.copy()
        link2[ContextGraphImport.K_TO] = "id2"
        link2[ContextGraphImport.K_FROM] = ["id0"]

        data[ContextGraphImport.K_LINKS] = [link0, link1, link2]

        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

    def test_OK_both(self):
        data = self.template_json.copy()

        link0 = self.template_link.copy()
        link0[ContextGraphImport.K_TO] = "id0"
        link0[ContextGraphImport.K_FROM] = ["id1"]
        link1 = self.template_link.copy()
        link1[ContextGraphImport.K_TO] = "id1"
        link1[ContextGraphImport.K_FROM] = ["id2"]
        link2 = self.template_link.copy()
        link2[ContextGraphImport.K_TO] = "id2"
        link2[ContextGraphImport.K_FROM] = ["id0"]

        ci0 = self.template_ci.copy()
        ci0[ContextGraphImport.K_ID] = "id0"
        ci1 = self.template_ci.copy()
        ci1[ContextGraphImport.K_ID] = "id1"
        ci2 = self.template_ci.copy()
        ci2[ContextGraphImport.K_ID] = "id2"

        data[ContextGraphImport.K_LINKS] = [link0, link1, link2]
        data[ContextGraphImport.K_CIS] = [ci0, ci1, ci2]

        self.store_import(data, self.uuid)
        try:
            self.ctx_import.import_context(self.uuid)
        except Exception as e:
            self.fail(self._desc_fail.format(repr(e)))

class ReportManager(TestCase):

    def setUp(self):
        self.import_storage = Middleware.get_middleware_by_uri(
            'storage-default-testgraphimport://'
        )
        self.template_report = {ImportKey.F_ID: None,
                                ImportKey.F_STATUS: None,
                                ImportKey.F_CREATION: None,
                                ImportKey.F_START: None,
                                ImportKey.F_EXECTIME: None,
                                ImportKey.F_STATS: {
                                    ImportKey.F_DELETED: None,
                                    ImportKey.F_UPDATED: None}}

        self.manager = Manager()
        self.manager[Manager.STORAGE] = self.import_storage

        self.uuid = "i-am-an-uuid"

    def tearDown(self):
        self.import_storage.remove_elements()

    def test_get_next_uuid(self):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"
        uuid_2 = self.uuid + "_2"
        uuid_3 = self.uuid + "_3"

        self.manager.create_import_status(uuid_0)
        self.manager.create_import_status(uuid_1)
        self.manager.create_import_status(uuid_2)
        self.manager.create_import_status(uuid_3)

        self.manager.update_status(uuid_0,
                                   {ImportKey.F_STATUS: ImportKey.ST_DONE})
        self.manager.update_status(uuid_2,
                                   {ImportKey.F_STATUS: ImportKey.ST_ONGOING})
        self.manager.update_status(uuid_3,
                                   {ImportKey.F_STATUS: ImportKey.ST_FAILED})

        uuid = self.manager.get_next_uuid()
        report = self.manager.get_import_status(uuid)
        self.assertEqual(report[ImportKey.F_ID], uuid_1)
        self.assertEqual(report[ImportKey.F_STATUS], ImportKey.ST_PENDING)
        try:
             time.strptime(report[ImportKey.F_CREATION], Manager.DATE_FORMAT)
        except ValueError:
            self.fail("Error while converting the creation time.")

    def test_get_next_uuid_no_uuid(self):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"
        uuid_2 = self.uuid + "_2"

        self.manager.create_import_status(uuid_0)
        self.manager.create_import_status(uuid_1)
        self.manager.create_import_status(uuid_2)

        self.manager.update_status(uuid_0,
                                   {ImportKey.F_STATUS: ImportKey.ST_DONE})
        self.manager.update_status(uuid_1,
                                   {ImportKey.F_STATUS: ImportKey.ST_ONGOING})
        self.manager.update_status(uuid_2,
                                   {ImportKey.F_STATUS: ImportKey.ST_FAILED})

        uuid = self.manager.get_next_uuid()
        self.assertEquals(uuid, None)

    def test_get_next_uuid_multiple(self):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"

        self.manager.create_import_status(uuid_0)
        time.sleep(2)
        self.manager.create_import_status(uuid_1)

        uuid = self.manager.get_next_uuid()
        self.assertEquals(uuid, uuid_0)

    def test_is_present(self):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"

        self.manager.create_import_status(uuid_0)

        self.assertTrue(self.manager.is_present(uuid_0))
        self.assertFalse(self.manager.is_present(uuid_1))

    def test_update_status_authorized(self):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"

        self.manager.create_import_status(uuid_0)
        self.manager.create_import_status(uuid_1)

        expected_1 = self.manager.get_import_status(uuid_1)
        expected_0 = self.manager.get_import_status(uuid_0)

        update = {ImportKey.F_STATUS: ImportKey.ST_DONE,
                  ImportKey.F_INFO: "Nothing wrong",
                  ImportKey.F_EXECTIME: "01:23:45",
                  ImportKey.F_START: "54:32:01",
                  ImportKey.F_STATS: "Stats"}

        for k in update.keys():
            expected_0[k] = update[k]

        self.manager.update_status(uuid_0, update)

        result_0 = self.manager.get_import_status(uuid_0)
        result_1 = self.manager.get_import_status(uuid_1)

        self.assertDictEqual(result_0, expected_0)
        self.assertDictEqual(result_1, expected_1)

    def test_update_status_authorized_no_uuid(self):
        with self.assertRaises(ValueError):
            self.manager.update_status("uuid", {})

    def test_update_status_not_authorized(self):
        uuid = self.uuid

        self.manager.create_import_status(uuid)
        expected = self.manager.get_import_status(uuid)

        update = {ImportKey.F_ID: ImportKey.ST_DONE,
                  ImportKey.F_CREATION: "01:23:45"}

        self.manager.update_status(uuid, update)

        result = self.manager.get_import_status(uuid)

        self.assertDictEqual(result, expected)

    def test_create_import(self):
        expected = {ImportKey.F_ID: self.uuid,
                    ImportKey.F_CREATION: time.asctime(),
                    ImportKey.F_STATUS: ImportKey.ST_PENDING}
        self.manager.create_import_status(self.uuid)
        result = self.manager.get_import_status(self.uuid)

        self.assertDictEqual(result, expected)

    def test_create_import_same_uuid(self):
        self.manager.create_import_status(self.uuid)
        desc = "An import status with the same uuid ({0}) already "\
               "exist.".format(self.uuid)

        # FIXME the message did not match the expected one. I don't know why
        # with self.assertRaisesRegexp(ValueError, des):
        #     self.manager.create_import_status(self.uuid)

        try:
            self.manager.create_import_status(self.uuid)
        except ValueError as e:
            self.assertEqual(desc, e.message)
        except:
            self.fail("An exception different of ValueError was raised")

    def _test_state_on_db(self, func, state, other_state):
        uuid_0 = self.uuid + "_0"
        uuid_1 = self.uuid + "_1"
        uuid_2 = self.uuid + "_2"

        self.manager.create_import_status(uuid_0)
        self.manager.create_import_status(uuid_1)
        self.manager.create_import_status(uuid_2)

        self.manager.update_status(uuid_0,
                                   {ImportKey.F_STATUS: ImportKey.ST_DONE})
        self.manager.update_status(uuid_1,
                                   {ImportKey.F_STATUS: ImportKey.ST_FAILED})
        self.manager.update_status(uuid_2,
                                   {ImportKey.F_STATUS: state})

        self.assertTrue(func())
        self.manager.update_status(uuid_2, {ImportKey.F_STATUS: other_state})
        self.assertFalse(func())

    def test_on_going_in_db(self):
        self._test_state_on_db(self.manager.on_going_in_db,
                               ImportKey.ST_ONGOING,
                               ImportKey.ST_DONE)

    def test_pending_in_db(self):
        self._test_state_on_db(self.manager.pending_in_db,
                               ImportKey.ST_PENDING,
                               ImportKey.ST_DONE)

    def test_check_id(self):
        self.assertFalse(self.manager.check_id(self.uuid))
        self.manager.create_import_status(self.uuid)
        self.assertTrue(self.manager.check_id(self.uuid))

    def test_get_import_status(self):
        self.manager.create_import_status(self.uuid)

        expected = self.manager.get_import_status(self.uuid)

        update = {ImportKey.F_STATUS: ImportKey.ST_DONE,
                  ImportKey.F_INFO: "Nothing wrong",
                  ImportKey.F_EXECTIME: "01:23:45",
                  ImportKey.F_START: "54:32:01",
                  ImportKey.F_STATS: "Stats"}

        for k in update.keys():
            expected[k] = update[k]

        self.manager.update_status(self.uuid, update)
        result = self.manager.get_import_status(self.uuid)

        self.assertDictEqual(result, expected)


if __name__ == '__main__':
    main()
