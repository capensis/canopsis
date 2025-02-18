import { keyBy } from 'lodash';

import { sortPinnedSearches } from '@/helpers/search/sorting';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetAdvancedSearchSavedItemsMixin = {
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return this.userPreference?.content?.searches ?? [];
    },

    newSearches() { // TODO: rename
      return this.userPreference?.content?.newSearches ?? [];
    },

    newSearchesById() {
      return keyBy(this.newSearches, '_id'); // TODO: rename
    },
  },
  methods: {
    updateNewSearch(search) {
      this.query = {
        ...this.query,

        page: 1,
        alarm_pattern: JSON.stringify(search.alarm_pattern),
        entity_pattern: JSON.stringify(search.entity_pattern),
        pbehavior_pattern: JSON.stringify(search.pbehavior_pattern),
      };

      let newSearches;

      if (this.newSearchesById[search._id]) {
        newSearches = this.newSearches.map(value => (value._id === search._id ? search : value));
      } else {
        newSearches = [...this.newSearches, search];
      }

      this.updateContentInUserPreference({
        newSearches: sortPinnedSearches(newSearches, search._id), // TODO: rename
      });
    },

    updateSearchInQuery(search) {
      this.query = {
        ...this.query,

        search,
        page: 1,
      };
    },

    addSearchIntoUserPreferences(search) {
      if (!search) {
        return;
      }

      const newSearches = [...this.searches, { search, pinned: false }];

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(newSearches, search),
      });
    },

    togglePinSearchInUserPreferences(search) {
      const searchItem = this.searches.find(item => item.search === search);

      if (!searchItem) {
        return;
      }

      const newSearches = this.searches.filter(item => item.search !== search);

      newSearches.push({ ...searchItem, pinned: !searchItem.pinned });

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(newSearches, search),
      });
    },

    removeSearchFromUserPreferences(search) {
      this.updateContentInUserPreference({
        searches: this.searches.filter(item => item.search !== search),
      });
    },
  },
};
