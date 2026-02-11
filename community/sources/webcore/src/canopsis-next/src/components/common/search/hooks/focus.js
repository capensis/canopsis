import { unref, provide } from 'vue';

/**
 * Composable for managing focus state
 *
 * @returns {Object} Object containing focus management methods
 */
export const useLastInputFocus = (key) => {
  let lastInputFocuses = [];

  const focusRegister = {
    register: focus => lastInputFocuses.push(focus),
    unregister: focus => lastInputFocuses = lastInputFocuses.filter(f => f !== focus),
    call: () => lastInputFocuses.at(-1)?.(),
  };

  provide(unref(key), focusRegister);

  return {
    focusRegister,
  };
};
