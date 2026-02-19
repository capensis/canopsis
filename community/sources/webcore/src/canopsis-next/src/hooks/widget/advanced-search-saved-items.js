import { computed, provide, unref } from 'vue';
import { keyBy } from 'lodash';

import { sortPinnedSearches } from '@/helpers/search/sorting';
import {
  createAdvancedSearchFromFieldValue,
  mergeSearchIntoSavedSearches,
  prepareQueryWithAdvancedSearch,
  prepareQueryWithoutAdvancedSearch,
  isEmptyAdvancedSearch,
} from '@/helpers/search/advanced-search';
import { uuid } from '@/helpers/uuid';

import { useWidgetUserPreference } from '@/hooks/store/modules/user-preference';

/**
 * Hook that provides advanced search saved items functionality for widgets.
 * Implements the same logic as widgetAdvancedSearchSavedItemsMixin.
 *
 * @param {Object} options - Options object
 * @param {string|Ref<string>} options.widgetId - ID of the widget to get preferences for
 * @param {Object|Ref<Object>} options.query - Current query object (required for updateSearch/resetSearch)
 * @param {boolean} [options.localWidget=false] - When true, user preference updates are stored locally only
 * @param {Function} emit - Emit function to send query updates (e.g. emit('update:query', query))
 * @returns {Object} Object containing search-related data and methods
 * @property {ComputedRef<Array>} searches - Computed reference to the saved searches array (with _id ensured)
 * @property {ComputedRef<Object>} searchesById - Computed reference to searches keyed by _id
 * @property {Function} updateSearch - Function to update search and persist to user preference
 * @property {Function} resetSearch - Function to clear search from query
 * @property {Function} togglePinSearch - Function to toggle pin status of a search by id
 * @property {Function} removeSearch - Function to remove a search by id
 * @property {Function} registerSelectAdvancedSearch - Provider for $registerSelectAdvancedSearch
 * @property {Function} selectAdvancedSearchField - Provider for $selectAdvancedSearchField
 */
export const useWidgetAdvancedSearchSavedItems = (
  { widgetId, query, localWidget = false },
  emit,
) => {
  let selectAdvancedSearch;

  const { userPreference, updateContentInUserPreference } = useWidgetUserPreference({
    widgetId,
    localWidget,
  });

  const searches = computed(() => (
    (userPreference.value?.content?.searches ?? []).map(search => (search?._id ? search : { ...search, _id: uuid() }))
  ));

  const searchesById = computed(() => keyBy(searches.value, '_id'));

  /**
   * Emits the updated query to the parent component.
   *
   * @param {Object} newQuery - The new query object to apply
   */
  const setQuery = newQuery => emit?.('update:query', newQuery);

  /**
   * Clears all advanced search-related fields from the query (search, alarm, entity, pbehavior)
   * and resets the page to 1.
   */
  const clearQuery = () => setQuery(prepareQueryWithoutAdvancedSearch(unref(query)));

  /**
   * Updates the current search in the query and persists it to user preferences.
   * If the search is empty, calls resetSearch instead.
   *
   * @param {Object} search - Advanced search object with search, alarm, entity, pbehavior fields
   */
  const updateSearch = (search) => {
    if (isEmptyAdvancedSearch(search)) {
      clearQuery();

      return;
    }

    setQuery(prepareQueryWithAdvancedSearch(unref(query), search));

    updateContentInUserPreference({
      searches: sortPinnedSearches(mergeSearchIntoSavedSearches(searches.value, search), search._id, '_id'),
    });
  };

  /**
   * Clears search-related fields from the query and emits the updated query.
   */
  const resetSearch = () => {
    if (isEmptyAdvancedSearch(unref(query))) {
      return;
    }

    clearQuery();
  };

  /**
   * Toggles the pinned status of a saved search by id.
   *
   * @param {string} id - The _id of the search to toggle
   */
  const togglePinSearch = (id) => {
    const updatedSearches = searches.value.map(search => (
      search._id === id ? { ...search, pinned: !search.pinned } : search
    ));

    updateContentInUserPreference({
      searches: sortPinnedSearches(updatedSearches, id, '_id'),
    });
  };

  /**
   * Removes a saved search from user preferences by id.
   *
   * @param {string} id - The _id of the search to remove
   */
  const removeSearch = id => updateContentInUserPreference({
    searches: searches.value.filter(search => search._id !== id),
  });

  /**
   * Registers the select function from c-advanced-search-field for programmatic search selection.
   *
   * @param {Function|null} selectFunc - Function that receives a search object, or null to unregister
   */
  const registerSelectAdvancedSearch = selectFunc => selectAdvancedSearch = selectFunc;

  /**
   * Creates a search from a field/value pair and passes it to the registered select function.
   * Used by alarm column cells to apply a filter from a chip click.
   *
   * @param {string} field - The advanced search field name (e.g. entity.component)
   * @param {*} value - The value to filter by
   */
  const selectAdvancedSearchField = (field, value) => {
    const search = createAdvancedSearchFromFieldValue(field, value);

    selectAdvancedSearch?.(search);
  };

  provide('$registerSelectAdvancedSearch', registerSelectAdvancedSearch);
  provide('$selectAdvancedSearchField', selectAdvancedSearchField);

  return {
    searches,
    searchesById,
    updateSearch,
    resetSearch,
    togglePinSearch,
    removeSearch,
    registerSelectAdvancedSearch,
    selectAdvancedSearchField,
    addSearchIntoUserPreferences: updateSearch,
    togglePinSearchInUserPreferences: togglePinSearch,
    removeSearchFromUserPreferences: removeSearch,
  };
};
