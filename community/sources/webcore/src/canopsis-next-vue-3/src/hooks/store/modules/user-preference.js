import { storeToRefs } from 'pinia';
import { computed, unref } from 'vue';

import { useUserPreferenceStore } from '@/stores/modules/user-preference';

import { setField } from '@/helpers/immutable';

/**
 * User preference store facade.
 *
 * @returns {Object}
 */
export const useUserPreference = () => {
  const userPreferenceStore = useUserPreferenceStore();
  const { items, itemsByWidget, pending, getItemByWidgetId } = storeToRefs(userPreferenceStore);

  return {
    userPreferences: items,
    userPreferencesByWidget: itemsByWidget,
    userPreferencesPending: pending,
    getUserPreferenceByWidgetId: getItemByWidgetId,

    fetchUserPreference: userPreferenceStore.fetchItem,
    fetchUserPreferenceWithoutStore: userPreferenceStore.fetchItemWithoutStore,
    updateUserPreference: userPreferenceStore.update,
    updateLocalUserPreference: userPreferenceStore.updateLocal,
  };
};

/**
 * Hook that provides user preference functionality for a specific widget.
 *
 * @param {Object} options - Options object
 * @param {string|import('vue').Ref<string>} options.widgetId - ID of the widget to get preferences for
 * @param {boolean|import('vue').Ref<boolean>} [options.localWidget=false] - When true, updates are stored locally only
 * @returns {Object} Object containing widget-specific user preference data and methods
 */
export const useWidgetUserPreference = ({ widgetId, localWidget = false }) => {
  const {
    getUserPreferenceByWidgetId,

    fetchUserPreference: fetchUserPreferenceItem,
    updateUserPreference,
    updateLocalUserPreference,
  } = useUserPreference();

  const userPreference = computed(() => getUserPreferenceByWidgetId.value(unref(widgetId)));

  const fetchUserPreference = () => {
    if (unref(localWidget)) {
      return Promise.resolve();
    }

    return fetchUserPreferenceItem({ id: unref(widgetId) });
  };

  const updateContentInUserPreference = (content = {}) => {
    const method = unref(localWidget) ? updateLocalUserPreference : updateUserPreference;

    return method({
      data: setField(userPreference.value, 'content', prevContent => ({ ...prevContent, ...content })),
    });
  };

  return {
    userPreference,

    fetchUserPreference,
    updateContentInUserPreference,
  };
};
