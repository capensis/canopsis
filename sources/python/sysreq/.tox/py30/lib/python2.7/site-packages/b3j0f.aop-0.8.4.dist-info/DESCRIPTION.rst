Description
-----------

This project is an Aspect Oriented Programming library for python with reflective concerns.

.. image:: https://img.shields.io/pypi/l/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: License

.. image:: https://img.shields.io/pypi/status/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: Development Status

.. image:: https://img.shields.io/pypi/v/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: Latest release

.. image:: https://img.shields.io/pypi/pyversions/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: Supported Python versions

.. image:: https://img.shields.io/pypi/implementation/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: Supported Python implementations

.. image:: https://img.shields.io/pypi/wheel/b3j0f.aop.svg
   :target: https://travis-ci.org/b3j0f/aop
   :alt: Download format

.. image:: https://travis-ci.org/b3j0f/aop.svg?branch=master
   :target: https://travis-ci.org/b3j0f/aop
   :alt: Build status

.. image:: https://coveralls.io/repos/b3j0f/aop/badge.png
   :target: https://coveralls.io/r/b3j0f/aop
   :alt: Code test coverage

.. image:: https://img.shields.io/pypi/dm/b3j0f.aop.svg
   :target: https://pypi.python.org/pypi/b3j0f.aop/
   :alt: Downloads

.. image:: https://readthedocs.org/projects/b3j0faop/badge/?version=master
   :target: https://readthedocs.org/projects/b3j0faop/?badge=master
   :alt: Documentation Status

.. image:: https://landscape.io/github/b3j0f/aop/master/landscape.svg?style=flat
   :target: https://landscape.io/github/b3j0f/aop/master
   :alt: Code Health

Links
-----

- `Homepage`_
- `PyPI`_
- `Documentation`_

Installation
------------

pip install b3j0f.aop

Features
--------

1. Free and unlimited access: no limits to sharing of ideas and knowledges with the license MIT.

2. Performance:

   - less memory consumption in using the __slots__ class property.
   - less time on (un-)weaving and advice application improvement with binary python encoding and in using constants var in code.
   - (dis/en)abling advices without remove them in using dedicated Advice class.

3. Easy to use:

   - joinpoint matching with function or regex.
   - distributed programming:

      + interception context sharing in order to ease behaviour sharing between advices.
      + uuid for advice identification in order to ease its use in a distributed context.

   - maintenable with well named variables and functions, comments and few lines.
   - extensible through pythonic code (PEP8), same logic to function code interception and concern modularisation with one module by joinpoint or advice.
   - respect of aspects vocabulary in order to ease its use among AOP users.
   - close to callable python objects in weaving all types of callable elements such as (built-in) functions, (built-in) class, (built-in) methods, callable objects, etc.
   - advices are callable objects.
   - Unit tests for all functions such as examples.

4. Benchmark:

   - speed execution

Limitations
-----------

- Do not weave advices on readonly instance methods (where class use __slots__ attribute).

Examples
--------

How to change the behaviour of min by max ?

>>> from b3j0f.aop import weave, is_intercepted
>>> double_advice = lambda joinpoint: joinpoint.proceed() * 2
>>> weave(target=min, advices=double_advice)
>>> min(6, 7)
12

How to check if a function is intercepted ?

>>> from b3j0f.aop import is_intercepted
>>> is_intercepted(min)
True

Ok, let's get back its previous behaviour ...

>>> from b3j0f.aop import unweave
>>> unweave(min)
>>> min(6, 7)
6
>>> is_intercepted(min)
False

And with an annotation ?

>>> from b3j0f.aop import weave_on
>>> weave_on(advices=double_advice)(min)
>>> min(6, 7)
12
>>> is_intercepted(min)
True
>>> unweave(min)  # do not forget to unweave if weaving is useless ;)

Enjoy ...

State of the art
----------------

Related to improving criteria points (1. Free and unlimited access, etc.), a state of the art is provided here.

+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| Library    | Url                          | License  | Execution | Use | Benchmark     | Compatibility |
+============+==============================+==========+===========+=====+===============+===============+
| b3j0f.aop  | https://github.com/b3j0f/aop | MIT      | 4/5       | 4/5 | 4/5           | 4/5 (>=2.6)   |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| pyaspects  | http://tinyurl.com/n7ccof5   | GPL 2    | 4/5       | 2/5 | 2/5           | 2/5           |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| aspects    | http://tinyurl.com/obp8t2v   | LGPL 2.1 | 2/5       | 2/5 | 2/5           | 2/5           |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| aspect     | http://tinyurl.com/lpd87bd   | BSD      | 2/5       | 1/5 | 1/5           | 1/5           |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| spring     | http://tinyurl.com/dmkpj3    | Apache   | 4/5       | 2/5 | 3/5           | 2/5           |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+
| pytilities | http://tinyurl.com/q49ulr5   | GPL 3    | 1/5       | 1/5 | 1/5           | 1/5           |
+------------+------------------------------+----------+-----------+-----+---------------+---------------+

pyaspects
#########

weaknesses
>>>>>>>>>>

- Not functional approach: Aspect class definition.
- Side effects: Not close to python API.
- Not optimized Weaving and Time execution: use classes and generic methods.
- Not maintenable: poor comments.
- open-source and use limitations: GPL 2.
- limited in weave filtering.

aspects
#######

weaknesses
>>>>>>>>>>

- open-source and use limitations: LGPL 2.1.
- more difficulties to understand code with no respect of the AOP vocabulary, packaged into one module.
- limited in weave filtering.

aspect
######

strengths
>>>>>>>>>

+ invert the AOP in decorating advices with joinpoint instead of weaving advices on joinpoint.
+ open-source and no use limitations: BSD.

weaknesses
>>>>>>>>>>

- Simple and functional approach with use of python tools.
- maintenable: commented in respect of the PEP8.
- limited in weave filtering.

spring
######

strengths
>>>>>>>>>

- a very powerful library dedicated to develop strong systems based on component based software engineering.
- unittests.
- huge community.

weaknesses
>>>>>>>>>>

- require to understand a lot of concepts and install an heavy library before doing a simple interception with AOP concerns.

pytilities
##########

strenghts
>>>>>>>>>

+ Very complex and full library for doing aspects and other things.

weaknesses
>>>>>>>>>>

- open-source and use limitations: GPL 3.
- not maintenable: missing documentations and not respect of the PEP8.
- Executon time is not optimized with several classes used with generic getters without using __slots__. The only one optimization comes from the yield which requires from users to use it in their own advices (which must be a class).

Perspectives
------------

- wait feedbacks during 6 months before passing it to a stable version.
- Cython implementation.

Donation
--------

.. image:: https://cdn.rawgit.com/gratipay/gratipay-badge/2.3.0/dist/gratipay.png
   :target: https://gratipay.com/b3j0f/
   :alt: I'm grateful for gifts, but don't have a specific funding goal.

.. _Homepage: https://github.com/b3j0f/aop
.. _Documentation: http://b3j0faop.readthedocs.org/en/develop/
.. _PyPI: https://pypi.python.org/pypi/b3j0f.aop/


