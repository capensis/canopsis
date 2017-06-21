# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.model import Parameter
from canopsis.configuration.configurable.decorator import add_category
from canopsis.event import forger

from canopsis.selector.links import build_all_links

import time

CONF_PATH = 'context_graph/manager.conf'
CATEGORY = 'CONTEXTGRAPH'
CONTENT = [
    Parameter('event_types', Parameter.array()),
    Parameter('extra_fields', Parameter.array()),
    Parameter('authorized_info_keys', Parameter.array())
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class ContextGraph(MiddlewareRegistry):
    """ContextGraph"""

    ENTITIES_STORAGE = 'entities_storage'
    ORGANISATIONS_STORAGE = 'organisations_storage'
    USERS_STORAGE = 'measurements_storage'
    NAME = 'name'

    RESOURCE = "resource"
    COMPONENT = "component"
    CONNECTOR = "connector"

    @property
    def event_types(self):
        """
        Array of event_type
        """

        if not hasattr(self, '_event_types'):
            self.event_types = None

        return self._event_types

    @event_types.setter
    def event_types(self, value):
        if value is None:
            value = []

        self._event_types = value

    @property
    def extra_fields(self):
        """
        Array of fields to save from event in entity.
        """

        if not hasattr(self, '_extra_fields'):
            self.extra_fields = None

        return self._extra_fields

    @extra_fields.setter
    def extra_fields(self, value):
        if value is None:
            value = []

        self._extra_fields = value

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
            id_ = event["component"]
        elif source_type == cls.RESOURCE:
            id_ = "{0}/{1}".format(event["resource"], event["component"])
        elif source_type == cls.CONNECTOR:
            id_ = "{0}/{1}".format(event["connector"], event["connector_name"])
        else:
            error_desc = ("Event type should be 'connector', 'resource' or"
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

    def keys_info_filter(self, info):
        """Remove non authorized key present in the info field. See the
        configuration path givent to the class, by default
        etc/context_graph/manager.conf
        """

        if not hasattr(self, "authorized_info_keys"):
            values = values = self.conf.get(CATEGORY)
            self.authorized_info_keys = values.get("authorized_info_keys").value

        for key in info.keys():
            if key not in self.authorized_info_keys:
                info.pop(key)

    def __init__(self, event_types=None, extra_fields=None, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraph, self).__init__(*args, **kwargs)
        self.collection_name = 'default_entities'

        if event_types is None:
            self.event_types = event_types

        if extra_fields is None:
            self.extra_fields = extra_fields

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
            self[ContextGraph.ENTITIES_STORAGE].get_elements(query=query))

        return result

    def _put_entities(self, entities):
        """
        Store entities into database. Do no use this function unless you know
        exactly what you are doing. This function does not update the
        impact and depends links between entities. If you want these feature,
        use create_entity.

        :param entities: the entities to store in database
        """
        if not isinstance(entities, list):
            entities = [entities]

        self[ContextGraph.ENTITIES_STORAGE].put_elements(entities)

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
        self[ContextGraph.ENTITIES_STORAGE].remove_elements(entities)

    def get_all_entities_id(self):
        """
        Get the ids of every stored entities.
        TODO: use an iterator instead

        :return type: a set with every entities id.
        """
        entities = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={}))
        ret_val = set([])
        for i in entities:
            ret_val.add(i['_id'])
        return ret_val

    @classmethod
    def _enable_entity(cls, entity):
        """Enable an entity. If the given entity does not have an infos.tate
        field create it with a the enable status then create an
        infos.enable_history field with the current timestamp.

        :param entity: an entity
        :type entity: a dict
        """
        if "enabled" not in entity["infos"]:
            entity["infos"]["enabled"] = True

            if "enable_history" not in entity["infos"]:
                entity["infos"]["enable_history"] = [int(time.time())]
            else:
                infos = entity["infos"]
                if not isinstance(entity["infos"]["enable_history"], list):
                    infos["enable_history"] = [infos["enable_history"]]

                infos["enable_history"].append(int(time.time()))

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

        entities = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': entity["_id"]}))

        if len(entities) > 0:
            desc = "An entity  id {0} already exist".format(entities[0]["_id"])
            raise ValueError(desc)

        # update depends/impact links
        status = {"insertions": entity["depends"],
                  "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"],
                                                      status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = {"insertions": entity["impact"],
                  "deletions": []}
        updated_entities = self.__update_dependancies(entity["_id"],
                                                      status, "impact")
        updated_entities.append(entity)
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # rebuild selectors links
        if entity['type'] != 'selector':
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

            return {"deletions": list(deletions),
                    "insertions": list(insertions)}

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
        updated_entities = self.__update_dependancies(entity["_id"],
                                                      status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = compare_change(old_entity["impact"], entity["impact"])
        updated_entities = self.__update_dependancies(entity["_id"],
                                                      status, "impact")

        updated_entities.append(entity)
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # rebuild selectors links
        if entity['type'] != 'selector':
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
            entity = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query={'_id': id_}))[0]
        except IndexError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        # update depends/impact links
        status = {"deletions": entity["depends"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(id_,
                                                      status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        # update impact/depends links
        status = {"deletions": entity["impact"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(id_,
                                                      status, "impact")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        self[ContextGraph.ENTITIES_STORAGE].remove_elements(ids=[id_])

        # rebuild selectors links
        build_all_links(self)

    def get_entities(self,
                     query={},
                     projection=None,
                     limit=0,
                     start=0,
                     sort=False,
                     with_count=False):
        """Retreives entities matching the query and the projection.
        """
        # TODO handle projection, limit, sort, with_count
        # TODO complete docstring

        if not isinstance(query, dict):
            raise TypeError("Query must be a dict")

        result = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query=query,
            limit=limit,
            skip=start,
            sort=sort,
            projection=projection,
            with_count=with_count
        ))

        return result

    def get_event(self, entity, event_type='check', **kwargs):
        """Get an event from an entity.

        :param dict entity: entity to convert to an event.
        :param str event_type: specific event_type. Default check.
        :param dict kwargs: additional fields to put in the event.
        :rtype: dict
        """

        # keys from entity that should not be in event
        delete_keys = ["_id", "depends", "impact", "type", "measurements",
                       "infos"]

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
        col = self[ContextGraph.ENTITIES_STORAGE]._backend
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
        col = self[ContextGraph.ENTITIES_STORAGE]._backend
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
