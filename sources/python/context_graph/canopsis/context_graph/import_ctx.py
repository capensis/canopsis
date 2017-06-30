#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
import jsonschema
import threading
import ijson
import time
import Queue

#FIXME : move the check of the element in the superficial check method

CONF_FILE = 'context_graph/manager.conf'
CATEGORY = "IMPORTCONTEXT"

def execution_time(exec_time):
    """Return from exec_time a human readable string that represent the
    execution time in a human readable format.
    :param exec_time: the execution time
    :type exec_time: a float"""

    exec_time = int(exec_time) # we do not care of everything under the second

    hours =  exec_time / 3600
    minutes = (exec_time - 3600 * hours) / 60
    seconds = exec_time - (hours * 3600) - (minutes * 60)

    return "{0}:{1}:{2}".format(str(hours).zfill(2),
                                str(minutes).zfill(2),
                                str(seconds).zfill(2))

class ExceptionThread(threading.Thread):
    """Wrapper to support handling exception"""

    def __init__(self, *args, **kwargs):
        super(ExceptionThread, self).__init__(*args, **kwargs)
        self.except_queue = Queue.Queue()

    def run(self):
        try:
            super(ExceptionThread, self).run()
        except Exception as e:
            self.except_queue.put_nowait(e)

class ImportKey:
    """Usefull values for the import."""

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

    EVT_IMPORT_UUID = "jobs_uuid"
    EVT_JOBID = "jobid"

    JOB_ID = "importctx_{0}"

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
                             "already exist.".format(uuid))

        new_status = {ImportKey.F_ID: uuid,
                      ImportKey.F_CREATION: time.asctime(),
                      ImportKey.F_STATUS: ImportKey.ST_PENDING}

        self[self.STORAGE].put_element(new_status)

    def on_going_in_db(self):
        """Check if there is an import on going
        :return: True if an import is on going
        """
        result = list(self[self.STORAGE].find_elements(
            query={ImportKey.F_STATUS: ImportKey.ST_ONGOING}))

        return len(result) == 1

    def pending_in_db(self):
        """Check if there is a pending import in database
        :return: True if a pending import in database
        """
        result = list(self[self.STORAGE].find_elements(
            query={ImportKey.F_STATUS: ImportKey.ST_PENDING}))

        return len(result) > 0

    def check_id(self, _id):
        """Check if an id is already taken.
        :param _id: the id to check
        :return: True if the id is in db. False otherwise
        """
        result = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_ID: _id}))

        return len(result) == 1

    def get_import_status(self, _id):
        """Return the state of an import.
        :param _id: the id of the import
        :return dict: the report.
        """
        status = list(self[self.STORAGE].get_elements(
            query={ImportKey.F_ID: _id}))

        if len(status) > 0:
            return status[0]

        return None


class ContextGraphImport(ContextGraph):
    """The manager in charge of an import of a context.
    """

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
    K_ENABLED = "enabled"

    # If you add an action, remember to add in the a_pattern string in method
    # import_checker
    A_DELETE = "delete"
    A_CREATE = "create"
    A_UPDATE = "update"
    A_DISABLE = "disable"
    A_ENABLE = "enable"


    __A_PATTERN = "^delete$|^create$|^update$|^disable$|^enable$"
    __T_PATTERN = "^resource$|^component$|^connector$|^watcher$"
    __CI_REQUIRED = [K_ID,
                     K_ACTION,
                     K_TYPE]
    __LINK_REQUIRED = [K_FROM,
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

            K_MEASUREMENTS: {"type": "array",
                             "items": {
                                 "type": "string"}},
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


    def __init__(self, logger=None, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraphImport, self).__init__(*args, **kwargs)

        if logger is not None:
            self.logger = logger

        self.entities_to_update = {}
        self.update = {}
        self.delete = []

    @classmethod
    def check_element(cls, element, type_):
        """Check an element with a schema schema specified by his type.
        :param element: the element to check
        :param type_: the expected type of the element
        :raise: ValueError if the type_ is not correct
        :raise: ValidationError if the element does not match the schema
        """

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

        If a ci or link does not match the schema a ValidationError is raised.
        :param json: the json with every actions required for the update
        :param rtype: a dict with the entity id as a key and the entity as
        a value
        """
        # a set so no duplicate ids without effort and lower time complexity
        ids_cis = set()
        ids_links = set()

        def __get_entities_to_update_links(file_):
            """Parse the file_ to extract every link"""
            fd = open(file_, 'r')
            for link in ijson.items(fd, "{0}.item".format(self.K_LINKS)):
                self.check_element(link, self.K_LINKS)
                for id_ in link[self.K_FROM]:
                    ids_links.add(id_)

                # ids_cis.add(ci[self.K_ID])
                ids_links.add(link[self.K_TO])
            fd.close()

        def __get_entities_to_update_cis(file_):
            """Parse the file_ to extract every CI"""
            fd = open(file_, 'r');
            for ci in ijson.items(fd, "{0}.item".format(self.K_CIS)):
                self.check_element(ci, self.K_CIS)
                ids_cis.add(ci[self.K_ID])

                # we need to retreive every related entity to update the links
                if ci[self.K_ACTION] == self.A_DELETE:
                    # FIXME do the get_entities_by_id in one call Then add all
                        # impacts depends
                    entity = self.get_entities_by_id(ci[self.K_ID])[0]

                    for id_ in entity["depends"] + entity["impact"]:
                        ids_cis.add(id_)
            fd.close()

        cis_thd = ExceptionThread(group=None,
                                   target=__get_entities_to_update_cis,
                                   name="cis_thread",
                                   args=(file_,))

        links_thd = ExceptionThread(group=None,
                                     target=__get_entities_to_update_links,
                                     name="links_thread",
                                     args=(file_,))

        threads = [cis_thd, links_thd]

        cis_thd.start()
        links_thd.start()

        cis_thd.join()
        links_thd.join()

        # Unqueue an raise existing exception
        for thread in threads:
            try:
                excep = thread.except_queue.get_nowait()
            except Queue.Empty:
                pass
            else:
                self.logger.error("Exception in {0}".format(thread.getName()))
                self.logger.exception(excep)
                raise excep

        ids = ids_links.union(ids_cis)
        result = self.get_entities_by_id(list(ids))
        ctx = {}

        # transform "depends" and "impact" list in set for improved performance
        for doc in result:
            doc[self.K_DEPENDS] = set(doc[self.K_DEPENDS])
            doc[self.K_IMPACT] = set(doc[self.K_IMPACT])
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
        for ent_id in entity[self.K_DEPENDS]:
            if ent_id in self.delete:
                # the entity of id ent_id is already deleted, skipping
                continue

            if ent_id not in self.update:
                self.update[ent_id] = self.entities_to_update[ent_id].copy()
            try:
                self.update[ent_id][self.K_IMPACT].remove(id_)
            except ValueError:
                raise ValueError("Try to remove {0} from impacts field of"\
                                 "entity {1}.".format(id_, ent_id))

        # Update the impact/depends link
        for ent_id in entity[self.K_IMPACT]:
            if ent_id in self.delete:
                # the entity of id ent_id is already deleted, skipping
                continue

            if ent_id not in self.update:
                self.update[ent_id] = self.entities_to_update[ent_id].copy()
            try:
                self.update[ent_id][self.K_DEPENDS].remove(id_)
            except ValueError:
                raise ValueError("Try to remove {0} from impacts field of"\
                                 "entity {1}.".format(id_, ent_id))

        if id_ in self.update:
            self.update.pop(id_)

        self.delete.append(id_)

    def __a_update_entity(self, ci):
        """Update the entity with the information stored into the ci and store
        the result into self.update.

        If the entity to be updated is not initially store in the context,
        a ValueError will be raised.

        :param ci: the ci (see the JSON specification).
        """

        if ci[self.K_ID] not in self.entities_to_update:
            desc = "The ci of id {0} does not match any existing"\
                   " entity.".format(ci[self.K_ID])
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
        if ci[self.K_ID] in self.entities_to_update:
            desc = "The ci of id {0} match an existing entity.".format(
                ci["_id"])
            raise ValueError(desc)

        # set default value for required fields
        if not ci.has_key(self.K_NAME):
            ci[self.K_NAME] = ci[self.K_ID]
        if not ci.has_key(self.K_DEPENDS):
            ci[self.K_DEPENDS] = set()
        else:
            ci[self.K_DEPENDS] = set(ci[self.K_DEPENDS])
        if not ci.has_key(self.K_IMPACT):
            ci[self.K_IMPACT] = set()
        else:
            ci[self.K_IMPACT] = set(ci[self.K_IMPACT])
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
        if state != self.A_DISABLE and state != self.A_ENABLE:
            raise ValueError("{0} is not a valid state.".format(state))

        id_ = ci[self.K_ID]

        if id_ not in self.entities_to_update:
            desc = "The ci of id {0} does not match any existing"\
                   " entity.".format(id_)
            raise ValueError(desc)

        # If the entity is not in the update dict, add it
        if id_ not in self.update:
            self.update[id_] = self.entities_to_update[id_].copy()

        if state == self.A_DISABLE:
            key_history = "disable_history"
            key = self.K_DISABLE
            self.update[id_][self.K_INFOS][self.K_ENABLED] = False
        else:
            key_history = "enable_history"
            key = self.K_ENABLE
            self.update[id_][self.K_INFOS][self.K_ENABLED] = True

        # Update entity the fields enable/disable of infos
        timestamp = ci[self.K_PROPERTIES][key]

        if not isinstance(timestamp, list):
            if timestamp is None:
                timestamp = []
            else:
                timestamp = [timestamp]

        if self.update[id_][self.K_INFOS].has_key(key_history):

            if self.update[id_][self.K_INFOS][key_history] is None:
                self.update[id_][self.K_INFOS][key_history] = timestamp
            else:
                self.update[id_][self.K_INFOS][key_history] += timestamp

        else:
            self.update[id_][self.K_INFOS][key_history] = timestamp

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

        for id_ in link[self.K_FROM]:
            if id_ not in self.update:
                self.update[id_] = self.entities_to_update[id_]

        if link[self.K_TO] not in self.update:
            self.update[link[self.K_TO]] = self.entities_to_update\
                                           [link[self.K_TO]]

        for id_ in link[self.K_FROM]:
            self.update[id_][self.K_IMPACT].remove(link[self.K_TO])
            self.update[link[self.K_TO]][self.K_DEPENDS].remove(id_)

    def __a_update_link(self, link):
        raise NotImplementedError()

    def __a_create_link(self, link):
        """Create a link between two entity and store the modify entities
        into self.udpate.

        :param link: the link that identify a link (see the JSON specification).
        """

        if link[self.K_TO] not in self.update:
            self.update[link[self.K_TO]] = self.entities_to_update\
                                           [link[self.K_TO]]

        for ci_id in link[self.K_FROM]:
            if ci_id not in self.update:
                self.update[ci_id] = self.entities_to_update[ci_id]

            self.update[ci_id][self.K_IMPACT].add(link[self.K_TO])
            self.update[link[self.K_TO]][self.K_DEPENDS].add(ci_id)

    def __a_disable_link(self, link):
        raise NotImplementedError()

    def __a_enable_link(self, link):
        raise NotImplementedError()

    @classmethod
    def __superficial_check(cls, fd):
        """Check if the cis and links field are a list. If not, raise a
        jsonschema.ValidationError. It move the cursor of the fd to back 0."""

        cis_end = False
        links_end = False

        parser = ijson.parse(fd)
        for prefix, event, value in parser:
            if prefix == "cis":
                if event == "end_array":
                    cis_end = True
            if prefix == "links":
                if event == "end_array":
                    links_end = True

        fd.seek(0)

        if (cis_end == True) and (links_end == True):
            return True
        elif (cis_end == False) and (links_end == False):
            raise jsonschema.ValidationError(
                "CIS and LINKS should be an array.")
        elif cis_end == False:
            raise jsonschema.ValidationError("CIS should be an array.")
        elif links_end == False:
            raise jsonschema.ValidationError("LINKS should be an array.")


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

        start = time.time()
        self.__superficial_check(fd)
        end = time.time()
        self.logger.debug("Import {0} : superficial"\
                          " check {1}.".format(uuid,
                                               execution_time(end - start)))

        start = time.time()
        self.entities_to_update = self.__get_entities_to_update(file_)
        end = time.time()

        self.logger.debug("Import {0} : get_entities_to"\
                          "_update {1}.".format(uuid,
                                                execution_time(end - start)))

        # Process cis list
        start = time.time()
        for ci in ijson.items(fd, "{0}.item".format(self.K_CIS)):
            self.logger.debug("Current ci : {0}".format(ci))
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
        end = time.time()
        self.logger.debug("Import {0} : update cis {1}.".format(uuid, execution_time(end - start)))

        fd.seek(0)

        # Process link list
        start = time.time()
        for link in ijson.items(fd, "{0}.item".format(self.K_LINKS)):
            self.logger.debug("Current link : {0}".format(link))
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
                raise ValueError("The action {0} is not recognized".format(
                    link[self.K_ACTION]))
        end = time.time()
        self.logger.debug("Import {0} : update links"\
                          " {1}.".format(uuid, execution_time(end - start)))

        for id_ in self.update:
            if id_ in self.delete:
                desc = "The entity {0} to be deleted is updated in "\
                       "the same import. Update aborted.".format(id_)
                raise ValueError(desc)

        updated_entities = len(self.update)
        deleted_entities = len(self.delete)

        for entity in self.update.values():
            entity[self.K_IMPACT] = list(entity[self.K_IMPACT])
            entity[self.K_DEPENDS] = list(entity[self.K_DEPENDS])

        start = time.time()
        self._put_entities(self.update.values())
        end = time.time()
        self.logger.debug("Import {0} : push updated"\
                          " entities {1}.".format(uuid,
                                                  execution_time(end - start)))

        start = time.time()
        self._delete_entities(self.delete)
        end = time.time()
        self.logger.debug("Import {0} : delete entities"\
                          " {1}.".format(uuid, execution_time(end - start)))

        self.clean_attributes()
        return updated_entities, deleted_entities
