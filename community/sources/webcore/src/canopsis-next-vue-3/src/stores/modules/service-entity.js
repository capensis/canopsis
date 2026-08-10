import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useServiceEntityStore = defineStore('serviceEntity', () => {
  const services = ref({});

  const getListByServiceId = computed(() => serviceId => services.value[serviceId]?.items ?? []);
  const getPendingByServiceId = computed(
    () => serviceId => services.value[serviceId]?.pending ?? false,
  );
  const getMetaByServiceId = computed(() => serviceId => services.value[serviceId]?.meta ?? {});

  const setServiceState = (id, patch) => {
    services.value = {
      ...services.value,
      [id]: {
        ...services.value[id],
        ...patch,
      },
    };
  };

  const fetchList = async ({ id, params = {} } = {}) => {
    try {
      setServiceState(id, { pending: true });

      const response = await request.get(`${API_ROUTES.weatherService}/${id}`, { params });

      setServiceState(id, {
        items: response?.data ?? response ?? [],
        meta: response?.meta ?? {},
        pending: false,
      });

      return response;
    } catch (err) {
      console.error(err);

      setServiceState(id, { error: err, pending: false });

      throw err;
    }
  };

  return {
    services,
    getListByServiceId,
    getPendingByServiceId,
    getMetaByServiceId,
    fetchList,
  };
});
