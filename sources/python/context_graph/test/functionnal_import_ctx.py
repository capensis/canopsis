import unittest
import requests
import json
import urllib
from http.cookiejar import CookieJar
import sys
import argparse
from time import sleep


def test(auth, serv):
    # auth
    cj = CookieJar()
    opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))
    r = opener.open('http://{0}/autologin/{1}'.format(serv, auth))
    charset = r.info().get_param('charset', 'utf8')
    try:
        response = json.loads(r.read().decode(charset))
    except Exception as err:
        print('bad response from server {0}'.format(err))
        sys.exit()

    if not response['success']:
        print('error: the provided authkey does not match any user')
        sys.exit()

    print('test entity creation')
    js = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"create"}],"links":[]}'
    params = urllib.parse.urlencode({'json': js})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    sleep(2)
    

    print('test update entity')
    update = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{"coucou":"bouh"},"action":"update"}],"links":[]}'
    params = urllib.parse.urlencode({'json': update})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    sleep(2)

    print('test entity deletion')
    deletion = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"delete"}],"links":[]}'
    params = urllib.parse.urlencode({'json': deletion})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    sleep(2)

    print('link test')
    print('link creation between 2 entities')
    link_create = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"create"},{"_id":"resource_1/host_1","name":"resource_1","impact":[],"depends":[],"type":"resource","infos":{},"action":"create"}],"links":[{"_id":"resource_1/host_1-to-host_1","from":"resource_1/host_1","to":"host_1","infos":{},"action":"create"}]}'
    params = urllib.parse.urlencode({'json': link_create})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    sleep(2)

    print('link deletion between 2 entities')
    link_delete = '{"cis":[],"links":[{"_id":"resource_1/host_1-to-host_1","from":"resource_1/host_1","to":"host_1","infos":{},"action":"delete"}]}'
    params = urllib.parse.urlencode({'json': link_delete})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    sleep(2)

    print('cleaning')
    clean = '{"cis":[{"_id":"host_1","name":"host_1","impact":[],"depends":[],"type":"component","infos":{},"action":"delete"},{"_id":"resource_1/host_1","name":"resource_1","impact":[],"depends":[],"type":"resource","infos":{},"action":"delete"}],"links":[]}'
    params = urllib.parse.urlencode({'json': clean})
    print('http://{0}/coucou/bouh?{1}'.format(serv, params))
    req = urllib.request.Request(url='http://{0}/coucou/bouh?{1}'.format(serv, params),method='PUT')
    r = opener.open(req)
    print(r.read())

    print('Done')

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-a', type=str, help='authkey')
    parser.add_argument('-s', type=str, help='server')
    args = parser.parse_args()
    serv = args.s
    auth = args.a
    test(auth, serv)
