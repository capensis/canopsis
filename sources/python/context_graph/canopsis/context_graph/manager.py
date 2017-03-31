# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.event import forger


@conf_paths('context_graph/manager.conf')
@add_category('CONTEXTGRAPH')
class ContextGraph(MiddlewareRegistry):
    """ContextGraph"""

    ENTITIES_STORAGE = 'entities_storage'
    ORGANISATIONS_STORAGE = 'organisations_storage'
    USERS_STORAGE = 'measurements_storage'

    RESOURCE = "resource"
    COMPONENT = "component"
    CONNECTOR = "connector"

    @classmethod
    def get_id(cls, event):
        """Return the id extracted from the event as a string
        :param event: the event from which we extract the id
        :return type: boolean a string
        """
        if '_id' in event:
            return event['_id']

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
            error_desc = "Event type should be 'connector', 'resource' or\
            'component' not {0}.".format(source_type)
            raise ValueError(error_desc)

        return id_

    @classmethod
    def is_equals(cls, entity1, entity2):
        """Compare two entities and return True if the 2 entity are equals.
        False otherwise"""
        keys1 = entity1.keys()
        keys2 = entity2.keys()

        if len(keys1) != len(keys2):
            return False

        sorted(keys1)
        sorted(keys2)

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

    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraph, self).__init__(*args, **kwargs)

    def get_entities_by_id(self, id):
        """
        Retreive the entity identified by his id. If id is a list of id,
        get_entity return every entities who match the ids present in the list

        :param id the id of the entity. id can be a list
        """
        query = {"_id": None}
        if isinstance(id, type([])):
            ids = []
            for i in id:
                ids.append(i)
            query["_id"] = {"$in": ids}
        else:
            query["_id"] = id

        result = list(
            self[ContextGraph.ENTITIES_STORAGE].get_elements(query=query))

        return result

    def put_entities(self, entities):
        """
        Store entities into database.
        """
        if not isinstance(entities, list):
            entities = [entities]
        self[ContextGraph.ENTITIES_STORAGE].put_elements(entities)

    def get_all_entities_id(self):
        """
            get all entities ids by types
        """
        entities = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={}))
        ret_val = set([])
        for i in entities:
            ret_val.add(i['_id'])
        return ret_val

    def create_entity(self, entity):
        """Create an entity in the contexte with the given entity."""
        # TODO add traitement to check if every required field are present

        entities = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query={'_id': entity["_id"]}))
        if len(entities) > 0:
            desc = "An entity  id {0} already exist".format(entities[0]["_id"])
            raise ValueError(desc)

        status = {"insertions": entity["depends"],
                  "deletions": []}
        updated_entities = self.__update_dependancies(entity, status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        status = {"insertions": entity["impact"],
                  "deletions": []}
        updated_entities = self.__update_dependancies(entity, status, "impact")
        updated_entities.append(entity)
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

    def __update_dependancies(self, from_entity, status, dependancy_type):
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

        for entity in entities:
            try:
                entity[to].remove(from_entity["_id"])
            except ValueError:
                # TODO add a log.debug
                pass

        updated_entities += entities

        # retreive the entities that was not link to from_entity to add a
        # reference of the 'from_entity'
        entities = self.get_entities_by_id(status["insertions"])

        if len(entities) != len(status["insertions"]):
            raise ValueError("Could not find some entity in database.")

        for entity in entities:
            entity[to].append(from_entity["_id"])

        updated_entities += entities

        return updated_entities

    def update_entity(self, entity):
        """Update an entity identified by id_ with the given entity."""

        def compare_change(old, new):
            """Retourn a dict with the insertion and the deletion"""

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

        status = compare_change(old_entity["depends"], entity["depends"])
        updated_entities = self.__update_dependancies(entity, status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        status = compare_change(old_entity["impact"], entity["impact"])
        updated_entities = self.__update_dependancies(entity, status, "impact")

        updated_entities.append(entity)
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

    def delete_entity(self, id_):
        """Delete an entity identified by id_ from the context."""
        try:
            entity = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query={'_id': id_}))[0]
        except IndexError:
            desc = "No entity found for the following id : {0}".format(id_)
            raise ValueError(desc)

        status = {"deletions": entity["depends"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(entity, status, "depends")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        status = {"deletions": entity["impact"],
                  "insertions": []}
        updated_entities = self.__update_dependancies(entity, status, "impact")
        self[ContextGraph.ENTITIES_STORAGE].put_elements(updated_entities)

        self[ContextGraph.ENTITIES_STORAGE].remove_elements(ids=[id_])

    def get_entities(self,
                     query={},
                     projection={},
                     limit=0,
                     start=0,
                     sort=False,
                     with_count=False):
        """Retreives entities matching the query and the projection.
        """
        # TODO handle projection, limit, sort, with_count

        if not isinstance(query, dict):
            raise TypeError("Query must be a dict")

        if not isinstance(projection, dict):
            raise TypeError("Projection must be a dict")

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
        delete_keys = ["_id", "depends", "impact", "type"]

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
