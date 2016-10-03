Description
===========

Utilities for Python.

.. image:: https://img.shields.io/pypi/l/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: License

.. image:: https://img.shields.io/pypi/status/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: Development Status

.. image:: https://img.shields.io/pypi/v/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: Latest release

.. image:: https://img.shields.io/pypi/pyversions/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: Supported Python versions

.. image:: https://img.shields.io/pypi/implementation/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: Supported Python implementations

.. image:: https://img.shields.io/pypi/wheel/b3j0f.utils.svg
   :target: https://travis-ci.org/b3j0f/utils
   :alt: Download format

.. image:: https://travis-ci.org/b3j0f/utils.svg?branch=master
   :target: https://travis-ci.org/b3j0f/utils
   :alt: Build status

.. image:: https://coveralls.io/repos/b3j0f/utils/badge.png
   :target: https://coveralls.io/r/b3j0f/utils
   :alt: Code test coverage

.. image:: https://img.shields.io/pypi/dm/b3j0f.utils.svg
   :target: https://pypi.python.org/pypi/b3j0f.utils/
   :alt: Downloads

.. image:: https://readthedocs.org/projects/b3j0futils/badge/?version=master
   :target: https://readthedocs.org/projects/b3j0futils/?badge=master
   :alt: Documentation Status

.. image:: https://landscape.io/github/b3j0f/utils/master/landscape.svg?style=flat
   :target: https://landscape.io/github/b3j0f/utils/master
   :alt: Code Health

Links
=====

- `Homepage`_
- `PyPI`_
- `Documentation`_

Installation
============

pip install b3j0f.utils

Features
========

This library provides a set of generic tools in order to ease development of projects in python >= 2.6.

Provided tools are:

- chaining: chain object methods calls in a dedicated Chaining object. Such method calls return the Chaining object itself, allowing multiple calls to object methods to be invoked in a concise statement.
- iterable: tools in order to manage iterable elements.
- path: python object path resolver, from object to absolute/relative path or the inverse.
- property: (un)bind/find properties in reflective and oop concerns.
- reflect: tools which ease development with reflective concerns.
- runtime: ease runtime execution (transform dynamic variable to static variable in function, provide safe eval/exec functions).
- proxy: create proxy (from design pattern) objects from a routine or an object which respects the signature and description of the proxified element.
- ut: improve unit tests.
- version: ease compatibility between python version (from 2.x to 3.x).

Examples
========

Chaining
--------

>>> # add characters to a string in one line
>>> from b3j0f.utils.chaining import Chaining, ListChaining
>>> c = Chaining("te").__iadd__("s").__iadd__("t")
>>> # display content of Chaining
>>> c._
test
>>> # call several strings operations on several strings and get operation results in one line
>>> ListChaining("Test", "Example").upper().lower()[:]
[["TEST", "EXAMPLE"], ["test", "example"]]

Iterable
--------

>>> from b3j0f.utils.iterable import is_iterable, first, last, itemat, sliceit, hashiter
>>> is_iterable(1)
False
>>> is_iterable('aze')
True
>>> is_iterable('aze', exclude=str)
False

>>> from b3j0f.utils.version import OrderedDict
>>> od = OrderedDict((('1', 2), ('3', 4), ('5', 6)))
>>> first(od)
'1'
>>> first({}, default='test')
'test'

>>> last(od)
'5'
>>> last('', default='test')
'test'

>>> itemat(od, -1)
'5'
>>> itemat(od, 1)
'3'

>>> sliceit(od, -2, -1)
['3']

>>> hashiter([1, 2])
8

Path
----

>>> from b3j0f.utils.path import lookup, getpath
>>> getpath(lookup)
"b3j0f.utils.path.lookup"
>>> getpath(lookup("b3j0f.utils.path.getpath"))
"b3j0f.utils.path.getpath"

Property
--------

>>> from b3j0f.utils.property import put_properties, get_properties, del_properties
>>> put_properties(min, {'test': True})
>>> assert get_properties(min) == {'test': True}
>>> del_properties(min)
>>> assert get_properties(min) is None

>>> from b3j0f.utils.property import addproperties
>>> def before(self, value, name):  # define a before setter
>>>     self.before = value if hasattr(self, 'after') else None
>>> def after(self, value, name):
>>>     self.after = value + 2  # define a after setter
>>> @addproperties(['test'], bfset=before, afset=after)  # add python properties
>>> class Test(object):
>>>     pass
>>> assert isinstance(Test.test, property)  # assert property is bound
>>> test = Test()
>>> test.test = 2
>>> assert test.update is None  # assert before setter
>>> assert test.test == test._test == 2  # assert default setter
>>> assert test.after == 4

Reflect
-------

>>> from b3j0f.utils.reflect import base_elts, is_inherited
>>> class BaseTest(object):
>>>     def test(self): pass
>>> class Test(BaseTest): pass
>>> class FinalTest(Test): pass
>>> base_elts(FinalTest().test, depth=1)[-1].im_class.__name__
Test
>>> base_elts(FinalTest().test)[-1].im_class.__name__
BaseTest

>>> is_inherited(FinalTest.test)
True
>>> is_inherited(BaseTest.test)
False

Proxy
-----

>>> from b3j0f.utils.proxy import get_proxy, proxified_elt
>>> l = lambda: 2
>>> proxy = get_proxy(l, lambda: 3)
>>> proxy()
3
>>> assert proxified_elt(proxy) is l
True
>>> proxified_elt(proxy)()
2
>>> proxy = get_proxy(l)
>>> proxy()
2
>>> assert proxy is not l
>>> assert proxified_elt(proxy) is l

Runtime
-------

>>> from b3j0f.utils.runtime import safe_eval
>>> try:
>>>     safe_eval('open')
>>> except NameError:
>>>     print('open does not exist')
open does not exist

Version
-------

>>> from b3j0f.utils.version import getcallargs
>>> # getcallargs is same function from python>2.7 for python2.6
>>> from b3j0f.utils.version import PY3, PY2, PY26, PY27
>>> # PY3 is True if python version is 3, etc.

UT
--

>>> from b3j0f.utils.ut import UTCase  # class which inherits from unittest.TestCase
>>> UTCase.assertIs and True  # all methods of python2/3 TestCase are implemented in the UTCase for python v>2.
True

Perspectives
============

- Cython implementation.

Donation
========

.. image:: https://cdn.rawgit.com/gratipay/gratipay-badge/2.3.0/dist/gratipay.png
   :target: https://gratipay.com/b3j0f/
   :alt: I'm grateful for gifts, but don't have a specific funding goal.

.. _Homepage: https://github.com/b3j0f/utils
.. _Documentation: http://b3j0futils.readthedocs.org/en/master/
.. _PyPI: https://pypi.python.org/pypi/b3j0f.utils/


