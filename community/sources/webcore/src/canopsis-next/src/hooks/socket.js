import { useComponentInstance } from './vue';

/**
 * Hook to access socket management methods and properties from the Vue instance.
 *
 * @returns {Object} An object containing socket management methods and properties.
 */
export const useSocket = () => {
  const vm = useComponentInstance();

  return vm.$socket;
};
