export const DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES_WITHOUT_UNKNOWN = {
  open: 1,
  assigned: 2,
  inProgress: 3,
  closed: 4,
};

export const DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES = {
  unknown: 0,

  ...DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES_WITHOUT_UNKNOWN,
};

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
