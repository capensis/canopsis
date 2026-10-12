import { ref } from 'vue';

/**
 * Creates a handler that tracks its pending state.
 *
 * This function wraps any given asynchronous handler to monitor its execution state. It provides a reactive `pending`
 * state that indicates whether the handler is currently executing. The `pending` state starts as `false` and is set to
 * `true` when the handler begins execution and set back to `false` once the handler completes or throws an error.
 *
 * @param {Function} handler - The asynchronous function to be wrapped.
 * @param {boolean} [initialPending=false] - The initial value of the pending state.
 * @returns {Object} An object containing the `pending` reactive state and the `wrappedHandler` as `handler`.
 */
export const usePendingHandler = (handler, initialPending = false) => {
  const pending = ref(initialPending);

  const wrappedHandler = async (...args) => {
    try {
      pending.value = true;

      return await handler(...args);
    } catch (err) {
      console.warn(err);

      throw err;
    } finally {
      pending.value = false;
    }
  };

  return {
    pending,
    handler: wrappedHandler,
  };
};
