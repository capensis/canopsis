# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid

CONF_PATH = 'event/event.conf'
CATEGORY = 'EVENT'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Event(MiddlewareRegistry):

    default_state = 0
    EVENT_STORAGE = 'event_storage'

    states_str = {
        0: 'info',
        1: 'minor',
        2: 'major',
        3: 'critical'
    }

    """
    Manage events in Canopsis
    """

    @staticmethod
    def get_rk(event):
        rk = '{0}.{1}.{2}.{3}.{4}'.format(
            event['connector'],
            event['connector_name'],
            event['event_type'],
            event['source_type'],
            event['component']
        )

        if event['source_type'] == 'resource':
            rk = '{0}.{1}'.format(rk, event['resource'])

        return rk

    def is_ack(self, event):
        """
        Define if an event is in ack state
        :param: event is the event to test
        """
        return event.get('ack', {}).get('isAck', False)

    def is_alert(self, state):
        """
            Define if a state is in alert
            allow progressive alert definition migration
        """
        result = None
        if state == 0:
            result = False
        if state in (1, 2, 3):
            result = True
        return result

    def get_last_state(self, event):
        """
            Retrieve last event state from database
            This is a subset information of a find query focused on state
        """

        existing_event = self.get(Event.get_rk(event), {})
        return existing_event.get('state', self.default_state)

    def get(self, rk, projection=None, default=None):
        result = self.find(query={'rk': rk}, limit=1, projection=projection)

        return result[0] if len(result) else default

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        query={},
        projection=None
    ):

        """
        Retrieve information from data sources

        :param ids: an id list for document to search
        :param limit: maximum record fetched at once
        :param skip: ordinal number where selection should start
        :param with_count: compute selection count when True
        """

        result = self[Event.EVENT_STORAGE].get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=query,
            with_count=with_count,
            projection=projection
        )
        return result
