import { useComponentInstance } from './vue';

/**
 * Hook to access popup management methods and properties from the Vue instance.
 *
 * @returns {Object} An object containing popup management methods and properties.
 */
export const usePopups = () => {
  const vm = useComponentInstance();

  return vm.$popups;
};
