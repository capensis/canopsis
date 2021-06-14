# -*- coding: utf-8 -*-
from hashlib import sha1
from time import time

from canopsis.webcore.services import rights as rights_module


def check(mode='authkey', user=None, password=None):
    def _check_shadow(user, key):
        if user and user['shadowpasswd'].upper() == key.upper():
            return user

        return None

    def _check_plain(user, key):
        shadowpasswd = sha1(key).hexdigest()
        return _check_shadow(user, shadowpasswd)

    def _check_crypted(user, key):
        if user:
            shadowpasswd = user['shadowpasswd'].upper()
            ts = str(int(time() / 10) * 10)
            tmpKey = '{0}{1}'.format(shadowpasswd, ts)

            cryptedKey = sha1(tmpKey).hexdigest().upper()

            if cryptedKey == key.upper():
                return user

        return None

    def _check_authkey(user, key):
        manager = rights_module.get_manager()

        if not user:
            user = manager.user_storage.get_elements(
                query={
                    'crecord_type': 'user',
                    'authkey': key
                }
            )

        if user and user[0]['authkey'] == key:
            return user[0]

        else:
            return None

    handlers = {
        'plain': _check_plain,
        'shadow': _check_shadow,
        'crypted': _check_crypted,
        'authkey': _check_authkey
    }

    if mode in handlers:
        return handlers[mode](user, str(password))

    return None
