import { omit, keyBy } from 'lodash';

import {
  ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS,
  PATTERN_CONDITIONS,
  ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES,
} from '@/constants';

import { uuid } from '@/helpers/uuid';
import { sortPinnedSearches } from '@/helpers/search/sorting';
import { advancedSearchToForm, isEmptyAlarmSearch, isEqualAlarmSearches } from '@/helpers/search/advanced-search';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetAlarmAdvancedSearchSavedItemsMixin = {
  provide() {
    return {
      $selectAdvancedSearchField: (field, value) => {
        const patternField = ALARM_ADVANCED_SEARCH_FIELDS_TO_PATTERNS[field];
        const preparedField = field.replace(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.entity, '')
          .replace(ALARM_ADVANCED_SEARCH_PATTERNS_PREFIXES.pbehavior, '');

        const pattern = [[{ field: preparedField, cond: { value, type: PATTERN_CONDITIONS.equal } }]];
        const search = {
          _id: uuid(),
          pinned: false,
          rules: advancedSearchToForm({
            positions: [patternField],
            [patternField]: pattern,
          }),
        };

        this.$refs.alarmAdvancedSearch.advancedSearchElement.select(search);
      },
    };
  },
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
        search: search.search,
        alarm_pattern: JSON.stringify(search.alarm_pattern),
        entity_pattern: JSON.stringify(search.entity_pattern),
        pbehavior_pattern: JSON.stringify(search.pbehavior_pattern),
      };

      let searches;
      let found = false;

      searches = this.searches.map((value) => {
        if (isEqualAlarmSearches(value, search)) {
          found = true;

          return {
            ...search,

            pinned: value.pinned || search.pinned,
          };
        }

        return value;
      });

      if (!found) {
        searches = [search, ...this.searches];
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
