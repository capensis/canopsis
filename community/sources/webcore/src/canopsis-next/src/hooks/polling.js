import { unref } from 'vue';

/**
 * Function to handle polling mechanism for a process
 *
 * @param {Object} options - Options for polling mechanism
 * @param {Function} options.startHandler - Function to start the polling process
 * @param {Function} options.processHandler - Function to process the polling
 * @param {Function} options.endHandler - Function to handle the end of polling
 * @param {number} [options.interval = 2000] - Interval in milliseconds for polling (default: 2000)
 * @returns {Object} Object containing the polling function
 */
export const usePolling = ({ startHandler, processHandler, endHandler, interval = 2000 }) => {
  let cancelWait = () => {};

  /**
   * Function to wait for the process to complete
   *
   * @param {Object} options - Options for the process
   * @returns {Promise} Promise that resolves when the process is completed
   */
  const wait = options => new Promise((resolve, reject) => {
    let finished = false;
    let cancelled = false;

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
    const startResponse = await startHandler(...args);

    const waitResponse = await wait({ ...startResponse, ...args });

    return endHandler(waitResponse);
  };

  const cancel = () => cancelWait();

  return {
    poll,
    cancel,
  };
};
