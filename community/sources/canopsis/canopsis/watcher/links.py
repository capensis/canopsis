# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from json import loads


def build_links(watcher_entity, context_graph):
    """
    Check and build links of a watcher.

    :param watcher_entity:
    :param context_graph:
    """
    jmfilter = watcher_entity.get('mfilter', None)

    if jmfilter is None:
        return

    mfilter = loads(jmfilter)

    entities = context_graph.get_entities(
        query=mfilter,
        projection={'_id': 1}
    )

    watcher_updated = False
    watcher_depends = set(watcher_entity.get('depends', []))

    for entity in entities:
        eid = entity['_id']
        if eid not in watcher_depends:
            watcher_updated = True
            watcher_entity['depends'].append(eid)

    # FIXIT: updating only when needed is required to avoid resetting
    # the watcher entity impact field.
    # Also it avoids useless updates in DB.
    if watcher_updated:
        context_graph.update_entity(watcher_entity)


def build_all_links(context_graph):
    """
    Check and rebuild links of all watcher entities

    :param context_graph:
    """
    watchers = context_graph.get_entities(query={'type': 'watcher'})
    for i in watchers:
        build_links(i, context_graph)
