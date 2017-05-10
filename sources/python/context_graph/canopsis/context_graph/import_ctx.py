#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
import jsonschema
import ijson
import time

class ImportKey:

    # Status
    ST_PENDING = "pending"
    ST_ONGOING = "ongoing"
    ST_FAILED = "failed"
    ST_DONE = "done"

    # Fields
    F_CREATION = "creation"
    F_ID = "_id"
    F_STATUS = "status"
    F_INFO = "info"
    F_EXECTIME = "exec_time"
    F_START = "start"
    F_STATS = "stats"
    F_DELETED = "deleted"
    F_UPDATED = "updated"

    # daemon PID file
    PID_FILE = "/opt/canopsis/var/run/importd.pid"
    # import file pattern
    IMPORT_FILE = "/opt/canopsis/tmp/import_{0}.json"

CONF_FILE = 'context_graph/manager.conf'
CATEGORY = "IMPORTCONTEXT"

@conf_paths(CONF_FILE)
@add_category(CATEGORY)
class Manager(MiddlewareRegistry):
    """The manager use to interact with the default_importgraph collecion."""

    STORAGE = 'import_storage'
    DATE_FORMAT = "%a %b %d %H:%M:%S %Y"

    def __init__(self, *args, **kwargs):
        """__init__
        :param *args:
        :param **kwargs:
        """
        super(Manager, self).__init__(*args, **kwargs)

    def get_next_uuid(self):
        """Retreive the uuid of the next import to process using his creation
        date.
        :return the uuid as a string or None if they are no new import
        to process.
        """
        imports = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_STATUS: ImportKey.ST_PENDING}))

        if len(imports) == 0:
            return None


        next_ = imports[0]
        next_[ImportKey.F_CREATION] = time.strptime(
            str(next_[ImportKey.F_CREATION]), self.DATE_FORMAT)

        for import_ in imports:

            if import_[ImportKey.F_CREATION] <\
               next_[ImportKey.F_CREATION]:
                next_ = import_

        return next_[ImportKey.F_ID]

    def is_present(self, uuid):
        """Return true if the given import uuid exist in database.
        :param uuid: the given import uuid
        :return type: boolean
        :return True if the uuid exist in database. False otherwise
        """

        imports = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_ID: uuid}))

        return True if len(imports) == 1 else False

    def update_status(self, uuid, infos):
        """Update an import status with the fields of kwargs.
        If a field present in kwargs is not intended to be updated, it will
        be silently ignored.
        If no import status are found for the uuid an exception ValueError
        will be raised.
        :param uuid: the uuid of the import to update:
        :param infos: the fields and values to update as a dict:
        """
        authorized_fields = [ImportKey.F_STATUS,
                             ImportKey.F_INFO,
                             ImportKey.F_EXECTIME,
                             ImportKey.F_START,
                             ImportKey.F_STATS]

        if not self.is_present(uuid):
            raise ValueError("No import with the given uuid"\
                             " ({0}).".format(uuid))

        new_status = {ImportKey.F_ID: uuid}

        for field in authorized_fields:
            if infos.has_key(field):
                new_status[field] = infos[field]

        self[self.STORAGE].put_element(new_status)

    def create_import_status(self, uuid):
        """Create a new import status. It will be create with the given uuid,
        the creation date corresponding to the call of this function and the
        PENDING status.
        :param uuid: the uuid as a string
        """

        if self.is_present(uuid):
            raise ValueError("An import status with the same uuid ({0}) "\
                             " already exist.")

        new_status = {ImportKey.F_ID: uuid,
                      ImportKey.F_CREATION: time.asctime(),
                      ImportKey.F_STATUS: ImportKey.ST_PENDING}

        self[self.STORAGE].put_element(new_status)

    def on_going_in_db(self):
        """
            check if there is an import on going

            :return: True if an import is on going
        """
        result = list(self[self.STORAGE].find_elements(
            query={ImportKey.F_STATUS: ImportKey.ST_ONGOING}))

        return len(result) == 1

    def pending_in_db(self):
        """
            check if there is a pending import in database

            :return: True if a pending import in database
        """
        result = list(self[self.STORAGE].find_elements(
            query={ImportKey.F_STATUS: ImportKey.ST_PENDING}))

        return len(result) > 0

    def check_id(self, _id):
        """
            check if an id is already taken
        """
        result = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_ID: _id}))

        return len(result) == 1

    def get_import_status(self, _id):
        """
        return the state of an import.
        :param _id: the id of the import
        :return type: a string containg one of the following value "pending",
        "ongoing","failed" or "done".
        """
        status = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_ID: _id}))[0]

        return status

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

    CIS_SCHEMA = {
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
            K_PROPERTIES: {"type": "object"}}}

    LINKS_SCHEMA = {
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
            K_PROPERTIES: {"type": "object"}}}


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

        if type_ == cls.K_LINKS:
            schema = cls.LINKS_SCHEMA
        elif type_ == cls.K_CIS:
            schema = cls.CIS_SCHEMA
        else:
            raise ValueError("Unknowed type {0}\n".format(type_))

        jsonschema.validate(element, schema)

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

    def import_context(self, uuid):
        """Import a new context.

        :param uuid: the uuid of the import to process
        :return type: a tuple (updated entities, deleted entities)
        """

        file_ = ImportKey.IMPORT_FILE.format(uuid)

        fd = open(file_, 'r')

        # In case the previous import failed and/or raise an exception, we\
            # clean now
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

        updated_entities = len(self.update)
        deleted_entities = len(self.delete)

        self._put_entities(self.update.values())
        self._delete_entities(self.delete)

        self.clean_attributes()
        return updated_entities, deleted_entities
