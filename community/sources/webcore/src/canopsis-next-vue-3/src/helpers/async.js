/**
 * Waits for a specified timeout before resolving a promise.
 * This function is useful for introducing delays in asynchronous workflows.
 *
 * @function
 * @param {number} delay - The time, in milliseconds, to wait before the promise resolves.
 * @returns {Promise} A promise that resolves after the specified delay.
 *
 * @example
 * // Example 1: Using promisedWait to introduce a delay
 * promisedWait(2000).then(() => {
 *   console.log('Waited for 2 seconds');
 * });
 *
 * @example
 * // Example 2: Using promisedWait in an async function
 * async function delayedExecution() {
 *   console.log('Start delay');
 *   await promisedWait(1000);
 *   console.log('End delay after 1 second');
 * }
 *
 * delayedExecution();
 */
export const promisedWait = (delay) => new Promise((resolve) => setTimeout(resolve, delay));
