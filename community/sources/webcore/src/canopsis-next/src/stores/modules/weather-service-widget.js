import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

/**
 * @typedef {Object} WeatherServiceWidgetState
 * @property {Array<Object>} [items] - Weather services list.
 * @property {boolean} [pending] - Whether a fetch request is in progress.
 * @property {Error|null} [error] - Last fetch error, if any.
 */

/**
 * Weather service widget store.
 *
 * Keeps weather services list, pending and error state per widget id.
 */
export const useWeatherServiceWidgetStore = defineStore('weatherServiceWidget', () => {
  /** @type {import('vue').Ref<Record<string, WeatherServiceWidgetState>>} */
  const widgets = ref({});

  /**
   * Returns weather services list for the given widget id.
   *
   * @type {import('vue').ComputedRef<(widgetId: string) => Array<Object>>}
   */
  const getListByWidgetId = computed(() => widgetId => widgets.value[widgetId]?.items ?? []);

  /**
   * Returns pending state for the given widget id.
   *
   * @type {import('vue').ComputedRef<(widgetId: string) => boolean>}
   */
  const getPendingByWidgetId = computed(
    () => widgetId => widgets.value[widgetId]?.pending ?? false,
  );

  /**
   * Returns last fetch error for the given widget id.
   *
   * @type {import('vue').ComputedRef<(widgetId: string) => Error|null>}
   */
  const getErrorByWidgetId = computed(() => widgetId => widgets.value[widgetId]?.error ?? null);

  /**
   * Merge patch into widget state.
   *
   * @param {string} widgetId
   * @param {Partial<WeatherServiceWidgetState>} patch
   */
  const setWidgetState = (widgetId, patch) => {
    widgets.value = {
      ...widgets.value,
      [widgetId]: {
        ...widgets.value[widgetId],
        ...patch,
      },
    };
  };

  /**
   * Fetch weather services list for a widget.
   *
   * @param {Object} [params={}]
   * @param {string} params.widgetId - Widget id used as store key.
   * @param {Object} [params.params={}] - Request query params.
   * @returns {Promise<Object|Array>}
   */
  const fetchList = async ({ widgetId, params = {} } = {}) => {
    try {
      setWidgetState(widgetId, { pending: true, error: null });

      const response = await request.get(API_ROUTES.weatherService, { params });
      const items = response?.data ?? response ?? [];

      setWidgetState(widgetId, { items, pending: false });

      return response;
    } catch (err) {
      console.error(err);

      setWidgetState(widgetId, { error: err, pending: false });

      throw err;
    }
  };

  return {
    widgets,
    getListByWidgetId,
    getPendingByWidgetId,
    getErrorByWidgetId,
    fetchList,
  };
});
