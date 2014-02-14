import ldap
import clogging

from caccount import caccount
from cstorage import get_storage

from account import create_account

storage = get_storage(namespace='object')

logger = clogging.getLogger()

OPT_NETWORK_TIMEOUT = 1
CONFIG=None
root = account=caccount(user='root')

def get_config():
	global CONFIG

	if not CONFIG:
		record = storage.get("ldap.config", account=root)
		CONFIG = record.dump()

	return CONFIG

def auth(user, password):
	logger.debug("Check auth via ldap")
	config = get_config()

	if not config.get('enable', False):
		logger.debug("Ldap is disable")
		return None

	logger.debug(" + uri: %s" % config['uri'])

	if not password or not user:
		return None

	if not config.get('user_dn', None) and config['domain']:
		dn = "%s@%s" % (user, config['domain'])
	else:
		try:
			dn = config['user_dn'] % user
		except:
			dn = config['user_dn']

	logger.debug(" + dn: %s" % dn)

	# Connect
	conn = ldap.initialize(config['uri'])
	conn.set_option(ldap.OPT_REFERRALS, 0)
	conn.set_option(ldap.OPT_NETWORK_TIMEOUT, OPT_NETWORK_TIMEOUT)

	# Login
	try:
		conn.simple_bind_s(dn, password)
		return True
	
	except ldap.INVALID_CREDENTIALS:
		logger.error(" + Invalid password") 
		return False

	except Exception as err:
		logger.error("%s: %s" % (type(err), err)) 

	return None

def prov(user, password):
	logger.debug("Check prov via ldap")
	config = get_config()

	logger.debug(" + uri: %s" % config['uri'])

	if not config.get('enable', False):
		logger.debug("Ldap is disable")
		return None

	if not password or not user:
		return None

	if not config.get('user_dn', None) and config['domain']:
		dn = "%s@%s" % (user, config['domain'])
	else:
		try:
			dn = config['user_dn'] % user
		except:
			dn = config['user_dn']


	logger.debug(" + dn: %s" % dn)

	# Connect
	conn = ldap.initialize(config['uri'])
	conn.set_option(ldap.OPT_REFERRALS, 0)
	conn.set_option(ldap.OPT_NETWORK_TIMEOUT, OPT_NETWORK_TIMEOUT)

	# Login
	try:
		conn.simple_bind_s(dn, password)
	
	except ldap.INVALID_CREDENTIALS:
		logger.error(" + Invalid password") 
		return None

	except Exception as err:
		logger.error("%s: %s" % (type(err), err)) 
		return None

	# Get informations
	attrs = [
		str(config['lastname']),
		str(config['firstname']),
		str(config['mail'])
	]

	user_filter = config['user_filter'] % user

	logger.debug(" + Filter: %s" % user_filter) 

	result = conn.search_s(config['base_dn'], ldap.SCOPE_SUBTREE, user_filter, attrs)

	if not len(result):
		logger.debug(" + Impossible to find user info")
		return None

	elif len(result) > 1:
		logger.error(" + Too many result")
		return None	

	else:
		(dn, data) = result[0]

		logger.debug(" + dn: %s, data: %s" % (dn, data))

		info = {}
		for field in ['lastname', 'firstname', 'mail']:
			value = data.get(config[field], None)

			if isinstance(value, list) and len(value):
				value = value[0]

			info[field] = value 

		info["lastname"] =  str(info["lastname"]).title()
		info["firstname"] =  str(info["firstname"]).title()
		info["user"] = user
		info["passwd"] = password
		info["external"] = True
		info["aaa_group"] = config.get("aaa_group", "group.Canopsis")

		try:
			info["mail"] = data[config.get("mail","")]
		except Exception, err:
			info["mail"] = "Please set your mail in active directory (field mail)"

		logger.debug(" + Info: %s" % info)

		return create_account(info)