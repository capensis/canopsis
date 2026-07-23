import { useHealthcheckStore } from '@/stores/modules/healthcheck';

export const useHealthcheck = () => {
  const healthcheckStore = useHealthcheckStore();

  return {
    fetchHealthcheckEnginesWithoutStore: healthcheckStore.fetchHealthcheckEnginesWithoutStore,
    fetchHealthcheckStatusWithoutStore: healthcheckStore.fetchHealthcheckStatusWithoutStore,
    fetchHealthcheckEnginesOrderWithoutStore:
      healthcheckStore.fetchHealthcheckEnginesOrderWithoutStore,
  };
};
