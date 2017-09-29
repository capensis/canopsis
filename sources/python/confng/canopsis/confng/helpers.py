#!/usr/bin/env python
# -*- coding: utf-8 -*-
from six import string_types


def cfg_to_array(value, sep=','):
    """
    Convert a configuration line to an array.

    :param str value: a value to parse
    :param str sep: separator
    :rtype: list
    """
    values = value.split(',')

    f_values = []

    for val in values:
        val = val.strip()
        if len(val) > 0:
            f_values.append(val)

    return f_values


def cfg_to_bool(value):
    """
    Convert a configuration line to a bool.

    :param str value: a value to parse
    :raises: ValueError
    :rtype: bool
    """
    if isinstance(value, bool):
        return value
    elif isinstance(value, string_types):
        value_cap = value.capitalize()
        if value_cap == 'True':
            return True
        elif value_cap == 'False':
            return False

    raise ValueError('Cannot parse to boolean: {}'.format(value))
