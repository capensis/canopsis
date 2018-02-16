#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
PBehavior object.
"""

from __future__ import unicode_literals

import time
from six import string_types


class PBehavior(object):

    """
    Representation of a pbehavior element.
    """

    _ID = '_id'
    NAME = 'name'
    FILTER = 'filter'
    COMMENTS = 'comments'
    TSTART = 'tstart'
    TSTOP = 'tstop'
    RRULE = 'rrule'
    ENABLED = 'enabled'
    EIDS = 'eids'
    CONNECTOR = 'connector'
    CONNECTOR_NAME = 'connector_name'
    AUTHOR = 'author'
    TYPE = 'type'
    REASON = 'reason'

    def __init__(self, _id, name, filter_, tstart, tstop, rrule, author,
                 connector, connector_name,
                 comments=None, eids=None, enabled=None,
                 type_=None, reason=None,
                 *args, **kwargs):
        if comments is None or not isinstance(comments, list):
            comments = []
        if eids is None or not isinstance(eids, list):
            eids = []
        if enabled is None or not isinstance(enabled, bool):
            now = int(time.now())
            enabled = tstart <= now <= tstop
        if type_ is None or not isinstance(type_, string_types):
            type_ = ''
        if reason is None or not isinstance(reason, string_types):
            reason = ''

        self._id = _id
        self.name = name
        self.filter_ = filter_
        self.tstart = tstart
        self.tstop = tstop
        self.rrule = rrule
        self.connector = connector
        self.connector_name = connector_name
        self.author = author
        self.comments = comments
        self.eids = eids
        self.enabled = enabled
        self.type_ = type_
        self.reason = reason

        if args is not None or kwargs is not None:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self.name)

    def __repr__(self):
        return '<PBehavior {}>'.format(self.__str__())

    @staticmethod
    def convert_keys(pbehavior_dict):
        """
        Convert keys from mongo pbehavior dict, to PBehavior attribute names.

        :param dict pbehavior_dict: a raw pbehavior dict from mongo
        :rtype: dict
        """
        if PBehavior.FILTER in pbehavior_dict:
            pbehavior_dict['filter_'] = pbehavior_dict[PBehavior.FILTER]
            del pbehavior_dict[PBehavior.FILTER]

        if PBehavior.TYPE in pbehavior_dict:
            pbehavior_dict['type_'] = pbehavior_dict[PBehavior.TYPE]
            del pbehavior_dict[PBehavior.TYPE]

        return pbehavior_dict

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dico = {
            self._ID: self._id,
            self.NAME: self.name,
            self.FILTER: self.filter_,
            self.TSTART: self.tstart,
            self.TSTOP: self.tstop,
            self.RRULE: self.rrule,
            self.CONNECTOR: self.connector,
            self.CONNECTOR_NAME: self.connector_name,
            self.AUTHOR: self.author,
            self.COMMENTS: self.comments,
            self.EIDS: self.eids,
            self.ENABLED: self.enabled,
            self.TYPE: self.type_,
            self.REASON: self.reason
        }

        return dico
