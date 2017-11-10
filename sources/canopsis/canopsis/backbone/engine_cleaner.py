#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from engine import Engine

from canopsis.backbone.event import EventUTF8Error


class EngineCleaner(Engine):

    def work(self, event):
        if not event.is_valid():
            self.drop(event, 'invalid event')
            return

        try:
            event.ensure_utf8_format()
        except EventUTF8Error:
            self.drop(event, 'utf8 error in event ')
            return

        self.logger.debug(event)
        return event

    def drop(self, event, reason=None):
        """
        Drop a single event.

        :param Event event: the event to drop
        :param str reason: the reason of the drop
        """
        self.logger.warning('Dropping event: {}'.format(
            event if event is not None else '')
        )
