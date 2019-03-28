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

from urllib import quote_plus


from bottle import request, static_file, redirect, response, HTTPError

from canopsis.common.ws import route
from canopsis.webcore.services import session as session_module
from canopsis.webcore.services import rights as rights_module
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_FORBIDDEN
from canopsis.auth.check import check
from canopsis.userinterface.manager import UserInterfaceManager
from canopsis.version import CanopsisVersionManager
from canopsis.common.mongo_store import MongoStore


def exports(ws):
    session = session_module
    rights = rights_module.get_manager()

    @ws.application.get('/api/internal/login/login_info')
    def get_internal_login_info():
        login_info = {}

        cservices = {
            'webserver': {provider: 1 for provider in ws.providers},
        }

        records = ws.db.find(
            {'crecord_name': {'$in': ['casconfig', 'ldapconfig']}},
            namespace='object'
        )

        for cservice in records:
            cservice = cservice.dump()
            cname = cservice['crecord_name']
            cservices[cname] = cservice

            ws.logger.info(u'found cservices type {}'.format(cname))

            if cname == 'casconfig':
                cservice['server'] = cservice['server'].rstrip('/')
                cservice['service'] = cservice['service'].rstrip('/')
                ws.logger.info(u'cas config : server {}, service {}'.format(
                    cservice['server'],
                    cservice['service'],
                ))

        if "canopsis_cat.webcore.services.saml2" in ws.webmodules:
            result = ws.db.find({'_id': "canopsis"}, namespace='default_saml2')

            cservices["saml2config"] = {
                "url": result[0].data["saml2"]["settings"]["idp"]["singleSignOnService"]["url"]}

        login_info["auth"] = cservices

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        user_interface = user_interface_manager.get()

        if user_interface is not None:
            login_info["user_interface"] = user_interface.to_dict()

        vStore = MongoStore.get_default()
        vCollection = \
            vStore.get_collection(name=CanopsisVersionManager.COLLECTION)
        vDocument = CanopsisVersionManager(vCollection).\
            find_canopsis_version_document()

        login_info[CanopsisVersionManager.VERSION_FIELD] = vDocument[CanopsisVersionManager.VERSION_FIELD]

        return gen_json(login_info)

    @ws.application.post('/api/internal/login/login_info')
    def update_internal_login_info():

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        return gen_json(user_interface_manager.get().to_dict())
