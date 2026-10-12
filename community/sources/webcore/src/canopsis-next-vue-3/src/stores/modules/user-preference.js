import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useUserPreferenceStore = defineStore('userPreference', () => {
  const itemsByWidget = ref({});
  const pending = ref(false);

  const items = computed(() => Object.values(itemsByWidget.value));
  const getItemByWidgetId = computed(
    () => widgetId =>
      itemsByWidget.value[widgetId] ?? {
        widget: widgetId,
        content: {},
      },
  );

  const setPreference = (preference = {}) => {
    if (!preference.widget) {
      return;
    }

    itemsByWidget.value = {
      ...itemsByWidget.value,
      [preference.widget]: preference,
    };
  };

  const fetchItem = async ({ id } = {}) => {
    try {
      pending.value = true;

      const preference = await request.get(`${API_ROUTES.userPreferences}/${id}`);

      setPreference(preference);

      return preference;
    } catch (err) {
      console.warn(err);

      return undefined;
    } finally {
      pending.value = false;
    }
  };

  const fetchItemWithoutStore = ({ id } = {}) => request.get(`${API_ROUTES.userPreferences}/${id}`);

  const update = async ({ data } = {}) => {
    try {
      const preference = await request.put(API_ROUTES.userPreferences, data);

      setPreference(preference ?? data);

      return preference;
    } catch (err) {
      console.warn(err);

      return undefined;
    }
  };

  const updateLocal = ({ data } = {}) => setPreference(data);

  return {
    itemsByWidget,
    items,
    pending,
    getItemByWidgetId,
    fetchItem,
    fetchItemWithoutStore,
    update,
    updateLocal,
  };
});
