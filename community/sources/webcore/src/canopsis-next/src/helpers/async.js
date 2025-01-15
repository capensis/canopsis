/**
 * Executes a callback function after a specified timeout and returns a promise that resolves with
 * the callback's result.
 *
 * @function
 * @param {Function} callback - The asynchronous function to be executed after the timeout. It should return a promise
 * or a value.
 * @param {number} [timeout=0] - The time, in milliseconds, to wait before executing the callback.
 * @returns {Promise<*>} A promise that resolves with the result of the callback function.
 *
 * @example
 * // Example 1: Using promisesTimeout with a simple callback
 * promisesTimeout(() => 'Hello, World!', 1000).then(result => {
 *   console.log(result); // Logs 'Hello, World!' after 1 second
 * });
 *
 * @example
 * // Example 2: Using promisesTimeout with an asynchronous callback
 * const asyncFunction = async () => {
 *   return new Promise(resolve => {
 *     setTimeout(() => {
 *       resolve('Async Result');
 *     }, 500);
 *   });
 * };
 *
 * promisedTimeout(asyncFunction, 1000).then(result => {
 *   console.log(result); // Logs 'Async Result' after 1 second
 * });
 *
 * @example
 * // Example 3: Using promisedTimeout with default timeout
 * promisesTimeout(() => 'Immediate Result').then(result => {
 *   console.log(result); // Logs 'Immediate Result' immediately
 * });
 */
export const promisedTimeout = (callback, timeout = 0) => new Promise(resolve => setTimeout(async () => {
  const result = await callback?.();

  return resolve(result);
}, timeout));

/**
 * Waits for a specified timeout before resolving a promise.
 * This function is useful for introducing delays in asynchronous workflows.
 *
 * @function
 * @param {number} timeout - The time, in milliseconds, to wait before the promise resolves.
 * @returns {Promise<void>} A promise that resolves after the specified timeout.
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
export const promisedWait = timeout => promisedTimeout(() => {}, timeout);
