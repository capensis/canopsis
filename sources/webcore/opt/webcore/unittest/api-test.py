import sys
sys.path.append("/opt/canopsis/opt/webcore/") 

import traceback

import wsgi_webserver
from webtest import TestApp

from caccount import caccount
from cstorage import get_storage

storage = get_storage(namespace='object')

user = 'root'
pwd = 'root'
shadow = caccount().make_shadow('root')
crypted = caccount().make_tmp_cryptedKey(shadow=shadow)
authkey = storage.get('account.%s' % user, account=caccount(user='root')).data['authkey']

app = TestApp(wsgi_webserver.app)

def quit(code=0):
	wsgi_webserver.unload_webservices()
	sys.exit(code)

def get(uri, args={}, status=None, params=None):
	print "Get %s" % uri
	resp = app.get(uri, status=status, params=params)
	print " + code: %s" % resp.status_int

	return resp

def test(func):
	def wrapper(*args, **kwargs):
		try:
			print "### %s" % func.__name__
			func(*args, **kwargs)
			return True
		except Exception as err:
			traceback.print_exc(file=sys.stdout)
			quit(1)
			return False

	return wrapper

@test
def logout():
	resp = get('/logout')

@test
def login_failed():
	resp = get('/auth/toto/dummy', status=[403])
	resp = get('/auth/%s/dummy' % user, status=[403])

@test
def login_plain():
	resp = get('/auth/%s/%s' % (user, pwd), status=[200])
	logout()

@test
def login_plain_ldap():
	storage.remove('account.toto', account=caccount(user='root'))
	resp = get('/auth/toto/aqzsedrftg123;', status=[200])
	logout()

	resp = get('/auth/toto/tata', status=[403])
	resp = get('/auth/toto/aqzsedrftg123;', status=[200])
	logout()

@test
def login_shadow():
	params={ 'shadow': 1 }
	resp = get('/auth/%s/%s' % (user, shadow), params=params, status=[200])
	logout()

@test
def login_crypted():
	params={ 'crypted': 1 }
	resp = get('/auth/%s/%s' % (user, crypted), params=params, status=[200])
	logout()

@test
def login_authkey():
	resp = get('/autoLogin/dummy', status=[403])
	resp = get('/autoLogin/%s' % authkey, status=[200])
	logout()

	params={ 'authkey': authkey }
	resp = get('/autoLogin/dummy', status=[403])
	resp = get('/account/me', params=params, status=[200])
	logout()

@test
def login_checkAuthPlugin():
	resp = get('/account/checkAuthPlugin1', status=[403])
	#resp = get('/canopsis/auth.html', status=[200])

	# Login
	resp = get('/autoLogin/%s' % authkey, status=[200])

	resp = get('/account/checkAuthPlugin1', status=[200])
	logout()

	resp = get('/auth/canopsis/canopsis', status=[200])
	resp = get('/ui/view', status=[200])
	resp = get('/account/checkAuthPlugin1', status=[200])
	resp = get('/account/checkAuthPlugin2', status=[403])

	logout()

## Execute test
#login_failed()
#login_plain()
login_plain_ldap()
#login_shadow()
#login_crypted()
#login_authkey()
#login_checkAuthPlugin()


quit()
