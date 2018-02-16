#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Entity object.
"""

from __future__ import unicode_literals

import time


class Entity(object):

    """
    Representation of an entity element.
    """

    # Keys as seen in db
    _ID = '_id'
    NAME = 'name'
    TYPE = 'type'
    DEPENDS = 'depends'
    IMPACTS = 'impacts'
    MEASUREMENTS = 'measurements'
    INFOS = 'infos'
    ENABLED = 'enabled'
    ENABLED_HISTORY = 'enabled_history'

    def __init__(self, _id, name, type_,
                 depends=None, impacts=None, measurements=None, infos=None,
                 enabled=None, enabled_history=None,
                 *args, **kwargs):
        if depends is None or not isinstance(depends, list):
            depends = []
        if impacts is None or not isinstance(impacts, list):
            impacts = []
        if enabled_history is None or not isinstance(enabled_history, list):
            enabled_history = []
        if enabled is None or not isinstance(enabled, bool):
            enabled = False
        if infos is None or not isinstance(infos, dict):
            infos = {}
        if measurements is None or not isinstance(measurements, dict):
            measurements = {}

        self._id = _id
        self.name = name
        self.type_ = type_
        self.depends = depends
        self.impacts = impacts
        self.measurements = measurements
        self.infos = infos
        self.enabled = enabled
        self.enabled_history = enabled_history

        if args is not None or kwargs is not None:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

        # Automatically enable the entity
        self._enable()

    def __str__(self):
        return '{}-{}'.format(self._id)

    def __repr__(self):
        return '<Entity {}>'.format(self.__str__())

    @staticmethod
    def convert_keys(entity_dict):
        """
        Convert keys from mongo entity dict, to Entity attribute names.

        :param dict entity_dict: a raw entity dict from mongo
        :rtype: dict
        """
        if Entity.TYPE in entity_dict:
            entity_dict['type_'] = entity_dict[Entity.TYPE]
            del entity_dict[Entity.TYPE]

        return entity_dict

    def _enable(self):
        """
        Enable the entity.
        """
        timestamp = int(time.time())

        self.enabled = True
        self.enable_history = self.enable_history + timestamp

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dico = {
            self._ID: self._id,
            self.TYPE: self.type_,
            self.NAME: self.name,
            self.DEPENDS: self.depends,
            self.IMPACT: self.impact,
            self.MEASUREMENTS: self.measurements,
            self.INFOS: self.infos,
            self.ENABLED: self.enabled,
            self.ENABLED_HISTORY: self.enabled_history
        }

        return dico
