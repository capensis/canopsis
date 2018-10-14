#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Action object.
"""

from __future__ import unicode_literals


class Action(object):

    """
    Representation of an action element.
    """

    # Keys as seen in db
    _ID = '_id'
    FIELD = 'field'
    REGEX = 'regex'
    PARAMETERS = 'parameters'

    def __init__(self, _id, field, regex, parameters=None, *args, **kwargs):
        """
        :param str _id: action id
        :param str field: targeted field
        :param str regex: regex matcher on the field
        :param dict parameters: variable parameters to apply
        """
        if parameters is None or not isinstance(parameters, dict):
            parameters = {}

        self._id = _id
        self.field = field
        self.regex = regex
        self.parameters = parameters

        if args not in [(), None] or kwargs not in [{}, None]:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<Action {}>'.format(self.__str__())

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
            self._ID: self._id,
            self.FIELD: self.field,
            self.REGEX: self.regex,
            self.PARAMETERS: self.parameters,
        }

        return dictionnary
