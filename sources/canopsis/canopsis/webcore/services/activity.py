from __future__ import unicode_literals

from canopsis.confng import Configuration, Ini
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.webcore.utils import gen_json
from canopsis.activity import ActivityManager


def exports(ws):

    conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)

    mongo_store = MongoStore(config=conf_store, cred_config=conf_store)
    ac_coll = MongoCollection(
        mongo_store.get_collection(ActivityManager.ACTIVITY_COLLECTION))
    ac_man = ActivityManager(ac_coll)

    @ws.application.get('/api/v2/activity/activities')
    def get_activities():
        activities = ac_man.get_all()
        jactivities = []

        for act in activities:
            jactivities.append(act.to_dict())

        return gen_json(jactivities)
