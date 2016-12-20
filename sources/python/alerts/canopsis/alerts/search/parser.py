# -*- coding: utf-8 -*-

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
