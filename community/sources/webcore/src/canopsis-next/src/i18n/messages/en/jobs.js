import {
  JOB_STATUS,
  JOB_LAST_RUN_STATUS,
  JOB_RULE_TYPE,
  JOB_RULE_TYPES,
  TICKET_STATUS_JOBS_TABS,
} from '@/constants';

export default {
  filterByStatus: 'Filter by status',
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
  data: {
    request: 'Request',
    response: 'Response',
    webhookFailedPrefix: 'Webhook is failed',
  },
  status: {
    [JOB_STATUS.running]: 'Running',
    [JOB_STATUS.paused]: 'Paused',
    [JOB_STATUS.stopped]: 'Stopped',
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
    [JOB_RULE_TYPES.declareTicket]: 'Ticket declaration rules',
  },
  actions: {
    editJob: 'Edit job',
    repeatJob: 'Repeat job',
    pauseJob: 'Pause job',
    startJob: 'Start job',
  },
  tabs: {
    [TICKET_STATUS_JOBS_TABS.ticketStatus]: 'Ticket status',
  },
  popups: {
    updated: 'Job updated',
    repeated: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> repeated | {count} jobs repeated',
    restarted: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> restarted | {count} jobs restarted',
    repeatFailed: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to repeat | {count} jobs failed to repeat',
    restartFailed: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to restart | {count} jobs failed to restart',
    paused: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> paused | {count} jobs paused',
    pauseFailed: 'Job for <strong>{ruleName}</strong> / <strong>Ticket number {ticketNumber}</strong> failed to pause | {count} jobs failed to pause',
  },
};
