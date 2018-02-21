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

from pyparsing import (
    Literal, CaselessLiteral, Word, Combine, Optional,
    ZeroOrMore, Forward, nums, alphas, ParseException, delimitedList)
from canopsis.common.utils import ensure_iterable
import math
import operator
import re


class Formulas(object):
    """Class that reads formulas and parse it using EBNF grammar"""
    # map operator symbols to corresponding arithmetic operations
    global epsilon
    epsilon = 1e-12
    opn = {
        "+": operator.add,
        "-": operator.sub,
        "*": operator.mul,
        "/": operator.truediv,
        "^": operator.pow
    }

    fn = {
        "sin": math.sin,
        "cos": math.cos,
        "tan": math.tan,
        "abs": abs,
        "trunc": lambda a: int(a),
        "round": round,
        "max": lambda l: max(float(i) for i in l),
        "min": lambda l: min(float(i) for i in ensure_iterable(l)),
        "sum": lambda l: sum(float(i) for i in l),
        "sgn": lambda a: abs(a) > epsilon and ((a > 0) - (a < 0)) or 0
    }

    def __init__(self, _dict=None):
        self.exprStack = []
        self._bnf = None
        self._dict = _dict  # The dictionnary value as dictionnary {'x':2}
        self.variables = _dict

    def push_first(self, strg, loc, toks):
        '''
        Define an action to apply on the matched tokens
        :param strg: is the original parse string
        :param loc: is the location in the string where matching started
        :param toks: is the list of the matched tokens
        '''
        self.exprStack.append(toks[0])

    def push_minus(self, strg, loc, toks):
        '''
        Define an action to apply on the matched tokens
        :param strg: is the original parse string.
        :param loc: is the location in the string where matching started.
        :param toks: is the list of the matched tokens.
        '''
        if toks and toks[0] == '-':
            self.exprStack.append('unary -')

    def _import(self, _dict):
        '''
        set variables data.
        :param _dict: variables and thier values.
        '''
        self._dict = _dict

    def reset(self):
        '''
        Reset the variables and thier values.
        '''
        self._dict = {}

    def bnf(self):
        '''
        The BNF grammar is defined bellow.
        expop   :: '^'
        multop  :: '*' | '/'
        addop   :: '+' | '-'
        integer :: ['+' | '-'] '0'..'9'+
        atom    :: PI | E | real | fn '(' expr ')' | '(' expr ')'
        factor  :: atom [ expop factor ]*
        term    :: factor [ multop factor ]*
        expr    :: term [ addop term ]*
        '''
        if not self._bnf:
            point = Literal(".")
            e = CaselessLiteral("E")
            fnumber = Combine(
                Word("+-"+nums, nums) +
                Optional(point + Optional(Word(nums))) +
                Optional(e + Word("+-" + nums, nums))
            )
            ident = Word(alphas, alphas + nums + "_$")
            minus = Literal("-")
            plus = Literal("+")
            div = Literal("/")
            mult = Literal("*")
            rpar = Literal(")").suppress()
            lpar = Literal("(").suppress()
            addop = plus | minus
            multop = mult | div
            expop = Literal("^")
            pi = CaselessLiteral("PI")

            expr = Forward()
            atom = (
                Optional("-") +
                (
                    pi |
                    e |
                    fnumber |
                    ident + lpar + delimitedList(expr) + rpar
                ).setParseAction(self.push_first) |
                (lpar + expr.suppress() + rpar)
            ).setParseAction(self.push_minus)

            # The right way to define exponentiation is -> 2^3^2 = 2^(3^2),
            # not (2^3)^2.
            factor = Forward()
            factor << atom + ZeroOrMore(
                (expop + factor).setParseAction(self.push_first)
            )

            term = factor + ZeroOrMore(
                (multop + factor).setParseAction(self.push_first)
            )
            expr << term + ZeroOrMore(
                (addop + term).setParseAction(self.push_first)
            )
            self._bnf = expr
        return self._bnf

    def evaluate_parsing(self, parsing_result):
        '''
        '''
        op = parsing_result.pop()
        if op == 'unary -':
            return -self.evaluate_parsing(parsing_result)

        if op in "+-*/^":
            op2 = self.evaluate_parsing(parsing_result)
            op1 = self.evaluate_parsing(parsing_result)
            return self.opn[op](op1, op2)

        elif op.lower() == "pi":
            return math.pi  # 3.1415926535

        elif op.lower() == "e":
            return math.e  # 2.718281828

        elif op.lower() in self.fn:
            t_op = op.lower()
            if t_op in ('max', 'min', 'sum'):
                if type(parsing_result) is list:
                    return self.fn[t_op](parsing_result)

                return self.fn[t_op](self.evaluate_parsing(parsing_result))

            return self.fn[op](self.evaluate_parsing(parsing_result))

        elif re.search('^[a-zA-Z][a-zA-Z0-9_]*$', op):
            if op in self._dict:
                return self._dict[op]

            else:
                return 0

        elif op[0].isalpha():
            return 0

        else:
            return float(op)

    def evaluate(self, formula):
        '''
        Evaluate the formula
        '''
        if self._dict is not None:
            for k, v in self._dict.iteritems():
                formula = formula.replace(str(k), str(v))

        self.exprStack = []  # reset the stack before each eval
        try:
            results = self.bnf().parseString(formula)
        except ParseException:
            results = ['Parse Failure', formula]
        if len(results) == 0 or results[0] == 'Parse Failure':
            return 'Parse Failure-{}'.format(formula)
        val = self.evaluate_parsing(self.exprStack[:])
        return val
