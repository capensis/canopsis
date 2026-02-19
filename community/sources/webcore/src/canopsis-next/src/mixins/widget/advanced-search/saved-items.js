import { keyBy } from 'lodash';

import {
  createAdvancedSearchFromFieldValue,
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
      $selectAdvancedSearchField: (field, value) => {
        const search = createAdvancedSearchFromFieldValue(field, value);

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
