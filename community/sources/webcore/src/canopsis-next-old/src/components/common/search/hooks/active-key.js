import { ref } from 'vue';

/**
 * Composable for managing active key state
 * @returns {Object} Object containing active key state and methods
 */
export const useActiveKey = () => {
  const activeKey = ref();

  /**
   * Sets the active key to the specified key.
   *
   * @param {string} key - The key to set as active.
   */
  const makeActive = key => activeKey.value = key;

  /**
   * Resets the active key if it matches the specified key.
   *
   * @param {string} key - The key to check against the active key.
   */
  const resetActive = (key) => {
    if (key === activeKey.value) {
      activeKey.value = null;
    }
  };

  return {
    activeKey,
    makeActive,
    resetActive,
  };
};
