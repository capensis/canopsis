from flask_restful import Resource as FlaskResource

class Resource(FlaskResource):

    _routes = []

    @classmethod
    def add_resources(cls, api):
        """
        Calls add_resource on api for every route defined in cls._routes.
        """
        api.add_resource(cls, *cls._routes)