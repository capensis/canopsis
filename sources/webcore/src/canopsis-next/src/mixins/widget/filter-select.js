import isEmpty from 'lodash/isEmpty';
import isBoolean from 'lodash/isBoolean';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import { prepareMainFilterToQueryFilter } from '@/helpers/filter/index';

export default {
  computed: {
    mainFilterCondition() {
      return this.userPreference.widget_preferences.mainFilterCondition;
    },

    mainFilter() {
      const mainFilter = this.userPreference.widget_preferences.mainFilter || this.widget.parameters.mainFilter;

      if (Array.isArray(mainFilter)) {
        return mainFilter;
      }

      return isEmpty(mainFilter) ? null : mainFilter;
    },

    viewFilters() {
      const viewFilters = this.userPreference.widget_preferences.viewFilters || this.widget.parameters.viewFilters;

      return isEmpty(viewFilters) ? [] : viewFilters;
    },
  },
  methods: {
    updateFilterFieldInWidgetPreferences(field, value) {
      const hasAccessToEditFilter = this.hasAccessToEditFilter || !isBoolean(this.hasAccessToEditFilter);

      if (hasAccessToEditFilter) {
        return this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,

          [field]: value,
        });
      }

      return Promise.resolve();
    },

    updateSelectedCondition(condition = FILTER_DEFAULT_VALUES.condition) {
      this.updateFilterFieldInWidgetPreferences('mainFilterCondition', condition);
      this.updateQueryBySelectedFilterAndCondition(this.mainFilter, condition);
    },

    updateSelectedFilter(filterObject) {
      this.updateFilterFieldInWidgetPreferences('mainFilter', filterObject || {});
      this.updateQueryBySelectedFilterAndCondition(filterObject, this.mainFilterCondition);
    },

    updateQueryBySelectedFilterAndCondition(filter, condition) {
      this.query = {
        ...this.query,

        filter: prepareMainFilterToQueryFilter(filter, condition),
      };
    },
  },
};
