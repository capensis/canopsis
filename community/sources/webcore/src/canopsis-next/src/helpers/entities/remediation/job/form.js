import { omit } from 'lodash';

import { REMEDIATION_JOB_EXECUTION_STATUSES } from '@/constants';

import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { durationWithEnabledToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} RemediationJob
 * @property {Object|string} config
 * @property {string} _id
 * @property {string} job_id
 * @property {string} name
 * @property {boolean} multiple_executions
 * @property {Duration} [job_wait_interval]
 * @property {string} [payload]
 * @property {Object} [query]
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
 * Check job status is running
 *
 * @param {RemediationJobExecution} job
 * @returns {boolean}
 */
export const isJobExecutionRunning = job => job.status === REMEDIATION_JOB_EXECUTION_STATUSES.running;

/**
 * Check job status is succeeded
 *
 * @param {RemediationJobExecution} job
 * @returns {boolean}
 */
export const isJobExecutionSucceeded = job => job.status === REMEDIATION_JOB_EXECUTION_STATUSES.succeeded;

/**
 * Check job status is failed
 *
 * @param {RemediationJobExecution} job
 * @returns {boolean}
 */
export const isJobExecutionFailed = job => job.status === REMEDIATION_JOB_EXECUTION_STATUSES.failed;

/**
 * Check job status is canceled
 *
 * @param {RemediationJobExecution} job
 * @returns {boolean}
 */
export const isJobExecutionCancelled = job => job.status === REMEDIATION_JOB_EXECUTION_STATUSES.canceled;

/**
 * Check job is finished
 *
 * @param {RemediationJobExecution} job
 * @returns {boolean}
 */
export const isJobFinished = job => [
  REMEDIATION_JOB_EXECUTION_STATUSES.canceled,
  REMEDIATION_JOB_EXECUTION_STATUSES.failed,
  REMEDIATION_JOB_EXECUTION_STATUSES.succeeded,
].includes(job.status);

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
  job_wait_interval: remediationJob.job_wait_interval
    ? durationWithEnabledToForm({ enabled: true, ...remediationJob.job_wait_interval })
    : { value: undefined, unit: undefined },
});

/**
 * Convert remediation job form object to API compatible object
 *
 * @param {RemediationJobForm} form
 * @param {RemediationConfigurationFormType} configType
 * @return {RemediationJob}
 */
export const formToRemediationJob = (form, configType) => {
  const { job_wait_interval: jobWaitInterval, config, payload, query, ...remediationJob } = form;

  if (configType?.with_body) {
    remediationJob.payload = payload;
  }

  if (configType?.with_query) {
    remediationJob.query = textPairsToObject(query);
  }

  if (jobWaitInterval?.enabled) {
    remediationJob.job_wait_interval = omit(jobWaitInterval, ['enabled']);
  }

  return {
    ...remediationJob,

    config: config._id,
  };
};

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
