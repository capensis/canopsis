import { computed } from 'vue';

import { sortPinnedSearches } from '@/helpers/search/sorting';

import { useWidgetUserPreference } from '@/hooks/store/modules/user-preference';

/**
 * Hook that provides advanced search saved items functionality for widgets
 *
 * @param {Object} options - Options object
 * @param {string|Ref<string>} options.widgetId - ID of the widget to get preferences for
 * @returns {Object} Object containing search-related data and methods
 * @property {ComputedRef<Array>} searches - Computed reference to the saved searches array
 * @property {Function} addSearchIntoUserPreferences - Function to add a search to user preferences
 * @property {Function} togglePinSearchInUserPreferences - Function to toggle pin status of a search
 * @property {Function} removeSearchFromUserPreferences - Function to remove a search from user preferences
 */
export const useWidgetAdvancedSearchSavedItems = ({ widgetId }) => {
  const { userPreference, updateContentInUserPreference } = useWidgetUserPreference({ widgetId });

  const searches = computed(() => userPreference.value?.content?.searches ?? []);

  /**
   * Adds a search into user preferences
   *
   * @param {string} search - Search string to add
   */
  const addSearchIntoUserPreferences = (search) => {
    if (!search) {
      return;
    }

    const newSearches = [...searches.value, { search, pinned: false }];

    updateContentInUserPreference({
      searches: sortPinnedSearches(newSearches, search),
    });
  };

  /**
   * Toggles pin status of a search in user preferences
   *
   * @param {string} search - Search string to toggle pin for
   */
  const togglePinSearchInUserPreferences = (search) => {
    const searchItem = searches.value.find(item => item.search === search);

    if (!searchItem) {
      return;
    }

    const newSearches = searches.value.filter(item => item.search !== search);

    newSearches.push({ ...searchItem, pinned: !searchItem.pinned });

    updateContentInUserPreference({
      searches: sortPinnedSearches(newSearches, search),
    });
  };

  /**
   * Removes a search from user preferences
   *
   * @param {string} search - Search string to remove
   */
  const removeSearchFromUserPreferences = (search) => {
    updateContentInUserPreference({
      searches: searches.value.filter(item => item.search !== search),
    });
  };

  return {
    searches,

    addSearchIntoUserPreferences,
    togglePinSearchInUserPreferences,
    removeSearchFromUserPreferences,
  };
};
