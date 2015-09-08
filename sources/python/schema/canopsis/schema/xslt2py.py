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

from canopsis.schema.codebuffer import Codebuffer


class xslt2py(object):
    """
    Converts xslt markups in python code. The method used to parse xslt
    and compile it in python code is based on this tutorial :
    http://lxml.de/parsing.html#the-target-parser-interface
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
        """
        cf http://lxml.de/parsing.html#the-target-parser-interface
        """
        if self.xls_ns in tag:
            self.action_mapping[tag][0](attrib)
        else:
            self.xmlnode(tag)

    def end(self, tag):
        """
        cf http://lxml.de/parsing.html#the-target-parser-interface
        """
        if self.xls_ns in tag:
            self.action_mapping[tag][1]()
        else:
            self.xmlnode(tag, close=True)

    def data(self, data):
        """
        cf http://lxml.de/parsing.html#the-target-parser-interface
        """
        stripped = data.strip('\n\t\r ')
        if stripped != '':
            self.xml_out('\'{}\''.format(stripped))

    def comment(self, text):
        """
        cf http://lxml.de/parsing.html#the-target-parser-interface
        """

    def close(self):
        """
        cf http://lxml.de/parsing.html#the-target-parser-interface
        """
        return self.py()

    def get_action_mapping(self):
        """
        In order to translate xslt to python, here are defined
        functions to call when a xslt node is detected (open and close)

        :return: dictionnary with xslt nodes in keys and mapped
           functions in values
        :rtype: dict
        """
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
        """
        Some xslt operators are different in xslt and python. This
        function returns a matching table between the two langages.

        :return: dictonnary with xslt operator tokens and their
           equivalent in python
        :rtype: dict
        """
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
        """
        Add a factored expression to transformed xml data

        :param str expression: some piece of xml (or format of the
           returned data)
        """
        self.py.addline(self.xout + ' += {}'.format(expression))

    def xmlnode(self, tag, close=False):
        """
        Provide an easy way to create xml nodes in transformed xml data

        :param str tag: name of the xml tag (text in '<' and '>')
        :param bool close: wether or not the node ends with '>' or '/>'
        """
        out = '<'
        if close is True:
            out += '/'
        out += tag
        out += '>'

        self.xml_out('\'{}\''.format(out))

    def xpath(self, *args):
        """
        An easy way to query a value in input xml data.

        :param *args: identificatier for node to query
        """
        xpath_query = self.xin + '.xpath(\''
        xpath_query += self.xpath_prefix
        for arg in args:
            xpath_query += '/*[local-name() = \\\'{}\\\']'.format(arg)
        xpath_query += '\')[0].text'

        return xpath_query

    def convert_test(self, expression):
        """
        Convert a xslt test expression in python test expression.

        :param str expression: xslt expression
        :return: test expression in python
        :rtype: str
        :raises SyntaxError: if expression syntax is not recognized
        """
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
        """
        Triggered when a `template` xslt node is detected.

        Loads xml prefix (`match` parameter of `template` node)
        """
        prefix_nodes = attrib['match'].split('/')
        for node in prefix_nodes:
            if node != '':
                self.xpath_prefix += '/*[local-name() = \\\'{}\\\']'.format(
                    node)

    def valueof(self, attrib):
        """
        Triggered when a `value-of` xslt node is detected.

        Copies the content in transformed data.
        """
        self.xml_out(self.xpath(attrib['select']))

    def when_open(self, attrib):
        """
        Triggered when a `when` xslt node is opened.

        Open a conditional indented bloc in factored code.
        """
        pyif = 'if ' + self.convert_test(attrib['test']) + ':'
        self.py.addline(pyif)

        self.py.indent()

    def when_close(self):
        """
        Triggered when a `when` node is closed.

        Dedent a bloc in factored code.
        """
        self.py.dedent()

    def otherwise_open(self, attrib):
        """
        Triggered when an `otherwise` node is opened.

        Open a conditional alternative indented bloc in factored code
        (else).
        """
        self.py.addline('else:')
        self.py.indent()

    def otherwise_close(self):
        """
        Triggered when an `otherwise` node is closed.

        Dedent a bloc in factored code.
        """
        self.py.dedent()
