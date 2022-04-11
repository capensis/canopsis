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
 * @typedef {Object} RemediationJobExecution
 * @property {string} _id
 * @property {string} job_id
 * @property {string} name
 * @property {string} payload
 * @property {number} started_at
 * @property {number} launched_at
 * @property {number} completed_at
 * @property {number} queue_number
 * @property {string} fail_reason
 * @property {string} output
 * @property {number | null} status
 * @property {Object | null} query
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

/**
 * Get empty job execution
 *
 * @returns {RemediationJobExecution}
 */
export const getEmptyRemediationJobExecution = () => ({
  _id: '',
  job_id: '',
  name: '',
  payload: '',
  status: null,
  fail_reason: '',
  output: '',
  query: null,
  started_at: 0,
  launched_at: 0,
  completed_at: 0,
  queue_number: 0,
});
