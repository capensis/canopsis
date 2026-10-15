import { COLORS } from '@/config';
import { SNMP_TEMPLATE_STATE_STATES } from '@/constants';

/**
 * Get color by snmp rule state
 *
 * @param {number} value
 * @returns {string}
 */
export const getSnmpRuleStateColor = value => ({
  [SNMP_TEMPLATE_STATE_STATES.info]: COLORS.state.ok,
  [SNMP_TEMPLATE_STATE_STATES.minor]: COLORS.state.minor,
  [SNMP_TEMPLATE_STATE_STATES.major]: COLORS.state.major,
  [SNMP_TEMPLATE_STATE_STATES.critical]: COLORS.state.critical,
}[value]);
