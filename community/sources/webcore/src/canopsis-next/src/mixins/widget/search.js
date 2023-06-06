import { MAX_SEARCH_ITEMS } from '@/constants';

import { sortPinnedSearches } from '@/helpers/search/sorting';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetSearchMixin = {
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return this.userPreference?.content?.searches ?? [];
    },
  },
  methods: {
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

    updateSearchesInUserPreferences(search) {
      if (!search) {
        return;
      }

      const searchItem = this.searches.find(item => item.search === search) ?? { search, pinned: false };
      const newSearches = this.searches.filter(item => item.search !== search);

      newSearches.push(searchItem);

      this.updateContentInUserPreference({
        searches: sortPinnedSearches(newSearches, search).slice(0, MAX_SEARCH_ITEMS),
      });
    },

    removeSearchFromUserPreferences(search) {
      this.updateContentInUserPreference({
        searches: this.searches.filter(item => item.search !== search),
      });
    },
  },
};
