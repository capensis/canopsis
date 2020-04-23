# -*- coding: utf-8 -*-

import pytz
from uuid import uuid4
from json import loads, dumps

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection, CollectionError
from canopsis.logger import Logger

class BaseMetaAlarmRule(dict):
    """
    Base MetaAlarmRule structure.
    """

    _FIELDS = ()
    _EDITABLE_FIELDS = ()

    def __init__(self, **kwargs):
        super(BaseMetaAlarmRule, self).__init__()
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


META_ALARM_CONFIG_FIELD = 'config'

class MetaAlarmRule(BaseMetaAlarmRule):
    """
    MetaAlarmRule class.
    """
    ID = "_id"
    NAME = 'name'
    TYPE = 'type'
    PATTERNS = 'patterns'
    CONFIG = 'config'

    _FIELDS = (NAME, TYPE, PATTERNS, CONFIG, ID)

    _EDITABLE_FIELDS = (NAME, TYPE, PATTERNS, CONFIG)

    def __init__(self, **kwargs):
        super(MetaAlarmRule, self).__init__(**kwargs)

    def update(self, **kwargs):
        super(MetaAlarmRule, self).update(**kwargs)


class MetaAlarmRuleManager(object):
    """
    MetaAlarmRule manager class
    """

    MA_RULE_COLLECTION = 'meta_alarm_rules'
    LOG_PATH = 'var/log/metaalarmrulemanager.log'
    LOG_NAME = 'metaalarmrulemanager'

    @classmethod
    def provide_default_basics(cls):
        """
        Provide the default configuration and logger objects
        for MetaAlarmRuleManager.

        Do not use those defaults for tests.

        :return: config, logger, storage
        :rtype: Union[dict, logging.Logger, canopsis.storage.core.Storage]
        """
        logger = Logger.get(cls.LOG_NAME, cls.LOG_PATH)
        mongo = MongoStore.get_default()
        collection = mongo.get_collection(cls.MA_RULE_COLLECTION)
        mongo_collection = MongoCollection(collection)
  
        return logger, mongo_collection

    def __init__(self, logger, ma_rule_collection):
        """
        :param dict config: configuration
        :param str ma_rule_collection: MetaAlarmRule collection name
        """
        super(MetaAlarmRuleManager, self).__init__()
 
        self.logger = logger
        self.collection = ma_rule_collection

    def create(self, name, rule_type, patterns=None, config=None, ma_rule_id=None):
        if ma_rule_id is None:
            ma_rule_id = str(uuid4())

        create_kwargs = {
            MetaAlarmRule.ID: ma_rule_id,
            MetaAlarmRule.NAME: name,
            MetaAlarmRule.TYPE: rule_type,
            MetaAlarmRule.PATTERNS: patterns,
            MetaAlarmRule.CONFIG: config,
        }

        data = MetaAlarmRule(**create_kwargs)

        config_dict = data.to_dict().get("config")

        if not config_dict is None and not isinstance(config_dict, dict):
            raise ValueError("config_dict {} from {}".format(config_dict, config))

        try:
            result = self.collection.insert(data.to_dict())
        except CollectionError:
            # when inserting already existing id
            raise ValueError("Trying to insert MetaAlarmRule with already existing _id")

        return result

    def read(self, id):
        return self.collection.find_one({"_id": id})

    def delete(self, id):
        def _check_response(response):
            ack = True if 'ok' in response and response['ok'] == 1 else False

            return {
                'acknowledged': ack,
                'deletedCount': response['n']
            }

        return _check_response(self.collection.remove({"_id": id}))

    def read_rules_with_names(self, ids):
        query = dict()
        if ids:
            query["_id"] = {"$in": ids}

        return dict(((rule[MetaAlarmRule.ID], rule[MetaAlarmRule.NAME]) for rule in self.collection.find(query)))