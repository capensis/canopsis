#!/usr/bin/env python
# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2014 Jonathan Labéjof <jonathan.labejof@gmail.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
# --------------------------------------------------------------------

from unittest import main

from time import sleep

from random import random

from six import PY2

from ..ut import UTCase
from ..property import (
    get_properties, put_properties, del_properties,
    get_local_properties, get_local_property,
    firsts, remove_ctx,
    get_first_property, get_first_properties,
    setdefault, put_property,
    __B3J0F__PROPERTIES__,
    find_ctx, addproperties, _protectedattrname
)


class FindCTXTest(UTCase):
    """
    Test find_ctx method
    """
    def test_notctx(self):
        """
        Test with elt without ctx
        """

        elt = 1

        ctx = find_ctx(elt=elt)

        self.assertIs(ctx, elt)

    def test_method(self):
        """
        Test method ctx
        """

        class A:
            def a(self):
                pass

        elt = A.a

        ctx = find_ctx(elt=elt)

        if PY2:
            self.assertIs(A, ctx)
        else:
            self.assertIs(elt, ctx)

    def test_instance_method(self):
        """
        Test instance method ctx
        """

        class A:
            def a(self):
                pass

        a = A()

        elt = a.a

        ctx = find_ctx(elt=elt)

        self.assertIs(a, ctx)


class RemoveCTXTest(UTCase):
    """
    Test remove_ctx function.
    """

    def test_no_properties(self):
        """
        Test if not properties exist.
        """

        properties = remove_ctx({})

        self.assertFalse(properties)

    def test_empty(self):
        """
        Test with empty properties.
        """

        properties = remove_ctx({'test': []})

        self.assertFalse(properties['test'])

    def test_one(self):
        """
        Test with one value property.
        """

        properties = {'test': [('elt', 0)]}

        properties = remove_ctx(properties)

        self.assertEqual(properties['test'], [0])

    def test_properties(self):
        """
        Test with many property values.
        """

        properties = {'test': [('elt', 0), ('elt', 1)]}

        properties = remove_ctx(properties)

        self.assertEqual(properties['test'], [0, 1])


class PropertyTest(UTCase):
    """
    Test scenarios of puting/getting/deleting properties.
    """

    def _assert_properties(self, elt, count=10):

        properties = get_properties(elt=elt)
        self.assertFalse(properties)
        local_properties = get_local_properties(elt=elt)
        self.assertFalse(local_properties)

        properties = dict((str(i), i) for i in range(count))
        put_properties(elt=elt, properties=properties)

        local_properties = get_local_properties(elt=elt)
        properties = get_properties(elt=elt)

        self.assertEqual(len(properties), count)
        self.assertEqual(len(local_properties), count)
        for index in range(count):
            name = str(index)
            self.assertIs(properties[name][0][0], elt)
            self.assertIs(properties[name][0][1], index)
            self.assertIs(local_properties[name], index)

        for index in range(count):
            name = str(index)
            properties = get_properties(elt=elt, keys=name)
            self.assertEqual(len(properties), 1)
            self.assertIs(properties[name][0][0], elt)
            self.assertIn(name, properties)
            self.assertIs(properties[name][0][1], index)
            local_properties = get_local_properties(elt=elt, keys=name)
            self.assertEqual(len(local_properties), 1)
            self.assertIn(name, local_properties)
            self.assertIs(local_properties[name], index)

        del_properties(elt=elt, keys='0')

        properties = get_properties(elt=elt)
        self.assertEqual(len(properties), count - 1)

        local_properties = get_local_properties(elt=elt)
        self.assertEqual(len(local_properties), count - 1)

        for i in range(1, count):
            name = str(i)
            self.assertIs(properties[name][0][1], i)
            self.assertIs(local_properties[name], i)

        del_properties(elt=elt)
        properties = get_properties(elt=elt)
        self.assertFalse(properties)
        local_properties = get_local_properties(elt=elt)
        self.assertFalse(local_properties)

        if hasattr(elt, '__dict__'):
            self.assertNotIn(__B3J0F__PROPERTIES__, elt.__dict__)

    def test_builtin(self):
        """
        Test lookup of builtin
        """

        self._assert_properties(min)

    def test_object(self):
        """
        Test scenario on an object
        """
        self._assert_properties(1)

    def test_none(self):
        """
        Test scenario on None
        """
        self._assert_properties(None)

    def test_lambda(self):
        """
        Test scenario on a lambda expression
        """
        self._assert_properties(lambda: None)

    def test_function(self):
        """
        Test scenario on a function.
        """
        def a():
            pass

        self._assert_properties(a)

    def _test_inheritance(self, first, second, ctx1=None, ctx2=None, count=5):
        """
        Test inherited properties between first and second elements.
        """

        properties = dict((str(i), i) for i in range(count))
        put_properties(elt=first, properties=properties, ctx=ctx1)

        properties = get_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(properties), count)
        self.assertEqual(first, properties['0'][0][0])

        local_properties = get_local_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(local_properties), 0)

        properties = dict((str(i), i) for i in range(count))
        put_properties(elt=second, properties=properties, ctx=ctx2)

        properties = get_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(properties), count)
        self.assertEqual(first, properties['0'][1][0])
        self.assertEqual(second, properties['0'][0][0])

        local_properties = get_local_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(local_properties), count)

        del_properties(elt=first, ctx=ctx1)
        properties = get_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(properties), count)
        self.assertEqual(len(properties['0']), 1)
        self.assertIs(second, properties['0'][0][0])
        local_properties = get_local_properties(elt=second, ctx=ctx2)
        self.assertEqual(len(local_properties), count)

        del_properties(elt=second, ctx=ctx2)
        properties = get_properties(elt=second, ctx=ctx2)
        self.assertFalse(properties)
        local_properties = get_local_properties(elt=second, ctx=ctx2)
        self.assertFalse(local_properties)

    def test_class(self):
        """
        Test scenario on a class.
        """
        class A(object):
            pass

        self._assert_properties(A)

        class B(A):
            pass

        self._test_inheritance(A, B)

    def test_instance(self):
        """
        Test scenario on an instance.
        """
        class A:
            pass

        a = A()

        self._assert_properties(a)

        self._test_inheritance(A, a)

    def test_namespace(self):
        """
        Test scenario on a namespace.
        """
        class A:
            pass

        class B(A):
            pass

        self._assert_properties(A)

        self._test_inheritance(A, B)

    def test_method(self):
        """
        Test scenario on a method.
        """
        class A:
            def a(self):
                pass

        class B(A):
            pass

        self._assert_properties(A.a)

        self._test_inheritance(A.a, B.a, ctx1=A, ctx2=B)

    def test_bound_method(self):
        """
        Test scenario on a bound method.
        """
        class A(object):
            def a(self):
                pass

        a = A()

        self._test_inheritance(A.a, a.a, ctx1=A)

    def test_module(self):
        """
        Test scenario on a module.
        """
        import b3j0f

        self._assert_properties(b3j0f)

    def test_property_module(self):
        """
        Test scenario on the property module.
        """
        import b3j0f.utils.property

        self._assert_properties(b3j0f.utils.property)

    def test_dict(self):
        """
        Test scenario on a dictionary.
        """
        elt = {}

        self._assert_properties(elt)

    def test_list(self):
        """
        Test scenario on a list.
        """
        elt = []

        self._assert_properties(elt)

    def test_inheritance(self):
        """
        Test scenario on inherited methods
        """

        key, a, b, c = 'test', 1, 2, 3

        class A:
            def test(self):
                pass

        class B(A):
            pass

        class C(B):
            def test(self):
                pass

        put_property(elt=A.test, ctx=A, key=key, value=a)
        put_property(elt=B.test, ctx=B, key=key, value=b)
        put_property(elt=C.test, ctx=C, key=key, value=c)

        properties = get_properties(elt=C.test, ctx=C)
        self.assertEqual(len(properties), 1)

        self.assertEqual(properties[key][2], (A.test, a))
        self.assertEqual(properties[key][1], (B.test, b))
        self.assertEqual(properties[key][0], (C.test, c))


class TTLTest(UTCase):
    """
    Test ttl parameters while puting property.
    """

    def tearDown(self):
        """
        Del properties at the end of tests.
        """
        del_properties(elt=self)
        self.assertNotIn(__B3J0F__PROPERTIES__, self.__dict__)

    def test_zero(self):
        """
        Test with ttl = 0
        """

        put_property(elt=self, ttl=0, key='name', value=1)

        sleep(0.1)

        properties = get_local_properties(elt=self)

        self.assertFalse(properties)

    def test_100(self):
        """
        Test with ttl = 100
        """

        ttl = 0.1

        put_property(elt=self, ttl=ttl, key='name', value=1)

        properties = get_local_properties(elt=self)

        self.assertTrue(properties)

        sleep(ttl + 0.2)

        properties = get_local_properties(elt=self)

        self.assertFalse(properties)


class OneTest(UTCase):
    """
    UT for get_local_property, get_first_property and get_first_properties.
    """

    def tearDown(self):
        """
        delete properties.
        """
        del_properties(elt=self)
        self.assertNotIn(__B3J0F__PROPERTIES__, self.__dict__)

    def test_first_none(self):
        """
        test first default.
        """
        _property = get_first_property(elt=self, key='a', default=2)

        self.assertEqual(_property, 2)

    def test_first(self):
        """
        test first on existing property
        """
        put_property(elt=self, key='a', value=1)

        _property = get_first_property(elt=self, key='a', default=2)

        self.assertEqual(_property, 1)

    def test_firsts_none(self):
        """
        test firsts default.
        """
        _properties = get_first_properties(elt=self, keys=['a', 'b'])

        self.assertFalse(_properties)

    def test_firsts(self):
        """
        test firsts on existing property
        """
        properties = {'a': 1, 'b': 2}
        put_properties(elt=self, properties=properties)

        _properties = get_first_properties(elt=self, keys='a')

        self.assertIn('a', _properties)
        self.assertEqual(_properties['a'], 1)

    def test_none_local(self):
        """
        test local default.
        """
        local_property = get_local_property(elt=self, key='a', default=2)

        self.assertEqual(local_property, 2)

    def test_local(self):
        """
        test local with existing property.
        """
        put_property(elt=self, key='a', value=1)

        local_property = get_local_property(elt=self, key='a', default=2)

        self.assertEqual(local_property, 1)


class UnifyTest(UTCase):
    """
    Test firsts function.
    """

    def test_empty(self):
        """
        Test empty properties.
        """
        properties = []

        unified_properties = firsts(properties=properties)

        self.assertFalse(unified_properties)

    def test_one_value(self):
        """
        Test one property
        """

        properties = {1: [(2, 3)]}

        unified_properties = firsts(properties=properties)

        self.assertIn(1, unified_properties)
        self.assertEqual(unified_properties[1], 3)

    def test_values(self):
        """
        Test several properties.
        """

        count = 10

        properties = {}

        for i in range(count):
            properties[i] = []
            for j in range(i + 1):
                properties[i].append((i, j))

        unified_properties = firsts(properties=properties)

        self.assertEqual(len(unified_properties), count)
        for i in range(count):
            self.assertIs(unified_properties[i], 0)


class SetDefaultTest(UTCase):
    """
    Test setdefault function
    """

    def setUp(self):
        """
        Set attributes to self key=test and new_value=2
        """
        self.key = 'test'

        self.new_value = 2

    def tearDown(self):
        """
        del properties
        """
        del_properties(elt=self)

        self.assertNotIn(__B3J0F__PROPERTIES__, self.__dict__)

    def test_exists(self):
        """
        Test with an existing property
        """

        put_properties(elt=self, properties={self.key: self.new_value + 1})

        value = setdefault(elt=self, key=self.key, default=self.new_value)

        self.assertNotEqual(value, self.new_value)

    def test_new(self):
        """
        Test on a missing property
        """

        value = setdefault(elt=self, key=self.key, default=self.new_value)

        self.assertEqual(value, self.new_value)


class TestAddProperties(UTCase):
    """Test the addproperties decorator."""

    def setUp(self):

        self.count = 5
        # generate count random property names
        self.names = [str(random()) for i in range(self.count)]

        self.getternames = set()
        self.setternames = {}
        self.deleternames = set()

    def getter(self, name=None):
        """Property getter."""
        if name is None:
            name = self.name
        self.getternames.add(name)
        return self

    def setter(self, value, name=None):
        """Property setter."""
        if name is None:
            name = self.name
        self.setternames[name] = value

    def deleter(self, name=None):
        """Property deleter."""
        if name is None:
            name = self.name
        self.deleternames.add(name)

    def test_empty_cls(self):
        """Test to add properties on an empty cls."""
        # add properties
        @addproperties(names=self.names)
        class Test(object):
            """Test class."""

        test = Test()

        for name in self.names:
            # assert properties are in Test
            prop = getattr(Test, name)
            self.assertTrue(isinstance(prop, property))
            # get protected attr name
            protectedattrname = _protectedattrname(name)
            # assert getter
            self.assertFalse(hasattr(test, protectedattrname))
            value = getattr(test, name)
            self.assertIsNone(value)
            setattr(test, protectedattrname, self)
            value = getattr(test, name)
            self.assertIs(value, self)
            delattr(test, protectedattrname)
            # assert setter
            self.assertFalse(hasattr(test, protectedattrname))
            setattr(test, name, self)
            self.assertTrue(hasattr(test, protectedattrname))
            # assert deleter
            self.assertTrue(hasattr(test, protectedattrname))
            delattr(test, name)
            self.assertFalse(hasattr(test, protectedattrname))

    def test_not_empty_cls(self):
        """Test to add properties with existing properties."""
        class Test(object):
            """Test class."""

        for name in self.names:
            prop = property(
                fget=self.getter, fset=self.setter, fdel=self.deleter
            )
            setattr(Test, name, prop)

        addproperties(names=self.names)(Test)

        test = Test()

        for name in self.names:
            protectedattrname = _protectedattrname(name)
            # assert existing getter
            self.name = name
            self.assertNotIn(name, self.getternames)
            self.assertFalse(hasattr(test, protectedattrname))
            value = getattr(test, name)
            self.assertIn(name, self.getternames)
            self.assertIs(value, self)
            # assert exising setter
            self.assertNotIn(name, self.setternames)
            setattr(test, name, value)
            self.assertIn(name, self.setternames)
            self.assertIs(self.setternames[name], self)
            # assert existing deleter
            self.assertNotIn(name, self.deleternames)
            delattr(test, name)
            self.assertFalse(hasattr(test, protectedattrname))
            self.assertIn(name, self.deleternames)

    def test_ab(self):
        """Test after/before getter/setter/deleter."""

        @addproperties(
            names=self.names,
            bfget=self.getter, afget=self.getter,
            bfset=self.setter, afset=self.setter,
            bfdel=self.deleter, afdel=self.deleter
        )
        class Test(object):
            """Test class."""

        test = Test()

        for name in self.names:
            # get protected attr name
            protectedattrname = _protectedattrname(name)
            # assert existing getter
            self.assertNotIn(name, self.getternames)
            self.assertFalse(hasattr(test, protectedattrname))
            value = getattr(test, name)
            self.assertIs(value, self)
            # assert setter
            setattr(test, protectedattrname, self)
            value = getattr(test, protectedattrname)
            self.assertIs(value, self)
            value = getattr(test, name)
            self.assertIn(name, self.getternames)
            self.assertIs(value, self)
            # assert exising setter
            self.assertNotIn(name, self.setternames)
            setattr(test, name, value)
            self.assertTrue(hasattr(test, protectedattrname))
            self.assertIn(name, self.setternames)
            self.assertIs(self.setternames[name], self)
            # assert existing deleter
            self.assertNotIn(name, self.deleternames)
            delattr(test, name)
            self.assertFalse(hasattr(test, protectedattrname))
            self.assertIn(name, self.deleternames)


if __name__ == '__main__':
    main()
