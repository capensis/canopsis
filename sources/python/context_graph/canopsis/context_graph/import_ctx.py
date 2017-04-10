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

    def __a_delete_entity(self, ci):
        id_ = ci["_id"]
        try:
            entity = self.entities_to_update.get(id_["_id"])
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
        pass

    def __a_create_entity(self, ci):
        pass

    def __a_disable_entity(self, ci):
        pass

    def __a_enable_entity(self, ci):
        pass

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


        for ci in json[self.CIS]:
            if ci[self.ACTION] == self.DELETE:
                self.__a_delete_entity(ci)
            if ci[self.ACTION] == self.CREATE:
                self.__a_delete_entity(ci)
            elif ci[self.ACTION] == self.UPDATE:
                self.__a_delete_entity(ci)
            elif ci[self.ACTION] == self.DISABLE:
                self.__a_delete_entity(ci)
            elif ci[self.ACTION] == self.ENABLE:
                self.__a_delete_entity(ci)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    ci[self.ACTION]))

        for link in json[self.LINKS]:
            if link[self.ACTION] == self.DELETE:
                self.__a_delete_link(link)
            if link[self.ACTION] == self.CREATE:
                self.__a_delete_link(link)
            elif link[self.ACTION] == self.UPDATE:
                self.__a_delete_link(link)
            elif link[self.ACTION] == self.DISABLE:
                self.__a_delete_link(link)
            elif link[self.ACTION] == self.ENABLE:
                self.__a_delete_link(link)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    link[self.ACTION]))

        self.put_entities(self.update.values())
