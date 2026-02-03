import { unref } from 'vue';

import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';

/**
 * Function to handle polling mechanism for a process
 *
 * @param {Object} options - Options for polling mechanism
 * @param {Function} options.startHandler - Function to start the polling process
 * @param {Function} options.processHandler - Function to process the polling
 * @param {Function} [options.endHandler = v => v] - Function to handle the end of polling (default: v => v)
 * @param {number} [options.interval = 2000] - Interval in milliseconds for polling (default: 2000)
 * @returns {Object} Object containing the polling function
 */
export const usePolling = ({ startHandler, processHandler, endHandler = v => v, interval = 2000 }) => {
  let cancelWait = () => {};
  let finished = false;
  let cancelled = false;

  /**
   * Function to wait for the process to complete
   *
   * @param {Object} options - Options for the process
   * @returns {Promise} Promise that resolves when the process is completed
   */
  const wait = options => new Promise((resolve, reject) => {
    const customResolve = (...args) => {
      finished = true;

      return resolve(...args);
    };

    cancelWait = () => cancelled = true;

    const customReject = (err) => {
      finished = true;

      return reject(err);
    };

    const callTimeout = async () => {
      try {
        await processHandler(options, customResolve, customReject);
      } catch (err) {
        customReject(err);

        return;
      }

      if (cancelled) {
        reject();
      }

      if (finished) {
        return;
      }

      setTimeout(callTimeout, unref(interval));
    };

    return callTimeout();
  });

  /**
   * Function to initiate the polling process
   *
   * @param {...any} args - Arguments for starting the polling process
   * @returns {Promise} Promise that resolves when the polling process is completed
   */
  const poll = async (...args) => {
    finished = false;
    cancelled = false;

    const startResponse = await startHandler(...args);

    if (cancelled) {
      return startResponse;
    }

    const waitResponse = await wait({ ...startResponse, ...args });

    return endHandler(waitResponse);
  };

  const cancel = () => cancelWait();

  return {
    poll,
    cancel,
  };
};

/**
 * Hook to handle polling with pending state tracking
 *
 * @param {Object} options - Options for polling with pending state tracking
 * @param {Function} options.startHandler - Function to start the polling process
 * @param {Function} options.processHandler - Function to process the polling
 * @param {Function} options.endHandler - Function to handle the end of polling
 * @param {number} options.interval - Interval in milliseconds for polling
 * @param {boolean} [options.initialPending = false] - Initial value for the pending state
 * @returns {Object} Object containing the polling function and the pending state
 * @property {Ref<boolean>} pending - Reactive reference to the pending state
 * @property {Function} poll - Function to start polling with pending state tracking
 */
export const usePollingWithPending = ({
  startHandler,
  processHandler,
  endHandler,
  interval,
  initialPending = false,
}) => {
  const { poll: basicPoll, cancel } = usePolling({ startHandler, processHandler, endHandler, interval });
  const { pending, handler: poll } = usePendingHandler(basicPoll, initialPending);

  return {
    pending,
    poll,
    cancel,
  };
};

/**
 * Function to handle a file polling (import or export)
 *
 * @param {Object} options - Options for polling a file
 * @param {Function} options.createHandler - Function to create polling
 * @param {Function} options.fetchHandler - Function to fetch polling status
 * @param {Function} [options.failedHandler = v => v] - Function to handle failed polling (default: v => v)
 * @param {Function} [options.endHandler = v => v] - Function to handle end of polling (default: v => v)
 * @param {number} [options.completedStatus = EXPORT_STATUSES.completed] - Status code for completed (default: 1)
 * @param {number} [options.failedStatus = EXPORT_STATUSES.failed] - Status code for failed (default: 2)
 * @param {number} [options.interval = EXPORT_FETCHING_INTERVAL] - Interval in milliseconds for polling (default: 2000)
 * @returns {Object} Object containing the function to start poll
 */
export const useFilePolling = ({
  createHandler,
  fetchHandler,
  endHandler = v => v,
  failedHandler = v => v,
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
  interval = EXPORT_FETCHING_INTERVAL,
}) => {
  /**
   * Function to process the export file
   *
   * @param {Object} options - Options for the export file
   * @param {string} options._id - ID of the file
   * @param {Object} rest - Additional data for the export
   * @param {Function} resolve - Function to resolve the export process
   * @param {Function} reject - Function to reject the export process
   * @returns {Promise} Promise that resolves when the export process is completed
   */
  const processHandler = async ({ _id: id, ...rest }, resolve, reject) => {
    const response = await fetchHandler({ id, ...rest });

    if (response.status === unref(completedStatus)) {
      return resolve(response);
    }

    if (response.status === unref(failedStatus)) {
      return reject(failedHandler(response));
    }

    return response;
  };

  const {
    poll,
    cancel,
  } = usePolling({
    interval,
    endHandler,
    processHandler,
    startHandler: createHandler,
  });

  return {
    poll,
    cancel,
  };
};

/**
 * Function to handle a file polling with pending state tracking
 *
 * @param {Object} options - Options for polling a file with pending state
 * @param {Function} options.createHandler - Function to create polling
 * @param {Function} options.fetchHandler - Function to fetch polling status
 * @param {Function} [options.endHandler] - Function to handle end of polling
 * @param {number} [options.completedStatus] - Status code for completed
 * @param {number} [options.failedStatus] - Status code for failed
 * @param {number} [options.interval] - Interval in milliseconds for polling
 * @param {boolean} [options.initialPending=false] - Initial value for the pending state
 * @returns {Object} Object containing polling functions and state
 * @property {Ref<boolean>} pending - Reactive reference to the pending state
 * @property {Function} poll - Function to start polling with pending state tracking
 * @property {Function} cancel - Function to cancel the polling process
 */
export const useFilePollingWithPending = ({
  createHandler,
  fetchHandler,
  endHandler,
  completedStatus,
  failedStatus,
  interval,
  initialPending,
}) => {
  const { poll: basicPoll, cancel } = useFilePolling({
    createHandler,
    fetchHandler,
    endHandler,
    completedStatus,
    failedStatus,
    interval,
  });

  const { pending, handler: poll } = usePendingHandler(basicPoll, initialPending);

  return {

    pending,
    poll,
    cancel,
  };
};
