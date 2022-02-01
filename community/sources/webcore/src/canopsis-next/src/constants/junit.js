import { COLORS } from '@/config';

export const TEST_SUITE_STATUSES = {
  passed: 0,
  skipped: 1,
  error: 2,
  failed: 3,

  /**
   * Special frontend value. We don't have this value on the backend side.
   */
  total: 4,
};

export const TEST_SUITE_COLORS = {
  [TEST_SUITE_STATUSES.passed]: COLORS.testSuiteStatuses.passed,
  [TEST_SUITE_STATUSES.error]: COLORS.testSuiteStatuses.error,
  [TEST_SUITE_STATUSES.failed]: COLORS.testSuiteStatuses.failed,
  [TEST_SUITE_STATUSES.skipped]: COLORS.testSuiteStatuses.skipped,
  [TEST_SUITE_STATUSES.total]: COLORS.testSuiteStatuses.skipped,
};

export const TEST_CASE_FILE_MASK = '%test_case%-hh-mm-ss-YYYY-MM-DD';

export const TEST_SUITE_HISTORICAL_DATA_MONTHS_DEFAULT_ITEMS = [1, 3, 6, 12];

export const JUNIT_ALARM_CONNECTOR = 'junit';
