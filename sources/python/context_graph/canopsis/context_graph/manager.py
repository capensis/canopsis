# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import copy
import jsonschema
import time

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.middleware.core import Middleware
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.link_builder.link_builder import HypertextLinkManager
from canopsis.event import forger
from canopsis.watcher.links import build_all_links

from canopsis.confng import Configuration, Ini

CONF_PATH = 'context_graph/manager.conf'
CTX_HYPERLINK = "hypertextlink_conf"

DEFAULT_SCHEMA_ID = "context_graph.filter_infos"


class InfosFilter(object):
    """Class use to clean the infos field of an entity"""

    OBJ_STORAGE = 'storage-default://'

    def __init__(self, logger=None):
        super(InfosFilter, self).__init__()
        self.obj_storage = Middleware.get_middleware_by_uri(
            self.OBJ_STORAGE, table='schemas')

        self.config = Configuration.load(CONF_PATH, Ini)
        self.logger = logger

        self.event_types = self.config.CONTEXTGRAPH.get(
            'event_types', []).split(',')
        self.extra_fields = self.config.CONTEXTGRAPH.get(
            'extra_fields', []).split(',')
        self.schema_id = self.config.INFOS_FILTER.get(
            'schema_id', DEFAULT_SCHEMA_ID)

        self.reload_schema()

    def reload_schema(self):
        """Reload the schema and regenerate the internal structure used to
        filter the infos dict."""

        try:
            schemas = self.obj_storage.get_elements(
                query={"_id": self.schema_id}, projection={"_id": 0})
        except IndexError:
            raise ValueError("No infos schema found in database.")

        if len(schemas) > 0:
            self._schema = schemas[0]
            if not isinstance(self._schema, dict):
                raise ValueError("The schema should be a dict not"
                                 " a {0}.".format(type(self._schema)))

    def __clean(self, infos, iteration_dict, schema):
        """Recursive method use to clean the infos dict following the given
        schema. If a key or a sub key of infos is not in schema, it will
        be deleted.

        :param dict infos: the info field to clean
        :pararm dict iteration_dict: a copy of infos used to iterate over every
        keys in infos
        :param dict schema: the schema used to clean select the key to delete.
        """

        for key in iteration_dict:
            if key not in schema:
                infos.pop(key)
            elif isinstance(iteration_dict[key], dict):
                self.__clean(infos[key], infos[key].copy(), schema[key])

    def filter(self, infos):
        """Filter the fieds in infos. If a key from infos did not exist in the
        schema, it will deleted. If a a field type did not match the expected
        one or a required field is missing, the error will logged and the
        filtering will be stopped.

        :param dict infos: the dict to parse
        """

        try:
            jsonschema.validate(infos, self._schema)
        except jsonschema.ValidationError as v_err:
            self.logger.warning(v_err.message)
        try:
            schema = self._schema["schema"]["properties"]
        except KeyError:
            self.logger.warning("No properties field")
        else:
            self.__clean(infos, copy.deepcopy(infos), schema)


class ContextGraph(MiddlewareRegistry):
    """ContextGraph"""

    ENTITIES_STORAGE = 'entities_storage'
    ORGANISATIONS_STORAGE = 'organisations_storage'
    MEASUREMENTS_STORAGE = 'measurements_storage'
    NAME = 'name'

    RESOURCE = "resource"
    COMPONENT = "component"
    CONNECTOR = "connector"

    @classmethod
    def get_id(cls, event):
        """Return the id extracted from the event as a string.
        If the event come from a component, the id will follow the pattern
        "component". If the event come from resource, the id will follow the
        pattern resource/component. If the event come from a connector, the
        id will follow the pattern "connector/connector_name".

        :param event: the event from which we extract the id.
        :return type: the id as a string
        """
        source_type = ''
        if 'source_type' in event:
            source_type = event['source_type']
        elif 'type' in event:
            source_type = event['type']

        id_ = ""

        if source_type == cls.COMPONENT:
            try:
                id_ = event["component"].encode('utf-8')
            except UnicodeEncodeError:
                id_ = event['component']

        elif source_type == cls.RESOURCE:
            try:
                id_ = "{0}/{1}".format(
                    event["resource"].encode('utf-8'),
                    event["component"].encode('utf-8')
                )
            except UnicodeEncodeError:
                id_ = "{0}/{1}".format(
                    event["resource"],
                    event["component"]
                )

        elif source_type == cls.CONNECTOR:
            try:
                id_ = "{0}/{1}".format(
                    event["connector"].encode('utf-8'),
                    event["connector_name"].encode('utf-8')
                )
            except UnicodeEncodeError:
                id_ = "{0}/{1}".format(
                    event["connector"],
                    event["connector_name"]
                )

        else:
            error_desc = ("Event type should be 'connector', 'resource' or "
                          "'component' not {}.".format(source_type))
            raise ValueError(error_desc)

        return id_

    @classmethod
    def is_equals(cls, entity1, entity2):
        """Compare two entities and return True if the 2 entity are equals.
        False otherwise.

        :param entity1: the first entity
        :param entity2: the first entity
        :return type: a boolean
        """
        keys1 = entity1.keys()
        keys2 = entity2.keys()

        if len(keys1) != len(keys2):
            return False

        keys1 = sorted(keys1)
        keys2 = sorted(keys2)

        if keys1 != keys2:
            return False

        for key in keys1:
            if isinstance(entity1[key], list):
                # copy and sorte the list
                list1 = sorted(entity1[key][:])
                list2 = sorted(entity2[key][:])

                if list1 != list2:
                    return False
            else:
                if entity1[key] != entity2[key]:
                    return False

        return True

    DEFAULT_EVENT_TYPES = []
    DEFAULT_EXTRA_FIELDS = []

    def __init__(self, event_types=None, extra_fields=None, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraph, self).__init__(*args, **kwargs)
        self.collection_name = 'default_entities'

        self.config = Configuration.load(CONF_PATH, Ini)
        self.storages = {
            self.ENTITIES_STORAGE: Middleware.get_middleware_by_uri(self.config.CONTEXTGRAPH.entities_storage_uri),
            self.ORGANISATIONS_STORAGE: Middleware.get_middleware_by_uri(self.config.CONTEXTGRAPH.organisations_storage_uri),
            self.MEASUREMENTS_STORAGE: Middleware.get_middleware_by_uri(self.config.CONTEXTGRAPH.measurements_storage_uri),
        }
        # FIXIT: as event_types and extra_fields are set from __init__
        # arguments... what's the point of this?
        self.event_types = self.config.CONTEXTGRAPH.get(
            'event_types', self.DEFAULT_EVENT_TYPES).split(',')
        self.extra_fields = self.config.CONTEXTGRAPH.get(
            'extra_fields', self.DEFAULT_EXTRA_FIELDS).split(',')

        if event_types is None:
            self.event_types = self.DEFAULT_EVENT_TYPES

        if extra_fields is None:
            self.extra_fields = self.DEFAULT_EXTRA_FIELDS

        # For links building
        self.at_manager = AssociativeTableManager(logger=self.logger)
        self.hypertextlink_conf = self.config.CONTEXTGRAPH.get(
            CTX_HYPERLINK, "")
        if self.hypertextlink_conf != "":
            atable = self.at_manager.get(self.hypertextlink_conf)
            if atable is not None:
                conf = atable.get_all()
                self.hlb_manager = HypertextLinkManager(conf, self.logger)

        self.filter_ = InfosFilter(logger=self.logger)

    def get_entities_by_id(self, _id):
        """
        Retreive the entity identified by an id. If id is a list of id,
        get_entities_by_id return every entities who match the ids present
        in the list

        :param id: the id of an entity. id can be a list
        :return type: a list of entity
        """

        query = {"_id": None}
        if isinstance(_id, list):
            query["_id"] = {"$in": _id}
        else:
            query["_id"] = _id

        result = list(
            self.storages[ContextGraph.ENTITIES_STORAGE].get_elements(query=query))

        return result

    def _put_entities(self, entities):
        """
        Store entities into database. Do no use this function unless you know
        exactly what you are doing. This function does not update the
        impact and depends links between entities. If you want these feature,
        use create_entity. Nor add/update the the enable, enable_history or
        disable_history.

        :param entities: the entities to store in database
        """
        if not isinstance(entities, list):
            entities = [entities]

        for entity in entities:
            self.filter_.filter(entity["infos"])

        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(entities)

        # rebuild selectors links
        build_all_links(self)

    def _delete_entities(self, entities):
        """
        Remove entities from database. Do no use this function unless you know
        exactly what you are doing. This function does not update the
        impact and depends links between entities. If you want these feature,
        use delete_entity.

        :param entities: the entities to remove from database
        """
        if not isinstance(entities, list):
            entities = [entities]
        self.storages[ContextGraph.ENTITIES_STORAGE].remove_elements(entities)

    def get_all_entities_id(self):
        """
        Get the ids of every stored entities.
        TODO: use an iterator instead

        :return type: a set with every entities id.
        """
        entities = list(
            self.storages[ContextGraph.ENTITIES_STORAGE].get_elements(query={}))
        ret_val = set([])
        for i in entities:
            ret_val.add(i['_id'])
        return ret_val

    @classmethod
    def _enable_entity(cls, entity, timestamp=None):
        """Enable an entity. If the given entity does not have an infos.enabled
        field create it with a the enable status then create an
        infos.enable_history field with the current timestamp.

        :param entity: an entity
        :type entity: a dict
        :param int timestamp: a timestamp. If None, the current timestamp will
        be use.
        """

        if timestamp is None:
            timestamp = int(time.time())

        entity["enabled"] = True

        if "enable_history" not in entity:
            entity["enable_history"] = [timestamp]
        else:
            if not isinstance(entity["enable_history"], list):
                entity["enable_history"] = [entity["enable_history"]]

            entity["enable_history"].append(timestamp)

    def create_entity(self, entity):
        """Create an entity in the context with the given entity. This will
        update the depends and impact links between entities. If they are
        one or several id in the fields depends and/or impact, this function
        will update every related entity by adding the entity id in the correct
        field.

        If the entity contains one or several unknowned entity id, a ValueError
        will be same id as

        If an entity is already store with same id as than entity, a ValueError
        will be raised with a description.

        If an entity does not have a field infos.state, it will enable
        the entity with the current timestamp. This will create two fields
        infos.state and a infos.enable_history. Infos.state will store the
        current state of the entity (enable) and infos.enable_history will
        store the enable time.

        Other exception maybe raised, see __update_dependancies.

        :param entity: the new entity.
        """

        # TODO add treatment to check if every required field are present
        self._enable_entity(entity)

        entities = list(self.storages[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': entity["_id"]}))

        if len(entities) > 0:
            desc = "An entity id {0} already exist".format(entities[0]["_id"])
            raise ValueError(desc)

        # update depends/impact links
        status = {"insertions": entity["depends"], "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "depends")
        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = {"insertions": entity["impact"], "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "impact")
        updated_entities.append(entity)

        for entity in updated_entities:
            self.filter_.filter(entity["infos"])

        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # rebuild selectors links
        if entity['type'] != 'watcher':
            build_all_links(self)

    def __update_dependancies(self, id_, status, dependancy_type):
        """Return the list of entities whom the depends or impact fields are
        updated by removing or adding the id_.

        Dependancy_type allow you to specify the type of links (impact/depends
        or depends/impact) link.

        Status must contain two keys : insertions and deletions. The key
        insertions "contains" a list of entities id that need to be updated by
        adding the from_entity id in the field impact/depends.The key
        deletions "contains" a list of entities id that need to be updated by
        removing the from_entity id in the field impact/depends.

        If some entities are not found with the ids from the list behind the
        keys "insertions and "deletion" a ValueError will be raised.

        :param id_: the entity to add or remove
        :param status: a dict with two keys "insertions" and "deletion"
        :param dependancy_type: a string. "depends" to update the depends/impact
        links and "impact" to update the impact/depends links

        :return type: a list of entities.
        """

        if dependancy_type == "depends":
            to = "impact"
        elif dependancy_type == "impact":
            to = "depends"
        else:
            raise ValueError(
                "Dependancy_type should be depends or impact not {0}.".format(
                    dependancy_type))

        updated_entities = []

        # retreive the entities that were link to from_entity to remove the
        # reference of the 'from_entity'
        entities = self.get_entities_by_id(status["deletions"])

        if len(entities) != len(status["deletions"]):
            raise ValueError("Could not find some entity in database.")

        # update the related entities
        for entity in entities:
            try:
                entity[to].remove(id_)
            except ValueError:
                desc = "No {0} id in the {1} of the entity of id {2}.".format(
                    id_, to, entity["_id"])
                self.logger.debug(desc)

        updated_entities += entities

        # retreive the entities that was not link to from_entity to add a
        # reference of the 'from_entity'
        entities = self.get_entities_by_id(status["insertions"])

        if len(entities) != len(status["insertions"]):
            raise ValueError("Could not find some entity in database.")

        # update the related entities
        for entity in entities:
            entity[to].append(id_)

        updated_entities += entities

        return updated_entities

    def update_entity(self, entity):
        """Update an entity identified by id_ with the given entity.
        If needed, the fields impact/depends of the related entity will be
        updated.

        If the entity does not exist exist in database, a ValueError will be
        raise.

        Other exception maybe raised, see __update_dependancies.

        :param entity: the entity updated
        """

        def compare_change(old, new):
            """Retourn a dict with the insertion and the deletion."""

            s_old = set(list(old))
            s_new = set(list(new))
            deletions = s_old.difference(s_new)
            insertions = s_new.difference(s_old)

            return {
                "deletions": list(deletions),
                "insertions": list(insertions)
            }

        try:
            old_entity = self.get_entities_by_id(entity["_id"])[0]
        except IndexError:
            raise ValueError(
                "The _id {0} does not match any entity in database.".format(
                    entity["_id"]))

        # check if the old entity differ from the updated one
        if self.is_equals(old_entity, entity):
            # nothing to do.
            return

        # update depends/impact links
        status = compare_change(old_entity["depends"], entity["depends"])
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "depends")
        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = compare_change(old_entity["impact"], entity["impact"])
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "impact")

        updated_entities.append(entity)
        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # rebuild selectors links
        if entity['type'] != 'watcher':
            build_all_links(self)

    def delete_entity(self, id_):
        """Delete an entity identified by id_ from the context.

        If needed, the fields impact/depends of the related entity will be
        updated.

        If the entity does not exist exist in database, a ValueError will be
        raise.

        Other exception maybe raised, see __update_dependancies.

        :param id_: the id of the entity to delete.
        """

        try:
            entity = list(self.storages[ContextGraph.ENTITIES_STORAGE].get_elements(
                query={'_id': id_}))[0]
        except IndexError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # update depends/impact links
        status = {"deletions": entity["depends"], "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "depends")
        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = {"deletions": entity["impact"], "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "impact")
        self.storages[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        self.storages[ContextGraph.ENTITIES_STORAGE].remove_elements(ids=[id_])

        build_all_links(self)

    def get_entities(self,
                     query={},
                     projection=None,
                     limit=0,
                     start=0,
                     sort=False,
                     with_count=False):
        """
        Retreives entities matching the query and the projection.

        :param dict query: set of couple of (field name, field value)
        :param int limit: max number of elements to get
        :param int start: first element index among searched list
        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order
        :type sort: list of {(str, {ASC, DESC}}), or str}
        :param bool with_count: If True (False by default), add count to
            the result

        :return: a list of entities
        :rtype: list of dict elements
        """

        if not isinstance(query, dict):
            raise TypeError("Query must be a dict")

        result = list(self.storages[ContextGraph.ENTITIES_STORAGE].get_elements(
            query=query,
            limit=limit,
            skip=start,
            sort=sort,
            projection=projection,
            with_count=with_count
        ))

        # Enrich each entity with http links
        for res in result:
            res['links'] = {}
            if hasattr(self, 'hlb_manager'):
                res['links'] = self.hlb_manager.links_for_entity(res)

        return result

    def get_event(self, entity, event_type='check', **kwargs):
        """Get an event from an entity.

        :param dict entity: entity to convert to an event.
        :param str event_type: specific event_type. Default check.
        :param dict kwargs: additional fields to put in the event.
        :rtype: dict
        """

        # keys from entity that should not be in event
        delete_keys = [
            "_id", "depends", "impact", "type", "measurements", "infos"
        ]

        kwargs['event_type'] = event_type

        # In some cases, name is present but is component in fact
        if 'name' in entity:
            if 'component' not in entity:
                entity['component'] = entity['name']
            entity.pop('name')

        # remove key that should not be in the event
        tmp = entity.copy()
        for key in delete_keys:
            try:
                tmp.pop(key)
            except KeyError:
                pass

        # fill kwargs with entity values
        for field in tmp:
            kwargs[field] = tmp[field]

        # forge the event
        result = forger(**kwargs)

        return result

    def get_graph_impact(self, _id, deepness=None):
        """
        """
        col = self.storages[ContextGraph.ENTITIES_STORAGE]._backend
        ag = []

        match = {'$match': {'_id': _id}}
        ag.append(match)

        glookup = {
            '$graphLookup': {
                'from': self.collection_name,
                'startWith': '$_id',
                'connectFromField': 'impact',
                'connectToField': '_id',
                'depthField': 'depth',
                'as': 'graph'
            }
        }
        if deepness is not None:
            glookup['$graphLookup']['maxDepth'] = deepness
        ag.append(glookup)

        res = col.aggregate(ag)
        return res['result'][0]

    def get_graph_depends(self, _id, deepness=None):
        """
        """
        col = self.storages[ContextGraph.ENTITIES_STORAGE]._backend
        ag = []

        match = {'$match': {'_id': _id}}
        ag.append(match)

        glookup = {
            '$graphLookup': {
                'from': self.collection_name,
                'startWith': '$_id',
                'connectFromField': 'depends',
                'connectToField': '_id',
                'depthField': 'depth',
                'as': 'graph'
            }
        }
        if deepness is not None:
            glookup['$graphLookup']['maxDepth'] = deepness
        ag.append(glookup)

        res = col.aggregate(ag)
        return res['result'][0]

    def get_leaves_impact(self, _id, deepness=None):
        """
        """
        graph = self.get_graph_impact(_id, deepness)
        ret_val = []

        if graph['graph'] == []:
            graph = graph.pop('graph')
            ret_val.append(graph)

        for i in graph['graph']:
            if i['impact'] == [] or i['depth'] == deepness:
                ret_val.append(i)

        return ret_val

    def get_leaves_depends(self, _id, deepness=None):
        """
        """
        graph = self.get_graph_depends(_id, deepness)
        ret_val = []

        if graph['graph'] == []:
            graph = graph.pop('graph')
            ret_val.append(graph)

        for i in graph['graph']:
            if i['depends'] == [] or i['depth'] == deepness:
                ret_val.append(i)

        return ret_val
