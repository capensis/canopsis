#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
UserInterface object.
"""

from __future__ import unicode_literals


class UserInterface(object):

    """
    Representation of an user interface element.
    """

    # Keys as seen in db
    _ID = '_id'
    APP_TITLE = 'app_title'
    FOOTER = 'footer'
    LOGIN_PAGE_DESCRIPTION = 'login_page_description'
    LOGO = 'logo'
    LANGUAGE = 'language'
    POPUP_TIMEOUT = 'popup_timeout'
    AllowChangeSeverityToInfo = 'allow_change_severity_to_info'

    def __init__(self, _id, app_title=None, footer=None, login_page_description=None, logo=None, language=None,
                 popup_timeout=None, allow_change_severity_to_info=False, *args, **kwargs):

        self._id = _id
        self.app_title = app_title
        self.footer = footer
        self.login_page_description = login_page_description
        self.logo = logo
        self.language = language
        self.popup_timeout = popup_timeout
        self.allow_change_severity_to_info = allow_change_severity_to_info

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
            self.LOGIN_PAGE_DESCRIPTION: self.login_page_description,
            self.LOGO: self.logo,
            self.LANGUAGE: self.language,
            self.POPUP_TIMEOUT: self.popup_timeout,
            self.AllowChangeSeverityToInfo: self.allow_change_severity_to_info
        }

        return dictionnary
