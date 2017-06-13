# -*- coding: utf-8 -*- 
from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from json import loads


def build_links(selector_entity):
    """
        check and build links of a selector,
    """
    from canopsis.context_graph.manager import ContextGraph
    context_graph = ContextGraph()
    mfilter = loads(selector_entity['infos']['mfilter'])
    dep = context_graph.get_entities(
        query=mfilter,
        projection={'_id':1}
    )
    for i in dep:
        if i['_id'] not in selector_entity['depends']:
            selector_entity['depends'].append(i['_id'])
    context_graph.update_entity(selector_entity)

def build_all_links():
    """
        check and rebuild links of all selector entities
    """
    selectors = context_graph.get_entities(query={'type': 'selector'})
    for i in selectors:
        build_links(i)

