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

from six import string_types
from calendar import timegm
from datetime import datetime
from dateutil.rrule import rrulestr
from json import loads, dumps
from uuid import uuid4

from canopsis.common.utils import singleton_per_scope
from canopsis.context_graph.manager import ContextGraph
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.middleware.registry import MiddlewareRegistry

# Might be useful when dealing with rrules
# from dateutil.rrule import rrulestr

CONF_PATH = 'pbehavior/pbehavior.conf'
CATEGORY = 'PBEHAVIOR'


class BasePBehavior(dict):
    _FIELDS = ()
    _EDITABLE_FIELDS = ()

    def __init__(self, **kwargs):
        for key, value in kwargs.iteritems():
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
        for key, value in kwargs.iteritems():
            if key in self._EDITABLE_FIELDS:
                self.__dict__[key] = value
        return self.__dict__

    def to_dict(self):
        return self.__dict__


class PBehavior(BasePBehavior):
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

    _FIELDS = (NAME, FILTER, COMMENTS, TSTART, TSTOP, RRULE, ENABLED, EIDS,
               CONNECTOR, CONNECTOR_NAME, AUTHOR)

    _EDITABLE_FIELDS = (NAME, FILTER, TSTART, TSTOP, RRULE, ENABLED,
                        CONNECTOR, CONNECTOR_NAME, AUTHOR)

    def __init__(self, **kwargs):
        if PBehavior.FILTER in kwargs:
            kwargs[PBehavior.FILTER] = dumps(kwargs[PBehavior.FILTER])
        super(PBehavior, self).__init__(**kwargs)

    def update(self, **kwargs):
        if PBehavior.FILTER in kwargs:
            kwargs[PBehavior.FILTER] = dumps(kwargs[PBehavior.FILTER])
        super(PBehavior, self).update(**kwargs)


class Comment(BasePBehavior):
    ID = '_id'
    AUTHOR = 'author'
    TS = 'ts'
    MESSAGE = 'message'

    _FIELDS = (ID, AUTHOR, TS, MESSAGE)
    _EDITABLE_FIELDS = (AUTHOR, MESSAGE)


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PBehaviorManager(MiddlewareRegistry):

    PBEHAVIOR_STORAGE = 'pbehavior_storage'

    _UPDATE_FLAG = 'updatedExisting'
    __TYPE_ERR = "id_ must be a list of string or a string"

    def __init__(self, *args, **kwargs):
        super(PBehaviorManager, self).__init__(*args, **kwargs)
        self.context = singleton_per_scope(ContextGraph)

    @property
    def pbehavior_storage(self):
        return self[PBehaviorManager.PBEHAVIOR_STORAGE]

    def get(self, _id, query=None):
        """Get pbehavior by id.
        :param str id: pbehavior id
        :param dict query: filtering options
        """
        return self.pbehavior_storage.get_elements(ids=_id, query=query)

    def create(
        self,
        name, filter, author,
        tstart, tstop, rrule='',
        enabled=True, comments=None,
        connector='canopsis', connector_name='canopsis'
    ):
        """
        Method creates pbehavior record
        :param str name: filtering options
        :param dict filter: a mongo filter that match entities from canopsis context
        :param str author: the name of the user/app that has generated the pbehavior
        :param timestamp tstart: timestamp that correspond to the start of the pbehavior
        :param timestamp tstop: timestamp that correspond to the end of the pbehavior
        :param str rrule: reccurent rule that is compliant with rrule spec
        :param bool enabled: boolean to know if pbhevior is enabled or disabled
        :param list of dict comments: a list of comments made by users
        :param str connector: a string representing the type of connector that has generated the pbehavior
        :param str connector_name:  a string representing the name of connector that has generated the pbehavior
        :return: created element eid
        :rtype: str
        """

        if enabled in [True, "True", "true"]:
            enabled = True
        elif enabled in [False, "False", "false"]:
            enabled = False
        else:
            raise ValueError("The enabled value does not match a boolean")

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

        kw = locals()
        kw.pop('self')

        data = PBehavior(**kw)
        if not data.comments or not isinstance(data.comments, list):
            data.update(comments=[])
        else:
            for c in data.comments:
                c.update({'_id': str(uuid4())})
        result = self.pbehavior_storage.put_element(element=data.to_dict())

        return result

    def get_pbehaviors_by_eid(self, id_):
        """Retreive from database every pbehavior that contains the given id_
        :param list,str: the id(s) as a str or a list of string
        :return list: a list of pbehavior
        """

        if not isinstance(id_, (list, str, unicode)):
            raise TypeError(self.__TYPE_ERR)

        if isinstance(id_, list):
            for element in id_:
                if not isinstance(element, (str, unicode)):
                    raise TypeError(self.__TYPE_ERR)
        else:
            id_ = [id_]

        cursor = self.pbehavior_storage.get_elements(
            query={PBehavior.EIDS: {"$in": id_}}
        )

        return list(cursor)

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
        :param dict kwargs: values pbehavior fields
        """
        pb_value = self.get(_id)

        if pb_value is None:
            raise ValueError("The id does not match any pebahvior")

        pbehavior = PBehavior(**self.get(_id))
        new_data = {k: v for k, v in kwargs.iteritems() if v is not None}
        pbehavior.update(**new_data)

        result = self.pbehavior_storage.put_element(
            element=new_data, _id=_id
        )

        if (PBehaviorManager._UPDATE_FLAG in result and
                result[PBehaviorManager._UPDATE_FLAG]):
            return pbehavior.to_dict()
        return None

    def delete(self, _id=None, _filter=None):
        """
        Delete pbehavior record
        :param str _id: pbehavior id
        """

        result = self.pbehavior_storage.remove_elements(
            ids=_id, _filter=_filter
        )

        return self._check_response(result)

    def _update_pbehavior(self, pbehavior_id, query):
        result = self.pbehavior_storage._update(
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
        :return:
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
        :return:
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

        result = self.pbehavior_storage._update(
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
        :return:
        """
        result = self.pbehavior_storage._update(
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
        pbehaviors = self.pbehavior_storage.get_elements(
            query={PBehavior.EIDS: {'$in': [entity_id]}},
            sort={PBehavior.TSTART: -1}
        )
        result = [PBehavior(**pb).to_dict() for pb in pbehaviors]

        return result

    def compute_pbehaviors_filters(self):
        """
        Compute all filters and update eids attributes.
        """
        pbehaviors = self.pbehavior_storage.get_elements(
            query={PBehavior.FILTER: {'$exists': True}}
        )

        for pb in pbehaviors:
            entities = self.context[ContextGraph.ENTITIES_STORAGE].get_elements(
                query=loads(pb[PBehavior.FILTER])
            )

            pb[PBehavior.EIDS] = [e['_id'] for e in entities]
            self.pbehavior_storage.put_element(element=pb)

    def check_pbehaviors(self, entity_id, list_in, list_out):
        """

        :param str entity_id:
        :param list list_in: list of pbehavior names
        :param list list_out: list of pbehavior names
        :return: bool if the entity_id is currently in list_in arg and out list_out arg
        """
        return (self._check_pbehavior(entity_id, list_in) and
                not self._check_pbehavior(entity_id, list_out))

    def _check_pbehavior(self, entity_id, pb_names):
        """

        :param str entity_id:
        :param list pb_names: list of pbehavior names
        :return: bool if the entity_id is currently in pb_names arg
        """
        try:
            entity = self.context.get_entities_by_id(entity_id)[0]
        except:
            self.logger.error('Unable to check_behavior on {} entity_id'
                              .format(entity_id))
            return None
        event = self.context.get_event(entity)

        pbehaviors = self.pbehavior_storage.get_elements(
            query={
                PBehavior.NAME: {'$in': pb_names},
                PBehavior.EIDS: {'$in': [entity_id]}
            }
        )

        names = []
        fromts = datetime.fromtimestamp
        for pb in pbehaviors:
            tstart = fromts(pb[PBehavior.TSTART])
            tstop = fromts(pb[PBehavior.TSTOP])

            dt_list = [tstart, tstop]
            if pb['rrule'] is not None:
                dt_list = list(
                    rrulestr(pb['rrule'], dtstart=tstart).between(
                        tstart, tstop, inc=True
                    )
                )

            if (len(dt_list) >= 2 and
                fromts(event['timestamp']) >= dt_list[0] and
                fromts(event['timestamp']) <= dt_list[-1]):
                names.append(pb[PBehavior.NAME])

        result = set(pb_names).isdisjoint(set(names))

        return not result

    def _check_response(self, response):
        ack = True if "ok" in response and response["ok"] == 1 else False
        return {"acknowledged": ack,
                "deletedCount": response["n"]}

    def get_active_pbehaviors(self, eids):
        """
        Return a list of active pbehaviors linked to some entites.

        :param list eids: the desired entities id
        :return: list of pbehaviors
        """
        result = []
        for eid in eids:
            pbhs = self.get_pbehaviors(eid)
            result = result + [x for x in pbhs if self._check_pbehavior(eid, [x['name']])]

        return result
