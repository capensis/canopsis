# -*- coding: utf-8 -*-
# --------------------------------
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

from canopsis.schema.codebuffer import Codebuffer


class xslt2py(object):
    """
    Converts xslt markups in python code
    """

    def __init__(self, xslt=None, xml_in='xin', xml_out='xout',
                 *args, **kwargs):
        """
        :param str xslt: XSLT object
        :param str xml_in: Name of the variable containing the xml
        :param str xml_out: Produced xml will be stored in this variable
        """

        super(xslt2py, self).__init__(*args, **kwargs)

        self.xls_ns = 'http://www.w3.org/1999/XSL/Transform'

        self.xpath_prefix = ''

        self.xin = xml_in
        self.xout = xml_out

        self.action_mapping = self.get_action_mapping()
        self.operator_mapping = self.get_operator_mapping()

        self.py = Codebuffer()
        self.py.addline(self.xout + ' = \'\'')

    def start(self, tag, attrib):
        if self.xls_ns in tag:
            self.action_mapping[tag][0](attrib)
        else:
            self.xmlnode(tag)

    def end(self, tag):
        if self.xls_ns in tag:
            self.action_mapping[tag][1]()
        else:
            self.xmlnode(tag, close=True)

    def data(self, data):
        stripped = data.strip('\n\t\r ')
        if stripped != '':
            self.xml_out('\'{}\''.format(stripped))

    def comment(self, text):
        pass

    def close(self):
        return self.py()

    def get_action_mapping(self):
        xs = '{' + self.xls_ns + '}'
        mapping = {
            xs + 'stylesheet': [lambda node: None, lambda: None],
            xs + 'output': [lambda node: None, lambda: None],
            xs + 'template': [self.template, lambda: None],
            xs + 'value-of': [self.valueof, lambda: None],
            xs + 'choose': [lambda node: None, lambda: None],
            xs + 'when': [self.when_open, self.when_close],
            xs + 'otherwise': [self.otherwise_open, self.otherwise_close],
        }
        return mapping

    def get_operator_mapping(self):
        mapping = {
            '=': '=',
            '!=': '!=',
            '<': '<',
            '<=': '<=',
            '>': '>',
            '>=': '>=',
        }
        return mapping

    def xml_out(self, expression):
        self.py.addline(self.xout + ' += {}'.format(expression))

    def xmlnode(self, tag, close=False):
        out = '<'
        if close is True:
            out += '/'
        out += tag
        out += '>'

        self.xml_out('\'{}\''.format(out))

    def xpath(self, *args):
        xpath_query = self.xin + '.xpath(\''
        xpath_query += self.xpath_prefix
        for arg in args:
            xpath_query += '/*[local-name() = \\\'{}\\\']'.format(arg)
        xpath_query += '\')[0].text'

        return xpath_query

    def convert_test(self, expression):
        for operator in self.operator_mapping:
            if operator in expression:
                break
        else:
            raise SyntaxError('No operator found in expression \'{}\''.format(
                expression))

        splited_expression = expression.split(operator)
        value = splited_expression[0].strip()
        compared = splited_expression[1].strip()

        pytest = self.xpath(value)
        pytest += ' '
        pytest += self.operator_mapping[operator]
        pytest += ' {}'.format(compared)

        return pytest

    def template(self, attrib):
        prefix_nodes = attrib['match'].split('/')
        for node in prefix_nodes:
            if node != '':
                self.xpath_prefix += '/*[local-name() = \\\'{}\\\']'.format(
                    node)

    def valueof(self, attrib):
        self.xml_out(self.xpath(attrib['select']))

    def when_open(self, attrib):
        pyif = 'if ' + self.convert_test(attrib['test']) + ':'
        self.py.addline(pyif)

        self.py.indent()

    def when_close(self):
        self.py.dedent()

    def otherwise_open(self, attrib):
        self.py.addline('else:')
        self.py.indent()

    def otherwise_close(self):
        self.py.dedent()
