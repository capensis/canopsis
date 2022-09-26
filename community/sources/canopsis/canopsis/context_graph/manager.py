# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.link_builder.link_builder import HypertextLinkManager
from canopsis.confng import Configuration, Ini
from canopsis.common.middleware import Middleware

CONF_PATH = 'etc/context_graph/manager.conf'


class ConfName(object):
    """List of values used for the configuration"""

    SECT_GCTX = "CONTEXTGRAPH"

    ENT_STORAGE = "entities_storage_uri"
    EVENT_TYPES = "event_types"
    EXTRA_FIELDS = "extra_fields"

    CTX_HYPERLINK = "hypertextlink_conf"


class ContextGraph(object):
    """ContextGraph"""

    NAME = 'name'

    RESOURCE = "resource"
    COMPONENT = "component"
    CONNECTOR = "connector"

    def __init__(self,
                 logger,
                 *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """

        parser = Configuration.load(CONF_PATH, Ini)
        section = parser.get(ConfName.SECT_GCTX)

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
                if 'val' in conf:
                    conf = conf['val']
                self.hlb_manager = HypertextLinkManager(conf, self.logger)

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

    def enrich_links_to_entity_with_alarm(self, entity, alarm):
        if hasattr(self, 'hlb_manager'):
            links = self.hlb_manager.links_for_entity(
                entity, options={'alarm': alarm})
            return links
        return {}

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
