import { JOB_STATE, JOB_RUN_STATUS, JOB_RULE_TYPE } from '@/constants/jobs';

export default {
  filterByStatus: 'Filter by status',
  ruleName: 'Rule name',
  authTokenName: 'Auth token name',
  ticketSystemName: 'Ticket system name',
  ruleType: 'Rule type',
  active: 'Active',
  statusLabel: 'Status',
  startDate: 'Start date',
  finishDate: 'Finish date',
  failReason: 'Fail reason',
  expirationDate: 'Expiration date',
  searchByRuleName: 'Search by rule name',
  webhookFailedPrefix: 'Webhook is failed',
  status: {
    [JOB_STATE.running]: 'Running',
    [JOB_STATE.paused]: 'Paused',
    [JOB_STATE.stopped]: 'Stopped',
    unknown: 'Unknown',
  },
  activeState: {
    [JOB_STATE.running]: 'Active',
    [JOB_STATE.paused]: 'Paused',
    [JOB_STATE.stopped]: 'Stopped',
  },
  runStatus: {
    [JOB_RUN_STATUS.succeed]: 'Succeed',
    [JOB_RUN_STATUS.failed]: 'Failed',
    inProgress: 'In progress',
  },
  ruleTypeValues: {
    [JOB_RULE_TYPE.ticketDeclarationRule]: 'Ticket declaration rule',
    [JOB_RULE_TYPE.scenario]: 'Scenario',
  },
  tabs: {
    instructions: 'Instructions',
    webhooks: 'Webhooks',
    ticketStatus: 'Ticket status',
    authToken: 'Auth token',
  },
  actions: {
    start: 'Start',
    stop: 'Stop',
    resume: 'Resume',
    pause: 'Pause',
    edit: 'Edit',
  },
};
