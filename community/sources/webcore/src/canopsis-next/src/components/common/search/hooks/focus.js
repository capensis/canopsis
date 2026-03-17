import { unref, provide } from 'vue';

/**
 * Composable for managing focus state
 *
 * @returns {Object} Object containing focus management methods
 */
export const useLastInputFocus = (key) => {
  let inputFocuses = [];

  const focusRegister = {
    register: focus => inputFocuses.push(focus),
    unregister: focus => inputFocuses = inputFocuses.filter(f => f !== focus),
    call: () => inputFocuses.at(-1)?.(),
  };

  provide(unref(key), focusRegister);

  return {
    focusRegister,
  };
};
