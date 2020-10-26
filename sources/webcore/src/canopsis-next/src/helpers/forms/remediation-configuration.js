import { REMEDIATION_CONFIGURATION_TYPES } from '@/constants';

/**
 * Convert remediation configuration entity to form object
 *
 * @typedef {Object} RemediationConfiguration
 * @property {string} name
 * @property {string} host
 * @property {string} token
 * @property {string} type
 * @param {RemediationConfiguration} remediationConfiguration
 * @return {RemediationConfiguration}
 */
export const remediationConfigurationToForm = (remediationConfiguration = {}) => ({
  name: remediationConfiguration.name || '',
  host: remediationConfiguration.host || '',
  type: remediationConfiguration.type || REMEDIATION_CONFIGURATION_TYPES.rundeck,
  auth_token: remediationConfiguration.auth_token || '',
});

/**
 * Convert remediation configuration form object to API compatible object
 *
 * @param {RemediationConfiguration} form
 * @return {RemediationConfiguration}
 */
export const formToRemediationConfiguration = form => ({ ...form });
