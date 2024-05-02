import { useComponentInstance } from './vue';

/**
 * Hook to access modal management methods and properties from the Vue instance.
 *
 * @returns {Object} An object containing modal management methods and properties.
 */
export const useModals = () => {
  const vm = useComponentInstance();

  return vm.$modals;
};
