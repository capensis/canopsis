export const JOB_STATUS = {
  running: 0,
  paused: 1,
  stopped: 2,
};

export const JOB_EXECUTION_STATUS = {
  running: 1,
  succeeded: 2,
  failed: 3,
  aborted: 4,
};

export const JOB_LAST_RUN_STATUS = {
  succeed: 0,
  failed: 1,
};

export const JOB_RULE_TYPE = {
  ticketDeclarationRule: 'ticket_declaration_rule',
  scenario: 'scenario',
};

export const JOB_RULE_TYPES = {
  scenario: 0,
  declareTicket: 1,
};

export const JOBS_TABS = {
  instructions: 'instructions',
  webhooks: 'webhooks',
  ticketStatus: 'ticketStatus',
  authToken: 'authToken',
};

export const JOBS_TABS_PERMISSIONS = {
  [JOBS_TABS.instructions]: 'remediationJobManagement',
  [JOBS_TABS.webhooks]: 'webhookJobManagement',
  [JOBS_TABS.ticketStatus]: 'ticketStatusJobManagement',
  [JOBS_TABS.authToken]: 'webhookTokenJobManagement',
};

export const JOB_ACTION_REFETCH_DELAY = 3000;
