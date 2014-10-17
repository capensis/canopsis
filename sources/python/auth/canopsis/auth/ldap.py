#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from bottle import request, HTTPError
import ldap

from canopsis.old.storage import get_storage
from canopsis.old.account import Account

from canopsis.webcore.services.account import create_account
from canopsis.auth.base import BaseBackend


class LDAPBackend(BaseBackend):
    name = 'LDAPBackend'

    def get_config(self):
        try:
            record = self.ws.db.get('cservice.ldapconfig')
            return record.dump()

        except KeyError:
            return None

    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = request.environ.get('beaker.session')

            if not s.get('auth_on', False):
                user = self.do_auth()

                if user and not self.install_account(user):
                    return HTTPError(403, 'Forbidden')

            return callback(*args, **kwargs)

        return decorated

    def do_auth(self):
        self.logger.debug('Fetch LDAP configuration from database')

        config = self.get_config()

        if not config:
            self.logger.error('LDAP configuration not found')
            return None

        user = request.params.get('username', default=None)
        passwd = request.params.get('password', default=None)

        dn = config.get('user_dn', None)

        if dn:
            try:
                dn = dn % user

            except TypeError:
                pass

        else:
            dn = '{0}@{1}'.format(user, config['domain'])

        self.logger.debug('Connecting to LDAP server: {0}'.format(
            config['uri']
        ))

        conn = ldap.initialize(config['uri'])
        conn.set_option(ldap.OPT_REFERRALS, 0)
        conn.set_option(ldap.OPT_NETWORK_TIMEOUT, ldap.OPT_NETWORK_TIMEOUT)

        try:
            self.logger.info('Authenticate user {0} to LDAP server'.format(
                user
            ))

            conn.simple_bind_s(dn, passwd)

        except ldap.INVALID_CREDENTIALS:
            self.logger.error('Invalid credentials for user {0}'.format(user))

            # Will try with the next backend
            return None

        try:
            self.logger.debug('Ensure user\'s presence in database: {}'.format(
                user
            ))

            record = self.rights.get_user(user)

        except KeyError:
            record = None

        if not record:
            self.logger.info(
                'Account {0} not found in database, create it'.format(user)
            )

            attrs = [
                config['firstname'],
                config['lastname'],
                config['mail']
            ]

            ufilter = config['user_filter'] % user

            result = conn.search_s(
                config['base_dn'],
                ldap.SCOPE_SUBTREE,
                ufilter,
                attrs
            )

            # TODO: Replace this with crecord.user
            if result:
                dn, data = result[0]

                for field in ['firstname', 'lastname', 'mail']:
                    val = data.get(config[field], None)

                    if val and isinstance(val, list):
                        val = val[0]

                    info[field] = val

                info['firstname'] = info['firstname'].title()
                info['lastname'] = info['lastname'].title()
                info['user'] = user
                info['passwd'] = passwd
                info['external'] = True
                info['aaa_group'] = 'group.Canopsis'
                info['groups'] = ['group.CPS_view']

                account = create_account(info)

            else:
                info = {
                    'user': user,
                    'passwd': passwd,
                    'firstname': user,
                    'lastname': '',
                    'mail': None,
                    'external': True,
                    'aaa_group': 'group.Canopsis',
                    'groups': ['group.CPS_view']
                }

                account = create_account(info)

        return account


def get_backend():
    return LDAPBackend()
