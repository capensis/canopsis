# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from canopsis.common import root_path
from os.path import join
from canopsis.alerts.search.interpreter import interpret
import json


_GRAMMAR_FILE = 'etc/alerts/search/grammar.bnf'
_GRAMMAR_FILE_PATH = join(root_path, _GRAMMAR_FILE)


def parse_search(search):
    """
    Parse a search expression to return a mongo filter and a search scope.

    :param str search: Search expression

    :return: Scope ('this' or 'all') and filter (dict)
    :rtype: tuple

    :raises ValueError: If search is not grammatically correct
    """

    if not search:
        return ('this', {})
    t, q = interpret(search, grammar_file=_GRAMMAR_FILE_PATH)
    if q and type(q) is str:
        q = json.loads(json.dumps(q).replace("\\\\'", "'").decode("string_escape"))
    return t, q


def interpret_search(bnf_search, translate_func):
    """
    translate bnf_filter dict
    :param bnf_search: bnf_filter dict
    :param translate_func: translation function
    :return:
    """
    if not translate_func:
        return bnf_search
    return translate_func(bnf_search)
