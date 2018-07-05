"""
Trace activity on context graph entities.
"""

from __future__ import unicode_literals

from canopsis.common.collection import CollectionError
from canopsis.common.enumerations import FastEnum
from canopsis.common.middleware import Emulator as Middleware


class TraceError(Exception):
    pass


class TraceNotFound(TraceError):
    pass


class TraceSetError(TraceError):
    pass


class Trace(FastEnum):

    TRIGGERED_BY = 'triggered_by'
    ID = '_id'
    EXTRA = 'extra'
    IMPACT_ENTITIES = 'impact_entities'


class TracerManager(object):

    CONF_PATH = 'etc/tracer/manager.conf'
    DEFAULT_STORAGE_URI = 'mongodb-default-tracer://'

    @classmethod
    def provide_default_storage(cls):
        return Middleware.get_middleware_by_uri(cls.DEFAULT_STORAGE_URI)

    @classmethod
    def provide_default_basics(cls):
        return (cls.provide_default_storage(), )

    def __init__(self, storage):
        """
        :param storage: MongoDB Storage
        """
        super(TracerManager, self).__init__()
        self.storage = storage

    def set_trace(self, _id, triggered_by, impact_entities=None, extra=None):
        """
        Creates a new trace or update existing one based on _id.

        :param str _id: trace id
        :param str triggered_by: who triggered this trace. free-form string
        :param list impact_entities: list of entity ids
        :param dict extra: free-form dict, additional and optional informations
        :raises TraceSetError: on put_element error
        :return dict res: result mongo
        """

        impact_entities = [] if impact_entities is None else impact_entities
        extra = {} if extra is None else extra

        trace = {
            Trace.ID: _id,
            Trace.IMPACT_ENTITIES: impact_entities,
            Trace.EXTRA: extra,
            Trace.TRIGGERED_BY: triggered_by
        }

        res = self.storage.put_element(trace)
        if res is None or res.get('ok', 0) != 1.0:
            raise TraceSetError('put result: {}'.format(res))

        return res

    def add_trace_entities(self, _id, impact_entities):
        """
        Update impact_entities of a trace.

        :param str _id: trace id
        :param list impact_entities: list of entity ids
        :return dict : result mongo
        """
        trace = self.get_by_id(_id)

        entities = set(trace[Trace.IMPACT_ENTITIES])

        len_before_update = len(entities)

        for entity_id in impact_entities:
            entities.add(entity_id)

        if len_before_update != len(entities):
            trace[Trace.IMPACT_ENTITIES] = list(entities)
            return self.storage.put_element(trace)

        return None

    def set_trace_extra(self, _id, extra):
        """
        Update extra informations of a trace. This a simple replace,
        no merge is processed.

        :param str _id: trace id
        :param dict extra: new extra informations
        :return dict : result mongo
        """
        trace = self.get_by_id(_id)

        trace[Trace.EXTRA] = extra

        return self.storage.put_element(trace)

    def del_trace(self, _id):
        return self.storage.remove_elements(filter_={Trace.ID: _id})

    def get_by_id(self, _id):
        """
        :param str _id: trace id
        :rtype dict: trace
        :raise TraceNotFound: no trace with given _id
        """
        query = {Trace.ID: _id}
        res = self.get(_filter=query)

        if len(res) == 0:
            raise TraceNotFound('no trace found with id {}'.format(_id))

        return res[0]

    def get(self, _filter):
        """
        :param dict _filter: mongo filter applied on the trace collection.
        :rtype list: list of documents.
        """
        return list(self.storage.get_elements(query=_filter))
