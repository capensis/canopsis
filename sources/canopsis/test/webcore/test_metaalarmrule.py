from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.webcore.services.metaalarmrule import RouteHandlerMetaAlarmRule
import unittest
from canopsis.common import root_path
import xmlrunner


class TestMetaAlarmRuleWebservice(unittest.TestCase):

    INVALID_RULE = {
        'name': 'test_bad_pb',
        'filter': 'bad filter',
        'enabled': True,
        'rule_type': None,
    }

    VALID_RULE = {
        'name': 'test_pb',
        'filter': '{"nokey": "novalue"}',
        'enabled': True,
        'rule_type': 'relation',
    }

    @classmethod
    def setUpClass(cls):
        ma_rule_manager = MetaAlarmRuleManager(*MetaAlarmRuleManager.provide_default_basics())
        cls.rh_ma_rule = RouteHandlerMetaAlarmRule(ma_rule_manager)

    def test_create_bad_rule(self):
        with self.assertRaises(ValueError):
            self.rh_ma_rule.create(**self.INVALID_RULE)

    def test_create_rule(self):
        self.rh_ma_rule.create(**self.VALID_RULE)

    def test_read_rule(self):
        ma_rule_id = self.rh_ma_rule.create(**self.VALID_RULE)
        self.assertIsInstance(ma_rule_id, str)

        ma_rule = self.rh_ma_rule.read(ma_rule_id)

        self.assertIsInstance(ma_rule, dict)
        self.assertEquals(ma_rule.get('name'), self.VALID_RULE.get('name'))

    def test_delete_rule(self):
        ma_rule_id = self.rh_ma_rule.create(**self.VALID_RULE)

        ma_rule = self.rh_ma_rule.read(ma_rule_id)

        self.assertEquals(ma_rule.get('name'), self.VALID_RULE.get('name'))

        delres = self.rh_ma_rule.delete(ma_rule_id)
        self.assertEquals(delres.get('deletedCount'), 1)

        ma_rule = self.rh_ma_rule.read(ma_rule_id)
        self.assertIsNone(ma_rule)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
