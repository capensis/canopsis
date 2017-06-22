# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from json import loads


def build_links(watcher_entity, context_graph):
    """
    Check and build links of a watcher.

    :param watcher_entity:
    :param context_graph:
    """
    mfilter = loads(watcher_entity['infos']['mfilter'])
    dep = context_graph.get_entities(
        query=mfilter,
        projection={'_id': 1}
    )
    for i in dep:
        if i['_id'] not in watcher_entity['depends']:
            watcher_entity['depends'].append(i['_id'])
    context_graph.update_entity(watcher_entity)


def build_all_links(context_graph):
    """
    Check and rebuild links of all watcher entities

    :param context_graph:
    """
    watchers = context_graph.get_entities(query={'type': 'watcher'})
    for i in watchers:
        build_links(i, context_graph)
