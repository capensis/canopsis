import argparse
import json
import sys
import urllib
from http.cookiejar import CookieJar
from unittest import TestCase, TestSuite, TextTestRunner

import pymongo


class Test(TestCase):
    def __init__(self, server, authkey):
        super(Test, self).__init__('test_graph_import')
        self.server = server
        self.authkey = authkey

    def test_graph_import(self):
        # auth
        cj = CookieJar()
        opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))
        r = opener.open('http://{0}:8082/autologin/{1}'.format(self.server, self.authkey))
        charset = r.info().get_param('charset', 'utf8')
        try:
            response = json.loads(r.read().decode(charset))
        except Exception as err:
            print('bad response from server {0}'.format(err))
            sys.exit()

        client = pymongo.MongoClient('mongodb://cpsmongo:canopsis@{0}:27017/canopsis'.format(self.server))
        db = client.canopsis
        col = db.default_entities

        if not response['success']:
            print('error: the provided authkey does not match any user')
            sys.exit()

        print('test entities')
        print('test entity creation')
        js = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"create"}],"links":[]}'
        params = urllib.parse.urlencode({'json': js})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        self.assertDictEqual(col.find_one({'_id': 'host_1'}),
                             {'impact': [], 'name': 'host_1', 'type': 'component', 'infos': {}, '_id': 'host_1',
                              'depends': []})

        print('test update entity')
        update = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{"coucou":"bouh"},"action":"update"}],"links":[]}'
        params = urllib.parse.urlencode({'json': update})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        self.assertDictEqual(col.find_one({'_id': 'host_1'}), {'impact': [], 'name': 'host_1', 'type': 'component',
                                                               'infos': {'coucou': 'bouh'}, '_id': 'host_1', 'depends':
                                                                   []})

        print('test entity deletion')
        deletion = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"delete"}],"links":[]}'
        params = urllib.parse.urlencode({'json': deletion})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        print('link test')
        print('link creation between 2 entities')
        link_create = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"create"},{"_id":"resource_1/host_1","name":"resource_1","impact":[],"depends":[],"type":"resource","infos":{},"action":"create"}],"links":[{"_id":"resource_1/host_1-to-host_1","from":"resource_1/host_1","to":"host_1","infos":{},"action":"create"}]}'
        params = urllib.parse.urlencode({'json': link_create})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        self.assertDictEqual(col.find_one({'_id': 'host_1'}), {'impact': [], 'name': 'host_1', 'type': 'component',
                                                               'infos': {}, '_id': 'host_1', 'depends':
                                                                   ['resource_1/host_1']})
        self.assertDictEqual(col.find_one({'_id': 'resource_1/host_1'}),
                             {'impact': ['host_1'], 'name': 'resource_1', 'type':
                                 'resource', 'infos': {}, '_id':
                                  'resource_1/host_1', 'depends': []})

        print('link deletion between 2 entities')
        link_delete = '{"cis":[],"links":[{"_id":"resource_1/host_1-to-host_1","from":"resource_1/host_1","to":"host_1","infos":{},"action":"delete"}]}'
        params = urllib.parse.urlencode({'json': link_delete})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        self.assertDictEqual(col.find_one({'_id': 'host_1'}), {'impact': [], 'name': 'host_1', 'type': 'component',
                                                               'infos': {}, '_id': 'host_1', 'depends': []})
        self.assertDictEqual(col.find_one({'_id': 'resource_1/host_1'}), {'impact': [], 'name': 'resource_1', 'type':
            'resource', 'infos': {}, '_id':
                                                                              'resource_1/host_1', 'depends': []})

        print('cleaning')
        clean = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"delete"},{"_id":"resource_1/host_1","name":"resource_1","impact":[],"depends":[],"type":"resource","infos":{},"action":"delete"}],"links":[]}'
        params = urllib.parse.urlencode({'json': clean})
        req = urllib.request.Request(url='http://{0}:8082/coucou/bouh?{1}'.format(self.server, params), method='PUT')
        opener.open(req)

        print('Done')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-a', type=str, help='authkey')
    parser.add_argument('-s', type=str, help='server')
    args = parser.parse_args()
    serv = args.s
    auth = args.a
    suite = TestSuite()
    suite.addTest(Test(serv, auth))
    t = TextTestRunner()
    t.run(suite)
