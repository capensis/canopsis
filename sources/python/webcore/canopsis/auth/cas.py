#!/usr/bin/env python
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
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from bottle import request, response, HTTPError, redirect
import xml.etree.ElementTree
import requests

from canopsis.auth.base import BaseBackend

from sys import version as PYVER

if PYVER < '3':
    from urllib import quote_plus

else:
    from urllib.parse import quote_plus


class CASBackend(BaseBackend):
    name = 'CASBackend'
    handle_logout = True

    def get_config(self):
        try:
            record = self.ws.db.get('cservice.casconfig')
            return record.dump()

        except KeyError:
            return None

    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = self.session.get()

            config = self.get_config()

            if not config:
                self.logger.error('CAS configuration not found')
                return callback(*args, **kwargs)

            server = config.get('server')
            service = config.get('service')
            default_role = config.get('default_role')

            service = '{0}/logged_in'.format(service.rstrip('/'))

            if not s.get('auth_on', False):
                res = self.do_auth(s, server, service, default_role)

                if isinstance(res, basestring):
                    return res

                elif not res:
                    return HTTPError(403, 'Forbidden')

            elif request.path in ['/logout', '/disconnect']:
                self.undo_auth(s, server, service)

            return callback(*args, **kwargs)

        return decorated

    def do_auth(self, session, cas_server, service_url, default_role):
        if 'ticket' in request.params:
            self.logger.info(
                'Received ticket from CAS server, start validation'
            )

            ticket = request.params.get('ticket')

            validate_url = '{0}/serviceValidate?ticket={1}&service={2}'.format(
                cas_server.rstrip('/'),
                quote_plus(ticket),
                quote_plus(service_url)
            )

            res = requests.get(validate_url, verify=False)

            if res.status_code != 200:
                self.logger.error('Impossible to validate ticket')
                return False

            self.logger.debug('Parsing CAS server response')

            user = None
            xmlroot = xml.etree.ElementTree.fromstring(res.content)

            for e in xmlroot.iter():
                if e.tag == '{http://www.yale.edu/tp/cas}user':
                    user = e.text

            if not user:
                self.logger.error(
                    'Impossible to find user in response: {0}'.format(
                        res.content
                    )
                )

                return False

            # Get user
            record = self.session.get_user(user)

            if not record:
                self.logger.info('Creating user {0} in database'.format(user))

                record = {
                    '_id': user,
                    'external': True,
                    'enable': True,
                    'contact': {
                        'firstname': user,
                        'lastname': '',
                        'mail': None
                    },
                    'role': default_role
                }

                record = self.rights.save_user(self.ws, record)

            record['_id'] = record.get('_id', user)

            self.logger.info(
                'Authentication validated by CAS server for user {0}'.format(
                    user
                )
            )

            session['auth_cas'] = True
            session.save()

            return self.install_account(user, record)

        else:
            self.logger.info(
                'Redirecting user to CAS server: {0} --> {1}'.format(
                    cas_server,
                    service_url
                )
            )

            url = '{0}/login?service={1}'.format(
                cas_server.rstrip('/'),
                quote_plus(service_url)
            )

            username = request.params.pop('username', default=None)
            password = request.params.pop('password', default=None)

            if username and password:
                response.status = 307
                response.content_type = 'application/x-www-form-urlencoded'
                response.set_header('Location', url)

                return 'username={0}&password={1}'.format(username, password)

            else:
                redirect(url)

    def undo_auth(self, session, cas_server, service_url):
        if session.get('auth_cas', False):
            self.logger.info(
                'Redirecting user to CAS server: {0} --> {1}'.format(
                    cas_server, service_url
                )
            )

            self.session.delete()

            url = '{0}/logout?service={1}&url={1}'.format(
                cas_server.rstrip('/'),
                quote_plus(service_url)
            )

            redirect(url)


def get_backend(ws):
    return CASBackend(ws)
