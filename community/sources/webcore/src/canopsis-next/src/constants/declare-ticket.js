export const DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES = {
  header: 'index .Header "%field_name%"',
  response: 'index .Response "%field_name%"',
  responseByStep: 'index .ResponseMap "%N%.%field_name%"',
};

export const DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES = {
  author: '.AdditionalData.Author',
  user: '.AdditionalData.User',
  alarmChangeType: '.AdditionalData.Trigger',
  initiator: '.AdditionalData.Initiator',
  output: '.AdditionalData.Output',
  ruleName: '.AdditionalData.RuleName',
};

export const DECLARE_TICKET_EXECUTION_STATUSES = {
  waiting: 0,
  running: 1,
  succeeded: 2,
  failed: 3,
};
