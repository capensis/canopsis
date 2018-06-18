# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

"""
Managing PBehavior.
"""

from calendar import timegm
from datetime import datetime
from json import loads, dumps
from time import time
from uuid import uuid4
from six import string_types
from dateutil.rrule import rrulestr
from pymongo import DESCENDING

import pytz

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.common.utils import singleton_per_scope
from canopsis.context_graph.manager import ContextGraph
from canopsis.logger import Logger
from canopsis.middleware.core import Middleware
from canopsis.pbehavior.utils import check_valid_rrule


class BasePBehavior(dict):
    """
    Base PBehaviorManager structure.
    """

    _FIELDS = ()
    _EDITABLE_FIELDS = ()

    def __init__(self, **kwargs):
        super(BasePBehavior, self).__init__()
        for key, value in kwargs.items():
            if key in self._FIELDS:
                self.__dict__[key] = value

    def __repr__(self):
        return repr(self.__dict__)

    def __setitem__(self, key, item):
        if key in self._EDITABLE_FIELDS:
            self.__dict__[key] = item

    def __getitem__(self, key):
        return self._get(key)

    def __getattr__(self, item):
        return self._get(item)

    def _get(self, item):
        if item in self._FIELDS and item in self.__dict__:
            return self.__dict__[item]
        return None

    def update(self, **kwargs):
        """
        Update the current instance with every kwargs arguments.

        :param kwargs: the argument to use to update the instance
        :returns: the updated representation of the current instance
        :rtype: dict
        """
        for key, value in kwargs.items():
            if key in self._EDITABLE_FIELDS:
                self.__dict__[key] = value
        return self.__dict__

    def to_dict(self):
        """
        Return the dict representation of the current instance

        :returns: return the dict representation of the current instance
        :rtype: dict
        """
        return self.__dict__


class PBehavior(BasePBehavior):
    """
    PBehavior class.
    """

    NAME = 'name'
    FILTER = 'filter'
    COMMENTS = 'comments'
    TSTART = 'tstart'
    TSTOP = 'tstop'
    RRULE = 'rrule'
    ENABLED = 'enabled'
    EIDS = 'eids'
    CONNECTOR = 'connector'
    CONNECTOR_NAME = 'connector_name'
    AUTHOR = 'author'
    TYPE = 'type_'
    REASON = 'reason'

    DEFAULT_TYPE = 'generic'

    _FIELDS = (NAME, FILTER, COMMENTS, TSTART, TSTOP, RRULE, ENABLED, EIDS,
               CONNECTOR, CONNECTOR_NAME, AUTHOR, TYPE, REASON)

    _EDITABLE_FIELDS = (NAME, FILTER, TSTART, TSTOP, RRULE, ENABLED,
                        CONNECTOR, CONNECTOR_NAME, AUTHOR, TYPE, REASON)

    def __init__(self, **kwargs):
        if PBehavior.FILTER in kwargs:
            kwargs[PBehavior.FILTER] = dumps(kwargs[PBehavior.FILTER])
        super(PBehavior, self).__init__(**kwargs)

    def update(self, **kwargs):
        if PBehavior.FILTER in kwargs:
            kwargs[PBehavior.FILTER] = dumps(kwargs[PBehavior.FILTER])
        super(PBehavior, self).update(**kwargs)


class Comment(BasePBehavior):
    """
    Comment class.
    """

    ID = '_id'
    AUTHOR = 'author'
    TS = 'ts'
    MESSAGE = 'message'

    _FIELDS = (ID, AUTHOR, TS, MESSAGE)
    _EDITABLE_FIELDS = (AUTHOR, MESSAGE)


class PBehaviorManager(object):
    """
    PBehavior manager class.
    """

    PB_STORAGE_URI = 'mongodb-default-pbehavior://'
    LOG_PATH = 'var/log/pbehaviormanager.log'
    LOG_NAME = 'pbehaviormanager'

    _UPDATE_FLAG = 'updatedExisting'
    __TYPE_ERR = "id_ must be a list of string or a string"

    @classmethod
    def provide_default_basics(cls):
        """
        Provide the default configuration and logger objects
        for PBehaviorManager.

        Do not use those defaults for tests.

        :return: config, logger, storage
        :rtype: Union[dict, logging.Logger, canopsis.storage.core.Storage]
        """
        logger = Logger.get(cls.LOG_NAME, cls.LOG_PATH)
        pb_storage = Middleware.get_middleware_by_uri(cls.PB_STORAGE_URI)

        return logger, pb_storage

    def __init__(self, logger, pb_storage):
        """
        :param dict config: configuration
        :param pb_storage: PBehavior Storage object
        """
        super(PBehaviorManager, self).__init__()
        kwargs = {"logger": logger}
        self.context = singleton_per_scope(ContextGraph, kwargs=kwargs)
        self.logger = logger
        self.pb_storage = pb_storage

        self.pb_store = MongoCollection(MongoStore.get_default().get_collection('default_pbehavior'))

        self.currently_active_pb = set()

    def get(self, _id, query=None):
        """Get pbehavior by id.

        :param str id: pbehavior id
        :param dict query: filtering options
        """
        return self.pb_storage.get_elements(ids=_id, query=query)

    def create(
            self,
            name, filter, author,
            tstart, tstop, rrule='',
            enabled=True, comments=None,
            connector='canopsis', connector_name='canopsis',
            type_=PBehavior.DEFAULT_TYPE, reason=''):
        """
        Method creates pbehavior record

        :param str name: filtering options
        :param dict filter: a mongo filter that match entities from canopsis
        context
        :param str author: the name of the user/app that has generated the
        pbehavior
        :param timestamp tstart: timestamp that correspond to the start of the
        pbehavior
        :param timestamp tstop: timestamp that correspond to the end of the
        pbehavior
        :param str rrule: reccurent rule that is compliant with rrule spec
        :param bool enabled: boolean to know if pbhevior is enabled or disabled
        :param list of dict comments: a list of comments made by users
        :param str connector: a string representing the type of connector that
            has generated the pbehavior
        :param str connector_name:  a string representing the name of connector
            that has generated the pbehavior
        :param str type_: associated type_ for this pbh
        :param str reason: associated reason for this pbh
        :raises ValueError: invalid RRULE
        :return: created element eid
        :rtype: str
        """

        if enabled in [True, "True", "true"]:
            enabled = True
        elif enabled in [False, "False", "false"]:
            enabled = False
        else:
            raise ValueError("The enabled value does not match a boolean")

        check_valid_rrule(rrule)

        if comments is not None:
            for comment in comments:
                if "author" in comment:
                    if not isinstance(comment["author"], string_types):
                        raise ValueError("The author field must be an string")
                else:
                    raise ValueError("The author field is missing")
                if "message" in comment:
                    if not isinstance(comment["message"], string_types):
                        raise ValueError("The message field must be an string")
                else:
                    raise ValueError("The message field is missing")

        pb_kwargs = {
            'name': name,
            'filter': filter,
            'author': author,
            'tstart': tstart,
            'tstop': tstop,
            'rrule': rrule,
            'enabled': enabled,
            'comments': comments,
            'connector': connector,
            'connector_name': connector_name,
            PBehavior.TYPE: type_,
            'reason': reason
        }
        if PBehavior.EIDS not in pb_kwargs:
            pb_kwargs[PBehavior.EIDS] = []

        data = PBehavior(**pb_kwargs)
        if not data.comments or not isinstance(data.comments, list):
            data.update(comments=[])
        else:
            for comment in data.comments:
                comment.update({'_id': str(uuid4())})
        result = self.pb_storage.put_element(element=data.to_dict())

        return result

    def get_pbehaviors_by_eid(self, id_):
        """Retreive from database every pbehavior that contains
        the given id_ in the PBehavior.EIDS field.

        :param list,str: the id(s) as a str or a list of string
        :returns: a list of pbehavior, with the isActive key in pbehavior is
            active when queried.
        :rtype: list
        """

        if not isinstance(id_, (list, string_types)):
            raise TypeError(self.__TYPE_ERR)

        if isinstance(id_, list):
            for element in id_:
                if not isinstance(element, string_types):
                    raise TypeError(self.__TYPE_ERR)
        else:
            id_ = [id_]

        cursor = self.pb_storage.get_elements(
            query={PBehavior.EIDS: {"$in": id_}}
        )

        pbehaviors = []

        now = int(time())

        for pb in cursor:
            if pb['tstart'] <= now and (pb['tstop'] is None or pb['tstop'] >= now):
                pb['isActive'] = True
            else:
                pb['isActive'] = False

            pbehaviors.append(pb)

        return pbehaviors

    def read(self, _id=None):
        """Get pbehavior or list pbehaviors.
        :param str _id: pbehavior id, _id may be equal to None
        """
        result = self.get(_id)

        return result if _id else list(result)

    def update(self, _id, **kwargs):
        """
        Update pbehavior record
        :param str _id: pbehavior id
        :param dict kwargs: values pbehavior fields. If a field is None, it will
            **not** be updated.
        :raises ValueError: invalid RRULE or no pbehavior with given _id
        """
        pb_value = self.get(_id)

        if pb_value is None:
            raise ValueError("The id does not match any pebahvior")

        check_valid_rrule(kwargs.get('rrule', ''))

        pbehavior = PBehavior(**self.get(_id))
        new_data = {k: v for k, v in kwargs.items() if v is not None}
        pbehavior.update(**new_data)

        result = self.pb_storage.put_element(
            element=new_data, _id=_id
        )

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return pbehavior.to_dict()
        return None

    def upsert(self, pbehavior):
        """
        Creates or update the given pbehavior.

        This function uses MongoStore/MongoCollection instead of Storage.

        :param canopsis.models.pbehavior.PBehavior pbehavior:
        :rtype: bool, dict
        :returns: success, update result
        """
        r = self.pb_store.update({'_id': pbehavior._id}, pbehavior.to_dict(), upsert=True)

        if r.get('updatedExisting', False) and r.get('nModified') == 1:
            return True, r
        elif r.get('updatedExisting', None) is False and r.get('nModified') == 0 and r.get('ok') == 1.0:
            return True, r
        else:
            return False, r

    def delete(self, _id=None, _filter=None):
        """
        Delete pbehavior record
        :param str _id: pbehavior id
        """

        result = self.pb_storage.remove_elements(
            ids=_id, _filter=_filter
        )

        return self._check_response(result)

    def _update_pbehavior(self, pbehavior_id, query):
        result = self.pb_storage._update(
            spec={'_id': pbehavior_id},
            document=query,
            multi=False, cache=False
        )
        return result

    def create_pbehavior_comment(self, pbehavior_id, author, message):
        """
        Сreate comment for pbehavior.

        :param str pbehavior_id: pbehavior id
        :param str author: author of the comment
        :param str message: text of the comment
        """
        comment_id = str(uuid4())
        comment = {
            Comment.ID: comment_id,
            Comment.AUTHOR: author,
            Comment.TS: timegm(datetime.utcnow().timetuple()),
            Comment.MESSAGE: message
        }

        query = {'$addToSet': {PBehavior.COMMENTS: comment}}

        result = self._update_pbehavior(pbehavior_id, query)

        if not result:
            result = self._update_pbehavior(
                pbehavior_id, {'$set': {PBehavior.COMMENTS: []}}
            )
            if not result:
                return None

            result = self._update_pbehavior(pbehavior_id, query)

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return comment_id
        return None

    def update_pbehavior_comment(self, pbehavior_id, _id, **kwargs):
        """
        Update the comment record.

        :param str pbehavior_id: pbehavior id
        :param str_id: comment id
        :param dict kwargs: values comment fields
        """
        pbehavior = self.get(
            pbehavior_id,
            query={PBehavior.COMMENTS: {'$elemMatch': {'_id': _id}}}
        )
        if not pbehavior:
            return None

        _comments = pbehavior[PBehavior.COMMENTS]
        if not _comments:
            return None

        comment = Comment(**_comments[0])
        comment.update(**kwargs)

        result = self.pb_storage._update(
            spec={'_id': pbehavior_id, 'comments._id': _id},
            document={'$set': {'comments.$': comment.to_dict()}},
            multi=False, cache=False
        )

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return comment.to_dict()
        return None

    def delete_pbehavior_comment(self, pbehavior_id, _id):
        """
        Delete comment record.

        :param str pbehavior_id: pbehavior id
        :param str _id: comment id
        """
        result = self.pb_storage._update(
            spec={'_id': pbehavior_id},
            document={'$pull': {PBehavior.COMMENTS: {'_id': _id}}},
            multi=False, cache=False
        )

        return self._check_response(result)

    def get_pbehaviors(self, entity_id):
        """
        Return all pbehaviors related to an entity_id, sorted by descending
        tstart.

        :param str entity_id: Id for which behaviors have to be returned

        :return: pbehaviors, with name, tstart, tstop, rrule and enabled keys
        :rtype: list of dict
        """
        res = list(
            self.pb_storage._backend.find(
                {PBehavior.EIDS: {'$in': [entity_id]}},
                sort=[(PBehavior.TSTART, DESCENDING)]
            )
        )

        return res

    def compute_pbehaviors_filters(self):
        """
        Compute all filters and update eids attributes.
        """
        pbehaviors = self.pb_storage.get_elements(
            query={PBehavior.FILTER: {'$exists': True}}
        )

        for pbehavior in pbehaviors:

            query = loads(pbehavior[PBehavior.FILTER])
            if not isinstance(query, dict):
                self.logger.error('compute_pbehaviors_filters(): filter is '
                                  'not a dict !\n{}'.format(query))
                continue

            entities = self.context.ent_storage.get_elements(
                query=query
            )

            pbehavior[PBehavior.EIDS] = [e['_id'] for e in entities]
            self.pb_storage.put_element(element=pbehavior)

    @staticmethod
    def check_active_pbehavior(timestamp, pbehavior):
        """
        For a given pbehavior (as dict) check that the given timestamp is active
        using either:

        the rrule, if any, from the pbehavior + tstart and tstop to define
        start and stop times.

        tstart and tstop alone if no rrule is available.

        All dates and times are computed using UTC timezone, so check that your
        timestamp was exported using UTC.

        :param int timestamp: timestamp to check
        :param dict pbehavior: the pbehavior
        :rtype bool:
        :returns: pbehavior is currently active or not
        """
        fromts = datetime.utcfromtimestamp
        tstart = pbehavior[PBehavior.TSTART]
        tstop = pbehavior[PBehavior.TSTOP]

        if not isinstance(tstart, (int, float)):
            return False
        if not isinstance(tstop, (int, float)):
            return False

        tz = pytz.UTC
        dtts = fromts(timestamp).replace(tzinfo=tz)
        dttstart = fromts(tstart).replace(tzinfo=tz)
        dttstop = fromts(tstop).replace(tzinfo=tz)

        dt_list = [dttstart, dttstop]
        rrule = pbehavior['rrule']
        if rrule:
            # compute the minimal date from which to start generating
            # dates from the rrule.
            # a complementary date (missing_date) is computed and added
            # at index 0 of the generated dt_list to ensure we manage
            # dates at boundaries.
            dt_tstart_date = dtts.date()
            dt_tstart_time = dttstart.time().replace(tzinfo=tz)
            dt_dtstart = datetime.combine(dt_tstart_date, dt_tstart_time)

            # dt_list at 0 and 1 indexes can be equal, so we generate
            # three dates to ensure [1] - [2] will give a non-zero timedelta
            # object.
            dt_list = list(
                rrulestr(rrule, dtstart=dt_dtstart).xafter(
                    dttstart, count=3, inc=True
                )
            )

            # compute the "missing date": a date before the first one
            # just in case that first date isn't taking in account
            # a pbehavior that just begun.
            missing_date = dt_list[0] - (dt_list[-1] - dt_list[-2])
            dt_list.insert(0, missing_date)

            delta = dttstop - dttstart
            for dt in dt_list:
                if dtts >= dt and dtts <= dt + delta:
                    return True

            return False

        else:
            if dtts >= dttstart and dtts <= dttstop:
                return True
            return False

        return False

    def check_pbehaviors(self, entity_id, list_in, list_out):
        """

        :param str entity_id:
        :param list list_in: list of pbehavior names
        :param list list_out: list of pbehavior names
        :returns: bool if the entity_id is currently in list_in arg and out list_out arg
        """
        return (self._check_pbehavior(entity_id, list_in) and
                not self._check_pbehavior(entity_id, list_out))

    def _check_pbehavior(self, entity_id, pb_names):
        """

        :param str entity_id:
        :param list pb_names: list of pbehavior names
        :returns: bool if the entity_id is currently in pb_names arg
        """
        try:
            entity = self.context.get_entities_by_id(entity_id)[0]
        except Exception:
            self.logger.error('Unable to check_behavior on {} entity_id'
                              .format(entity_id))
            return None
        event = self.context.get_event(entity)

        pbehaviors = self.pb_storage.get_elements(
            query={
                PBehavior.NAME: {'$in': pb_names},
                PBehavior.EIDS: {'$in': [entity_id]}
            }
        )

        names = []
        fromts = datetime.fromtimestamp
        for pbehavior in pbehaviors:
            tstart = pbehavior[PBehavior.TSTART]
            tstop = pbehavior[PBehavior.TSTOP]
            if not isinstance(tstart, (int, float)):
                self.logger.error('Cannot parse tstart value: {}'
                                  .format(pbehavior))
                continue
            if not isinstance(tstop, (int, float)):
                self.logger.error('Cannot parse tstop value: {}'
                                  .format(pbehavior))
                continue
            tstart = fromts(tstart)
            tstop = fromts(tstop)

            dt_list = [tstart, tstop]
            if pbehavior['rrule'] is not None:
                dt_list = list(
                    rrulestr(pbehavior['rrule'], dtstart=tstart).between(
                        tstart, tstop, inc=True
                    )
                )

            if (len(dt_list) >= 2
                    and fromts(event['timestamp']) >= dt_list[0]
                    and fromts(event['timestamp']) <= dt_list[-1]):
                names.append(pbehavior[PBehavior.NAME])

        result = set(pb_names).isdisjoint(set(names))

        return not result

    @staticmethod
    def _check_response(response):
        ack = True if 'ok' in response and response['ok'] == 1 else False

        return {
            'acknowledged': ack,
            'deletedCount': response['n']
        }

    def get_active_pbehaviors(self, eids):
        """
        Return a list of active pbehaviors linked to some entites.

        :param list eids: the desired entities id
        :returns: list of pbehaviors
        """
        result = []
        for eid in eids:
            pbhs = self.get_pbehaviors(eid)
            result = result + [x for x in pbhs if self._check_pbehavior(
                eid, [x['name']]
            )]

        return result

    def get_all_active_pbehaviors(self):
        """
        Return all pbehaviors currently active using
        self.check_active_pbehavior
        """
        now = int(time())
        query = {}

        ret_val = list(self.pb_storage.get_elements(query=query))

        results = []

        for pb in ret_val:
            if self.check_active_pbehavior(now, pb):
                results.append(pb)

        return results

    def get_active_pbehaviors_from_type(self, types=None):
        """
        Return pbehaviors currently active, with a specific type,
        using self.check_active_pbehavior
        """
        if types is None:
            types = []
        now = int(time())
        query = {PBehavior.TYPE: {'$in': types}}

        ret_val = list(self.pb_storage.get_elements(query=query))

        results = []

        for pb in ret_val:
            if self.check_active_pbehavior(now, pb):
                results.append(pb)

        return results

    def get_varying_pbehavior_list(self):
        """
        get_varying_pbehavior_list

        :returns: list of PBehavior id activated since last check
        :rtype: list
        """
        active_pbehaviors = self.get_all_active_pbehaviors()
        active_pbehaviors_ids = set()
        for active_pb in active_pbehaviors:
            active_pbehaviors_ids.add(active_pb['_id'])

        varying_pbs = active_pbehaviors_ids.symmetric_difference(self.currently_active_pb)
        self.currently_active_pb = active_pbehaviors_ids

        return list(varying_pbs)

    def launch_update_watcher(self, watcher_manager):
        """
        launch_update_watcher update watcher when a pbehavior is active

        :param object watcher_manager: watcher manager
        :returns: number of watcher updated
        retype: int
        """
        new_pbs = self.get_varying_pbehavior_list()
        new_pbs_full = list(self.pb_storage._backend.find(
            {'_id': {'$in': new_pbs}}
        ))

        merged_eids = []
        for pbehaviour in new_pbs_full:
            merged_eids = merged_eids + pbehaviour['eids']

        watchers_ids = set()
        for watcher in self.get_wacher_on_entities(merged_eids):
            watchers_ids.add(watcher['_id'])
        for watcher_id in watchers_ids:
            watcher_manager.compute_state(watcher_id)

        return len(list(watchers_ids))

    def get_wacher_on_entities(self, entities_ids):
        """
        get_wacher_on_entities.

        :param entities_ids: entity id
        :returns: list of watchers
        :rtype: list
        """
        query = {
            '$and': [
                {'depends': {'$in': entities_ids}},
                {'type': 'watcher'}
            ]
        }
        watchers = self.context.get_entities(query=query)

        return watchers
