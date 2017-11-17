from __future__ import unicode_literals

import arrow

from bottle import request
from six import string_types
from time import time as now_ts

from canopsis.confng import Configuration, Ini
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.webcore.utils import gen_json, gen_json_error
from canopsis.activity.activity import Activity, ActivityAggregate
from canopsis.activity.manager import ActivityManager, ActivityAggregateManager
from canopsis.activity.pbehavior import PBehaviorGenerator


class RouteHandler(object):

    def __init__(self, ac_man, acag_man, pb_gen):
        self.ac_man = ac_man
        self.acag_man = acag_man
        self.pb_gen = pb_gen

    def get_activities(self):
        activities = self.ac_man.get_all()

        return map(Activity.to_dict, activities)

    def set_activities(self, document):
        """
        :raises ValueError: some activity is invalid
        """
        if not isinstance(document, dict):
            raise ValueError('document must be a dict')

        if 'aggregate_name' not in document:
            raise ValueError('missing aggregate_name parameter')

        if 'activities' not in document:
            raise ValueError('missing activities parameter')

        aggregate_name = document['aggregate_name']
        activities = document['activities']

        if not isinstance(activities, list):
            raise ValueError('activities is not an array')
        aggregate = ActivityAggregate(aggregate_name)
        for doc in activities:
            aggregate.add(Activity(**doc))

        ids = self.acag_man.store(aggregate)

        result = {
            'inserted': len(ids)
        }

        return result

    def _generate_pbs_register(self, pbehaviors):
        pass

    def _generate_pbs_return(self, aggregate_names):
        dict_pbs = {}

        for agname in aggregate_names:
            if not isinstance(agname, string_types):
                raise ValueError(
                    'aggregate name must be a string, got {}'.format(agname))

            acag = self.acag_man.get(agname)
            now = arrow.get(int(now_ts()))
            pbehaviors = self.pb_gen.activities_to_pbehaviors(acag, now)
            dict_pbs[agname] = [pb.to_dict() for pb in pbehaviors]

        return dict_pbs

    def generate_pbs(self, document):
        """
        :param document list[string]: list of aggregate names
        :param register_pb bool: register pbehaviors
        """

        if not isinstance(document, dict):
            raise ValueError('document must be a dict')

        register_pb = request.json.get('register', False)
        aggregate_names = request.json.get('aggregate_names')

        if not isinstance(aggregate_names, list):
            raise ValueError(
                'aggregate_names must be a list of string'
                ': [aggregate_name, ...]')

        if not isinstance(register_pb, bool):
            raise ValueError(
                'register_pb must be a bool, got {}'.format(register_pb))

        dict_pbs = self._generate_pbs_return(aggregate_names)

        if register_pb:
            self._generate_pbs_register(dict_pbs)

        return dict_pbs


def exports(ws):

    conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)

    mdbstore = MongoStore(config=conf_store, cred_config=conf_store)
    ac_coll = MongoCollection(
        mdbstore.get_collection(ActivityManager.ACTIVITY_COLLECTION))

    ac_man = ActivityManager(ac_coll)
    acag_man = ActivityAggregateManager(ac_man)
    pb_gen = PBehaviorGenerator()

    route_handler = RouteHandler(ac_man, acag_man, pb_gen)

    @ws.application.get('/api/v2/activity/activities')
    def get_activities():
        return gen_json(route_handler.get_activities())

    @ws.application.post('/api/v2/activity/set_aggregate')
    def set_activities():
        """
        The JSON document must contain parameters of the Activity object.

        You can pass an Array of Dicts and Activities will be inserted.
        """
        try:
            return gen_json(route_handler.set_activities(request.json))

        except (ValueError, TypeError) as exc:
            return gen_json_error(str(exc), 400)

    @ws.application.post('/api/v2/activity/generate_pbehaviors')
    def generate_pbehaviors():
        try:
            return gen_json(route_handler.generate_pbs(request.json))

        except (ValueError,) as exc:
            return gen_json_error(str(exc), 400)
