from engine import Engine
from canopsis.backbone.event import EventUTF8Error

ESSENTIAL_PARAMETERS = [
    'connector',
    'connector_name',
    'event_type',
    'source_type',
    'component'
]


class EngineCleaner(Engine):

    def work(self, event, message):
        if not event.is_valid():
            self.drop(event, 'invalid event')
            return
        try:
            event.ensure_utf8_format()
        except EventUTF8Error:
            self.drop(event, 'utf8 error in event ')
            return
        return event

    def drop(self, event, reason=None):
        self.logger.warning('droping event: {}'.format(
            event if event is not None else '')
        )

