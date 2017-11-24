from unittest import main, TestCase

from canopsis.event import get_routingkey


class TestEventFunctions(TestCase):

    refrk = 'keyboard.cherry.keypress.resource.mx.brown'
    refrk_component = 'keyboard.cherry.keypress.component.mx'

    @classmethod
    def get_ref_event(cls):
        return {
            'connector': 'keyboard',
            'connector_name': 'cherry',
            'component': 'mx',
            'resource': 'brown',
            'source_type': 'resource',
            'event_type': 'keypress',
        }

    def test_get_routingkey(self):
        event = self.get_ref_event()

        rk = get_routingkey(event)
        self.assertEqual(rk, self.refrk)

        event = self.get_ref_event()
        event['source_type'] = 'caps'

        rk = get_routingkey(event)
        self.assertEqual(rk, self.refrk)

    def test_get_routingkey_raise(self):
        event = {}

        with self.assertRaises(KeyError):
            get_routingkey(event)

    def test_get_routingkey_overrides_source_type(self):
        event = self.get_ref_event()
        del event['source_type']

        rk = get_routingkey(event)
        self.assertEqual(rk, self.refrk)
        self.assertEqual(event['source_type'], 'resource')

        del event['resource']
        rk = get_routingkey(event)
        self.assertEqual(rk, self.refrk_component)
        self.assertNotIn('resource', event)


if __name__ == '__main__':
    main()
