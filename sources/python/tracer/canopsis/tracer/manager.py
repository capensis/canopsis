from __future__ import unicode_literals

from pymongo.errors import PyMongoError
from bson.errors import BSONError

from canopsis.confng import Configuration, Ini
from canopsis.middleware.registry import MiddlewareRegistry

from canopsis.common.enumerations import FastEnum


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


class TracerManager(MiddlewareRegistry):

    CONF_PATH = 'etc/tracer/manager.conf'
    DEFAULT_STORAGE_URI = 'mongodb-default-tracer://'

    def __init__(self, storage_uri=None, *args, **kwargs):
        super(TracerManager, self).__init__(*args, **kwargs)

        self.config = Configuration.load(self.CONF_PATH, Ini)

        if storage_uri is None:
            storage_uri = self.DEFAULT_STORAGE_URI

        storage_uri = self.config.get('tracer_storage', storage_uri)
        self.storage = self.get_middleware_by_uri(storage_uri)

    def set_trace(self, _id, triggered_by, impact_entities=[], extra={}):
        """
        Creates a new trace or update existing one based on _id.

        :param str _id: trace id
        :param str triggered_by: who triggered this trace. free-form string
        :param list impact_entities: list of entity ids
        :param dict extra: free-form dict, additional and optional informations
        :raises TraceSetError: on put_element error
        """
        trace = {
            Trace.ID: _id,
            Trace.IMPACT_ENTITIES: impact_entities,
            Trace.EXTRA: extra,
            Trace.TRIGGERED_BY: triggered_by
        }

        try:
            res = self.storage.put_element(trace)
        except BSONError, ex:
            raise TraceSetError('document error: {}'.format(ex))
        except PyMongoError, ex:
            raise TraceSetError('pymongo error: {}'.format(ex))

        if res.get('ok', 0) != 1.0:
            raise TraceSetError('put result: {}'.format(res))

        return res

    def add_trace_entities(self, _id, impact_entities):
        """
        Update impact_entities of a trace.

        :param str _id: trace id
        :param list impact_entities: list of entity ids
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
        :raises TraceNotFound: no trace with given _id
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
