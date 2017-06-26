from unittest import TestCase, main
from mock import Mock

from canopsis.alerts.manager import Alerts
from canopsis.alerts.selector import StateCalculator, Selector
from canopsis.context.manager import Context
from canopsis.middleware.core import Middleware

from canopsis.old.record import Record
from canopsis.old.cfilter import Filter


class SelectorTest(TestCase):
    def setUp(self):
        alarm_storage = Middleware.get_middleware_by_uri(
            'storage-default-testalarm://'
        )

        context_storage = Middleware.get_middleware_by_uri(
            'storage-default-testcontexct://'
        )

        self.alert_manager = Alerts()
        self.ctx_manager = Context()

        self.alert_manager[Alerts.ALARM_STORAGE] = alarm_storage
        self.ctx_manager[Context.CTX_STORAGE] = context_storage

        alarm_id_1 = '/fake/alarm1/id'
        alarm_id_2 = '/fake/alarm2/id'
        alarm_id_3 = '/fake/alarm3/id'

        self.entity_id_1 = '/component/collectd/pbehavior/test1/'
        self.entity_id_2 = '/component/collectd/pbehavior/test2/'
        self.entity_id_3 = '/component/collectd/pbehavior/test3/'

        self.alarms = [{
                '_id': alarm_id_1,
                'entity_id': self.entity_id_1,
                'connector': 'canopsis-test-connector',
                'connector_name': 'canopsis-test',
                'component': 'canopsis-demo-sakura-debian7-pbehaviors',
                'timestamp': 0,
                'state': 0,
                'status': 3,
                'event_type': 'check',
                'source_type': 'component',
                'ack': {
                    'isAck': False
                }
            }, {
                '_id': alarm_id_2,
                'entity_id': self.entity_id_1,
                'connector': 'canopsis-test-connector',
                'connector_name': 'canopsis-test',
                'component': 'canopsis-demo-sakura-debian7-pbehaviors',
                'timestamp': 0,
                'state': 1,
                'status': 3,
                'event_type': 'check',
                'source_type': 'component'
            }, {
                '_id': alarm_id_3,
                'entity_id': self.entity_id_1,
                'connector': 'canopsis-test-connector',
                'connector_name': 'canopsis-test',
                'component': 'canopsis-demo-sakura-debian7-pbehaviors',
                'timestamp': 0,
                'state': 1,
                'status': 3,
                'event_type': 'check',
                'source_type': 'component'
            }
        ]

        self.entities = [{
            '_id': self.entity_id_1,
            'name': 'engine-test1',
            'type': 'metric-test',
            'connector': 'canopsis-test-connector',
            'connector_name': 'canopsis-test',
        }, {
            '_id': self.entity_id_2,
            'name': 'big-engine-test2',
            'type': 'metric-test',
            'connector': 'canopsis-test-connector',
            'connector_name': 'canopsis-test',
        }, {
            '_id': self.entity_id_3,
            'name': 'test_context3',
            'type': 'resource-test',
            'connector': 'nagios-test-connector',
            'connector_name': 'nagios-test',
        }]

        self.alert_manager[Alerts.ALARM_STORAGE].put_elements(self.alarms)
        self.ctx_manager[Context.CTX_STORAGE].put_elements(self.entities)

    def tearDown(self):
        self.alert_manager[Alerts.ALARM_STORAGE].remove_elements()
        self.ctx_manager[Context.CTX_STORAGE].remove_elements()

    def test_get_states(self):
        alarm_storage = self.alert_manager[Alerts.ALARM_STORAGE]

        self.cstate = StateCalculator(
            alarm_storage, Mock()
        )

        mfilter = {
            "connector": {"$eq": "canopsis-test-connector"}
        }

        entities = self.ctx_manager[Context.CTX_STORAGE].get_elements(
            query=mfilter)

        mfilter = {'entity_id': {'$in': [e['_id'] for e in entities]}}

        mfilter = Filter().make_filter(mfilter)

        states = self.cstate.get_states(mfilter)

        self.assertDictEqual(states, {0: 1, 1: 2})

        worst_state = self.cstate.get_worst_state(states)

        self.assertEqual(1, worst_state)

        wstate_for_ack = self.cstate.get_worst_state_for_ack(mfilter)

        self.assertEqual(0, wstate_for_ack)

        infobagot = self.cstate.get_infobagot(mfilter)

        self.assertEqual(1, infobagot)

    def test_get_alert(self):
        storage = self.alert_manager[Alerts.ALARM_STORAGE]
        event = {
            'mfilter': {
                "connector": {"$eq": "canopsis-test-connector"}
            }
        }

        selector = Selector(
            storage=storage,
            record=Record(event)
        )

        alert = selector.alert()

        self.assertIsNotNone(alert)
        self.assertIn('state', alert)
        self.assertEqual(alert['state'], 0)


if __name__ == '__main__':
    main()
