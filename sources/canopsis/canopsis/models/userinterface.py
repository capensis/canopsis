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
    APP_TITLE = 'user_app_title'
    FOOTER = 'user_footer'
    LOGO = 'user_logo'

    def __init__(self, _id, app_title=None, footer=None, logo=None, *args, **kwargs):

        self._id = _id
        self.user_app_title = app_title
        self.user_footer = footer
        self.user_logo = logo

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
            self._ID: self._id,
            self.APP_TITLE: self.user_app_title,
            self.FOOTER: self.user_footer,
            self.LOGO: self.user_logo,
        }

        return dictionnary
