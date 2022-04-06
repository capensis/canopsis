import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} RemediationJob
 * @property {Object|string} config
 * @property {string} _id
 * @property {string} job_id
 * @property {string} name
 * @property {string} payload
 * @property {boolean} multiple_executions
 * @property {Object} query
 */
/**
 * @typedef {RemediationJob} RemediationJobForm
 * @property {TextPairObject[]} query
 */

/**
 * Convert remediation job entity to form object
 *
 * @param {RemediationJob} remediationJob
 * @return {RemediationJobForm}
 */
export const remediationJobToForm = (remediationJob = {}) => ({
  config: remediationJob.config ?? '',
  job_id: remediationJob.job_id ?? '',
  name: remediationJob.name ?? '',
  payload: remediationJob.payload ?? '',
  multiple_executions: remediationJob.multiple_executions ?? false,
  query: remediationJob.query ? objectToTextPairs(remediationJob.query) : [],
});

/**
 * Convert remediation job form object to API compatible object
 *
 * @param {RemediationJobForm} form
 * @return {RemediationJob}
 */
export const formToRemediationJob = form => ({
  ...form,

  config: form.config._id,
  query: textPairsToObject(form.query),
});
