# -*- coding: utf-8 -*-

from redis.sentinel import Sentinel
from canopsis.confng import Configuration, Ini


class RedisSentinelStore(object):

    CONF_PATH = 'etc/common/redis_sentinel_store.conf'
    CONF_CAT = 'DATABASE'
    DEFAULT_DB_HOST = 'localhost'
    DEFAULT_DB_PORT = '6379'
    DEFAULT_DB_NUM = '0'
    DEFAULT_MASTER_NAME = 'canopsis'

    @classmethod
    def get_default(cls):
        """
        Get default redis connection using the default configuration file.
        """
        config = Configuration.load(cls.CONF_PATH, Ini)
        return RedisSentinelStore(config)

    def __init__(self, config):
        self.config = config
        conf = self.config.get(self.CONF_CAT, {})
        self.db_hosts = conf.get('hosts', self.DEFAULT_DB_HOST)
        self.db_pass = conf.get('pwd')
        self.master_name = conf.get('master', self.DEFAULT_MASTER_NAME)
        self.sentinel = None
        self.conn = None

        self._connect()

    def parse_hosts(self):
        hosts_ports = self.db_hosts.split(';')
        ret_val = []
        for host_port in hosts_ports:
            split_host = host_port.split(':')
            ret_val.append((split_host[0],int(split_host[1])))
        return ret_val

    def _connect(self):
        self.sentinel = Sentinel(self.parse_hosts(), socket_timeout=0.1, password=self.db_pass)
        self.conn = self.sentinel.master_for(self.master_name)

    def get_new_master(self):
        return self.sentinel.master_for(self.master_name)

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
        return self.conn.delete(name)
