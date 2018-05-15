import importlib
import logging

from flask import Flask
from flask_restful import Api

from canopsis.confng import Configuration, Ini
from canopsis.webcore.apps.flask.helpers import Resource

app = Flask(__name__)

def _auto_import(app, api, exports_funcname='exports_v3', configuration='/opt/canopsis/etc/webserver.conf'):
    conf = Configuration.load(configuration, Ini)
    webservices = conf.get('webservices')

    for webservice, enabled in webservices.items():
        if int(enabled) == 1:
            wsmod = importlib.import_module('{}'.format(webservice))

            if hasattr(wsmod, exports_funcname):
                app.logger.info('webservice {}: loading'.format(webservice))
                getattr(wsmod, exports_funcname)(app, api)
                app.logger.info('webservice {}: loaded'.format(webservice))
            else:
                app.logger.debug('webservice {}: {} unavailable'.format(webservice, exports_funcname))
        else:
            app.logger.debug('webservice {}: skipped'.format(webservice))

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
    logfile_handler = logging.FileHandler('/opt/canopsis/var/log/webserver.log')
    app.logger.addHandler(logfile_handler)
    app.logger.setLevel(logging.INFO)

    from beaker.middleware import SessionMiddleware
    from flask import session
    from flask.sessions import SessionInterface
    from canopsis.old.account import Account
    from canopsis.old.storage import get_storage

    db = get_storage(account=Account(user='root', group='root'))

    session_opts = {
        'session.type': 'mongodb',
        'session.cookie_expires': 300,
        'session.url': '{0}.beaker'.format(db.uri),
        'session.secret': 'canopsis',
        'session.lock_dir': '~/tmp/webcore_cache',
    }

    class BeakerSessionInterface(SessionInterface):
        def open_session(self, app, request):
            session = request.environ['beaker.session']
            return session

        def save_session(self, app, session, response):
            session.save()

    app.wsgi_app = SessionMiddleware(app.wsgi_app, session_opts)
    app.session_interface = BeakerSessionInterface()

    api = Api(app)

    _auto_import(app, api)

    return app, api

from flask import session
from flask_restful import reqparse

class APIRoot(Resource):

    resource_routes = ['/api/v3/']

    def get(self):
        self._app.logger.info(session)
        return {'message': 'authenticate with /auth | get v3 routes with /api/v3/routes/all | get other routes with /api/v2/rule/them/all/'}

def exports_v3(app, api):
    APIRoot.init(app, api)

app, api = _init(app)
exports_v3(app, api)
