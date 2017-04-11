#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from canopsis.context_graph.manager import ContextGraph
import ast


class ContextGraphImport(ContextGraph):

    K_CIS = "cis"
    K_LINKS = "links"
    K_FROM = "from"
    K_TO = "to"
    K_ACTION = "action"
    K_TYPE = "type"
    K_INFOS = "infos"
    K_ACTION = "action"
    K_ID = "_id"
    K_NAME = "name"

    A_DELETE = "delete"
    A_CREATE = "create"
    A_UPDATE = "update"
    A_DISABLE = "disable"
    A_ENABLE = "enable"

    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraphImport, self).__init__(*args, **kwargs)

        self.entities_to_update = {}
        self.update = {}
        self.delete = []

    def __get_entities_to_update(self, json):
        # a set so no duplicate ids without effort and low time complexity
        ids = set()

        for ci in json[self.K_CIS]:
            ids.add(ci[self.K_ID])
            if ci[self.K_ACTION] == self.A_DELETE:
                entity = self.get_entities_by_id(ci[self.K_ID])[0]

                for id_ in entity["depends"] + entity["impact"]:
                    ci.add(id_)

        for link in json[self.K_LINKS]:
            ids.add(link[self.K_FROM])
            ids.add(link[self.K_TO])

        ctx = {}
        result = self.get_entities_by_id(list(ids))

        for doc in result:
            ctx[doc[self.K_ID]] = doc

        return ctx

    def __a_delete_entity(self, ci):
        # TODO rewrite this function
        id_ = ci[self.K_ID]
        try:
            entity = self.entities_to_update.get(id_[self.K_ID])
        except IndexError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # update depends/impact links
        # Est ce que c'est utile ?
        status = {"deletions": entity["depends"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "depends")
        for entity in updated_entities:
            if self.update.has_key(entity["_id"]):
                self.update[entity["_id"]]['impact'].remove(entity["_id"])
            else:
                self.update[entity["_id"]] = entity

        # update impact/depends links
        status = {"deletions": entity["impact"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "impact")
        for entity in updated_entities:
            if self.update.has_key(entity["_id"]):
                self.update[entity["_id"]]['depends'].remove(entity["_id"])
            else:
                self.update[entity["_id"]] = entity


    def __a_update_entity(self, ci):
        if not self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} does not match any existing entity.".format(
                ci[self.K_ID])
            raise KeyError(desc)

        entity = self.entities_to_update[ci[self.K_ID]]

        fields_to_update = [self.K_NAME, self.K_INFOS,self.K_TYPE]

        for field in fields_to_update:
            entity[field] = ci[field]

        self.update[ci[self.K_ID]] = entity


    def __a_create_entity(self, ci):
        if self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} match an existing entity.".format(
                ci["_id"])
            raise ValueError(desc)

        entity = {'_id': ci[self.K_ID],
                  'type': ci[self.K_TYPE],
                  'name': ci[self.K_NAME],
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': ci[self.K_INFOS]}

        self.update[ci[self.K_ID]] = entity

    def __a_disable_entity(self, ci):
        if not self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} does not match any existing entity.".format(
                ci[self.K_ID])
            raise KeyError(desc)

        if not self.update.has_key(ci[self.K_ID]):
            self.update[ci[self.K_ID]] = self.entities_to_update.copy()
        # self.update[ci[self.K_ID]][self.K_INFOS][]

    def __a_enable_entity(self, ci):
        if not self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} does not match any existing entity.".format(
                ci[self.K_ID])
            raise KeyError(desc)

        entity = self.entities_to_update.copy()
        # store the enable timestamp
        self.update[ci[self.K_ID]] = entity

    def __a_delete_link(self, link):
        pass

    def __a_update_link(self, link):
        raise NotImplementedError()

    def __a_create_link(self, link):
        pass

    def __a_disable_link(self, link):
        pass

    def __a_enable_link(self, link):
        pass


    def import_context(self, json):
        if (not isinstance(json, dict) or isinstance(json, str)):
            raise ValueError("Json should a string or a dict")

        if isinstance(json, str):
            json = ast.literal_eval(json)

        self.entities_to_update = self.__get_entities_to_update(json)


        for ci in json[self.K_CIS]:
            if ci[self.K_ACTION] == self.A_DELETE:
                self.__a_delete_entity(ci)
            if ci[self.K_ACTION] == self.A_CREATE:
                self.__a_delete_entity(ci)
            elif ci[self.K_ACTION] == self.A_UPDATE:
                self.__a_delete_entity(ci)
            elif ci[self.K_ACTION] == self.A_DISABLE:
                self.__a_delete_entity(ci)
            elif ci[self.K_ACTION] == self.A_ENABLE:
                self.__a_delete_entity(ci)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    ci[self.K_ACTION]))

        for link in json[self.K_LINKS]:
            if link[self.K_ACTION] == self.A_DELETE:
                self.__a_delete_link(link)
            if link[self.K_ACTION] == self.A_CREATE:
                self.__a_delete_link(link)
            elif link[self.K_ACTION] == self.A_UPDATE:
                self.__a_delete_link(link)
            elif link[self.K_ACTION] == self.A_DISABLE:
                self.__a_delete_link(link)
            elif link[self.K_ACTION] == self.A_ENABLE:
                self.__a_delete_link(link)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    link[self.K_ACTION]))
