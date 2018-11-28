# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import copy
import jsonschema
import time

from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.link_builder.link_builder import HypertextLinkManager
from canopsis.common.utils import normalize_utf8
from canopsis.confng import Configuration, Ini
from canopsis.event import forger
from canopsis.logger import Logger, OutputFile
from canopsis.common.middleware import Middleware
from canopsis.watcher.links import build_all_links


CONF_PATH = 'etc/context_graph/manager.conf'


class ConfName(object):
    """List of values used for the configuration"""

    SECT_GCTX = "CONTEXTGRAPH"
    SECT_FILTER = "INFOS_FILTER"

    ENT_STORAGE = "entities_storage_uri"
    ORG_STORAGE = "organisations_storage_uri"
    USR_STORAGE = "users_storage_uri"
    MEASURMNT_STORAGE = "measurements_storage_uri"
    IMPORT_STORAGE = "import_storage_uri"
    EVENT_TYPES = "event_types"
    EXTRA_FIELDS = "extra_fields"

    SCHEMA_ID = "schema_id"
    CTX_HYPERLINK = "hypertextlink_conf"


class InfosFilter:
    """Class use to clean the infos field of an entity"""

    OBJ_STORAGE = "OBJECT_STORAGE"
    SCHEMA_ID = "schema_infos"
    LOG_NAME = "InfosFilter"
    LOG_PATH = "~/var/log/infos_filter.log"

    def __init__(self, config=None, logger=None):

        if logger is None:
            self.logger = Logger.get(self.LOG_NAME,
                                     self.LOG_PATH,
                                     output_cls=OutputFile)
        else:
            self.logger = logger

        if config is None:
            self.config = Configuration.load(CONF_PATH, Ini)
        else:
            self.config = config

        self.obj_storage = Middleware.get_middleware_by_uri(
            'storage-default://', table='schemas')

        section = self.config.get(ConfName.SECT_GCTX)
        self._event_types = section[ConfName.EVENT_TYPES]
        self._extra_fields = section[ConfName.EXTRA_FIELDS]

        section = self.config.get(ConfName.SECT_FILTER)
        self._schema_id = section[ConfName.SCHEMA_ID]

        self.reload_schema()

    def reload_schema(self):
        """Reload the schema and regenerate the internal structure used to
        filter the infos dict."""

        try:
            self._schema = self.obj_storage.get_elements(
                query={"_id": self._schema_id}, projection={"_id": 0})[0]
        except IndexError as exc:
            raise ValueError(
                "No infos schema found in database: {}".format(exc))

        if isinstance(self._schema, list):
            self._schema = self._schema[0]
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
            pass
        else:
            self.__clean(infos, copy.deepcopy(infos), schema)


class ContextGraph(object):
    """ContextGraph"""

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
        source_field = ''
        source_type = ''
        if 'source_type' in event:
            source_type = event['source_type']
            source_field = 'source_type'
        elif 'type' in event:
            source_type = event['type']
            source_field = 'type'

        connector = normalize_utf8(event.get('connector'))
        connector_name = normalize_utf8(event.get('connector_name'))
        component = normalize_utf8(event.get('component'))
        id_ = ""

        if source_type == cls.COMPONENT:
            id_ = component

        elif source_type == cls.RESOURCE:
            resource = normalize_utf8(event.get('resource'))
            id_ = "{}/{}".format(resource, component)

        elif source_type == cls.CONNECTOR:
            id_ = "{}/{}".format(connector, connector_name)

        else:
            error_desc = (
                "{} should be one of 'connector'"
                ", 'resource' or 'component' not {}.".format(
                    source_field, source_type
                )
            )
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

    @classmethod
    def create_entity_dict(cls,
                           id,
                           name,
                           etype,
                           depends=[],
                           impact=[],
                           measurements={},
                           infos={},
                           **kwargs):
        """
        Create an entity with following information and put it state at enable.

        :param id: the entity id
        :type id: a string
        :param name: the entity name
        :type name: a string
        :param etype: the entity type
        :type etype: a string
        :param depends: every entities that depends of the current entity
        :type depends: a list
        :param impact: every entities that depends of the current entity
        :type impact: a list
        :param measurements: measurements link to the current entity
        :type measurements: a dict
        :param infos: information related to the entity
        :type infos: a dict

        :return: a dict
        """
        ent = {
            '_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
            'infos': infos
        }

        for key in kwargs:
            ent[key] = kwargs[key]

        if etype == 'component':
            ent.pop("measurements")

        cls._enable_entity(ent)

        return ent

    def __init__(self,
                 logger,
                 *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """

        parser = Configuration.load(CONF_PATH, Ini)
        section = parser.get(ConfName.SECT_GCTX)

        self.collection_name = 'default_entities'

        self.at_storage = Middleware.get_middleware_by_uri(
            AssociativeTableManager.STORAGE_URI
        )

        self.ent_storage = Middleware.get_middleware_by_uri(
            section.get(ConfName.ENT_STORAGE)
        )

        self.logger = logger

        # For links building
        at_collection = self.at_storage._backend
        self.at_manager = AssociativeTableManager(logger=self.logger,
                                                  collection=at_collection)

        hypertextlink_conf = section.get(ConfName.CTX_HYPERLINK, "")
        self.event_types = section.get(ConfName.EVENT_TYPES, [])
        self.extra_fields = section.get(ConfName.EXTRA_FIELDS, [])

        if hypertextlink_conf != "":
            atable = self.at_manager.get(hypertextlink_conf)
            if atable is not None:
                conf = atable.get_all()
                self.hlb_manager = HypertextLinkManager(conf, self.logger)

        self.filter_ = InfosFilter(logger=self.logger)

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

        self.ent_storage.put_elements(entities)

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
        self.ent_storage.remove_elements(entities)

    def get_all_entities_id(self):
        """
        Get the ids of every stored entities.
        TODO: use an iterator instead

        :return type: a set with every entities id.
        """
        entities = list(
            self.ent_storage.get_elements(query={}))
        ret_val = set([])
        for i in entities:
            ret_val.add(i['_id'])

        return ret_val

    def get_entities_with_open_alarms(self, query, limit, offset):
        """
        Get a list of entities enhaced with open alarms data found with
        a given mongo filter.

        :param query: Custom mongodb filter for entities
        :type query: dict

        :return type: an array of entities including the open alarms.
        """

        pipeline = []
        col = self.ent_storage._backend
        match_query = {
            '$match': query
        }
        pipeline.append(match_query)

        join_alarms = {
            '$lookup': {
                'from': 'periodical_alarm',
                'localField': '_id',
                'foreignField': 'd',
                'as': 'alarm'
            }
        }
        pipeline.append(join_alarms)

        total_count = list(col.aggregate(
            pipeline + [{'$count': 'total_count'}]))[0]["total_count"]

        if offset > 0:
            set_offset = {
                '$skip': offset
            }
            pipeline.append(set_offset)

        if limit > 0:
            set_limit = {
                '$limit': limit
            }
            pipeline.append(set_limit)

        ignore_terminated_alarms = {
            '$addFields': {
                'alarm': {
                    '$filter': {
                        'input': '$alarm',
                        'as': 'alarm',
                        'cond': {
                            "$not": "$$alarm.v.resolved"
                        }
                    }
                }
            }
        }
        pipeline.append(ignore_terminated_alarms)

        transform_alarm_array_to_field = {
            '$addFields': {
                'alarm': {
                    "$cond": [
                        {"$eq": [{"$size": "$alarm"}, 1]},
                        {'$arrayElemAt': ["$alarm", 0]},
                        None
                    ]
                }
            }
        }
        pipeline.append(transform_alarm_array_to_field)

        entities = list(col.aggregate(pipeline))

        return {"total_count": total_count,
                "count": len(entities),
                "data": entities}

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

        entities = list(self.ent_storage.get_elements(
            query={'_id': entity["_id"]}))

        if len(entities) > 0:
            desc = "An entity id {0} already exist".format(entities[0]["_id"])
            raise ValueError(desc)

        # update depends/impact links
        status = {"insertions": entity["depends"], "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "depends")
        self.ent_storage.put_elements(updated_entities)

        # update impact/depends links
        status = {"insertions": entity["impact"], "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "impact")
        updated_entities.append(entity)

        for entity in updated_entities:
            self.filter_.filter(entity["infos"])

        self.ent_storage.put_elements(updated_entities)

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

        :raises: ValueError
        :return type: a list of entities.
        """

        if dependancy_type == "depends":
            field = "impact"
        elif dependancy_type == "impact":
            field = "depends"
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
                entity[field].remove(id_)
            except ValueError:
                desc = "No {0} id in the {1} of the entity of id {2}.".format(
                    id_, field, entity["_id"])
                self.logger.debug(desc)

        updated_entities += entities

        # retreive the entities that was not link to from_entity to add a
        # reference of the 'from_entity'
        entities = self.get_entities_by_id(status["insertions"])

        if len(entities) != len(status["insertions"]):
            raise ValueError("Could not find some entity in database.")

        # update the related entities
        for entity in entities:
            entity[field].append(id_)

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
        self.ent_storage.put_elements(updated_entities)

        # update impact/depends links
        status = compare_change(old_entity["impact"], entity["impact"])
        updated_entities = self.__update_dependancies(entity["_id"], status,
                                                      "impact")

        updated_entities.append(entity)
        self.ent_storage.put_elements(updated_entities)

        # rebuild selectors links
        if entity['type'] != 'watcher':
            build_all_links(self)

    def update_entity_body(self, entity):
        """Update the document corresponding to an entity.
        The fields impact/depends of the related entity will NOT be updated.

        :param entity: the entity updated
        """
        self.ent_storage.put_elements([entity])

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
            entity = list(self.ent_storage.get_elements(
                query={'_id': id_}))[0]
        except IndexError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # update depends/impact links
        status = {"deletions": entity["depends"], "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "depends")
        self.ent_storage.put_elements(updated_entities)

        # update impact/depends links
        status = {"deletions": entity["impact"], "insertions": []}
        updated_entities = self.__update_dependancies(id_, status, "impact")
        self.ent_storage.put_elements(updated_entities)

        self.ent_storage.remove_elements(ids=[id_])

        build_all_links(self)

    def get_entities(self,
                     query=None,
                     projection=None,
                     limit=0,
                     start=0,
                     sort=False,
                     with_count=False,
                     with_links=False):
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
        :param bool with_links: If True (False by default), add builded links

        :return: a list of entities
        :rtype: list of dict elements
        """

        if query is None:
            query = {}
        elif not isinstance(query, dict):
            raise TypeError("Query must be a dict")

        result = self.ent_storage.get_elements(
            query=query,
            limit=limit,
            skip=start,
            sort=sort,
            projection=projection,
            with_count=with_count
        )

        if with_count:
            count = result[1]
            # Don't invert those two lines to avoid duplicate results
            result = list(result[0])
        else:
            result = list(result)

        # Enrich each entity with http links
        for res in result:
            res['links'] = {}
            if with_links and hasattr(self, 'hlb_manager'):
                links = self.hlb_manager.links_for_entity(res)
                res['links'] = links

        if with_count:
            return result, count
        else:
            return result

    def get_entities_by_id(self, _id, with_links=False):
        """
        Retreive the entity identified by an id. If id is a list of id,
        get_entities_by_id return every entities who match the ids present
        in the list

        :param id: the id of an entity. id can be a list
        :param bool with_links: If True (False by default), add builded links
        :returns: a list of entity
        """

        query = {"_id": None}
        if isinstance(_id, list):
            query["_id"] = {"$in": _id}
        else:
            query["_id"] = _id

        result = self.get_entities(query=query, with_links=with_links)

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
            "_id", "depends", "impact", "type", "measurements", "infos",
            "last_state_change"
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
        """Return the impact graph of the entity design by _id.
        :param _id: the _id of the entity from the graph start
        :param deepness: the max depth of the graph.
        :return dict: the graph.
        """
        col = self.ent_storage._backend
        aggregate = []

        match = {'$match': {'_id': _id}}
        aggregate.append(match)

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
        aggregate.append(glookup)

        return list(col.aggregate(aggregate))[0]

    def get_graph_depends(self, _id, deepness=None):
        """Return the depends graph from the entity design by _id.
        :param _id: the _id of the entity from the graph start
        :param deepness: the max depth of the graph.
        :return dict: the graph.
        """
        col = self.ent_storage._backend
        aggregate = []

        match = {'$match': {'_id': _id}}
        aggregate.append(match)

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
        aggregate.append(glookup)

        return list(col.aggregate(aggregate))[0]

    def get_leaves_impact(self, _id, deepness=None):
        """Return the entities at the end of the impact graph from the entity
        design by _id.
        :param _id: the _id of the entity from the graph start
        :param deepness: the max depth of the graph.
        :return dict: the graph.
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
        """Return the entities at the end of the depends graph from the entity
        design by _id.
        :param _id: the _id of the entity from the graph start
        :param deepness: the max depth of the graph.
        :return dict: the graph.
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
