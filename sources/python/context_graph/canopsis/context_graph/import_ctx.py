#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from canopsis.context_graph.manager import ContextGraph
import ast


class ContextGraphImport(ContextGraph):

    # TODO add a feature to restore the context if an error occured during while
    # is pushed into the database

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
    K_ENABLE = "enable"
    K_DISABLE = "disable"

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
        """Return every entities id required for the update
        :param json: the json with every actions required for the update
        :param rtype: a dict with the entity id as a key and the entity as
        a value
        """
        # a set so no duplicate ids without effort and low time complexity
        ids = set()

        for ci in json[self.K_CIS]:
            ids.add(ci[self.K_ID])
            if ci[self.K_ACTION] == self.A_DELETE:
                # we need to retreive every related entity to update the links
                entity = self.get_entities_by_id(ci[self.K_ID])[0]

                for id_ in entity["depends"] + entity["impact"]:
                    ids.add(id_)

        for link in json[self.K_LINKS]:
            ids.add(link[self.K_FROM])
            ids.add(link[self.K_TO])

        ctx = {}
        result = self.get_entities_by_id(list(ids))

        for doc in result:
            ctx[doc[self.K_ID]] = doc

        return ctx

    def __a_delete_entity(self, ci):
        """Update the entities related with the entity to be deleted disigned
        by ci and store them into self.update. Add the id of entity to be
        deleted into self.delete.

        If the entity to be deleted is not initially store in the context,
        a ValueError will be raised.

        :param ci: the ci (see the JSON specification).
        """

        id_ = ci[self.K_ID]

        try:
            entity = self.entities_to_update[ci[self.K_ID]]
        except KeyError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # Update the depends/impact link
        for ent_id in entity["depends"]:
            self.update[ent_id] = self.entities_to_update[ent_id].copy()
            self.update[ent_id]["impact"].remove(id_)

        # Update the impact/depends link
        for ent_id in entity["impact"]:
            self.update[ent_id] = self.entities_to_update[ent_id].copy()
            self.update[ent_id]["depends"].remove(id_)

        self.delete.append(id_)

    def __a_update_entity(self, ci):
        """Update the entity with the information stored into the ci and store
        the result into self.update.

        If the entity to be updated is not initially store in the context,
        a ValueError will be raised.

        :param ci: the ci (see the JSON specification).
        """


        if not self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} does not match any existing entity.".format(
                ci[self.K_ID])
            raise ValueError(desc)

        entity = self.entities_to_update[ci[self.K_ID]]

        fields_to_update = [self.K_NAME, self.K_INFOS,self.K_TYPE]

        for field in fields_to_update:
            entity[field] = ci[field]

        self.update[ci[self.K_ID]] = entity


    def __a_create_entity(self, ci):
        """Create an entity with a ci and store it into self.update

        If the new entity is initially store in the context, a ValueError will
        be raised.

        :param ci: the ci (see the JSON specification).
        """
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

    def __change_state_entity(self, ci, state):
        """Change the state (enable/disable) of an entity and store the result
        into self.update.

        If state does not match enable or disable, a ValueError will be raised.

        :param ci: the ci (see the JSON specification).
        :param state: if the state is "disable", the timestamp of the
        deactivation of the entity will be store into the fields infos.disable.
        Same behaviour with "enable" but the timestamp will be store into
        infos.enable.
        """
        if state != self.K_DISABLE and state != self.K_ENABLE:
            raise ValueError("{0} is not a valid state.".format(state))

        id_ = ci[self.K_ID]

        if not self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} does not match any existing entity.".format(
                id_)
            raise ValueError(desc)

        # If the entity is not in the update dict, add it
        if not self.update.has_key(id_):
            self.update[id_] = self.entities_to_update[id_].copy()

        # Update entity the fields enable/disable of infos
        timestamp = ci[self.K_INFOS][state]

        if not isinstance(timestamp, list):
            timestamp = [timestamp]

        if self.update[id_][self.K_INFOS].has_key(state):

            if self.update[id_][self.K_INFOS][state] is None:
                self.update[id_][self.K_INFOS][state] = timestamp
            else:
                self.update[id_][self.K_INFOS][state] += timestamp

        else:
            self.update[id_][self.K_INFOS][state] = timestamp

    def __a_disable_entity(self, ci):
        """Disable an entity defined by ci. For more information, see
        __change_state.

        :param ci: the ci (see the JSON specification).
        """
        self.__change_state_entity(ci, self.K_DISABLE)

    def __a_enable_entity(self, ci):
        """Enable an entity defined by ci. For more information, see
        __change_state.

        :param ci: the ci (see the JSON specification).
        """
        self.__change_state_entity(ci, self.K_ENABLE)

    def __a_delete_link(self, link):
        """Delete a link between two entity and store the modify entities
        into self.udpate.

        :param link: the link that identify a link (see the JSON specification).
        """

        if link[self.K_FROM] not in self.update.keys():
            self.update[link[self.K_FROM]] = self.entities_to_update[link[self.K_FROM]]

        if link[self.K_TO] not in self.update.keys():
            self.update[link[self.K_TO]] = self.entities_to_update[link[self.K_TO]]

        self.update[link[self.K_FROM]]['impact'].remove(link[self.K_TO])
        self.update[link[self.K_TO]]['depends'].remove(link[self.K_FROM])

    def __a_update_link(self, link):
        raise NotImplementedError()

    def __a_create_link(self, link):
        """Create a link between two entity and store the modify entities
        into self.udpate.

        :param link: the link that identify a link (see the JSON specification).
        """

        if link[self.K_FROM] not in self.update.keys():
            self.update[link[self.K_FROM]] = self.entities_to_update[link[self.K_FROM]]

        if link[self.K_TO] not in self.update.keys():
            self.update[link[self.K_TO]] = self.entities_to_update[link[self.K_TO]]

        self.update[link[self.K_FROM]]['impact'].append(link[self.K_TO])
        self.update[link[self.K_TO]]['depends'].append(link[self.K_FROM])

    def __a_disable_link(self, link):
        raise NotImplementedError()

    def __a_enable_link(self, link):
        raise NotImplementedError()

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

        for id_ in self.update:
            if id_ in self.delete:
                desc = "The entity {0} to be deleted is updated in "\
                       "the same import. Update aborted.".format(id_)
                raise ValueError(desc)

        self.put_entities(self.update.values())
        self._delete_entity(self.delete)
