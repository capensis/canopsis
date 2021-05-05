/**
 * @typedef {Object} RemediationJob
 * @property {Object|string} config
 * @property {string} _id
 * @property {string} job_id
 * @property {string} name
 * @property {string} payload
 */
/**
 * @typedef {RemediationJob} RemediationJobForm
 */

/**
 * Convert remediation job entity to form object
 *
 * @param {RemediationJob} remediationJob
 * @return {RemediationJobForm}
 */
export const remediationJobToForm = (remediationJob = {}) => ({
  config: remediationJob.config || '',
  job_id: remediationJob.job_id || '',
  name: remediationJob.name || '',
  payload: remediationJob.payload || '{}',
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
  payload: form.payload,
});
