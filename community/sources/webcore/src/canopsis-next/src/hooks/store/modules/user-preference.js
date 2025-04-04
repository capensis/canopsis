import { computed, unref } from 'vue';

import { setField } from '@/helpers/immutable';

import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for interacting with the userPreference store module
 *
 * @returns {Object} Object containing store module hooks
 */
export const useUserPreferenceStoreModule = () => useStoreModuleHooks('userPreference');

/**
 * Hook that provides access to user preference getters and actions
 *
 * @returns {Object} Object containing user preference getters and actions
 * @property {Object} getUserPreferenceByWidgetId - Getter for retrieving user preference by widget ID
 * @property {Function} fetchUserPreferenceItem - Action to fetch a user preference item
 * @property {Function} fetchUserPreferenceWithoutStore - Action to fetch user preference without storing it
 * @property {Function} updateUserPreference - Action to update user preference
 */
export const useUserPreference = () => {
  const { useGetters, useActions } = useUserPreferenceStoreModule();

  const getters = useGetters({
    getUserPreferenceByWidgetId: 'getItemByWidgetId',
  });

  const actions = useActions({
    fetchUserPreferenceItem: 'fetchItem',
    fetchUserPreferenceWithoutStore: 'fetchItemWithoutStore',
    updateUserPreference: 'update',
  });

  return {
    ...getters,
    ...actions,
  };
};

/**
 * Hook that provides user preference functionality for a specific widget
 *
 * @param {Object} options - Options object
 * @param {string|Ref<string>} options.widgetId - ID of the widget to get preferences for
 * @returns {Object} Object containing widget-specific user preference data and methods
 * @property {ComputedRef<Object>} userPreference - Computed reference to the user preference for the widget
 * @property {Function} fetchUserPreference - Function to fetch user preference for the widget
 * @property {Function} updateContentInUserPreference - Function to update content in the user preference
 */
export const useWidgetUserPreference = ({ widgetId }) => {
  const {
    getUserPreferenceByWidgetId,

    fetchUserPreferenceItem,
    updateUserPreference,
  } = useUserPreference();

  const userPreference = computed(() => getUserPreferenceByWidgetId.value(unref(widgetId)));

  const fetchUserPreference = () => fetchUserPreferenceItem({ id: unref(widgetId) });

  const updateContentInUserPreference = (content = {}) => updateUserPreference({
    data: setField(userPreference.value, 'content', prevContent => ({ ...prevContent, ...content })),
  });

  return {
    userPreference,

    fetchUserPreference,
    updateContentInUserPreference,
  };
};
