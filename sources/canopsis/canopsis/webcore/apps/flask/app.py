# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import importlib
import logging
import os

from flask import Flask
from flask_restful import Api

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.webcore.apps.flask.helpers import Resource

app = Flask(__name__)


def _auto_import(app, api, configuration, exports_funcname='exports_v3'):
    for webservice, enabled in configuration.items():
        if int(enabled) == 1:
            try:
                wsmod = importlib.import_module('{}'.format(webservice))
            except ImportError:
                app.logger.error('webservice v3 {}: unable to import module'
                                 .format(webservice))

            if hasattr(wsmod, exports_funcname):
                app.logger.info('webservice v3 {}: loading'.format(webservice))
                getattr(wsmod, exports_funcname)(app, api)
                app.logger.info('webservice v3 {}: loaded'.format(webservice))
            else:
                app.logger.debug('webservice v3 {}: {} unavailable'
                                 .format(webservice, exports_funcname))
        else:
            app.logger.debug('webservice v3 {}: skipped'.format(webservice))


def _init(app):
    """
    For each configured webservice, run exports_v3 if function exists.

    Expected configuration:

    [webservices]
    wsname=0|1
    other_wsname=0|1

    0: skip webservice
    1: load webservice
    """
    logfile_handler = logging.FileHandler(os.path.join(root_path, 'var/log/oldapi.log'))
    app.logger.addHandler(logfile_handler)
    app.logger.setLevel(logging.INFO)

    configuration = os.path.join(root_path, 'etc/oldapi.conf')
    conf = Configuration.load(configuration, Ini)
    webservices = conf.get('webservices')

    from beaker.middleware import SessionMiddleware
    from flask.sessions import SessionInterface
    from canopsis.old.account import Account
    from canopsis.old.storage import get_storage

    db = get_storage(account=Account(user='root', group='root'))

    cfg_session = conf.get('session', {})
    session_opts = {
        'session.type': 'mongodb',
        'session.cookie_expires': int(cfg_session.get('cookie_expires', 300)),
        'session.url': '{0}.beaker'.format(db.uri),
        'session.secret': cfg_session.get('secret', 'canopsis'),
        'session.lock_dir': cfg_session.get('data_dir'),
    }

    class BeakerSessionInterface(SessionInterface):
        def open_session(self, app, request):
            return request.environ['beaker.session']

        def save_session(self, app, session, response):
            session.save()

    app.wsgi_app = SessionMiddleware(app.wsgi_app, session_opts)
    app.session_interface = BeakerSessionInterface()

    api = Api(app)

    _auto_import(app, api, webservices)

    return app, api

from flask import session


class APIRoot(Resource):

    resource_routes = ['/api/v3/']

    def get(self):
        self._app.logger.info(session)
        return {'message': 'authenticate with /auth | get v3 routes with /api/v3/routes/all | get other routes with /api/v2/rule/them/all/'}


def exports_v3(app, api):
    APIRoot.init(app, api)

app, api = _init(app)
exports_v3(app, api)
