/**
 * Builds a bulk stop payload for job-management endpoints.
 *
 * @param {Object[]} jobs
 * @param {Object<string, boolean>} stopPendingJobsById
 * @param {boolean} withStopPendingJobs
 * @returns {Object[]}
 */
export const prepareJobsStopRequest = (
  jobs = [],
  stopPendingJobsById = {},
  withStopPendingJobs = true,
) => jobs.map((job) => {
  const request = { _id: job._id };

  if (withStopPendingJobs) {
    request.stop_pending_jobs = job.pending_jobs > 0
      ? Boolean(stopPendingJobsById[job._id])
      : false;
  }

  return request;
});

/**
 * Splits a bulk API response into successful and failed items.
 *
 * @param {Object[]} response
 * @returns {{ successItems: Object[], failedItems: Object[] }}
 */
export const splitJobsBulkResponse = (response = []) => response.reduce(
  (acc, item) => {
    acc[item.error ? 'failedItems' : 'successItems'].push(item);

    return acc;
  },
  { successItems: [], failedItems: [] },
);
