import { useHealthcheckParametersStore } from '@/stores/modules/healthcheck-parameters';

/**
 * Healthcheck parameters store facade.
 *
 * @returns {Object}
 */
export const useHealthcheckParameters = () => {
  const healthcheckParametersStore = useHealthcheckParametersStore();

  return {
    fetchHealthcheckParametersWithoutStore: healthcheckParametersStore.fetchItemWithoutStore,
    updateHealthcheckParameters: healthcheckParametersStore.update,
  };
};
