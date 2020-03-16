from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.webcore.services.metaalarmrule import RouteHandlerMetaAlarmRule
import unittest
from canopsis.common import root_path
import xmlrunner


class TestMetaAlarmRuleWebservice(unittest.TestCase):

    INVALID_RULES = [{
        'name': 'test_bad_ma',
        'patterns': 'bad patterns',
        'config': 'bad config',
        'rule_type': None,
    }, {
        'name': 'test_ma',
        'config': '{"time_interval": 3}',
        'rule_type': 'attribute',
        'patterns': None,
    }, {
        'name': 'test_ma_1',
        'config': '{"attribute_patterns": [{"v": {"state": {"val": 3} } }]}',
        'rule_type': 'time',
        'patterns': None,
    }]

    VALID_RULES = [{
        'name': 'test_valid_ma1',
        'patterns': '{"nokey": "novalue"}',
        'config': '{"time_interval": 3}',
        'rule_type': 'time',
    }, {
        'name': 'test_valid_ma2',
        'patterns': None,
        'config': '{"attribute_patterns": [{"v": {"state": {"val": 3} } }]}',
        'rule_type': 'attribute',
    }, {
        'name': 'test_valid_ma3',
        'patterns': None,
        'config': '{"attribute_patterns": [{"v": {"state": {"val": 3} } }], "threshold_count": 3, "time_interval": 10}',
        'rule_type': 'complex',
    }]

    @classmethod
    def setUpClass(cls):
        ma_rule_manager = MetaAlarmRuleManager(
            *MetaAlarmRuleManager.provide_default_basics())
        cls.rh_ma_rule = RouteHandlerMetaAlarmRule(ma_rule_manager)

    def test_create_bad_rule(self):
        for rule in self.INVALID_RULES:
            with self.assertRaises(ValueError):
                self.rh_ma_rule.create(**rule)

    def test_create_rule(self):
        for rule in self.VALID_RULES:
            self.rh_ma_rule.create(**rule)

    def test_read_rule(self):
        for rule in self.VALID_RULES:
            ma_rule_id = self.rh_ma_rule.create(**rule)
            self.assertIsInstance(ma_rule_id, str)

            ma_rule = self.rh_ma_rule.read(ma_rule_id)

            self.assertIsInstance(ma_rule, dict)
            self.assertEquals(ma_rule.get('name'), rule.get('name'))

            self.assertIsInstance(ma_rule.get("config"), dict)

    def test_delete_rule(self):
        for rule in self.VALID_RULES:
            ma_rule_id = self.rh_ma_rule.create(**rule)

            ma_rule = self.rh_ma_rule.read(ma_rule_id)

            self.assertEquals(ma_rule.get('name'), rule.get('name'))

            delres = self.rh_ma_rule.delete(ma_rule_id)
            self.assertEquals(delres.get('deletedCount'), 1)

            ma_rule = self.rh_ma_rule.read(ma_rule_id)
            self.assertIsNone(ma_rule)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
