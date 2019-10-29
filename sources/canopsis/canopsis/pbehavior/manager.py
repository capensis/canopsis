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
import re
from calendar import timegm
from datetime import datetime, date
from dateutil import tz, rrule
from json import loads, dumps
from time import time
from uuid import uuid4
from six import string_types
from pymongo import DESCENDING

import pytz

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection, CollectionError
from canopsis.common.utils import singleton_per_scope
from canopsis.confng import Configuration, Ini
from canopsis.context_graph.manager import ContextGraph
from canopsis.logger import Logger
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
    ID = "_id"
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
    TIMEZONE = 'timezone'
    EXDATE = 'exdate'

    DEFAULT_TYPE = 'generic'

    _FIELDS = (NAME, FILTER, COMMENTS, TSTART, TSTOP, RRULE, ENABLED, EIDS,
               CONNECTOR, CONNECTOR_NAME, AUTHOR, TYPE, REASON, TIMEZONE,
               EXDATE, ID)

    _EDITABLE_FIELDS = (NAME, FILTER, TSTART, TSTOP, RRULE, ENABLED,
                        CONNECTOR, CONNECTOR_NAME, AUTHOR, TYPE, REASON,
                        TIMEZONE, EXDATE)

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

    PB_COLLECTION = 'default_pbehavior'
    LOG_PATH = 'var/log/pbehaviormanager.log'
    LOG_NAME = 'pbehaviormanager'
    CONF_PATH = 'etc/pbehavior/manager.conf'
    PBH_CAT = "PBEHAVIOR"

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
        mongo = MongoStore.get_default()
        collection = mongo.get_collection(cls.PB_COLLECTION)
        mongo_collection = MongoCollection(collection)
        config = Configuration.load(PBehaviorManager.CONF_PATH, Ini)

        return config, logger, mongo_collection

    def __init__(self, config, logger, pb_collection):
        """
        :param dict config: configuration
        :param pb_storage: PBehavior Storage object
        """
        super(PBehaviorManager, self).__init__()
        kwargs = {"logger": logger}
        self.context = singleton_per_scope(ContextGraph, kwargs=kwargs)
        self.logger = logger
        self.config = config
        self.config_data = self.config.get(self.PBH_CAT, {})
        self.default_tz = self.config_data.get("default_timezone",
                                               "Europe/Paris")
        # this line allow us to raise an exception pytz.UnknownTimeZoneError,
        # if the timezone defined in the pbehabior configuration file is wrong
        pytz.timezone(self.default_tz)
        self.collection = pb_collection
        self.currently_active_pb = set()

    def get(self, _id, search=None, limit=None, skip=None):
        """Get pbehavior by id.

        When _id is None, all the pbehaviors are returned. This behavior
        should be considered deprecated, and is only kept for backward
        compatibility. You probably want to use the get_enabled_pbehaviors
        method instead.

        :param str id: pbehavior id
        :param dict query: filtering options
        """
        pipeline = []
        if _id is None:
            if search is not None:
                or_query = [
                    {"name": re.compile(str(search), re.IGNORECASE)},
                    {"reason": re.compile(str(search), re.IGNORECASE)},
                    {"author": re.compile(str(search), re.IGNORECASE)},
                    {"type_": re.compile(str(search), re.IGNORECASE)},
                    {"eids": {"$elemMatch": {
                        "$regex": ".*{}.*".format(str(search)), '$options': 'i'}}}
                ]
                pipeline.append({"$match": {"$or": or_query}})
            else:
                pipeline.append({"$match": {}})
        else:
            pipeline.append({"$match": {"_id": _id}})

        total_count_data = list(self.collection.aggregate(
            pipeline + [{'$count': 'total_count'}]))

        if(len(total_count_data) == 1):
            try:
                total_count = total_count_data[0]["total_count"]
            except (IndexError, KeyError):
                self.logger.error(
                    "Exception while trying to reach total_count")
                return {"total_count": 0, "count": 0, "data": []}
        else:
            self.logger.error(
                "The aggregate returned unexpected data about total_count")
            return {"total_count": 0, "count": 0, "data": []}

        if _id is None:
            if skip is not None:
                pipeline.append({"$skip": skip})
            if limit is not None:
                pipeline.append({"$limit": limit})

        pbhs = list(self.collection.aggregate(pipeline))

        return {"total_count": total_count,
                "count": len(pbhs),
                "data": pbhs}

    def create(
            self,
            name, filter, author,
            tstart, tstop, rrule='',
            enabled=True, comments=None,
            connector='canopsis', connector_name='canopsis',
            type_=PBehavior.DEFAULT_TYPE, reason='', timezone=None,
            exdate=None, pbh_id=None):
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
        :param str timezone: the timezone of the new pbehabior. If no timezone
        are given, use the default one. See the pbehavior documentation
        for more information.
        :param list of str| str exdate: a list of string representation of a date
        following this pattern "YYYY/MM/DD HH:MM:00 TIMEZONE". The hour use the
        24 hours clock system and the timezone is the name of the timezone. The
        month, the day of the month, the hour, the minute and second are
        zero-padded.
        :param str pbh_id: Optional id for pbh. If not specified or none, a
        random id will be generated
        :raises ValueError: invalid RRULE
        :raises pytz.UnknownTimeZoneError: invalid timezone
        :return: created element eid
        :rtype: str
        """

        if timezone is None:
            timezone = self.default_tz

        if exdate is None:
            exdate = []

        # this line allow us to raise an exception pytz.UnknownTimeZoneError,
        # if the timezone defined in the pbehabior configuration file is wrong
        pytz.timezone(timezone)

        if enabled in [True, "True", "true"]:
            enabled = True
        elif enabled in [False, "False", "false"]:
            enabled = False
        else:
            raise ValueError("The enabled value does not match a boolean")

        if not isinstance(exdate, list):
            exdate = [exdate]

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

        if pbh_id is None:
            pbh_id = str(uuid4())

        pb_kwargs = {
            PBehavior.ID: pbh_id,
            PBehavior.NAME: name,
            PBehavior.FILTER: filter,
            PBehavior.AUTHOR: author,
            PBehavior.TSTART: tstart,
            PBehavior.TSTOP: tstop,
            PBehavior.RRULE: rrule,
            PBehavior.ENABLED: enabled,
            PBehavior.COMMENTS: comments,
            PBehavior.CONNECTOR: connector,
            PBehavior.CONNECTOR_NAME: connector_name,
            PBehavior.TYPE: type_,
            PBehavior.REASON: reason,
            PBehavior.TIMEZONE: timezone,
            PBehavior.EXDATE: exdate,
            PBehavior.EIDS: []
        }

        data = PBehavior(**pb_kwargs)
        if not data.comments or not isinstance(data.comments, list):
            data.update(comments=[])
        else:
            for comment in data.comments:
                # Add a unique id to each comment, so that it can be
                # manipulated with the /pbehavior/comment API
                comment['_id'] = str(uuid4())
        try:
            result = self.collection.insert(data.to_dict())
        except CollectionError:
            # when inserting already existing id
            raise ValueError("Trying to insert PBehavior with already existing _id")

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

        cursor = self.collection.find({PBehavior.EIDS: {"$in": id_}})

        pbehaviors = []

        now = int(time())

        for pb in cursor:
            if pb['tstart'] <= now and (pb['tstop'] is None or pb['tstop'] >= now):
                pb['isActive'] = True
            else:
                pb['isActive'] = False

            pbehaviors.append(pb)

        return pbehaviors

    def read(self, _id=None, search=None, limit=None, skip=None):
        """Get pbehavior or list pbehaviors.
        :param str _id: pbehavior id, _id may be equal to None
        """
        result = self.get(_id, search=search, limit=limit, skip=skip)

        return result

    def update(self, _id, **kwargs):
        """
        Update pbehavior record
        :param str _id: pbehavior id
        :param dict kwargs: values pbehavior fields. If a field is None, it will
            **not** be updated.
        :raises ValueError: invalid RRULE or no pbehavior with given _id
        """
        data = self.__get_and_check_pbehavior(_id, **kwargs)
        data["new_data"]["_id"] = _id
        result = self.collection.update(
            {PBehavior.ID: _id},
            {'$set': data["new_data"]})

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return data["pbehavior"].to_dict()
        return None

    def update_v2(self, _id, **kwargs):
        """
        Update pbehavior record
        :param str _id: pbehavior id
        :param dict kwargs: values pbehavior fields. If a field is None, it will
            **not** be updated.
        :raises ValueError: invalid RRULE or no pbehavior with given _id
        """
        pbehavior = self.__get_and_check_pbehavior(_id, **kwargs)["pbehavior"]

        result = self.collection.update(
            {'_id': pbehavior._id or _id}, pbehavior.to_dict(), upsert=False)

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return pbehavior.to_dict()
        return None

    def __get_and_check_pbehavior(self, _id, **kwargs):
        try:
            pb_value = self.get(_id).get('data')[0]
        except (TypeError, KeyError, IndexError):
            raise ValueError("The id does not match any pbehavior")

        check_valid_rrule(kwargs.get('rrule', ''))

        pbehavior = PBehavior(**pb_value)
        new_data = {k: v for k, v in kwargs.items() if v is not None}
        pbehavior.update(**new_data)

        return {"pbehavior": pbehavior, "new_data": new_data}

    def upsert(self, pbehavior):
        """
        Creates or update the given pbehavior.

        This function uses MongoStore/MongoCollection instead of Storage.

        :param canopsis.models.pbehavior.PBehavior pbehavior:
        :rtype: bool, dict
        :returns: success, update result
        """
        r = self.collection.update(
            {'_id': pbehavior._id}, pbehavior.to_dict(), upsert=True)

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

        if _id is None and _filter is None:
            raise ValueError("_id and _filter is None, this will erase every"
                             "pbehaviors.")
        filter_ = {}
        if _filter is not None:
            filter_ = _filter

        if _id is not None:
            filter_["_id"] = _id

        result = self.collection.remove(filter_)

        return self._check_response(result)

    def _update_pbehavior(self, pbehavior_id, query):
        result = self.collection.update(
            {'_id': pbehavior_id}, query, multi=False
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
            {PBehavior.COMMENTS: {'$elemMatch': {'_id': _id}}}
        )
        if not pbehavior:
            return None

        pbehavior = pbehavior.get('data')[0]
        if not pbehavior:
            return None

        _comments = pbehavior[PBehavior.COMMENTS]
        if not _comments:
            return None

        comment = Comment(**_comments[0])
        comment.update(**kwargs)

        result = self.collection.update(
            {'_id': pbehavior_id, 'comments._id': _id},
            {'$set': {'comments.$': comment.to_dict()}},
            multi=False
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
        result = self.collection.update(
            {'_id': pbehavior_id},
            {'$pull': {PBehavior.COMMENTS: {'_id': _id}}},
            multi=False
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
            self.collection.find(
                {PBehavior.EIDS: {'$in': [entity_id]}},
                sort=[(PBehavior.TSTART, DESCENDING)]
            )
        )

        return res

    def compute_pbehaviors_filters(self):
        """
        Compute all filters and update eids attributes.
        """
        pbehaviors = self.collection.find(
            {PBehavior.FILTER: {'$exists': True}})

        for pbehavior in pbehaviors:

            query = loads(pbehavior[PBehavior.FILTER])
            if not isinstance(query, dict):
                self.logger.error('compute_pbehaviors_filters(): filter is '
                                  'not a dict !\n{}'.format(query))
                continue

            entities = self.context.ent_storage.get_elements(
                query=query
            )

            eids = [e['_id'] for e in entities]
            self.collection.update({"_id": pbehavior[PBehavior.ID]},
                                   {"$set": {PBehavior.EIDS: eids}},
                                   upsert=False, multi=False)

    def _check_active_simple_pbehavior(self, timestamp, pbh):
        """ Check if a pbehavior without a rrule is active at the given time.

        :param int timestamp: the number a second this 1970/01/01 00:00:00
        :param dict pbehavior: a pbehavior as a dict.
        :return bool: True if the boolean is active, false otherwise
        """
        if pbh[PBehavior.TSTART] <= timestamp <= pbh[PBehavior.TSTOP]:
            return True

        return False

    @staticmethod
    def __convert_timestamp(timestamp, timezone):
        """Convert a pbehavior timestamp defined in the timezone to a datetime
        in the same timezone.
        :param timestamp:"""

        return datetime.fromtimestamp(timestamp, tz.gettz(timezone))

    def _get_recurring_pbehavior_rruleset(self, pbehavior):
        """ Gets the rec_set for a recurring pbehavior

        :param Dict[str, Any] pbehavior: the recurring pbehavior
        :rtype: rruleset
        """
        tz_name = pbehavior.get(PBehavior.TIMEZONE, self.default_tz)

        rec_set = rrule.rruleset()

        start = self.__convert_timestamp(pbehavior[PBehavior.TSTART], tz_name)

        if PBehavior.EXDATE in pbehavior and\
           isinstance(pbehavior[PBehavior.EXDATE], list):
            for date in pbehavior[PBehavior.EXDATE]:
                exdate = self.__convert_timestamp(date, tz_name)
                rec_set.exdate(exdate)

        rec_set.rrule(rrule.rrulestr(pbehavior[PBehavior.RRULE],
                                     dtstart=start))
        return rec_set

    def _check_active_recurring_pbehavior(self, timestamp, pbehavior):
        """ Check if a pbehavior with a rrule is active at the given time.

        :param int timestamp: the number a second this 1970/01/01 00:00:00
        :param dict pbehavior: a pbehavior as a dict.
        :return bool: True if the boolean is active, false otherwise
        :raise ValueError: if the pbehavior.exdate is invalid. Or if the
        date of an occurence of the pbehavior is not a valid date.
        """

        tz_name = pbehavior.get(PBehavior.TIMEZONE, self.default_tz)

        rec_set = self._get_recurring_pbehavior_rruleset(pbehavior)

        # convert the timestamp to a datetime in the pbehavior's timezone
        now = self.__convert_timestamp(timestamp, tz_name)

        start = self.__convert_timestamp(pbehavior[PBehavior.TSTART], tz_name)
        stop = self.__convert_timestamp(pbehavior[PBehavior.TSTOP], tz_name)
        duration = stop - start  # pbehavior duration

        rec_start = rec_set.before(now)

        self.logger.debug("Recurence start : {}".format(rec_start))
        # No recurrence found
        if rec_start is None:
            return False

        self.logger.debug("Timestamp       : {}".format(now))
        self.logger.debug("Recurence stop  : {}".format(rec_start + duration))

        if rec_start <= now <= rec_start + duration:
            return True

        return False

    def check_active_pbehavior(self, timestamp, pbehavior):
        """ Check if a pbehavior is active at the given time.

        :param int timestamp: the number a second this 1970/01/01 00:00:00
        :param dict pbehavior: a pbehavior as a dict.
        :return bool: True if the boolean is active, false otherwise
        :raise ValueError: if the pbehavior.exdate is invalid. Or if the
        date of an occurence of the pbehavior is not a valid date.
        """
        if PBehavior.RRULE not in pbehavior or\
           pbehavior[PBehavior.RRULE] is None or\
           pbehavior[PBehavior.RRULE] == "":
            return self._check_active_simple_pbehavior(timestamp, pbehavior)
        else:
            if PBehavior.EXDATE not in pbehavior:
                pbehavior[PBehavior.EXDATE] = []
            return self._check_active_recurring_pbehavior(timestamp, pbehavior)

    def check_pbehaviors(self, entity_id, list_in, list_out):
        """
        !!!! DEPRECATED !!!!
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
        self.logger.critical("_check_pbehavior is DEPRECATED !!!!")
        try:
            entity = self.context.get_entities_by_id(entity_id)[0]
        except Exception:
            self.logger.error('Unable to check_behavior on {} entity_id'
                              .format(entity_id))
            return None
        event = self.context.get_event(entity)

        pbehaviors = self.collection.find(
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
                    rrule.rrulestr(pbehavior['rrule'], dtstart=tstart).between(
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

    def get_active_pbehaviors_on_entities(self, entity_ids):
        """
        Yields the pbehaviors that are currently active on a list of entities,
        given their ids.

        :param List[str] entity_ids: the ids of the entities
        :returns Iterator[Dict[str, Any]]: an iterator on the active pbehaviors
        """
        now = int(time())

        pbehaviors = self.collection.find({
            "eids": {"$in": entity_ids}
        })
        for pbehavior in pbehaviors:
            try:
                if self.check_active_pbehavior(now, pbehavior):
                    yield pbehavior
            except ValueError as exept:
                self.logger.exception(
                    "Can't check if the pbehavior is active.")

    def get_all_active_pbehaviors(self):
        """
        Return all pbehaviors currently active using
        self.check_active_pbehavior
        """
        now = int(time())

        ret_val = list(self.collection.find({}))

        results = []

        for pb in ret_val:
            try:
                if self.check_active_pbehavior(now, pb):
                    results.append(pb)
            except ValueError as exept:
                self.logger.exception(
                    "Can't check if the pbehavior is active.")

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

        ret_val = list(self.collection.find(query))

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

        varying_pbs = active_pbehaviors_ids.symmetric_difference(
            self.currently_active_pb)
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
        new_pbs_full = list(self.collection.find({'_id': {'$in': new_pbs}}))

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

    @staticmethod
    def get_active_intervals(after, before, pbehavior):
        """
        Return all the time intervals between after and before during which the
        pbehavior was active.

        The intervals are returned as a list of tuples (start, end), ordered
        chronologically. start and end are UTC timestamps, and are always
        between after and before.

        :param int after: a UTC timestamp
        :param int before: a UTC timestamp
        :param Dict[str, Any] pbehavior:
        :rtype: List[Tuple[int, int]]
        """
        rrule_str = pbehavior[PBehavior.RRULE]
        tstart = pbehavior[PBehavior.TSTART]
        tstop = pbehavior[PBehavior.TSTOP]

        if not isinstance(tstart, (int, float)):
            return
        if not isinstance(tstop, (int, float)):
            return

        # Convert the timestamps to datetimes
        tz = pytz.UTC
        dttstart = datetime.utcfromtimestamp(tstart).replace(tzinfo=tz)
        dttstop = datetime.utcfromtimestamp(tstop).replace(tzinfo=tz)
        delta = dttstop - dttstart

        dtafter = datetime.utcfromtimestamp(after).replace(tzinfo=tz)
        dtbefore = datetime.utcfromtimestamp(before).replace(tzinfo=tz)

        if not rrule_str:
            # The only interval where the pbehavior is active is
            # [dttstart, dttstop]. Ensure that it is included in
            # [after, before], and convert the datetimes to timestamps.
            if dttstart < dtafter:
                dttstart = dtafter
            if dttstop > dtbefore:
                dttstop = dtbefore
            yield (
                timegm(dttstart.timetuple()),
                timegm(dttstop.timetuple())
            )
        else:
            # Get all the intervals that intersect with the [after, before]
            # interval.
            interval_starts = rrule.rrulestr(rrule_str, dtstart=dttstart).between(
                dtafter - delta, dtbefore, inc=False)
            for interval_start in interval_starts:
                interval_end = interval_start + delta
                # Ensure that the interval is included in [after, before], and
                # datetimes to timestamps.
                if interval_start < dtafter:
                    interval_start = dtafter
                if interval_end > dtbefore:
                    interval_end = dtbefore
                yield (
                    timegm(interval_start.timetuple()),
                    timegm(interval_end.timetuple())
                )

    def get_intervals_with_pbehaviors_by_eid(self, after, before, entity_id):
        """
        Yields intervals between after and before with a boolean indicating if
        a pbehavior affects the entity during this interval.

        The intervals are returned as a list of tuples (start, end, pbehavior),
        ordered chronologically. start and end are UTC timestamps, and are
        always between after and before, pbehavior is a boolean indicating if a
        pbehavior affects the entity during this interval. None of the
        intervals overlap.

        :param int after: a UTC timestamp
        :param int before: a UTC timestamp
        :param str entity_id: the id of the entity
        :rtype: Iterator[Tuple[int, int, bool]]
        """
        return self.get_intervals_with_pbehaviors(
            after, before, self.get_pbehaviors(entity_id))

    def get_intervals_with_pbehaviors(self, after, before, pbehaviors):
        """
        Yields intervals between after and before with a boolean indicating if
        one of the pbehaviors is active during this interval.

        The intervals are returned as a list of tuples (start, end, pbehavior),
        ordered chronologically. start and end are UTC timestamps, and are
        always between after and before, pbehavior is a boolean indicating if a
        pbehavior affects the entity during this interval. None of the
        intervals overlap.

        :param int after: a UTC timestamp
        :param int before: a UTC timestamp
        :param List[Dict[str, Any]] pbehaviors: a list of pbehabiors
        :rtype: Iterator[Tuple[int, int, bool]]
        """
        intervals = []

        # Get all the intervals where a pbehavior is active
        for pbehavior in pbehaviors:
            for interval in self.get_active_intervals(after, before, pbehavior):
                intervals.append(interval)

        if not intervals:
            yield (after, before, False)
            return

        # Order them chronologically (by start date)
        intervals.sort(key=lambda a: a[0])

        # Yield the first interval without any active pbehavior
        merged_interval_start, merged_interval_end = intervals[0]
        yield (
            after,
            merged_interval_start,
            False
        )

        # At this point intervals is a list of intervals where a pbehavior is
        # active, ordered by start date. Some of those intervals may be
        # overlapping. This merges the overlapping intervals.
        for interval_start, interval_end in intervals[1:]:
            if interval_end < merged_interval_end:
                # The interval is included in the merged interval, skip it.
                continue

            if interval_start > merged_interval_end:
                # Since the interval starts after the end of the merged
                # interval, they cannot be merged. Yield the merged interval,
                # and move to the new one.
                yield (
                    merged_interval_start,
                    merged_interval_end,
                    True
                )
                yield (
                    merged_interval_end,
                    interval_start,
                    False
                )
                merged_interval_start = interval_start

            merged_interval_end = interval_end

        yield (
            merged_interval_start,
            merged_interval_end,
            True
        )
        yield (
            merged_interval_end,
            before,
            False
        )

    def get_enabled_pbehaviors(self):
        """
        Yields all the enabled pbehaviors.

        :rtype: Iterator[Dict[str, Any]]
        """
        return self.collection.find({PBehavior.ENABLED: True})

    def _get_last_tstop(self, pbh, now):
        """
        Returns last pbehavior stop timestamp before now.

        Warning : this method might return a timestamp greater than the now
                  timestamp, which means the pbehavior is currently running
                  It can also return 0 when the pbh hasn't started running
                  yet


        :param Dict[str, Any] pbh: a pbehavior
        :param datetime now: datetime corresponding to now
        :rtype: int
        """
        tz_name = pbh.get(PBehavior.TIMEZONE, self.default_tz)
        start = self.__convert_timestamp(pbh[PBehavior.TSTART], tz_name)
        if start > now:
            # when pbh hasn't started yet, we return 0 in order to exclude pbh
            return 0
        if PBehavior.RRULE not in pbh or\
                pbh[PBehavior.RRULE] is None or\
                pbh[PBehavior.RRULE] == "":
            #pbh is simple
            pbh_last_tstop = pbh[PBehavior.TSTOP]
        else:
            # convert the timestamp to a datetime in the pbehavior's timezone
            start = self.__convert_timestamp(pbh[PBehavior.TSTART], tz_name)
            stop = self.__convert_timestamp(pbh[PBehavior.TSTOP], tz_name)

            duration = stop - start  # pbehavior duration
            rec_set = self._get_recurring_pbehavior_rruleset(pbh)
            last_tstart = rec_set.before(now)
            # when the pbh is recurrent but hasn't started running yet
            # we return 0, which ensures this pbh isn't used
            if last_tstart is None:
                return 0
            last_tstop_dt = last_tstart + duration
            pbh_last_tstop = int(
                (last_tstop_dt - datetime(1970, 1, 1, tzinfo=tz.UTC)).total_seconds())
        return pbh_last_tstop

    def get_ok_ko_timestamp(self, entity_id):
        """
        Get the timestamp corresponding either to current day at midnight,
        or to the last pbehavior stop for that entity_id

        :param str entity_id: the entity id needing the ok ko timestamp
        :rtype: int
        """
        # get today at midnight timestamp as base return timestamp
        # because each alarm ok ko counter is soft-reseted at midnight
        # midnight at local timezone
        ret_timestamp = self.get_last_tstop_from_eid(entity_id)
        today_at_midnight = date.today()
        tam_timestamp = int(today_at_midnight.strftime("%s"))
        if ret_timestamp < tam_timestamp:
            return tam_timestamp
        return ret_timestamp

    def get_last_tstop_from_eid(self, entity_id):
        """
        Get the timestamp corresponding to
        the last pbehavior stop for that entity_id
        If pbh is active, then now timestamp is returned
        If no pbh is found, then 0 is returned

        :param str entity_id: the entity id needing the last pbh timestamp
        :rtype: int
        """
        now = int(time())
        ret_timestamp = 0
        for pbh in self.get_pbehaviors(entity_id):
            tz_name = pbh.get(PBehavior.TIMEZONE, self.default_tz)
            now_dt = self.__convert_timestamp(now, tz_name)
            if self.check_active_pbehavior(now, pbh):
                # if a pbh is active, then the ok ko counter
                # is supposed to be inactive
                return now

            pbh_last_tstop = self._get_last_tstop(pbh, now_dt)
            if now > pbh_last_tstop > ret_timestamp:
                # keeping the most recent timestamp that still is in the past
                ret_timestamp = pbh_last_tstop
        return ret_timestamp
