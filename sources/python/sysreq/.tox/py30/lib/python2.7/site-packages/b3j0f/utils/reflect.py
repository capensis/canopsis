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

"""Python reflection tools."""

from __future__ import unicode_literals, absolute_import

from .iterable import ensureiterable

from inspect import isclass, isroutine, ismethod, getmodule, getmro

from six import PY2, get_method_self, get_method_function
from six.moves import range

__all__ = ['isoldstyle', 'base_elts', 'find_embedding', 'is_inherited']


def isoldstyle(cls):
    """Return True if cls is an old style class (does not inherits from object
    in python2."""

    return getmro(cls)[-1] is not object


def base_elts(elt, cls=None, depth=None):
    """Get bases elements of the input elt.

    - If elt is an instance, get class and all base classes.
    - If elt is a method, get all base methods.
    - If elt is a class, get all base classes.
    - In other case, get an empty list.

    :param elt: supposed inherited elt.
    :param cls: cls from where find attributes equal to elt. If None,
        it is found as much as possible. Required in python3 for function
        classes.
    :type cls: type or list
    :param int depth: search depth. If None (default), depth is maximal.
    :return: elt bases elements. if elt has not base elements, result is empty.
    :rtype: list
    """

    result = []

    elt_name = getattr(elt, '__name__', None)

    if elt_name is not None:

        cls = [] if cls is None else ensureiterable(cls)

        elt_is_class = False

        # if cls is None and elt is routine, it is possible to find the cls
        if not cls and isroutine(elt):

            if hasattr(elt, '__self__'):  # from the instance

                instance = get_method_self(elt)  # get instance

                if instance is None and PY2:  # get base im_class if PY2
                    cls = list(elt.im_class.__bases__)

                else:  # use instance class
                    cls = [instance.__class__]

        # cls is elt if elt is a class
        elif isclass(elt):
            elt_is_class = True
            cls = list(elt.__bases__)

        if cls:  # if cls is not empty, find all base classes

            index_of_found_classes = 0  # get last visited class index
            visited_classes = set(cls)  # cache for visited classes
            len_classes = len(cls)

            if depth is None:  # if depth is None, get maximal value
                depth = -1  # set negative value

            while depth != 0 and index_of_found_classes != len_classes:
                len_classes = len(cls)

                for index in range(index_of_found_classes, len_classes):
                    _cls = cls[index]

                    for base_cls in _cls.__bases__:
                        if base_cls in visited_classes:
                            continue

                        else:
                            visited_classes.add(base_cls)
                            cls.append(base_cls)
                index_of_found_classes = len_classes
                depth -= 1

            if elt_is_class:
                # if cls is elt, result is classes minus first class
                result = cls

            elif isroutine(elt):

                # get an elt to compare with found element
                if ismethod(elt):
                    elt_to_compare = get_method_function(elt)
                else:
                    elt_to_compare = elt

                for _cls in cls:  # for all classes
                    # get possible base elt
                    b_elt = getattr(_cls, elt_name, None)

                    if b_elt is not None:
                        # compare funcs
                        if ismethod(b_elt):
                            bec = get_method_function(b_elt)
                        else:
                            bec = b_elt
                        # if matching, add to result
                        if bec is elt_to_compare:
                            result.append(b_elt)

    return result


def is_inherited(elt, cls=None):
    """True iif elt is inherited in a base class.

    :param elt: elt to check such as an inherited element.
    :param type cls: base cls where find the base elt.
    :return: true if elt is an inherited element.
    :rtype: bool
    """

    return base_elts(elt, cls=cls, depth=1)


def find_embedding(elt, embedding=None):
    """Try to get elt embedding elements.

    :param embedding: embedding element. Must have a module.

    :return: a list of [module [,class]*] embedding elements which define elt.
    :rtype: list
    """

    result = []  # result is empty in the worst case

    # start to get module
    module = getmodule(elt)

    if module is not None:  # if module exists

        visited = set()  # cache to avoid to visit twice same element

        if embedding is None:
            embedding = module

        # list of compounds elements which construct the path to elt
        compounds = [embedding]

        while compounds:  # while compounds elements exist
            # get last compound
            last_embedding = compounds[-1]
            # stop to iterate on compounds when last embedding is elt
            if last_embedding == elt:
                result = compounds  # result is compounds
                break

            else:
                # search among embedded elements
                for name in dir(last_embedding):
                    # get embedded element
                    embedded = getattr(last_embedding, name)

                    try:  # check if embedded has already been visited
                        if embedded not in visited:
                            visited.add(embedded)  # set it as visited

                        else:
                            continue

                    except TypeError:
                        pass

                    else:

                        # get embedded module
                        embedded_module = getmodule(embedded)
                        # and compare it with elt module
                        if embedded_module is module:
                            # add embedded to compounds
                            compounds.append(embedded)
                            # end the second loop
                            break

                else:
                    # remove last element if no coumpound element is found
                    compounds.pop(-1)

    return result
