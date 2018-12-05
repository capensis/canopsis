import isEmpty from 'lodash/isEmpty';
import isBoolean from 'lodash/isBoolean';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import { prepareMainFilterToQueryFilter } from '@/helpers/filter';

export default {
  computed: {
    mainFilterCondition() {
      return this.userPreference.widget_preferences.mainFilterCondition;
    },

    mainFilter() {
      const mainFilter = this.userPreference.widget_preferences.mainFilter || this.widget.parameters.mainFilter;

      return isEmpty(mainFilter) ? null : mainFilter;
    },

    viewFilters() {
      const viewFilters = this.userPreference.widget_preferences.viewFilters || this.widget.parameters.viewFilters;

      return isEmpty(viewFilters) ? [] : viewFilters;
    },
  },
  methods: {
    updateSelectedCondition(condition = FILTER_DEFAULT_VALUES.condition) {
      if (this.hasAccessToAddFilter || !isBoolean(this.hasAccessToAddFilter)) {
        this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,

          mainFilterCondition: condition,
        });
      }

      this.query = { ...this.query, filter: prepareMainFilterToQueryFilter(this.mainFilter, condition) };
    },

    updateSelectedFilter(filterObject) {
      if (this.hasAccessToAddFilter || !isBoolean(this.hasAccessToAddFilter)) {
        this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,

          mainFilter: filterObject || {},
        });
      }

      this.query = { ...this.query, filter: prepareMainFilterToQueryFilter(filterObject, this.mainFilterCondition) };
    },
  },
};
