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

from canopsis.common.ws import route
from bottle import HTTPError

import polib
import sys
import os


def exports(ws):
    @route(ws.application.get)
    def i18n(lang='en'):
        lang_file = os.path.join(sys.prefix, 'locale', lang, 'ui_lang.po')
        translations = {}

        if os.path.isfile(lang_file):
            try:
                po = polib.pofile(lang_file)
                # When language file is properly loaded
                for entry in po:
                    translations[entry.msgid] = entry.msgstr

            except Exception as e:
                return HTTPError(
                    500,
                    'Unable to load po file for lang {}: {} (file: {}'.format(
                        lang, e, lang_file
                    )
                )

        return translations
