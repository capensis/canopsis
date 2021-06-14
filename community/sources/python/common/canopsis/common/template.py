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

from canopsis.common.utils import ensure_unicode
from pybars import Compiler, _compiler
from datetime import datetime
from time import time

class Template(Compiler):
    def __init__(self, source, *args, **kwargs):
        super(Template, self).__init__(*args, **kwargs)

        self.source = ensure_unicode(source)
        self.vars = {}

        self.register_helper(u'foreach', self._helper_foreach)
        self.register_helper(u'ifnot', self._helper_ifnot)
        self.register_helper(u'ifeq', self._helper_ifeq)
        self.register_helper(u'set', self._helper_set)
        self.register_helper(u'get', self._helper_get)
        self.register_helper(u'increment', self._helper_increment)
        self.register_helper(u'compact', self._helper_compact)
        self.register_helper(u'formattedDate', self._helper_formatdate)
        self.register_helper(u'strip-newlines', self._help_strip_nl)
        self.register_helper(u'today', self._helper_today)
        self.register_helper(u'state2text', self._helper_state2text)

    def register_helper(self, name, handler):
        name = ensure_unicode(name)

        _compiler._pybars_['helpers'][name] = handler

    def __call__(self, context):
        compiled = self.compile(self.source)
        return u''.join(compiled(context))

    def _helper_foreach(self, this, options, items, sortKey=None):
        result = []
        index = 0

        if items:
            if sortKey:
                items = sorted(items, key=lambda k: k[sortKey])

            for item in items:
                item['@index'] = index
                item['@first'] = (index == 0)
                item['@last'] = (index == len(items) - 1)

                result.extend(options['fn'](item))
                index += 1

        return result

    def _helper_ifnot(self, this, options, context):
        if callable(context):
            context = context(this)

        if not context:
            return options['fn'](this)

    def _helper_ifeq(self, this, options, op1, op2):
        if callable(op1):
            op1 = op1(this)

        if callable(op2):
            op2 = op2(this)

        if op1 == op2:
            return options['fn'](this)

    def _helper_set(self, this, var, val):
        self.vars[var] = val

    def _helper_get(self, this, var):
        return self.vars.get(var, u'')

    def _helper_increment(self, this, var):
        if var in self.vars:
            self.vars[var] += 1

        else:
            self.vars[var] = 1

    def _helper_compact(self, this, options):
        result = u''

        for item in options['fn'](this):
            result = u'{0}{1}'.format(result, item)

        return result.replace('\n', ' ').replace('\r', '')

    def _helper_formatdate(self, this, dtformat, ts):
        return datetime.fromtimestamp(ts).strftime(dtformat)

    def _help_strip_nl(self, this, text):
        if not text:
            text = ''

        return text.replace('\r', '').replace('\n', '')

    def _helper_today(self, this, dtformat='%Y%m%d'):
        return self._helper_formatdate(None, dtformat, time())

    def _helper_state2text(self, this, value, state="OK,WARNING,WARNING,CRITICAL,UNKNOW"):
	if len(state.split(','))-1 < int(value):
		return "UNKNOW"
	else:
		return state.split(',')[int(value)]	
