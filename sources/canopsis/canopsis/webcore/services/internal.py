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


from bottle import request

from canopsis.webcore.services import session as session_module
from canopsis.webcore.services import rights as rights_module
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.userinterface.manager import UserInterfaceManager
from canopsis.version import CanopsisVersionManager
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import CollectionError

VALID_USER_INTERFACE_PARAMS = [
    'app_title', 'footer', 'login_page_description', 'logo', 'language', 'popup_timeout'
]

VALID_POPUP_UNIT = {
    's', 'h', 'm'
}

VALID_POPUP_PARAMS = {
    'unit', 'interval'
}

VALID_CANOPSIS_EDITIONS = [
    'cat', 'core'
]

VALID_CANOPSIS_STACKS = [
    'go', 'python'
]

VALID_CANOPSIS_LANGUAGES = [
    'en', 'fr'
]


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
    document = CanopsisVersionManager(version_collection). \
        find_canopsis_document()

    if document is not None:
        return {
            CanopsisVersionManager.EDITION_FIELD: document.get(CanopsisVersionManager.EDITION_FIELD, ""),
            CanopsisVersionManager.STACK_FIELD: document.get(CanopsisVersionManager.STACK_FIELD, ""),
            CanopsisVersionManager.VERSION_FIELD: document.get(
                CanopsisVersionManager.VERSION_FIELD, "")
        }

    return {
        CanopsisVersionManager.EDITION_FIELD: "",
        CanopsisVersionManager.STACK_FIELD: "",
        CanopsisVersionManager.VERSION_FIELD: ""
    }


def get_login_config(ws):
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

    return {"login_config": login_config}


def check_values(ws, edition, stack):
    if edition is not None and edition not in VALID_CANOPSIS_EDITIONS:
        ws.logger.error("edition is an invalid value : {}".format(edition))
        return False

    if stack is not None and stack not in VALID_CANOPSIS_STACKS:
        ws.logger.error("stack is an invalid value : {}".format(stack))
        return False

    return True


def sanitize_popup_timeout(popup_setting):
    if not isinstance(popup_setting, dict):
        return {
            'unit': 's',
            'interval': 10
        }

    else:
        if 'unit' not in popup_setting or popup_setting['unit'] not in VALID_POPUP_UNIT:
            popup_setting['unit'] = 's'
        if 'interval' not in popup_setting or not isinstance(popup_setting['interval'], int) or \
            popup_setting['interval'] < 0:
            popup_setting['interval'] = 10
    # remove redundant keys in popup_timeout
    for k in popup_setting.keys():
        if k not in VALID_POPUP_PARAMS:
            popup_setting.pop(k)
    return popup_setting


def exports(ws):
    session = session_module
    rights = rights_module.get_manager()

    @ws.application.get('/api/internal/login/login_info', skip=ws.skip_login)
    def get_internal_login_info():
        cservices = {}
        cservices.update(get_login_config(ws))
        cservices.update(get_user_interface())
        cservices.update(get_version())
        return gen_json(cservices)

    @ws.application.post('/api/internal/properties')
    def update_canopsis_properties():
        try:
            doc = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        store = MongoStore.get_default()
        version_collection = store.get_collection(
            name=CanopsisVersionManager.COLLECTION)

        try:
            ok = check_values(
                ws, doc.get("edition"), doc.get("stack"))
            if ok:
                success = CanopsisVersionManager(version_collection). \
                    put_canopsis_document(
                    doc.get("edition"), doc.get("stack"), None)

                if not success:
                    return gen_json_error({'description': 'failed to update edition/stack'},
                                          HTTP_ERROR)
                return gen_json({})
            else:
                err = 'Invalid value(s).'
                ws.logger.error(err)
                return gen_json_error(
                    {'description': err},
                    HTTP_ERROR
                )

        except CollectionError as ce:
            ws.logger.error('Update edition/stack error: {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while updating edition/stack values'},
                HTTP_ERROR
            )

    @ws.application.get('/api/internal/app_info')
    def get_internal_app_info():
        cservices = {}
        user_interface = get_user_interface().get("user_interface", None)
        if user_interface is not None:
            for key in user_interface.keys():
                if key not in ['app_title', 'logo', 'language', 'popup_timeout']:
                    user_interface.pop(key)
            cservices.update(user_interface)
        ws.logger.error(get_version())
        cservices.update(get_version())
        ws.logger.warning(cservices)

        return gen_json(cservices)

    @ws.application.post('/api/internal/user_interface')
    @ws.application.put('/api/internal/user_interface')
    def update_internal_interface():
        try:
            interface = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        if interface is None:
            return gen_json_error(
                {'description': 'nothing to insert'},
                HTTP_ERROR
            )

        for key in interface.keys():
            if key not in VALID_USER_INTERFACE_PARAMS:
                interface.pop(key)
            elif key == 'popup_timeout':
                # set default value for popup_timeout
                interface[key]['info'] = sanitize_popup_timeout(interface[key].get('info'))
                interface[key]['error'] = sanitize_popup_timeout(interface[key].get('error'))

        # set default value for popup_timeout
        if 'popup_timeout' not in interface.keys():
            interface['popup_timeout'] = dict(info={
                'unit': 's',
                'interval': 10
            }, error={
                'unit': 's',
                'interval': 10
            })

        language = interface.get('language', None)
        if language is not None and language not in VALID_CANOPSIS_LANGUAGES:
            ws.logger.error(
                "language is an invalid value : {}".format(language))
            return gen_json_error(
                {'description': "language is an invalid value : {}".format(
                    language)},
                HTTP_ERROR
            )

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        if len(interface) > 0:
            return gen_json(user_interface_manager.update(interface))
        return gen_json_error(
            {'description': 'nothing to insert'},
            HTTP_ERROR
        )

    @ws.application.delete('/api/internal/user_interface')
    def delete_internal_interface():
        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        return gen_json(user_interface_manager.delete())
