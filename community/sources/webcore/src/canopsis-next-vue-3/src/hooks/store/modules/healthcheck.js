import { useHealthcheckStore } from '@/stores/modules/healthcheck';

export const useHealthcheck = () => {
  const healthcheckStore = useHealthcheckStore();

  return {
    fetchHealthcheckStatusWithoutStore: healthcheckStore.fetchHealthcheckStatusWithoutStore,
  };
};
