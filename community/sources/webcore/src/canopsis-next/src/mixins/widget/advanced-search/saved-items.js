import { keyBy } from 'lodash';

import {
  createAdvancedSearchFromAlarmFieldValue,
  mergeSearchIntoSavedSearches,
  prepareQueryWithAdvancedSearch,
  prepareQueryWithoutAdvancedSearch,
  isEmptyAdvancedSearch,
} from '@/helpers/search/advanced-search';
import { sortPinnedSearches } from '@/helpers/search/sorting';
import { uuid } from '@/helpers/uuid';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetAdvancedSearchSavedItemsMixin = {
  provide() {
    return {
      $registerSelectAdvancedSearch: selectFunc => this.$selectAdvancedSearch = selectFunc,

      /**
       * Creates a search from a field/value pair and passes it to the registered select function.
       * Used by ONLY alarm column cells to apply a filter from a chip click.
       *
       * @param {string} field - The advanced search field name (e.g. entity.component)
       * @param {*} value - The value to filter by
       */
      $selectAdvancedSearchField: (field, value) => {
        const search = createAdvancedSearchFromAlarmFieldValue(field, value);

        this.$selectAdvancedSearch?.(search);
      },
    };
  },
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return (this.userPreference?.content?.searches ?? []).map(search => (search?._id ? search : {
        ...search,

        _id: uuid(),
      }));
    },

    searchesById() {
      return keyBy(this.searches, '_id');
    },
  },
  methods: {
    updateSearch(search) {
      if (isEmptyAdvancedSearch(search)) {
        this.resetSearch();

        return;
      }
      this.query = prepareQueryWithAdvancedSearch(this.query, search);

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(mergeSearchIntoSavedSearches(this.searches, search), search._id, '_id'),
      });
    },

    resetSearch() {
      if (isEmptyAdvancedSearch(this.query)) {
        return;
      }

      this.query = prepareQueryWithoutAdvancedSearch(this.query);
    },

    togglePinSearch(id) {
      const searches = this.searches.map(search => (
        search._id === id ? { ...search, pinned: !search.pinned } : search
      ));

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(searches, id, '_id'),
      });
    },

    removeSearch(id) {
      this.updateContentInUserPreference({
        searches: this.searches.filter(search => search._id !== id),
      });
    },
  },
};
