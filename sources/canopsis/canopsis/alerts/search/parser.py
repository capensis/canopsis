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

from sys import modules
from imp import new_module
from six import exec_, raise_from

from grako.tool import genmodel
from grako.exceptions import GrakoException
from grako.codegen.python import codegen as pythoncg


def get_parser(grammar_file, **kwargs):
    """
    Get a runtime generated grako parser from a bnf file.

    :param grammar_file: Grako-bnf grammar file.
    :param dict kwargs: Custom parameters that will be given to parser
      constructor.

    :return: Grako parser
    :rtype: grako.parsing.Parser

    :raises ValueError: if parser generation failed (bad grammar)
    :raises IOError: if filename cannot be read
    :raises OSError: if filename cannot be read
    :raises FileNotFoundError: if filename does not exist
    """

    with open(grammar_file) as f:
        grammar = f.read()

        try:
            model = genmodel('Search', grammar, filename=grammar_file)
            code = pythoncg(model)

        except GrakoException as e:
            msg = ("Error trying to generate a grako parser. Is your grammar "
                   "{} correct ?".format(grammar_file))
            raise_from(ValueError(msg), e)

    dynamic_name = 'canopsis.alerts.search.dynamic_parser'
    module = new_module(dynamic_name)
    exec_(code, module.__dict__)
    modules[dynamic_name] = module

    from canopsis.alerts.search.dynamic_parser import SearchParser

    return SearchParser(**kwargs)
