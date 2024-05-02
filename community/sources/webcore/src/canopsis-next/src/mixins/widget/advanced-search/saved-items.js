import { sortPinnedSearches } from '@/helpers/search/sorting';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetAdvancedSearchSavedItemsMixin = {
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return this.userPreference?.content?.searches ?? [];
    },
  },
  methods: {
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
