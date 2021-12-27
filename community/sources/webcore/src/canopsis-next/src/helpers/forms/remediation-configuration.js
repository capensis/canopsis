import { pick } from 'lodash';

import { REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME } from '@/constants';

/**
 * @typedef {Object} RemediationConfiguration
 * @property {string} name
 * @property {string} host
 * @property {string} type
 * @property {string} auth_token
 * @property {string} auth_username
 */

/**
 * @typedef {Object} RemediationConfigurationFormType
 * @property {string} name
 * @property {string} auth_type
 */

/**
 * @typedef {RemediationConfiguration} RemediationConfigurationForm
 * @property {RemediationConfigurationFormType | string} type
 */

/**
 * Convert remediation configuration entity to form object
 *
 * @param {RemediationConfiguration} remediationConfiguration
 * @return {RemediationConfigurationForm}
 */
export const remediationConfigurationToForm = (remediationConfiguration = {}) => ({
  name: remediationConfiguration.name || '',
  host: remediationConfiguration.host || '',
  type: remediationConfiguration.type || '',
  auth_token: remediationConfiguration.auth_token || '',
  auth_username: remediationConfiguration.auth_username || '',
});

/**
 * Convert remediation configuration form object to API compatible object
 *
 * @param {RemediationConfigurationForm} form
 * @return {RemediationConfiguration}
 */
export const formToRemediationConfiguration = (form) => {
  const remediationConfiguration = pick(form, ['name', 'host', 'auth_token']);
  const type = form.type || {};

  remediationConfiguration.type = type.name;

  if (REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME.includes(type.auth_type)) {
    remediationConfiguration.auth_username = form.auth_username;
  }

  return remediationConfiguration;
};
