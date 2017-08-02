from pkgutil import extend_path
__path__ = extend_path(__path__, __name__)

import importlib
import logging

from flask import Flask
from flask_restful import Api

from canopsis.common.ini_parser import IniParser
from canopsis.webcore.apps.flask.helpers import Resource

app = Flask(__name__)
api = Api(app)

class APIRoot(Resource):

    def get(self):
        return {'message': 'authenticate with /api/v3/auth | get routes with /api/v3/rule/them/all/'}

def exports_v3(app, api):
    api.add_resource(APIRoot, '/api/v3/')

def init(app, api, exports_funcname='exports_v3', configuration='/opt/canopsis/etc/webserver.conf'):
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

    conf = IniParser(configuration, app.logger)
    webservices = conf.get('webservices')

    for webservice, enabled in webservices.items():
        if int(enabled) == 1:
            wsmod = importlib.import_module('canopsis.webcore.services.{}'.format(webservice))

            if hasattr(wsmod, exports_funcname):
                app.logger.info('webservice {}: loading'.format(webservice))
                getattr(wsmod, exports_funcname)(app, api)
                app.logger.info('webservice {}: loaded'.format(webservice))
            else:
                app.logger.info('webservice {}: {} unavailable'.format(webservice, exports_funcname))
        else:
            app.logger.info('webservice {}: skipped'.format(webservice))

init(app, api)
exports_v3(app, api)