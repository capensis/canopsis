import { ref, unref, provide } from 'vue';

/**
 * Composable for managing focus state
 *
 * @returns {Object} Object containing focus management methods
 */
export const useLastInputFocus = (key) => {
  const lastInputFocus = ref(() => {});

  /**
   * Registers a new focus function
   *
   * @param {Function} focus - The function to set as the last input focus.
   */
  const registerLastInputFocus = focus => lastInputFocus.value = focus;

  provide(unref(key), registerLastInputFocus);

  return {
    lastInputFocus,
    registerLastInputFocus,
  };
};
