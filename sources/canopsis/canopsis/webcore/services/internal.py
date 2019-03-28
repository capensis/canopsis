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

    def get_login_config():
        login_config = {
            'webserver': {provider: 1 for provider in ws.providers},
        }

        records = ws.db.find(
            {'crecord_name': {'$in': ['casconfig', 'ldapconfig']}},
            namespace='object'
        )

        for login_service in records:
            login_service = login_service.dump()
            login_service_name = login_service['crecord_name']
            login_config[login_service_name] = login_service

            ws.logger.info(
                u'found cservices type {}'.format(login_service_name))

            if login_service_name == 'casconfig':
                login_service['server'] = login_service['server'].rstrip('/')
                login_service['service'] = login_service['service'].rstrip('/')
                ws.logger.info(u'cas config : server {}, service {}'.format(
                    login_service['server'],
                    login_service['service'],
                ))

        if "canopsis_cat.webcore.services.saml2" in ws.webmodules:
            result = ws.db.find({'_id': "canopsis"}, namespace='default_saml2')

            login_config["saml2config"] = {
                "url": result[0].data["saml2"]["settings"]["idp"]["singleSignOnService"]["url"]}

        return login_config

    def get_user_interface():
        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        user_interface = user_interface_manager.get()

        if user_interface is not None:
            return {"user_interface": user_interface.to_dict()}

        return {"user_interface": None}

    def get_version():
        store = MongoStore.get_default()
        version_collection = \
            store.get_collection(name=CanopsisVersionManager.COLLECTION)
        document = CanopsisVersionManager(version_collection).\
            find_canopsis_version_document()

        return {CanopsisVersionManager.VERSION_FIELD: document[CanopsisVersionManager.VERSION_FIELD]}

    @ws.application.get('/api/internal/login/login_info')
    def get_internal_login_info():
        cservices = {}

        cservices.update(get_login_config())

        cservices.update(get_user_interface())

        cservices.update(get_version())

        return gen_json(cservices)

    @ws.application.post('/api/internal/login/login_info')
    def update_internal_login_info():

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        return gen_json(user_interface_manager.get().to_dict())
