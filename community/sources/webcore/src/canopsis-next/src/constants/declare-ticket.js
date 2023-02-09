export const DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES = {
  header: '.Header.%field_name%',
  response: '.Response.%field_name%',
  headerByStep: 'index .ResponseMap "%N%.%field_name%"',
};

export const DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES = {
  author: '.AdditionalData.Author',
  user: '.AdditionalData.User',
  alarmChangeType: '.AdditionalData.AlarmChangeType',
  initiator: '.AdditionalData.Initiator',
  output: '.AdditionalData.Output',
};

export const DECLARE_TICKET_EXECUTION_STATUSES = {
  running: 1,
  succeeded: 2,
  failed: 3,
};
