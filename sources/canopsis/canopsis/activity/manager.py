# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from .activity import Activity


class ActivityAggregateManager(object):

    ACAGG_COLLECTION = 'default_activityaggregate'

    def __init__(self, acag_collection, activity_manager):
        """
        :type activity_manager: ActivityManager
        :type acag_collection: canopsis.common.collection.MongoCollection
        """
        self._coll = acag_collection
        self._activity_manager = activity_manager

    def store(self, aggregate):
        """
        Store an aggregate and attached activities.

        :type aggregate: ActivityAggregate
        """
        if self._coll.insert(aggregate.to_dict()) == aggregate.name:
            self._activity_manager.store(aggregate.activities)


class ActivityManager(object):
    """
    Store/get activities in/from database. Aggregates are never stored,
    they are only used to add a field in the activity so you can query
    activities grouped by aggregate.
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
        act.pop('_id')

        return Activity(**act)

    def get_all(self):
        """
        :rtype: [Activity]
        """
        cursor = self._coll.find({})
        activities = []
        for act in cursor:
            activities.append(Activity(**act))

        return activities

    def get_by_aggregate_name(self, aggregate_name):
        """
        :param str aggregate_name:
        :rtype: list[Activity]
        """
        activities = []
        res = self._coll.find({'aggregate_name': aggregate_name})

        for r in list(res):
            r.pop('_id')
            activities.append(Activity(**r))

        return activities
