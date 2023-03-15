import { SNMP_TEMPLATE_STATE_STATES } from '@/constants';

export default {
  oid: 'oid',
  module: 'Select a mib module',
  output: 'output',
  resource: 'resource',
  component: 'component',
  connectorName: 'connector_name',
  toCustom: 'To custom',
  defineVar: 'Define matching snmp var',
  writeTemplate: 'Write template',
  state: 'severity',
  moduleMibObjects: 'Snmp vars match field',
  regex: 'Regex',
  formatter: 'Format (capture group with \\x)',
  uploadMib: 'Upload MIB',
  addSnmpRule: 'Add SNMP rule',
  uploadedMibPopup:
    'File was uploaded.\nNotifications: {notification}\nObjects: {object}'
    + '|Files were uploaded.\nNotifications: {notification}\nObjects: {object}',
  states: {
    [SNMP_TEMPLATE_STATE_STATES.info]: 'Info',
    [SNMP_TEMPLATE_STATE_STATES.minor]: 'Minor',
    [SNMP_TEMPLATE_STATE_STATES.major]: 'Major',
    [SNMP_TEMPLATE_STATE_STATES.critical]: 'Critical',
  },
};
