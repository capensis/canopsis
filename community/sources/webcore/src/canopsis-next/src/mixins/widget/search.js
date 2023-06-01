import { MAX_SEARCH_ITEMS } from '@/constants';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetSearchMixin = {
  mixins: [entitiesUserPreferenceMixin],
  computed: {
    searches() {
      return this.userPreference?.content?.searches ?? [];
    },
  },
  methods: {
    updateSearchesInUserPreferences(search) {
      const newSearches = [
        search,
        ...this.searches.filter(item => item !== search),
      ].slice(0, MAX_SEARCH_ITEMS);

      return this.updateContentInUserPreference({
        searches: newSearches,
      });
    },
  },
};
