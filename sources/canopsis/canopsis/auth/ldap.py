# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
# along with Canopsis.  If not, see &lt;http://www.gnu.org/licenses/&gt;.
# ---------------------------------

from __future__ import absolute_import
from bottle import request, HTTPError
import ldap
import json

from canopsis.auth.base import BaseBackend


class LDAPBackend(BaseBackend):
    name = "LDAPBackend"

    def get_config(self):
        try:
            record = self.ws.db.get("cservice.ldapconfig")
            return record.dump()

        except KeyError:
            return None

    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = self.session.get()

            if not s.get("auth_on", False) and not self.do_auth(s):
                self.logger.error(u'Impossible to authenticate user')
                return HTTPError(403, "Forbidden")

            return callback(*args, **kwargs)

        return decorated

    def do_auth(self, session):
        self.logger.debug("Fetch LDAP configuration from database")
        mgr = self.rights.get_manager()

        config = self.get_config()

        if not config:
            self.logger.error("LDAP configuration not found")
            return False

        user = request.params.get("username", default=None)
        password = request.params.get("password", default=None)


        if config.get("ldap_uri"):
            self.logger.debug("Connecting to LDAP URI: {0}".format(
                config["ldap_uri"]
            ))
            conn = ldap.initialize(config["ldap_uri"])
        else:
            self.logger.debug("Connecting to LDAP server: {0}:{1}".format(
                config["host"], config["port"]
            ))
            conn = ldap.open(config["host"], config["port"])

        conn.set_option(ldap.OPT_REFERRALS, 0)
        conn.set_option(ldap.OPT_NETWORK_TIMEOUT, ldap.OPT_NETWORK_TIMEOUT)

        if not conn:
            self.logger.error("LDAP server unreachable: {0}:{1}".format(
                config["host"], config["port"]
            ))

            # Will try with the next backend
            return False

        try:
            self.logger.info("Authenticate to LDAP server")
            conn.simple_bind_s(config["admin_dn"], config["admin_passwd"])

        except ldap.INVALID_CREDENTIALS as err:
            self.logger.error("Invalid credentials: {0}".format(err))

            # Will try with the next backend
            return False

        self.logger.info("Authenticate user {0} to LDAP Server".format(user))

        username_attr = config.get('username_attr')
        attrs = [a.encode('utf-8') for a in config["attrs"].values()]
        if username_attr:
            attrs.append(username_attr.encode('utf-8'))
        ufilter = config["ufilter"] % user

        result = conn.search_s(
            config["user_dn"],
            ldap.SCOPE_SUBTREE,
            ufilter,
            attrs
        )

        if not result:
            self.logger.error("No match found for user: {0}".format(user))
            return False

        elif len(result) > 1:
            self.logger.warning("User matched multiple DN: {0}".format(
                json.dumps([dn for dn, _ in result])
            ))

        dn, data = result[0]

        try:
            conn.simple_bind_s(dn, password)

        except ldap.INVALID_CREDENTIALS as err:
            self.logger.error("Invalid credentials: {0}".format(err))
            return False

        username = user
        if username_attr:
            username = data.get(username_attr) or user
            if isinstance(username, list):
                username = username[0]

        info = mgr.get_user(username)

        if not info:
            info = {
                "_id": username,
                "external": True,
                "enable": True,
                "contact": {},
                "role": config["default_role"]
            }

        for field in config["attrs"].keys():
            val = data.get(config["attrs"][field], None)

            if val and isinstance(val, list):
                val = val[0]

            info["contact"][field] = val
            info[field] = val

        account = self.rights.save_user(self.ws, info)
        account['_id'] = username

        session['auth_ldap'] = True
        session.save()

        return self.install_account(username, account)


def get_backend(ws):
    return LDAPBackend(ws)

