/**
 * Debounce function
 *
 * @param {function} fn
 * @param {number} wait
 * @param {Object} token
 * @param {boolean} token.cancelled
 * @return {function}
 */
export const debounce = (fn, wait = 0, token = { cancelled: false }) => {
  if (wait === 0) {
    return fn;
  }

  let timeout;

  return (...args) => {
    const later = () => {
      timeout = null;

      // check if the fn call was cancelled.
      if (!token.cancelled) fn(...args);
    };

    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (!timeout) fn(...args);
  };
};
