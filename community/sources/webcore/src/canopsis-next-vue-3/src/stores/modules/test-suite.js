import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useTestSuiteStore = defineStore('testSuite', () => {
  const itemsById = ref({});
  const widgets = ref({});
  const histories = ref({});

  const getItem = computed(() => id => itemsById.value[id]);
  const getListByWidgetId = computed(
    () => widgetId =>
      (widgets.value[widgetId]?.allIds ?? []).map(id => itemsById.value[id]).filter(Boolean),
  );
  const getPendingByWidgetId = computed(
    () => widgetId => widgets.value[widgetId]?.pending ?? false,
  );
  const getHistoryListByTestSuiteId = computed(
    () => testSuiteId => histories.value[testSuiteId]?.items ?? [],
  );
  const getHistoryPendingByTestSuiteId = computed(
    () => testSuiteId => histories.value[testSuiteId]?.pending ?? false,
  );
  const getHistoryMetaByTestSuiteId = computed(
    () => testSuiteId => histories.value[testSuiteId]?.meta ?? {},
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

  const setHistoryState = (id, patch) => {
    histories.value = {
      ...histories.value,
      [id]: {
        ...histories.value[id],
        ...patch,
      },
    };
  };

  const setItems = (items = []) => {
    itemsById.value = items.reduce(
      (acc, item) => {
        acc[item._id] = item;

        return acc;
      },
      { ...itemsById.value },
    );
  };

  const validateDirectory = ({ data } = {}) => request.post(API_ROUTES.junit.directory, data);

  const fetchList = async ({ widgetId, params } = {}) => {
    setWidgetState(widgetId, { pending: true, error: null });

    const response = await request.get(`${API_ROUTES.junit.widget}/${widgetId}`, { params });
    const items = response?.data ?? response ?? [];

    setItems(items);
    setWidgetState(widgetId, { allIds: items.map(({ _id }) => _id), pending: false });

    return response;
  };

  const fetchSummaryWithoutStore = ({ id } = {}) =>
    request.get(`${API_ROUTES.junit.testSuites}/${id}/summary`);
  const fetchItemGanttIntervalsWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.junit.testSuites}/${id}/gantt`, { params });
  const fetchItemDetailsWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.junit.testSuites}/${id}/details`, { params });
  const fetchEntityGanttIntervalsWithoutStore = ({ id, params } = {}) =>
    request.get(API_ROUTES.junit.entityGantt, {
      params: { ...params, _id: id },
    });

  const fetchHistoryList = async ({ id, params = {} } = {}) => {
    try {
      setHistoryState(id, { pending: true });

      const response = await request.get(`${API_ROUTES.junit.history}/${id}`, { params });

      setHistoryState(id, {
        items: response?.data ?? response ?? [],
        meta: response?.meta ?? {},
        pending: false,
      });

      return response;
    } catch (err) {
      setHistoryState(id, { pending: false });

      throw err;
    }
  };

  return {
    itemsById,
    widgets,
    histories,
    getItem,
    getListByWidgetId,
    getPendingByWidgetId,
    getHistoryListByTestSuiteId,
    getHistoryPendingByTestSuiteId,
    getHistoryMetaByTestSuiteId,
    validateDirectory,
    fetchList,
    fetchSummaryWithoutStore,
    fetchItemGanttIntervalsWithoutStore,
    fetchItemDetailsWithoutStore,
    fetchEntityGanttIntervalsWithoutStore,
    fetchHistoryList,
  };
});
