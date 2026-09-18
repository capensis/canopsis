import { pick } from 'lodash';

import {
  REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME,
  REMEDIATION_CONFIGURATION_TYPES,
} from '@/constants';

/**
 * @typedef {Object} RemediationConfiguration
 * @property {string} name
 * @property {string} host
 * @property {string} type
 * @property {string} [auth_token]
 * @property {string} [auth_refresh_token]
 * @property {string} [client_id]
 * @property {string} [client_secret]
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
 * @property {boolean} enable_refresh_token
 */

/**
 * Check is job has a username
 *
 * @param {RemediationConfigurationFormType} [type]
 * @returns {boolean}
 */
export const isJobTypeIncludesUserName = (type = {}) =>
  REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME.includes(type.auth_type);

/**
 * Check is configuration type awx
 *
 * @param {RemediationConfigurationFormType | string} [type]
 * @returns {boolean}
 */
export const isAwxConfigurationType = (type = {}) =>
  (type?.name ?? type) === REMEDIATION_CONFIGURATION_TYPES.awx;

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
  auth_token: '',
  enable_refresh_token: Boolean(remediationConfiguration.client_id),
  auth_refresh_token: '',
  client_id: remediationConfiguration.client_id ?? '',
  client_secret: '',
  auth_username: remediationConfiguration.auth_username ?? '',
  skip_verify: remediationConfiguration.skip_verify ?? false,
});

/**
 * Resolve AWX refresh-token fields for the API payload.
 * `enable_refresh_token` is form-only and is never sent.
 *
 * @param {RemediationConfigurationForm} form
 * @return {Object}
 */
const getAwxRefreshTokenFields = (form) => {
  if (!form.enable_refresh_token) {
    return {
      client_id: '',
      auth_refresh_token: '',
      client_secret: '',
    };
  }

  const fields = {
    client_id: form.client_id ?? '',
  };

  if (form.auth_refresh_token) {
    fields.auth_refresh_token = form.auth_refresh_token;
  }

  if (form.client_secret) {
    fields.client_secret = form.client_secret;
  }

  return fields;
};

/**
 * Convert remediation configuration form object to API compatible object
 *
 * @param {RemediationConfigurationForm} form
 * @param {Object} [options]
 * @return {RemediationConfiguration}
 */
export const formToRemediationConfiguration = (form) => {
  const remediationConfiguration = {
    ...pick(form, ['name', 'host', 'skip_verify']),
    type: form.type?.name ?? form.type,
  };

  if (form.auth_username) {
    remediationConfiguration.auth_username = form.auth_username;
  }

  if (form.auth_token) {
    remediationConfiguration.auth_token = form.auth_token;
  }

  if (!isAwxConfigurationType(remediationConfiguration.type)) {
    return remediationConfiguration;
  }

  return {
    ...remediationConfiguration,
    ...getAwxRefreshTokenFields(form),
  };
};
