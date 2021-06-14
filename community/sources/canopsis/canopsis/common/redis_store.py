# -*- coding: utf-8 -*-

import redis

from canopsis.confng import Configuration, Ini


class RedisStore(object):

    CONF_PATH = 'etc/common/redis_store.conf'
    CONF_CAT = 'DATABASE'
    DEFAULT_DB_HOST = 'localhost'
    DEFAULT_DB_PORT = '6379'
    DEFAULT_DB_NUM = '0'

    @classmethod
    def get_default(cls):
        """
        Get default redis connection using the default configuration file.
        """
        config = Configuration.load(cls.CONF_PATH, Ini)
        return RedisStore(config)

    def __init__(self, config):
        self.config = config
        conf = self.config.get(self.CONF_CAT, {})
        self.db_host = conf.get('host', self.DEFAULT_DB_HOST)
        self.db_port = int(conf.get('port', self.DEFAULT_DB_PORT))
        self.db_num = int(conf.get('dbnum', self.DEFAULT_DB_NUM))
        self.db_pass = conf.get('pwd')
        self.conn = None

        self._connect()

    def _connect(self):
        self.conn = redis.StrictRedis(
            host=self.db_host,
            port=self.db_port,
            password=self.db_pass,
        )

    def echo(self, message):
        """
        :returns: the echo-ed essage
        :rtype: string
        """
        return self.conn.echo(message)

    def exists(self, name):
        """
        :returns: True if key name exists, False otherwise
        :rtype: bool
        """
        return self.conn.exists(name)

    def get(self, name):
        """
        :returns: None if key does not exists, value otherwise
        """
        return self.conn.get(name)

    def set(self, name, value, ex=None, px=None, nx=False, xx=False):
        """
        Set the key name to given value. See redis.StrictRedis.set for other
        parameters.
        """
        return self.conn.set(name, value, ex=ex, px=px, nx=nx, xx=xx)

    def remove(self, name):
        """
        remove object to free alert cash
        """
        return self.conn.remove(name)
