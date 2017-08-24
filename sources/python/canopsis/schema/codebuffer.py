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


class Codebuffer(object):

    def __init__(self, indentation=0, indent_string='    ',
                 first_lines='#! -*- coding: utf-8 -*-\n', *args, **kwargs):
        super(Codebuffer, self).__init__(*args, **kwargs)

        self.code = '' #'#! -*- coding: utf-8 -*-\n'

        self.indentation = indentation
        self.indent_string = indent_string

    def __call__(self):
        return self.code

    def indent(self, n=1):
        for i in range(n):
            self.indentation += 1

    def dedent(self, n=1):
        for i in range(n):
            if self.indentation >= 1:
                self.indentation -= 1

    def indent_code(self):
        for i in range(self.indentation):
            self.code += self.indent_string

    def carriage_return(self):
        self.code += '\n'
        self.indent_code()

    def add(self, code):
        self.code += code

    def addline(self, code):
        self.carriage_return()
        self.add(code)
