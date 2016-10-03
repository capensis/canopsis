
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..expr import Expression, Function, FuncName

from numbers import Number

from six import iteritems


class FuncNameTest(UTCase):

    def test_contains(self):

        for member in FuncName.__members__.values():

            self.assertTrue(FuncName.contains(member.value))


class ExpressionTest(UTCase):

    def test_init(self):

        name = 'test'

        expr = Expression(name=name)

        self.assertEqual(name, expr.name)
        self.assertIsNone(expr.alias)

    def test_getattr(self):

        expr = Expression.A.B.C

        self.assertEqual(expr.name, 'A.B.C')

    def test_getattr_(self):

        expr = Expression.A.name_

        self.assertEqual(expr.name, 'A.name')

    def _assertfunc(self, func, funcname, *params):

        self.assertIsInstance(func, Function)
        self.assertEqual(funcname.value, func.name)
        self.assertEqual(list(params), func.params)

    def _asserttype(self, expr, method, other=None):

        self.assertRaises(TypeError, getattr(expr, method), other)

    def test__and__(self):

        expr0, expr1 = Expression(name=''), Expression(name='')

        func = expr0 & expr1

        self._assertfunc(func, FuncName.AND, expr0, expr1)

    def test__rand__(self):

        expr0, expr1 = True, Expression(name='')

        func = expr0 & expr1

        self._assertfunc(func, FuncName.AND, expr1, expr0)

    def test__or__(self):

        expr0, expr1 = Expression(name=''), Expression(name='')

        func = expr0 | expr1

        self._assertfunc(func, FuncName.OR, expr0, expr1)

    def test__ror__(self):

        expr0, expr1 = True, Expression(name='')

        func = expr0 | expr1

        self._assertfunc(func, FuncName.OR, expr1, expr0)

    def test__xor__(self):

        expr0, expr1 = Expression(name=''), Expression(name='')

        func = expr0 ^ expr1

        self._assertfunc(func, FuncName.XOR, expr0, expr1)

    def test__rxor__(self):

        expr0, expr1 = True, Expression(name='')

        func = expr0 ^ expr1

        self._assertfunc(func, FuncName.XOR, expr1, expr0)

    def test__gt__(self):

        expr = Expression(name='')

        func = expr > 1

        self._assertfunc(func, FuncName.GT, expr, 1)

        self._asserttype(expr, '__gt__')

    def test__ge__(self):

        expr = Expression(name='')

        func = expr >= 1

        self._assertfunc(func, FuncName.GE, expr, 1)
        self._asserttype(expr, '__ge__')

    def test__lt__(self):

        expr = Expression(name='')

        func = expr < 1

        self._assertfunc(func, FuncName.LT, expr, 1)
        self._asserttype(expr, '__lt__')

    def test__le__(self):

        expr = Expression(name='')

        func = expr <= 1

        self._assertfunc(func, FuncName.LE, expr, 1)
        self._asserttype(expr, '__le__')

    def test__eq__(self):

        expr0, expr1 = Expression(name=''), Expression(name='')

        func = expr0 == expr1

        self._assertfunc(func, FuncName.EQ, expr0, expr1)

    def test__ne__(self):

        expr0, expr1 = Expression(name=''), Expression(name='')

        func = expr0 != expr1

        self._assertfunc(func, FuncName.NE, expr0, expr1)

    def test__add__(self):

        expr = Expression(name='')

        func = expr + 1

        self._assertfunc(func, FuncName.ADD, expr, 1)
        self._asserttype(expr, '__add__')

    def test_concat(self):

        expr = Expression(name='')

        func = expr + ''

        self._assertfunc(func, FuncName.CONCAT, expr, '')
        self._asserttype(expr, '__add__')

    def test__iadd__(self):

        expr = Expression(name='')

        func = expr.__iadd__(1)

        self._assertfunc(func, FuncName.IADD, expr, 1)
        self._asserttype(expr, '__iadd__')

    def test_iconcat(self):

        expr = Expression(name='')

        func = expr.__iadd__('')

        self._assertfunc(func, FuncName.ICONCAT, expr, '')
        self._asserttype(expr, '__iadd__')

    def test__sub__(self):

        expr = Expression(name='')

        func = expr - 1

        self._assertfunc(func, FuncName.SUB, expr, 1)
        self._asserttype(expr, '__sub__')

    def test__isub__(self):

        expr = Expression(name='')

        func = expr.__isub__(1)

        self._assertfunc(func, FuncName.ISUB, expr, 1)
        self._asserttype(expr, '__isub__')

    def test__mul__(self):

        expr = Expression(name='')

        func = expr * 1

        self._assertfunc(func, FuncName.MUL, expr, 1)
        self._asserttype(expr, '__mul__')

    def test__imul__(self):

        expr = Expression(name='')

        func = expr.__imul__(1)

        self._assertfunc(func, FuncName.IMUL, expr, 1)
        self._asserttype(expr, '__imul__')

    def test__truediv__(self):

        expr = Expression(name='')

        func = expr / 1

        self._assertfunc(func, FuncName.DIV, expr, 1)
        self._asserttype(expr, '__truediv__')

    def test__itruediv__(self):

        expr = Expression(name='')

        func = expr.__itruediv__(1)

        self._assertfunc(func, FuncName.IDIV, expr, 1)
        self._asserttype(expr, '__itruediv__')

    def test__invert__(self):

        expr = Expression(name='')

        func = ~expr

        self._assertfunc(func, FuncName.INVERT, expr)

    def test__neg__(self):

        expr = Expression(name='')

        func = -expr

        self._assertfunc(func, FuncName.NEG, expr)

    def test__abs__(self):

        expr = Expression(name='')

        func = abs(expr)

        self._assertfunc(func, FuncName.ABS, expr)

    def test__mod__(self):

        expr = Expression(name='')

        func = expr % r''

        self._assertfunc(func, FuncName.LIKE, expr, '')
        self._asserttype(expr, '__mod__')

    def test__pow__(self):

        expr = Expression(name='')

        func = expr ** 1

        self._assertfunc(func, FuncName.POW, expr, 1)
        self._asserttype(expr, '__pow__')

    def test__ipow__(self):

        expr = Expression(name='')

        func = expr.__ipow__(1)

        self._assertfunc(func, FuncName.IPOW, expr, 1)
        self._asserttype(expr, '__ipow__')

    def test__rand__(self):

        expr = Expression(name='')

        func = 1 & expr

        self._assertfunc(func, FuncName.AND, 1, expr)

    def test__ror__(self):

        expr = Expression(name='')

        func = 1 | expr

        self._assertfunc(func, FuncName.OR, 1, expr)

    def test__rxor__(self):

        expr = Expression(name='')

        func = 1 ^ expr

        self._assertfunc(func, FuncName.XOR, 1, expr)

    def test__radd__(self):

        expr = Expression(name='')

        func = 1 + expr

        self._assertfunc(func, FuncName.ADD, 1, expr)
        self._asserttype(expr, '__radd__')

    def test__rsub__(self):

        expr = Expression(name='')

        func = 1 - expr

        self._assertfunc(func, FuncName.SUB, 1, expr)
        self._asserttype(expr, '__rsub__')

    def test__rmul__(self):

        expr = Expression(name='')

        func = 1 * expr

        self._assertfunc(func, FuncName.MUL, 1, expr)
        self._asserttype(expr, '__rmul__')

    def test__rfloordiv__(self):

        expr = Expression(name='')

        func = 1 // expr

        self._assertfunc(func, FuncName.FLOORDIV, 1, expr)
        self._asserttype(expr, '__rfloordiv__')

    def test__rtruediv__(self):

        expr = Expression(name='')

        func = 1 / expr

        self._assertfunc(func, FuncName.DIV, 1, expr)
        self._asserttype(expr, '__rtruediv__')

    def test__rmod__(self):

        expr = Expression(name='')

        func = expr % expr

        self._assertfunc(func, FuncName.LIKE, expr, expr)
        self._asserttype(expr, '__rmod__')

    def test__rpow__(self):

        expr = Expression(name='')

        func = 1 ** expr

        self._assertfunc(func, FuncName.POW, 1, expr)
        self._asserttype(expr, '__rpow__')

    def test__getslice__(self):

        expr = Expression(name='')

        func = expr[1:2:3]

        self._assertfunc(func, FuncName.GETSLICE, expr, 1, 2, 3)
        self._asserttype(expr, '__getslice__')

    def test__setslice__(self):

        expr = Expression(name='')

        func = expr.__setslice__(1, 2, [])

        self._assertfunc(func, FuncName.SETSLICE, expr, 1, 2, [])
        self._asserttype(expr, '__setslice__')

    def test__delslice__(self):

        expr = Expression(name='')

        func = expr.__delslice__(1, 2)

        self._assertfunc(func, FuncName.DELSLICE, expr, 1, 2)
        self._asserttype(expr, '__delslice__')

    def test__getitem__(self):

        expr = Expression(name='')

        func = expr[1]

        self._assertfunc(func, FuncName.GETITEM, expr, 1)

    def test__setitem__(self):

        expr = Expression(name='')

        func = expr.__setitem__(1, True)

        self._assertfunc(func, FuncName.SETITEM, expr, 1, True)
        self._asserttype(expr, '__setitem__')

    def test__delitem__(self):

        expr = Expression(name='')

        func = expr.__delitem__(1)

        self._assertfunc(func, FuncName.DELITEM, expr, 1)

    def test__rshift__(self):

        expr = Expression(name='')

        func = expr >> 1

        self._assertfunc(func, FuncName.RSHIFT, expr, 1)

    def test__irshift__(self):

        expr = Expression(name='')

        func = expr.__irshift__(1)

        self._assertfunc(func, FuncName.IRSHIFT, expr, 1)

    def test__rrshift__(self):

        expr = Expression(name='')

        func = 1 >> expr

        self._assertfunc(func, FuncName.RSHIFT, 1, expr)

    def test__lshift__(self):

        expr = Expression(name='')

        func = expr << 1

        self._assertfunc(func, FuncName.LSHIFT, expr, 1)

    def test__ilshift__(self):

        expr = Expression(name='')

        func = expr.__ilshift__(1)

        self._assertfunc(func, FuncName.ILSHIFT, expr, 1)

    def test__rlshift__(self):

        expr = Expression(name='')

        func = 1 << expr

        self._assertfunc(func, FuncName.LSHIFT, 1, expr)

    def test_copy(self):

        expr = Expression(name='name', alias='alias')

        cexpr = expr.copy()

        self.assertIsNot(expr, cexpr)
        self.assertEqual(expr.name, cexpr.name)
        self.assertEqual(expr.alias, cexpr.alias)

    def test__call__(self):

        expr = Expression(name='test')
        params = (1,)
        func = expr(*params)

        self.assertIsInstance(func, Function)
        self.assertEqual(func.name, expr.name)
        self.assertEqual(func.params, params)


class FunctionTest(UTCase):

    def test_init_default(self):

        func = Function(name='')

        self.assertEqual(func.params, [])

    def test_init(self):

        params = [1, 2]

        func = Function(name='name', params=params)

        self.assertEqual(func.params, params)

    def test_call(self):

        params = [1, 2]

        func = Function(name='name')(1, 2)

        self.assertEqual(func.params, params)

if __name__ == '__main__':
    main()
