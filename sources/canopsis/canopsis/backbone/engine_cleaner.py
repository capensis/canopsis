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
        print(event)
        return event

    def drop(self, event, reason=None):
        self.logger.warning('Dropping event: {}'.format(
            event if event is not None else '')
        )

