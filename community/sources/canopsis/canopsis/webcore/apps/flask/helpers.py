# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import flask_restful

from flask import session
from flask_restful import Resource as FlaskResource


def authenticate(func):
    def wrapper(*args, **kwargs):
        # compatibility with bottle and /auth route
        if session.get('auth_on', False):
            return func(*args, **kwargs)

        flask_restful.abort(401)

    return wrapper


class Resource(FlaskResource):
    """
    Define routes in cls.resource_routes respecting
    http://flask-restful.readthedocs.io/en/0.3.5/quickstart.html#endpoints
    """

    resource_routes = []
    method_decorators = [authenticate]

    @classmethod
    def init(cls, app, api):
        cls._app = app
        cls._api = api
        cls.add_resources()

    @classmethod
    def add_resources(cls):
        """
        Calls add_resource on api for every route defined in cls.resource_routes list.
        """
        cls._api.add_resource(cls, *cls.resource_routes)
