import {
  JOB_STATUS,
  JOB_EXECUTION_STATUS,
  JOB_LAST_RUN_STATUS,
  JOB_RULE_TYPE,
  JOB_RULE_TYPES,
  JOBS_TABS,
} from '@/constants';

export default {
  filterByStatus: 'Filter by status',
  filterByRuleType: 'Filter by rule type',
  filterByLastStatus: 'Filter by last status',
  filterByActiveState: 'Filter by active state',
  ruleName: 'Rule name',
  ruleNameTooltip: 'Rule name at the time of job creation',
  authTokenName: 'Auth token name',
  ticketSystemName: 'Ticket system name',
  ticketNumber: 'Ticket number',
  ruleType: 'Rule type',
  activeState: 'Active state',
  lastStatus: 'Last status',
  startDate: 'Start date',
  finishDate: 'Finish date',
  failReason: 'Fail reason',
  expirationDate: 'Expiration date',
  searchByRuleName: 'Search by rule name, ticket system name or ticket number',
  searchByRuleNameOnly: 'Search by rule name',
  searchByAuthTokenName: 'Search by auth token name',
  data: {
    request: 'Request',
    response: 'Response',
    webhookFailedPrefix: 'Webhook is failed',
    statusPrefix: {
      [JOBS_TABS.instructions]: {
        [JOB_EXECUTION_STATUS.running]: 'Instruction job is running',
        [JOB_EXECUTION_STATUS.succeeded]: 'Instruction job executed successfully',
        [JOB_EXECUTION_STATUS.failed]: 'Instruction job is failed',
        [JOB_EXECUTION_STATUS.aborted]: 'Instruction job is aborted',
      },
      [JOBS_TABS.webhooks]: {
        [JOB_EXECUTION_STATUS.running]: 'Webhook is running',
        [JOB_EXECUTION_STATUS.succeeded]: 'Webhook executed successfully',
        [JOB_EXECUTION_STATUS.failed]: 'Webhook is failed',
        [JOB_EXECUTION_STATUS.aborted]: 'Webhook is aborted',
      },
      [JOBS_TABS.authToken]: {
        [JOB_EXECUTION_STATUS.running]: 'Auth token job is running',
        [JOB_EXECUTION_STATUS.succeeded]: 'Auth token job executed successfully',
        [JOB_EXECUTION_STATUS.failed]: 'Auth token job is failed',
        [JOB_EXECUTION_STATUS.aborted]: 'Auth token job is aborted',
      },
    },
  },
  executionStatus: {
    [JOB_EXECUTION_STATUS.running]: 'In progress',
    [JOB_EXECUTION_STATUS.succeeded]: 'Succeeded',
    [JOB_EXECUTION_STATUS.failed]: 'Failed',
    [JOB_EXECUTION_STATUS.aborted]: 'Aborted',
  },
  status: {
    [JOB_STATUS.running]: 'Running',
    [JOB_STATUS.paused]: 'Paused',
    [JOB_STATUS.stopped]: 'Stopped',
  },
  statusTooltip: {
    [JOB_STATUS.stopped]: 'Stopped\n(Status received from external API)',
  },
  lastRunStatus: {
    [JOB_LAST_RUN_STATUS.succeed]: 'Succeed',
    [JOB_LAST_RUN_STATUS.failed]: 'Failed',
  },
  ruleTypeValues: {
    [JOB_RULE_TYPE.ticketDeclarationRule]: 'Ticket declaration rule',
    [JOB_RULE_TYPE.scenario]: 'Scenario',
  },
  types: {
    [JOB_RULE_TYPES.scenario]: 'Scenario',
    [JOB_RULE_TYPES.declareTicket]: 'Ticket declaration rule',
  },
  actions: {
    editJob: 'Edit job',
    repeatJob: 'Repeat job',
    pauseJob: 'Pause job',
    startJob: 'Start job',
    stopJob: 'Stop job',
  },
  tabs: {
    [JOBS_TABS.instructions]: 'Instructions',
    [JOBS_TABS.webhooks]: 'Webhooks',
    [JOBS_TABS.ticketStatus]: 'Ticket status',
    [JOBS_TABS.authToken]: 'Auth token',
  },
  popups: {
    updated: 'Job updated',
    repeated:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> repeated | {count} jobs repeated',
    restarted:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> restarted | {count} jobs restarted',
    repeatFailed:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to repeat | {count} jobs failed to repeat',
    restartFailed:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to restart | {count} jobs failed to restart',
    paused:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> paused | {count} jobs paused',
    pauseFailed:
      'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to pause | {count} jobs failed to pause',
    stopped: 'Job <strong>{jobName}</strong> stopped | {count} jobs stopped',
    stopFailed: 'Job failed to stop | {count} jobs failed to stop',
  },
};
