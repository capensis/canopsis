# -*- coding: utf-8 -*-

from canopsis.common.ws import route
from canopsis.linklist.manager import Linklist
from canopsis.entitylink.manager import Entitylink

link_list_manager = Linklist(*Linklist.provide_default_basics())
entity_link_manager = Entitylink(*Entitylink.provide_default_basics())

DEFAULT_ROUTE = 'linklist'


def exports(ws):

    @route(ws.application.delete,
           payload=['ids'],
           name=DEFAULT_ROUTE)
    def delete(ids):
        link_list_manager.remove(ids)
        ws.logger.info(u'Delete : {}'.format(ids))
        return True

    @route(
        ws.application.post,
        payload=['document'],
        name='linklist/put'
    )
    def put(document):
        ws.logger.debug({
            'document': document,
            'type': type(document)
        })

        link_list_manager.put(document)

        return True

    @route(ws.application.post,
           payload=['limit', 'start', 'sort', 'filter'],
           name=DEFAULT_ROUTE)
    def get(limit=0, start=0, sort=None, filter={}):
        result = link_list_manager.find(
            limit=limit,
            skip=start,
            _filter=filter,
            sort=sort,
            with_count=True
        )
        return result

    @route(
        ws.application.post,
        payload=['event'],
        name='entitylink'
    )
    def get_entitylink(event):
        return entity_link_manager.get_links_from_event(event)
