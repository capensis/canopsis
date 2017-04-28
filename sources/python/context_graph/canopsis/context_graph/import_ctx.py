#!/usr/bin/env/python
# -*- coding: utf-8 -*-

import jsonschema
import ijson
from canopsis.context_graph.manager import ContextGraph


class ContextGraphImport(ContextGraph):

    # TODO add a feature to restore the context if an error occured during while
    # is pushed into the database

    K_LINKS = "links"
    K_FROM = "from"
    K_TO = "to"
    K_CIS = "cis"
    K_ID = "_id"
    K_NAME = "name"
    K_TYPE = "type"
    K_DEPENDS = "depends"
    K_IMPACT = "impact"
    K_MEASUREMENTS = "measurements"
    K_INFOS = "infos"
    K_ACTION = "action"
    K_ENABLE = "enable"
    K_DISABLE = "disable"
    K_PROPERTIES = "action_properties"

    # If you add an action, remember to add in the a_pattern string in method
    # import_checker
    A_DELETE = "delete"
    A_CREATE = "create"
    A_UPDATE = "update"
    A_DISABLE = "disable"
    A_ENABLE = "enable"


    __A_PATTERN = "^delete$|^create$|^update$|^disable$|^enable$"
    __T_PATTERN = "^resource$|^component$|^connector"
    __CI_REQUIRED = [K_ID,
                     K_ACTION,
                     K_TYPE]
    __LINK_REQUIRED = [K_ID,
                       K_FROM,
                       K_TO,
                       K_ACTION]

    SCHEMA = {
        "type": "object",
        "$schema": "http://json-schema.org/draft-04/schema#",
        "properties": {
            K_CIS: {
                "$schema": "http://json-schema.org/draft-04/schema#",
                "type": "array",
                "items": {
                    "$schema": "http://json-schema.org/draft-04/schema#",
                    "type": "object",
                    "required": __CI_REQUIRED,
                    "uniqueItems": True,
                    "properties": {
                        K_ID: {"type": "string"},
                        K_NAME: {"type": "string"},
                        K_DEPENDS: {"type": "array",
                                         "items": {
                                             "type": "string"}},
                        K_IMPACT: {"type": "array",
                                        "items": {
                                            "type": "string"}},
                        K_MEASUREMENTS: {"type": "array",
                                              "items":
                                              {"type":
                                               "string"}},
                        K_INFOS: {"type": "object"},
                        K_ACTION: {"type": "string",
                                        "pattern": __A_PATTERN},
                        K_TYPE: {"type": "string",
                                      "pattern": __T_PATTERN},
                        K_PROPERTIES: {"type": "object"}}}},
            K_LINKS: {
                "$schema": "http://json-schema.org/draft-04/schema#",
                "type": "array",
                "items": {
                    "$schema": "http://json-schema.org/draft-04/schema#",
                    "uniqueItems": True,
                    "type": "object",
                    "required": __LINK_REQUIRED,
                    "properties": {
                        K_ID: {"type": "string"},
                        K_FROM: {"type": "array",
                                      "items": {
                                          "type": "string"}},
                        K_TO: {"type": "string"},
                        K_INFOS: {"type": "object"},
                        K_ACTION: {"type": "string",
                                        "pattern": __A_PATTERN},
                        K_PROPERTIES: {"type": "object"}}}}}}



    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraphImport, self).__init__(*args, **kwargs)

        self.entities_to_update = {}
        self.update = {}
        self.delete = []

    @classmethod
    def check_element(cls, element, type_):

        jsonschema.validate(element, cls.SCHEMA["properties"][type_]["items"])

        state = element[ContextGraphImport.K_ACTION]
        if state == ContextGraphImport.A_DISABLE\
           or state == ContextGraphImport.A_ENABLE:
            element[ContextGraphImport.K_PROPERTIES][state]

    def clean_attributes(self):
        self.entities_to_update.clear()
        self.update.clear()
        del self.delete[:]


    def __get_entities_to_update(self, file_):
        """Return every entities id required for the update
        :param json: the json with every actions required for the update
        :param rtype: a dict with the entity id as a key and the entity as
        a value
        """
        # a set so no duplicate ids without effort and low time complexity
        ids = set()

        for ci in ijson.items(file_, "{0}.item".format(self.K_CIS)):
            ids.add(ci[self.K_ID])

            # we need to retreive every related entity to update the links
            if ci[self.K_ACTION] == self.A_DELETE:
            # FIXME do the get_entities_by_id in one call Then add all impacts\
                # depends
                entity = self.get_entities_by_id(ci[self.K_ID])[0]

                for id_ in entity["depends"] + entity["impact"]:
                    ids.add(id_)

        file_.seek(0)

        for link in ijson.items(file_, "{0}.item".format(self.K_LINKS)):
            for id_ in link[self.K_FROM]:
                ids.add(id_)

            ids.add(link[self.K_ID])

        result = self.get_entities_by_id(list(ids))
        ctx = {}

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
            entity = self.entities_to_update[id_]
        except KeyError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # Update the depends/impact link
        for ent_id in entity["depends"]:
            if ent_id in self.delete:
                # the entity of id ent_id is already deleted, skipping
                continue

            self.update[ent_id] = self.entities_to_update[ent_id].copy()
            try:
                self.update[ent_id]["impact"].remove(id_)
            except ValueError:
                raise ValueError("Try to remove {0} from impacts field of"\
                                 "entity {1}.".format(id_, ent_id))

        # Update the impact/depends link
        for ent_id in entity["impact"]:
            if ent_id in self.delete:
                # the entity of id ent_id is already deleted, skipping
                continue

            self.update[ent_id] = self.entities_to_update[ent_id].copy()
            try:
                self.update[ent_id]["depends"].remove(id_)
            except ValueError:
                raise ValueError("Try to remove {0} from impacts field of"\
                                 "entity {1}.".format(id_, ent_id))

        if id_ in self.update.keys():
            self.update.pop(id_)

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

        fields_to_update = [
            self.K_NAME,
            self.K_TYPE,
            self.K_MEASUREMENTS,
            self.K_INFOS]

        # if a a fields is missing we assume we did not need to update it
        for field in fields_to_update:
            try:
                entity[field] = ci[field]
            except KeyError:
                pass

        self.update[ci[self.K_ID]] = entity


    def __a_create_entity(self, ci):
        """Create an entity with a ci and store it into self.update

        If the new entity is initially store in the context, a ValueError will
        be raised.

        :param ci: the ci (see the JSON specification).
        """
        # TODO handle the creation of the name if needed and if the id
        # match the id scheme used in canopsis
        if self.entities_to_update.has_key(ci[self.K_ID]):
            desc = "The ci of id {0} match an existing entity.".format(
                ci["_id"])
            raise ValueError(desc)

        # set default value for required fields
        if not ci.has_key(self.K_NAME):
            ci[self.K_NAME] = ci[self.K_ID]
        if not ci.has_key(self.K_DEPENDS):
            ci[self.K_DEPENDS] = []
        if not ci.has_key(self.K_IMPACT):
            ci[self.K_IMPACT] = []
        if not ci.has_key(self.K_MEASUREMENTS):
            ci[self.K_MEASUREMENTS] = []
        if not ci.has_key(self.K_INFOS):
            ci[self.K_INFOS] = {}

        entity = {'_id': ci[self.K_ID],
                  'type': ci[self.K_TYPE],
                  'name': ci[self.K_NAME],
                  'depends': ci[self.K_DEPENDS],
                  'impact': ci[self.K_IMPACT],
                  'measurements': ci[self.K_MEASUREMENTS],
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
        timestamp = ci[self.K_PROPERTIES][state]

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

        if link[self.K_TO] not in self.update.keys():
            self.update[link[self.K_TO]] = self.entities_to_update[link[self.K_TO]]

        for ci_id in link[self.K_FROM]:
            if ci_id not in self.update.keys():
                self.update[ci_id] = self.entities_to_update[ci_id]

            self.update[ci_id]['impact'].append(link[self.K_TO])
            self.update[link[self.K_TO]]['depends'].append(ci_id)

    def __a_disable_link(self, link):
        raise NotImplementedError()

    def __a_enable_link(self, link):
        raise NotImplementedError()

    def import_context(self, file_, id_):
        fd = open(file_, 'r')

        # In case the previous import failed and raise an exception, we clean\
            #    now
        self.clean_attributes()

        self.entities_to_update = self.__get_entities_to_update(fd)

        fd.seek(0)

        for ci in ijson.items(fd, "{0}.item".format(self.K_CIS)):

            # add function to check if the element is correct and if the state is right
            self.check_element(ci, self.K_CIS)
            if ci[self.K_ACTION] == self.A_DELETE:
                self.__a_delete_entity(ci)
            elif ci[self.K_ACTION] == self.A_CREATE:
                self.__a_create_entity(ci)
            elif ci[self.K_ACTION] == self.A_UPDATE:
                self.__a_update_entity(ci)
            elif ci[self.K_ACTION] == self.A_DISABLE:
                self.__a_disable_entity(ci)
            elif ci[self.K_ACTION] == self.A_ENABLE:
                self.__a_enable_entity(ci)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    ci[self.K_ACTION]))

        fd.seek(0)

        for link in ijson.items(fd, "{0}.item".format(self.K_LINKS)):
            self.check_element(link, self.K_LINKS)
            if link[self.K_ACTION] == self.A_DELETE:
                self.__a_delete_link(link)
            elif link[self.K_ACTION] == self.A_CREATE:
                self.__a_create_link(link)
            elif link[self.K_ACTION] == self.A_UPDATE:
                self.__a_update_link(link)
            elif link[self.K_ACTION] == self.A_DISABLE:
                self.__a_disable_link(link)
            elif link[self.K_ACTION] == self.A_ENABLE:
                self.__a_enable_link(link)
            else:
                raise ValueError("The action {0} is not recognized\n".format(
                    link[self.K_ACTION]))

        for id_ in self.update:
            if id_ in self.delete:
                desc = "The entity {0} to be deleted is updated in "\
                       "the same import. Update aborted.".format(id_)
                raise ValueError(desc)

        self._put_entities(self.update.values())
        self._delete_entities(self.delete)

        self.clean_attributes()

    def on_going_in_db(self):
        """
            check if there is an import on going

            :return: True if an import is on going
        """
        return len(list(self[self.IMPORT_STORAGE].find_elements(query={'state': 'on going'}))) == 1

    def check_id(self, _id):
        """
            check if an id is already taken
        """
        return len(list(self[self.IMPORT_STORAGE].get_elements(query={'_id': _id}))) == 1

    def put_import(self, doc):
        """
            put an inport in db
        """
        return self[self.IMPORT_STORAGE].put_elements(doc)

    def update_import(self, doc):
        """
            update import in db
        """
        check = len(list(self[self.IMPORT_STORAGE].get_elements(query={'_id': doc['_id']}))) == 1
        if check:
            return self[self.IMPORT_STORAGE].put_elements(doc)
        else:
            raise ValueError:
                self.logger.error('not in db')
            
    def get_import_status(self, _id):
        """
            return the status of import
        """
        return list(self[self.IMPORT_STORAGE].get_elements(query={'_id': _id}))[0]['state']
