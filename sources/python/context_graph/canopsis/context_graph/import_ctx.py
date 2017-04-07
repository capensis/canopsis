#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from canopsis.context_graph.manager import ContextGraph
import ast


class ContextGraphImport(ContextGraph):

    CIS = "cis"
    LINKS = "links"
    FROM = "from"
    TO = "to"
    ACTION = "action"
    TYPE = "type"
    EXTRAS = "extras"
    ACTION = "action"
    ID = "_id"
    NAME = "name"

    DELETE = "delete"
    CREATE = "create"
    UPDATE = "update"
    DISABLE = "disable"
    ENABLE = "enable"

    JSON_FIELDS = [CIS, LINKS]

    __REQUIRED = 0
    __OPTIONAL = 1

    CIS_DICT = {ID: (__REQUIRED, str),
                NAME: (__REQUIRED, str),
                TYPE: (__REQUIRED, str),
                EXTRAS: (__REQUIRED, dict),
                ACTION: (__REQUIRED, str)}

    LINKS_DICT = {FROM: (__REQUIRED, str),
                  TO: (__REQUIRED, str),
                  EXTRAS: (__REQUIRED, dict),
                  ACTION: (__REQUIRED, str)}

    def __get_entities_to_update(self, json):
        # a set so no duplicate ids without effort and low time complexity
        ids = set()

        for ci in json[self.CIS]:
            ids.add(ci["_id"])

        for link in json[self.LINKS]:
            ids.add(link[self.FROM])
            ids.add(link[self.TO])

        ctx = {}
        result = self.get_entities_by_id(list(ids))

        for doc in result:
            ctx[doc["_id"]] = doc

        return ctx

    def a_delete_entity(self, ci, entities_to_update, update):
        pass

    def a_update_entity(self, ci, entities_to_update, update):
        pass

    def a_create_entity(self, ci, entities_to_update, update):
        pass

    def a_disable_entity(self, ci, entities_to_update, update):
        pass

    def a_enable_entity(self, ci, entities_to_update, update):
        pass

    def a_delete_link(self, link, entities_to_update, update):
        pass

    def a_update_link(self, link, entities_to_update, update):
        pass

    def a_create_link(self, link, entities_to_update, update):
        pass

    def a_disable_link(self, link, entities_to_update, update):
        pass

    def a_enable_link(self, link, entities_to_update, update):
        pass

    def _check_schema(self, json):
        # TODO add better description in the exception.
        # TODO add checks for list of expected values.
        # TODO add check for optional/required fields

        def check_sub_elt(schema, dict_):
            if len(schema.keys()) > len(dict_.keys()):
                raise ValueError(
                    "Some key/keys are missing in a sub element of the json.")
            if len(schema.keys()) > len(dict_.keys()):
                raise ValueError(
                    "Too many key/lkeys in a sub element of the json.")

            for key in schema.keys():
                if not schema[key][1] == type(dict_[key][1]):
                    raise ValueError(
                        "A value did not match the expected value.")

        if not self.JSON_FIELDS == json.keys():
            if len(self.JSON_FIELDS) > len(json.keys()):
                raise ValueError("Some keys are missing in the json.")
            if len(self.JSON_FIELDS) < len(json.keys()):
                raise ValueError("Too many keys are present in the json.")

        # check CIS fields
        for ci in json[self.CIS]:
            check_sub_elt(self.CIS_DICT, ci)

        # check CIS fields
        for link in json[self.LINKS]:
            check_sub_elt(self.LINK_DICT, link)

    def import_context(self, json):
        if (not isinstance(json, dict) or isinstance(json, str)):
            raise ValueError("Json should a string or a dict")

        if isinstance(json, str):
            json = ast.literal_eval(json)

        self._check_schema(json)

        entities_to_update = self.__get_entities_to_update(json)

        update = {}

        for ci in json[self.CIS]:
            if ci[self.ACTION] == self.DELETE:
                self.a_delete_entity(ci, entities_to_update, update)
            if ci[self.ACTION] == self.CREATE:
                self.a_delete_entity(ci, entities_to_update, update)
            elif ci[self.ACTION] == self.UPDATE:
                self.a_delete_entity(ci, entities_to_update, update)
            elif ci[self.ACTION] == self.DISABLE:
                self.a_delete_entity(ci, entities_to_update, update)
            elif ci[self.ACTION] == self.ENABLE:
                self.a_delete_entity(ci, entities_to_update, update)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    ci[self.ACTION]))

        for link in json[self.LINKS]:
            if link[self.ACTION] == self.DELETE:
                self.a_delete_link(link, entities_to_update, update)
            if link[self.ACTION] == self.CREATE:
                self.a_delete_link(link, entities_to_update, update)
            elif link[self.ACTION] == self.UPDATE:
                self.a_delete_link(link, entities_to_update, update)
            elif link[self.ACTION] == self.DISABLE:
                self.a_delete_link(link, entities_to_update, update)
            elif link[self.ACTION] == self.ENABLE:
                self.a_delete_link(link, entities_to_update, update)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    link[self.ACTION]))

        self.put_entities(update.values())
