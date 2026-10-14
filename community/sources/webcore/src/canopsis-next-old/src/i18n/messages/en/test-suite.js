import { TEST_SUITE_STATUSES } from '@/constants';

export default {
  xmlFeed: 'XML feed',
  hostname: 'Host name',
  lastUpdate: 'Last update',
  totalTests: 'Total tests',
  disabledTests: 'Tests disabled',
  copyMessage: 'Copy system message',
  systemError: 'System error',
  systemErrorMessage: 'System error message',
  systemOut: 'System out',
  systemOutMessage: 'System out message',
  compareWithHistorical: 'Compare with historical data',
  className: 'Classname',
  line: 'Line',
  failureMessage: 'Failure message',
  noData: 'No system messages found in XML',
  tabs: {
    globalMessages: 'Global messages',
    gantt: 'Gantt',
    details: 'Details',
    screenshots: 'Screenshots',
    videos: 'Videos',
  },
  statuses: {
    [TEST_SUITE_STATUSES.passed]: 'Passed',
    [TEST_SUITE_STATUSES.skipped]: 'Skipped',
    [TEST_SUITE_STATUSES.error]: 'Error',
    [TEST_SUITE_STATUSES.failed]: 'Failed',
    [TEST_SUITE_STATUSES.total]: 'Total time taken',
  },
  popups: {
    systemMessageCopied: 'System message copied to clipboard',
  },
};
