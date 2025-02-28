import { omit, keyBy } from 'lodash';

import { sortPinnedSearches } from '@/helpers/search/sorting';
import { isEmptyAlarmSearch } from '@/helpers/search/alarm-advanced-search';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetAlarmAdvancedSearchSavedItemsMixin = {
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return this.userPreference?.content?.searches ?? [];
    },

    searchesById() {
      return keyBy(this.searches, '_id');
    },
  },
  methods: {
    updateSearch(search) {
      if (isEmptyAlarmSearch(search)) {
        this.resetSearch();

        return;
      }

      this.query = {
        ...this.query,

        page: 1,
        search: search.text,
        alarm_pattern: JSON.stringify(search.alarm_pattern),
        entity_pattern: JSON.stringify(search.entity_pattern),
        pbehavior_pattern: JSON.stringify(search.pbehavior_pattern),
      };

      let searches;

      if (this.searchesById[search._id]) {
        searches = this.searches.map(value => (value._id === search._id ? search : value));
      } else {
        searches = [...this.searches, search];
      }

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(searches, search._id, '_id'),
      });
    },

    resetSearch() {
      if (isEmptyAlarmSearch(this.query)) {
        return;
      }

      this.query = {
        ...omit(this.query, ['search', 'alarm_pattern', 'entity_pattern', 'pbehavior_pattern']),

        page: 1,
      };
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
