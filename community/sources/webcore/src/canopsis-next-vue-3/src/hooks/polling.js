import axios from 'axios';
import { unref } from 'vue';

import { EXPORT_FETCHING_INTERVAL } from '@/config';
import { EXPORT_STATUSES } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';

/**
 * Creates a polling runner with cancel support.
 *
 * @param {Object} options
 * @param {Function} options.startHandler
 * @param {Function} options.processHandler
 * @param {Function} [options.endHandler]
 * @param {number|import('vue').Ref<number>} [options.interval]
 * @returns {{ poll: Function, cancel: Function }}
 */
export const usePolling = ({
  startHandler,
  processHandler,
  endHandler = v => v,
  interval = 2000,
}) => {
  let cancelWait = () => {};
  let finished = false;
  let cancelled = false;

  const wait = (options, cancelToken) =>
    new Promise((resolve, reject) => {
      const customResolve = (...args) => {
        finished = true;

        return resolve(...args);
      };

      const customReject = (err) => {
        finished = true;

        return reject(err);
      };

      const callTimeout = async () => {
        try {
          await processHandler(options, customResolve, customReject, cancelToken);
        } catch (err) {
          customReject(err);

          return;
        }

        if (cancelled) {
          reject();
        }

        if (!finished) {
          setTimeout(callTimeout, unref(interval));
        }
      };

      return callTimeout();
    });

  const poll = async (...args) => {
    const source = axios.CancelToken.source();
    const { token: cancelToken } = source;

    finished = false;
    cancelled = false;

    cancelWait = () => {
      cancelled = true;
      source.cancel();
    };

    const startResponse = await startHandler(...args, cancelToken);

    if (cancelled) {
      return startResponse;
    }

    const waitResponse = await wait({ ...startResponse, ...args }, cancelToken);

    return endHandler(waitResponse);
  };

  return {
    poll,
    cancel: () => cancelWait(),
  };
};

/**
 * Creates polling runner with pending state.
 *
 * @param {Object} options
 * @returns {{ pending: import('vue').Ref<boolean>, poll: Function, cancel: Function }}
 */
export const usePollingWithPending = ({
  startHandler,
  processHandler,
  endHandler,
  interval,
  initialPending = false,
}) => {
  const { poll: basicPoll, cancel } = usePolling({
    startHandler,
    processHandler,
    endHandler,
    interval,
  });
  const { pending, handler: poll } = usePendingHandler(basicPoll, initialPending);

  return {
    pending,
    poll,
    cancel,
  };
};

export const useFilePolling = ({
  createHandler,
  fetchHandler,
  endHandler = v => v,
  failedHandler = v => v,
  completedStatus = EXPORT_STATUSES.completed,
  failedStatus = EXPORT_STATUSES.failed,
  interval = EXPORT_FETCHING_INTERVAL,
}) => {
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

  return usePolling({
    startHandler: createHandler,
    processHandler,
    endHandler,
    interval,
  });
};

/**
 * Creates file polling runner with pending state.
 *
 * @param {Object} options
 * @returns {{ pending: import('vue').Ref<boolean>, poll: Function, cancel: Function }}
 */
export const useFilePollingWithPending = ({
  createHandler,
  fetchHandler,
  endHandler,
  failedHandler,
  completedStatus,
  failedStatus,
  interval,
  initialPending,
}) => {
  const { poll: basicPoll, cancel } = useFilePolling({
    createHandler,
    fetchHandler,
    endHandler,
    failedHandler,
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
