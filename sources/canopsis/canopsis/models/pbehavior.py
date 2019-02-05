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
    SOURCE = 'source'
    EXDATE = 'exdate'
    TIMEZONE = 'timezone'

    def __init__(self, _id, name, filter_, tstart, tstop, rrule, author,
                 connector=DEFAULT_CONNECTOR_VALUE,
                 connector_name=DEFAULT_CONNECTOR_NAME_VALUE,
                 comments=None, eids=None, type_=None, reason=None,
                 enabled=True, source=None, exdate=None,
                 timezone="UTC", *args, **kwargs):
        """
        :param str _id: pbehavior id
        :param str name: pbehavior name
        :param dict filter_: matched entites, as mongo filter dicts
        :param int tstart: starting timestamp of the pbehavior
        :param int tstop: stopping timestamp of the pbehavior. Can be None if this pbehavior does not expire
        :param str rrule: reccurent rule that is compliant with rrule spec
        :param str connector: connector that has generated the pbehavior
        :param str connector_name: connector_name that has generated the pbehavior
        :param str author: creator's name
        :param list comments: list of comments
        :param list eids: impacted entity ids
        :param str type_: particuliar type (editable trough ui ; pause...)
        :param str reason: explanation on pbehavior creation
        :param bool enabled: allow this pbehavior to be used. This is NOT the same as the is_active property.
        :param str source: if None, pbehavior was created from canopsis. if anything else, it was created from an external data like an event from Nagios or so.
        :param list int exdate: a list of exclusion date as a timestamp
        :param str timezone: a timezone name
        """
        if filter_ is None:
            filter_ = {}
        elif not isinstance(filter_, dict):
            raise TypeError('filter_ must be a dict, got {}'.format(type(filter_)))

        if comments is None:
            comments = []
        elif not isinstance(comments, list):
            raise TypeError('comments must be a list, got {}'.format(type(comments)))

        if eids is None:
            eids = []
        elif not isinstance(eids, list):
            raise TypeError('eids must be a list, got {}'.format(type(comments)))

        if type_ is None:
            type_ = ''
        elif not isinstance(type_, string_types):
            raise TypeError('type_ must be a string_type, got {}'.format(type(type_)))

        if reason is None:
            reason = ''
        elif not isinstance(reason, string_types):
            raise TypeError('reason must be a string_type, got {}'.format(type(reason)))

        if source is not None and not isinstance(source, string_types):
            raise TypeError('source must be None or a string, got {}'.format(type(source)))

        if exdate is None:
            exdate = []
        elif not isinstance(exdate, list):
            raise TypeError('exdate must be None or a list, got {}'.format(type(source)))
        else:
            for date in exdate:
                if not isinstance(date, int):
                    raise TypeError('The date inside exdate must be an int, got {}'.format(type(source)))

        self._id = _id
        self.enabled = enabled
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
        self.type_ = type_
        self.reason = reason
        self.source = source
        self.exdate = exdate
        self.timezone = timezone

        if args not in [(), None] or kwargs not in [{}, None]:
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

    @property
    def is_active(self):
        """
        :returns: is this pbehavior currently active.
        :rtype: bool
        """
        now = int(time.time())

        if self.tstart <= now and self.tstop is None:
            return True

        return self.tstart <= now <= self.tstop

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
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
            self.REASON: self.reason,
            self.SOURCE: self.source,
            self.EXDATE: self.exdate,
            self.TIMEZONE: self.timezone
        }

        return dictionnary
