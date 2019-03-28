#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
UserInterface object.
"""

from __future__ import unicode_literals


class UserInterface(object):

    """
    Representation of a ticketApiConfig element.
    """

    # Keys as seen in db
    _ID = '_id'
    APP_TITLE = 'app_title'
    FOOTER = 'footer'
    LOGO = 'logo'

    def __init__(self, _id, app_title=None, footer=None, logo=None, *args, **kwargs):

        self._id = _id
        self.app_title = app_title
        self.footer = footer
        self.logo = logo

        if args not in [(), None] or kwargs not in [{}, None]:
            print('Ignored values on creation: {} // {}'.format(args, kwargs))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<UserInterface {}>'.format(self.__str__())

    def to_dict(self):
        """
        Give a dict representation of the object.

        :rtype: dict
        """
        dictionnary = {
            self.APP_TITLE: self.app_title,
            self.FOOTER: self.footer,
            self.LOGO: self.logo,
        }

        return dictionnary
