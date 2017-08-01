from flask import Flask
from flask_restful import Resource
from flask_restful import Api

app = Flask(__name__)
api = Api(app)

class Root(Resource):

    def get(self):
        return 'coucou'

api.add_resource(Root, '/api/v3/')