import { useComponentInstance } from './vue';

/**
 * Hook to access sidebar management methods and properties from the Vue instance.
 *
 * @returns {Object} An object containing sidebar management methods and properties.
 */
export const useSidebar = () => {
  const vm = useComponentInstance();

  return vm.$sidebar;
};
