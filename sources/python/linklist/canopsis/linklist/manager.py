# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.context.manager import Context

CONF_PATH = 'linklist/linklist.conf'
CATEGORY = 'LINKLIST'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Linklist(MiddlewareRegistry):

    LINKLIST_STORAGE = 'linklist_storage'
    CONTEXT_CONFIGURABLE = 'configurable_linklist'
    TYPE = 'linklist'
    """
    Manage linklist information in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Linklist, self).__init__(*args, **kwargs)
        self.context = Context()

    def extra_delete(self, ids):
        """
        Remove extra fields persisted in a default storage.
        """
        self[Linklist.LINKLIST_STORAGE].remove_elements(ids=ids)

    def extra_put(self, entity_id, extra_fields, cache=False):
        """
        Allow persistance of a custom dict known as extra_fields

        :param entity_id: the identifier for the entity
        :param extra_fields: the dict that contains the fields to persist
        """
        self[Linklist.LINKLIST_STORAGE].put_element(
            _id=entity_id, element=extra_fields, cache=cache
        )

    def extra_find(self, selection):
        """
        Method getting extra fields from default storage
        It will merge the exra fields together to produce a single document.
        """

        # Gets ids from selection documents
        ids = [x['_id'] for x in selection]

        # find extra information for these ids
        self.logger.debug('ids for extra find is {}'.format(ids))
        if not ids:
            extras = []
        else:
            extras = self[Linklist.LINKLIST_STORAGE].get_elements(
                ids=ids
            )
            extras = list(extras)

        # Merge documents
        cache = {}
        for extra in extras:
            cache[extra['_id']] = extra

        # Build back information
        # between selection source and default storage
        for element in selection:
            extra_keys = {}
            if element['_id'] in cache:
                extra_keys = cache[element['_id']]
            element.update(extra_keys)

        return selection

    def find(
        self,
        limit=None,
        skip=None,
        _filter={},
        sort=None,
        with_count=False
    ):

        """
        Retrieve information from data sources

        :param search: a mongo filter like allowing accurate selection
        :param limit: maximum record fetched at once
        :param skip: ordinal number where selection should start
        :param with_count: compute selection count when True
        """

        _filter[Linklist.CONTEXT_CONFIGURABLE] = True

        result = self.context.find(
            _type=Linklist.TYPE,
            _filter=_filter,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count
        )

        if with_count:
            result, count = result
            if result is not None:
                result = self.extra_find(list(result))
            return (result, count)
        else:
            if result is not None:
                result = self.extra_find(list(result))
            return result

    def put(
        self,
        context,
        extra_keys,
    ):
        """
        Persistance layer for upsert operations

        :param context: contains data identifiers
        :param extra_keys: documents extra information depending on specific
        schema.
        """

        entity = {
            Linklist.CONTEXT_CONFIGURABLE: True,
            Context.NAME: Linklist.CONTEXT_CONFIGURABLE
        }

        extra_keys['last_update_date'] = int(time())

        self.context.put(
            _type=Linklist.TYPE,
            entity=entity,
            context=context,
        )

        _id = self.get_id_from_context(context)

        self.extra_put(_id, extra_keys)

    def get_id_from_context(self, context):
        """
        Retrieve the id of the document from the context collection.

        :param context: Context allowing collection identification
        """

        context_element = self.context.find(
            _type=Linklist.TYPE, context=context
        )

        context = list(context_element)

        if len(context) == 0:
            return None
        else:
            return context[0]['_id']

    def remove(
        self,
        context
    ):
        """
        Allow remove operation on records

        :param element_id: identifier for the document to remove
        """

        self.context.remove(
            context=context,
            _type=Linklist.TYPE,
        )

        self.extra_delete([context['_id']])
