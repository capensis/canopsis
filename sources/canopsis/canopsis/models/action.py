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
    TYPE = 'type'
    HOOK = 'hook'
    FIELDS = 'fields'
    REGEX = 'regex'
    PARAMETERS = 'parameters'

    def __init__(self, _id, type_, hook, fields, regex, parameters=None,
                 *args, **kwargs):
        """
        :param str _id: action id
        :param str type_: type of the action (pbehavior, ...)
        :param list fields: targeted fields
        :param str regex: regex matcher on the fields
        :param dict parameters: variable parameters to apply
        """
        if fields is None or not isinstance(fields, list):
            fields = []
        if parameters is None or not isinstance(parameters, dict):
            parameters = {}
        if not isinstance(hook, dict):
            hook = None

        self._id = _id
        self.type_ = type_
        self.hook = hook
        self.fields = fields
        self.regex = regex
        self.parameters = parameters

        if args not in [(), None] or kwargs not in [{}, None]:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<Action {}>'.format(self.__str__())

    @staticmethod
    def convert_keys(action_dict):
        """
        Convert keys from mongo action dict, to Action attribute names.

        :param dict action_dict: a raw action dict from mongo
        :rtype: dict
        """
        new_action_dict = action_dict.copy()

        if Action.TYPE in new_action_dict:
            new_action_dict['type_'] = new_action_dict[Action.TYPE]
            del new_action_dict[Action.TYPE]

        return new_action_dict

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
            self._ID: self._id,
            self.TYPE: self.type_,
            self.HOOK: self.hook,
            self.FIELDS: self.fields,
            self.REGEX: self.regex,
            self.PARAMETERS: self.parameters,
        }

        return dictionnary
