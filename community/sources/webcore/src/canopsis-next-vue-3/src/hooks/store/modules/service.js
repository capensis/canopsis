import { useServiceStore } from '@/stores/modules/service';

/**
 * Hook for accessing service store actions
 *
 * @returns {Object} Service store actions
 */
export const useService = () => {
  const serviceStore = useServiceStore();

  return {
    ...serviceStore,
  };
};
