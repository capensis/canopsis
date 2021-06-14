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
    IMPACTS = 'impact'
    MEASUREMENTS = 'measurements'
    INFOS = 'infos'
    ENABLED = 'enabled'
    ENABLE_HISTORY = 'enable_history'
    LAST_STATE_CHANGE = 'last_state_change'

    def __init__(self, _id, name, type_,
                 depends=None, impacts=None, measurements=None, infos=None,
                 enabled=True, enable_history=None, last_state_change=None,
                 *args, **kwargs):
        """
        :param str _id: entity id
        :param str name: entity name
        :param str type_: entity type ()
        :param list depends: dependency list of this entity
        :param list impacts: impact list of this entity
        :param list measurements: measurements list of this entity
        :param dic infos: extra informations
        :param bool enabled: his this entity enabled ?
        :param list enable_history: list of timestamp when the entity has been activated
        :param int last_state_change: timestamp of the last time the entity's
        state changed
        """
        if depends is None or not isinstance(depends, list):
            depends = []
        if impacts is None or not isinstance(impacts, list):
            impacts = []
        if enable_history is None or not isinstance(enable_history, list):
            enable_history = []
        if not isinstance(enabled, bool):
            enabled = True
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
        self.last_state_change = last_state_change

        self.enable_history = enable_history  # before enabled !!
        self.enabled = enabled

        if args not in [(), None] or kwargs not in [{}, None]:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<Entity {}>'.format(self.__str__())

    @staticmethod
    def convert_keys(entity_dict):
        """
        Convert keys from mongo entity dict, to Entity attribute names.

        :param dict entity_dict: a raw entity dict from mongo
        :rtype: dict
        """
        new_entity_dict = entity_dict.copy()
        if Entity.TYPE in new_entity_dict:
            new_entity_dict['type_'] = new_entity_dict[Entity.TYPE]
            del new_entity_dict[Entity.TYPE]

        if Entity.IMPACTS in new_entity_dict:
            new_entity_dict['impacts'] = new_entity_dict[Entity.IMPACTS]
            del new_entity_dict[Entity.IMPACTS]

        return new_entity_dict

    @property
    def enabled(self):
        return self._enabled

    @enabled.setter
    def enabled(self, value):
        self._enabled = value
        if self._enabled:
            timestamp = int(time.time())
            self.enable_history = self.enable_history + [timestamp]

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
            self._ID: self._id,
            self.TYPE: self.type_,
            self.NAME: self.name,
            self.DEPENDS: self.depends,
            self.IMPACTS: self.impacts,
            self.MEASUREMENTS: self.measurements,
            self.INFOS: self.infos,
            self.ENABLED: self.enabled,
            self.ENABLE_HISTORY: self.enable_history,
            self.LAST_STATE_CHANGE: self.last_state_change
        }

        return dictionnary
