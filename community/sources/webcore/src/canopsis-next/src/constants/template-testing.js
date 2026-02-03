import { USER_PERMISSIONS } from './permission';

export const TEMPLATE_TESTING_TABS = {
  data: 'data',
  tests: 'tests',
};

export const TEMPLATE_TESTING_DATA_TYPES = {
  event: 0,
  response: 1,
};

export const TEMPLATE_TESTING_TEST_TYPES = {
  eventFilter: 0,
  linkRule: 1,
  scenario: 2,
  widget: 3,
  declareTicketRule: 4,
  dynamicInfo: 5,
  instruction: 6,
  job: 7,
  metaAlarmRule: 8,
  externalAuthToken: 9,
};

export const TEMPLATE_TESTING_TESTS_TYPES_TO_PERMISSIONS = {
  [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: USER_PERMISSIONS.technical.exploitation.eventFilter,
  [TEMPLATE_TESTING_TEST_TYPES.linkRule]: USER_PERMISSIONS.technical.exploitation.linkRule,
  [TEMPLATE_TESTING_TEST_TYPES.scenario]: USER_PERMISSIONS.technical.exploitation.scenario,
  [TEMPLATE_TESTING_TEST_TYPES.widget]: USER_PERMISSIONS.technical.view,
  [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: USER_PERMISSIONS.technical.exploitation.declareTicketRule,
  [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: USER_PERMISSIONS.technical.exploitation.dynamicInfo,
  [TEMPLATE_TESTING_TEST_TYPES.instruction]: USER_PERMISSIONS.technical.remediationInstruction,
  [TEMPLATE_TESTING_TEST_TYPES.job]: USER_PERMISSIONS.technical.remediationJob,
  [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: USER_PERMISSIONS.technical.exploitation.metaAlarmRule,
};

export const TEMPLATE_TESTING_DATA_EVENT_PRE_FILLED_TEMPLATE = JSON.stringify({
  connector: 'example_connector',
  connector_name: 'example_connectorname',
  source_type: 'resource',
  event_type: 'check',
  component: 'example_component',
  resource: 'example_resource_1',
  state: 1,
  output: 'example alarm',
  author: 'root',
}, null, 2);

export const TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES = {
  event: 'event',
  response: 'response',
  alarm: 'alarm',
  user: 'user',
  entity: 'entity',
};

export const TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES_TO_DATA_TYPE = {
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event]: TEMPLATE_TESTING_DATA_TYPES.event,
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response]: TEMPLATE_TESTING_DATA_TYPES.response,
};

export const TEMPLATE_TESTING_TEST_VARIABLES_EDITOR_LINE_HEIGHT = 19;

export const TEMPLATE_TESTING_TEST_VARIABLES_MIN_EDITOR_LINES = 2;
