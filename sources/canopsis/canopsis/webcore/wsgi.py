# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from routes import Mapper


class AppMapper():

    def __init__(self, app_mappings):
        """
        :param app_mappings: list of mappings:
            [
                ['app_name', 'regex', wsgi_app],
            ]

            Read documentation of routes.Mapper.connect()
        """
        self.map = Mapper()
        for name, match, app in app_mappings:
            self.map.connect(name, match, app=app)

    def __call__(self, environ, start_response):
        match = self.map.routematch(environ=environ)
        if not match:
            return self.error404(environ, start_response)
        return match[0]['app'](environ, start_response)

    def error404(self, environ, start_response):
        html = b"""\
        <html>
          <head>
            <title>404 - Not Found</title>
          </head>
          <body>
            <h1>404 - Not Found</h1>
            <p>Route matched no application.</p>
          </body>
        </html>
        """
        headers = [
            ('Content-Type', 'text/html'),
            ('Content-Length', str(len(html)))
        ]
        start_response('404 Not Found', headers)
        return [html]


def get_default_app():
    from canopsis.webcore.apps.bottle.app import get_default_app as get_default_app_v1
    from canopsis.webcore.apps.flask.app import app as app_v3

    app_v1 = get_default_app_v1(logger=app_v3.logger)

    app_mappings = []
    app_mappings.append(['app_v3', '/api/v3{path:.*}', app_v3])
    app_mappings.append(['app_v2', '/api/v2{path:.*}', app_v1])
    app_mappings.append(['app_v1', '/{path:.*}', app_v1])

    app = AppMapper(app_mappings)

    return app
