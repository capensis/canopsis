# -*- coding: utf-8 -*-

from six import raise_from

from grako.model import ModelBuilderSemantics
from grako.exceptions import ParseError

from canopsis.common.utils import singleton_per_scope

from canopsis.alerts.search.walker import Walker
from canopsis.alerts.search.parser import get_parser


def interpret(condition, **kwargs):
    """
    Wrapper for Interpreter.interpret, allowing to interpret search filters
    without having to instanciate a Interpreter class.

    An Intepreter is instantiated as a singleton (per scope).

    :param str condition: Search condition to interpret
    :param dict kwargs: Given to Interpreter constructor the first time it is
      instantiated

    :return: Tuple with corresponding mongo filter and a string carrying
      informations about scope ('all' or 'this')
    :rtype: tuple

    :raises ValueError: If expression is not correct or if parser generation
      failed (bad grammar)
    """

    return singleton_per_scope(Interpreter, kwargs=kwargs).interpret(condition)


class Interpreter(object):
    def __init__(
            self,
            grammar_file,
            parser=None, walker=None,
    ):
        """
        Caches a parser and a walker.

        :param str grammar_file: Generate a parser at runtime with this
          Grako-bnf file.

        :param parser: Parser used to transform expressions to a grako AST
        :type parser: grako.parsing.Parser (or subclass) or None
        :param walker: Walker used to browse AST and serialize an object
          representing the expression
        :type walker: grako.model.Walker (or subclass) or None

        :param dict pkwargs: Additional parameters given to each parser.parse
          call (mostly for start_rule configurability)

        :raises ValueError: if parser generation failed (bad grammar)
        :raises IOError: if grammar_file cannot be read
        :raises OSError: if filename grammar_file be read
        :raises FileNotFoundError: if grammar_file does not exist
        """

        if parser:
            self.parser = parser
        else:
            self.parser = get_parser(
                grammar_file=grammar_file,
                semantics=ModelBuilderSemantics()
            )

        if walker:
            self.walker = walker
        else:
            self.walker = Walker()

    def interpret(self, condition):
        """
        Interpret a condition via a grako.model.NodeWalker subclass.

        :param str condition: Condition to interpret

        :return: Tuple with corresponding mongo filter and a string carrying
          informations about scope ('all' or 'this')
        :rtype: tuple

        :raises ValueError: If condition is not correct
        """

        try:
            model = self.parser.parse(condition, rule_name='start')

        except ParseError as e:
            raise_from(ValueError("Failed to parse"), e)

        res = self.walker.walk(model)

        return res
