#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
PBehavior object.
"""

from __future__ import unicode_literals

import json
import time
from six import string_types

DEFAULT_CONNECTOR_VALUE = 'canopsis'
DEFAULT_CONNECTOR_NAME_VALUE = 'canopsis'


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
    TYPE = 'type_'
    REASON = 'reason'

    def __init__(self, _id, name, filter_, tstart, tstop, rrule, author,
                 connector=DEFAULT_CONNECTOR_VALUE,
                 connector_name=DEFAULT_CONNECTOR_NAME_VALUE,
                 comments=None, eids=None, enabled=True,
                 type_=None, reason=None,
                 *args, **kwargs):
        """
        :param str _id: pbehavior id
        :param str name: pbehavior name
        :param dict filter_: matched entites, as mongo filter dicts
        :param int tstart: starting timestamp of the pbehavior
        :param int tstop: stopping timestamp of the pbehavior
        :param str rrule: reccurent rule that is compliant with rrule spec
        :param str connector: connector that has generated the pbehavior
        :param str connector_name: connector_name that has generated the pbehavior
        :param str author: creator's name
        :param list comments: list of comments
        :param list eids: impacted entity ids
        :param bool enabled: his this pbehavior actived ?
        :param str type_: particuliar type (editable trough ui ; pause...)
        :param str reason: explanation on pbehavior creation
        """
        if filter_ is None or not isinstance(filter_, dict):
            filter_ = {}
        if comments is None or not isinstance(comments, list):
            comments = []
        if eids is None or not isinstance(eids, list):
            eids = []
        if not isinstance(enabled, bool):
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
        new_pbh_dict = pbehavior_dict.copy()

        if isinstance(new_pbh_dict[PBehavior.FILTER], string_types):
            filter_ = new_pbh_dict[PBehavior.FILTER]
            new_pbh_dict[PBehavior.FILTER] = json.loads(filter_)

        if PBehavior.FILTER in new_pbh_dict:
            new_pbh_dict['filter_'] = new_pbh_dict[PBehavior.FILTER]
            del new_pbh_dict[PBehavior.FILTER]

        # if PBehavior.TYPE in new_pbh_dict:
        #    new_pbh_dict['type_'] = new_pbh_dict[PBehavior.TYPE]
        #    del new_pbh_dict[PBehavior.TYPE]

        return new_pbh_dict

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dico = {
            self._ID: self._id,
            self.NAME: self.name,
            self.FILTER: json.dumps(self.filter_),
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
