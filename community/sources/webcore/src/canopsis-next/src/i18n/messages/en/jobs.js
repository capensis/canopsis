import {
  JOB_STATUS,
  JOB_LAST_RUN_STATUS,
  JOB_RULE_TYPE,
  JOB_RULE_TYPES,
  TICKET_STATUS_JOBS_TABS,
} from '@/constants';

export default {
  filterByStatus: 'Filter by status',
  filterByActiveState: 'Filter by active state',
  ruleName: 'Rule name',
  authTokenName: 'Auth token name',
  ticketSystemName: 'Ticket system name',
  ticketNumber: 'Ticket number',
  ruleType: 'Rule type',
  active: 'Active',
  lastStatus: 'Last status',
  startDate: 'Start date',
  finishDate: 'Finish date',
  failReason: 'Fail reason',
  expirationDate: 'Expiration date',
  searchByRuleName: 'Search by rule name',
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
    [TICKET_STATUS_JOBS_TABS.instructions]: 'Instructions',
    [TICKET_STATUS_JOBS_TABS.webhooks]: 'Webhooks',
    [TICKET_STATUS_JOBS_TABS.ticketStatus]: 'Ticket status',
    [TICKET_STATUS_JOBS_TABS.authToken]: 'Auth token',
  },
};
