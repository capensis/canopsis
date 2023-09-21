import { COLORS } from '@/config';

export const SNMP_STATE_TYPES = {
  simple: 'simple',
  template: 'template',
};

export const SNMP_TEMPLATE_STATE_STATES = {
  info: 'info',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const SNMP_TEMPLATE_STATE_STATES_COLORS = {
  [SNMP_TEMPLATE_STATE_STATES.info]: COLORS.state.ok,
  [SNMP_TEMPLATE_STATE_STATES.minor]: COLORS.state.minor,
  [SNMP_TEMPLATE_STATE_STATES.major]: COLORS.state.major,
  [SNMP_TEMPLATE_STATE_STATES.critical]: COLORS.state.critical,
};
