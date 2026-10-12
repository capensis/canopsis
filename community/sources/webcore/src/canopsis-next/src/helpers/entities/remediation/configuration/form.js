import { pick } from 'lodash';

import { REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME } from '@/constants';

/**
 * @typedef {Object} RemediationConfiguration
 * @property {string} name
 * @property {string} host
 * @property {string} type
 * @property {string} auth_token
 * @property {string} auth_username
 * @property {boolean} skip_verify
 */

/**
 * @typedef {Object} RemediationConfigurationFormType
 * @property {string} name
 * @property {string} auth_type
 * @property {boolean} with_query
 * @property {boolean} with_body
 */

/**
 * @typedef {RemediationConfiguration} RemediationConfigurationForm
 * @property {RemediationConfigurationFormType | string} type
 */

/**
 * Check is job has a username
 *
 * @param {RemediationConfigurationFormType} [type]
 * @returns {boolean}
 */
export const isJobTypeIncludesUserName = (type = {}) => REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME
  .includes(type.auth_type);

/**
 * Convert remediation configuration entity to form object
 *
 * @param {RemediationConfiguration} remediationConfiguration
 * @return {RemediationConfigurationForm}
 */
export const remediationConfigurationToForm = (remediationConfiguration = {}) => ({
  name: remediationConfiguration.name ?? '',
  host: remediationConfiguration.host ?? '',
  type: remediationConfiguration.type ?? '',
  auth_token: remediationConfiguration.auth_token ?? '',
  auth_username: remediationConfiguration.auth_username ?? '',
  skip_verify: remediationConfiguration.skip_verify ?? false,
});

/**
 * Convert remediation configuration form object to API compatible object
 *
 * @param {RemediationConfigurationForm} form
 * @return {RemediationConfiguration}
 */
export const formToRemediationConfiguration = (form) => {
  const remediationConfiguration = pick(form, [
    'type',
    'name',
    'host',
    'auth_token',
    'skip_verify',
  ]);

  if (form.auth_username) {
    remediationConfiguration.auth_username = form.auth_username;
  }

  return remediationConfiguration;
};
