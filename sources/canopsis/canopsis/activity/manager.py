# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from .activity import Activity, ActivityAggregate


class ActivityAggregateManager(object):

    AG_COLLECTION = 'default_activity_aggregate'

    def __init__(self, ag_coll, pb_coll, activity_manager):
        """
        :param ag_coll MongoCollection: activity aggregate mongo collection
        :param pb_coll MongoCollection: pbehavior mongo collection
        :type activity_manager: ActivityManager
        """
        self._activity_manager = activity_manager
        self._ag_coll = ag_coll
        self._pb_coll = pb_coll

    def store(self, aggregate):
        """
        Store an aggregate and attached activities.

        All existing activities are dropped before insert.

        :type aggregate: ActivityAggregate
        """
        self._activity_manager.del_by_aggregate_name(aggregate.name)
        self._ag_coll.insert(aggregate.to_dict())
        res = self._activity_manager.store(aggregate.activities)
        self.attach_pbehaviors(aggregate)
        return res

    def get(self, aggregate_name):
        activities = self._activity_manager.get_by_aggregate_name(
            aggregate_name
        )

        aggregate = self._ag_coll.find_one({'_id': aggregate_name})

        if aggregate is None:
            return None

        acag = ActivityAggregate(
            aggregate['_id'],
            aggregate['entity_filter'],
            pb_ids=aggregate['pb_ids']
        )

        for ac in activities:
            acag.add(ac)

        return acag

    def delete(self, aggregate):
        """
        Remove all activities linked to this aggregate and pbehaviors too.
        :param aggregate ActivityAggregate:
        """
        self._pb_coll.remove({'_id': {'$in': aggregate.pb_ids}})
        self._ag_coll.remove({'_id': aggregate.name})

    def attach_pbehaviors(self, aggregate, pb_ids=None):
        """
        Replace all pb_ids in database with those from
        the aggregate plus those in pb_ids parameter.
        :param pb_ids list: list of pb_ids to attach. They extend those already
            in the aggregate object.
        """
        if pb_ids is not None:
            aggregate.pb_ids.extend(pb_ids)

        aggregate.pb_ids = list(set(aggregate.pb_ids))

        self._ag_coll.update(
            {'_id': aggregate.name},
            {
                '$set': {
                    'pb_ids': aggregate.pb_ids
                }
            }
        )


class ActivityManager(object):
    """
    Store/get activities in/from database. Aggregates are never stored,
    they are only used to add a field in the activity so you can query
    activities grouped by aggregate.

    Warning: prefer the use of ActivityAggregateManager to guarantee
    consistency across activities.
    """

    ACTIVITY_COLLECTION = 'default_activity'

    def __init__(self, activity_collection):
        """
        :param activity_collection: MongoCollection
        :type activity_collection: canopsis.common.collection.MongoCollection
        """
        self._coll = activity_collection

    def store(self, activities):
        """
        :type activities: list[Activity]
        """
        activities = [ac.to_dict() for ac in activities]

        return self._coll.insert(activities)

    def del_by_aggregate_name(self, aggregate_name):
        """
        :type aggregate_name: str
        """
        return self._coll.remove({'aggregate_name': aggregate_name})

    def get(self, _id):
        """
        :rtype: Activity
        """
        act = self._coll.find_one({'_id': _id})
        act[Activity.DBID] = act.pop('_id')

        return Activity(**act)

    def get_all(self):
        """
        :rtype: [Activity]
        """
        cursor = self._coll.find({})
        activities = []
        for act in cursor:
            act[Activity.DBID] = act.pop('_id')
            activities.append(Activity(**act))

        return activities

    def get_by_aggregate_name(self, aggregate_name):
        """
        :param str aggregate_name:
        :rtype: list[Activity]
        """
        activities = []
        res = self._coll.find({'aggregate_name': aggregate_name})

        for act in list(res):
            act[Activity.DBID] = act.pop('_id')
            activities.append(Activity(**act))

        return activities
