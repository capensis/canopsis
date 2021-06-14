# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from urlparse import urlparse

from canopsis.logger import Logger

def parse_scheme(uri):
    """
    uri examples:
    mongodb-default://
    mongodb-periodical-alarm://
    mongodb-default-alarmfilter://
    """

    protocol = None
    data_type = None
    data_scope = 'default'
    parsed_url = urlparse(uri)

    protocol = parsed_url.scheme

    if protocol and '-' in protocol:
        splitted_scheme = protocol.split('-')

        if splitted_scheme:
            protocol = splitted_scheme[0]

        if len(splitted_scheme) > 1:
            data_type = splitted_scheme[1]

        if len(splitted_scheme) > 2:
            data_scope = splitted_scheme[2]

    return protocol, data_type, data_scope

class Middleware(dict):
    """
    Emulator is to make other classes inherit from it instead of the
    regular Middleware class.

    And yes, we derive from dict because… middleware.
    """

    def __init__(self, *args, **kwargs):
        clsname = self.__class__.__name__
        self.logger = Logger.get(clsname, 'var/log/{}.log'.format(clsname))

    @property
    def log_lvl(self):
        return self.logger.level

    @property
    def safe(self):
        return True

    def _connect(self):
        raise NotImplementedError('empty middleware')

    def reconnect(self):
        return self._connect()

    @staticmethod
    def get_middleware_by_uri(uri, table=None):
        """
        table overrides data_scope
        """

        protocol, data_type, data_scope = parse_scheme(uri)

        storage = None

        if protocol == 'mongodb' or protocol == 'storage':
            if data_type == 'periodical':
                from canopsis.mongo.periodical import MongoPeriodicalStorage as msc
            else:
                from canopsis.mongo.core import MongoStorage as msc

            storage = msc()

        else:
            raise Exception('Unknown storage: {}'.format(protocol))

        storage.protocol = protocol
        storage.data_type = data_type if table is None else None
        storage.data_scope = data_scope if table is None else None
        storage.table = table
        storage.logger = Logger.get(
            'storage-{}'.format(protocol),
            'var/log/storage-{}.log'.format(protocol)
        )
        storage._connect()
        storage._backend = storage._get_backend(backend=storage.get_table())

        return storage


class SetSameSiteCookie(object):
    def __init__(self, app, secure):
        self.app = app
        self.secure = secure

    def __call__(self, environ, start_response):
        def session_start_response(status, headers, exc_info=None):
            for header in headers:
                if len(header) == 2 and header[0] == 'Set-cookie':
                    value = header[1].strip()
                    if value.startswith("beaker.session.id="):
                        cookie = header[1]
                        if not value.endswith('SameSite=Lax'):
                            headers.remove(('Set-cookie', cookie))
                            cookie = cookie + '; ' + 'SameSite=Lax'
                            headers.append(('Set-cookie', cookie))
                        if self.secure:
                            headers.remove(('Set-cookie', cookie))
                            cookie = cookie + '; ' + 'Secure'
                            headers.append(('Set-cookie', cookie))

            return start_response(status, headers, exc_info)
        return self.app(environ, session_start_response)
