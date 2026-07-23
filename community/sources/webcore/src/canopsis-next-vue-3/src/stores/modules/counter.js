import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useCounterStore = defineStore('counterWidget', () => {
  const widgets = ref({});

  const getListByWidgetId = computed(() => (widgetId) => widgets.value[widgetId]?.counters ?? []);
  const getPendingByWidgetId = computed(
    () => (widgetId) => widgets.value[widgetId]?.pending ?? false,
  );

  const setWidgetState = (widgetId, patch) => {
    widgets.value = {
      ...widgets.value,
      [widgetId]: {
        ...widgets.value[widgetId],
        ...patch,
      },
    };
  };

  const fetchList = async ({ widgetId, params = {}, filters = [] } = {}) => {
    try {
      setWidgetState(widgetId, { pending: true });

      const counters = await Promise.all(
        filters.map((filter) =>
          request.get(API_ROUTES.counter, {
            params: {
              ...params,
              filters: [filter],
            },
          }),
        ),
      );

      setWidgetState(widgetId, { counters, pending: false });

      return counters;
    } catch (err) {
      console.error(err);

      setWidgetState(widgetId, { pending: false });

      throw err;
    }
  };

  return {
    widgets,
    getListByWidgetId,
    getPendingByWidgetId,
    fetchList,
  };
});
