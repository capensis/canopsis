#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
TicketApi object.
"""

from __future__ import unicode_literals


class TicketApi(object):

    """
    Representation of a ticketApiConfig element.
    """

    # Keys as seen in db
    _ID = '_id'
    TYPE = 'type'
    FIELDS = 'fields'
    REGEX = 'regex'
    API = 'api'
    PARAMETERS = 'parameters'
    PAYLOAD = 'payload'
    API_URL = 'base_url'
    API_USER = 'username'
    API_PWD = 'password'

    def __init__(self, _id, type_, fields, regex,
                 api=None, parameters=None, payload=None,
                 *args, **kwargs):
        """
        :param str _id: ticketapi id
        :param str type_: type of api (snom, ...)
        :param list fields: targeted fields
        :param str regex: regex matcher on the field
        :param dict api: api base configuration
        :param dict parameters: other specific parameters (for ex, paths, static values)
        :param str payload: description, value by value, of the json sended to the api
        """
        if fields is None or not isinstance(fields, list):
            fields = []
        if api is None or not isinstance(api, dict):
            api = {}
        if parameters is None or not isinstance(parameters, dict):
            parameters = {}
        if payload is None or not isinstance(payload, dict):
            payload = {}

        for key in [self.API_URL, self.API_USER, self.API_PWD]:
            if key not in api:
                print('Missing {} in api dict: {}'.format(key, api))

        self._id = _id
        self.type_ = type_
        self.fields = fields
        self.regex = regex
        self.api = api
        self.parameters = parameters
        self.payload = payload

        if args not in [(), None] or kwargs not in [{}, None]:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<TicketApiConfig {}>'.format(self.__str__())

    @staticmethod
    def convert_keys(ticketapi_dict):
        """
        Convert keys from mongo ticketapi config dict, to TicketApi
        attribute names.

        :param dict ticketapi_dict: a raw ticketapi config dict from mongo
        :rtype: dict
        """
        new_ticketapi_dict = ticketapi_dict.copy()

        if TicketApi.TYPE in new_ticketapi_dict:
            new_ticketapi_dict['type_'] = new_ticketapi_dict[TicketApi.TYPE]
            del new_ticketapi_dict[TicketApi.TYPE]

        return new_ticketapi_dict

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
            self._ID: self._id,
            self.TYPE: self.type_,
            self.FIELDS: self.fields,
            self.REGEX: self.regex,
            self.API: self.api,
            self.PARAMETERS: self.parameters,
            self.PAYLOAD: self.payload,
        }

        return dictionnary
