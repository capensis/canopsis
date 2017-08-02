from flask_restful import Resource as FlaskResource

class Resource(FlaskResource):

    resource_routes = []

    @classmethod
    def init(cls, app, api):
        cls._app = app
        cls._api = api
        cls.add_resources()

    @classmethod
    def add_resources(cls):
        """
        Calls add_resource on api for every route defined in cls._routes.
        """
        cls._api.add_resource(cls, *cls.resource_routes)