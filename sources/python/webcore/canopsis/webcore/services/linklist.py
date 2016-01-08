# -*- coding: utf-8 -*-

from canopsis.common.ws import route
from canopsis.linklist.manager import Linklist
from canopsis.entitylink.manager import Entitylink

link_list_manager = Linklist()
entity_link_manager = Entitylink()


def exports(ws):

    @route(ws.application.delete, payload=['ids'])
    def linklist(ids):
        link_list_manager.remove(ids)
        ws.logger.info('Delete : {}'.format(ids))
        return True

    @route(
        ws.application.post,
        payload=['document'],
        name='linklist/put'
    )
    def linklist(document):
        ws.logger.debug({
            'document': document,
            'type': type(document)
        })

        link_list_manager.put(document)

        return True

    @route(ws.application.post, payload=['limit', 'start', 'sort', 'filter'])
    def linklist(limit=0, start=0, sort=None, filter={}):
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
    def linklist(event):
        ws.logger.debug({'received event': event})
        return entity_link_manager.get_links_from_event(event)
