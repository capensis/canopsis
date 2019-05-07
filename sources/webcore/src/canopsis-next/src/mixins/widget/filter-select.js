import { isEmpty, isBoolean } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import prepareMainFilterToQueryFilter from '@/helpers/filter';

export default {
  computed: {
    mainFilterCondition() {
      return this.userPreference.widget_preferences.mainFilterCondition || this.widget.parameters.mainFilterCondition;
    },

    mainFilter() {
      const mainFilter = this.userPreference.widget_preferences.mainFilter || this.widget.parameters.mainFilter;

      if (Array.isArray(mainFilter)) {
        return mainFilter;
      }

      return isEmpty(mainFilter) ? null : mainFilter;
    },

    viewFilters() {
      return this.userPreference.widget_preferences.viewFilters || [];
    },

    widgetViewFilters() {
      return this.widget.parameters.viewFilters || [];
    },
  },
  methods: {
    updateFieldsInWidgetPreferences(fields = {}) {
      const hasAccessToEditFilter = this.hasAccessToEditFilter || !isBoolean(this.hasAccessToEditFilter);

      if (hasAccessToEditFilter) {
        return this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,
          ...fields,
        });
      }

      return Promise.resolve();
    },

    updateFilters(viewFilters, mainFilter = this.mainFilter) {
      this.updateFieldsInWidgetPreferences({ viewFilters, mainFilter });
      this.updateQueryBySelectedFilterAndCondition(mainFilter, this.mainFilterCondition);
    },

    updateSelectedCondition(condition = FILTER_DEFAULT_VALUES.condition) {
      this.updateFieldsInWidgetPreferences({ mainFilterCondition: condition });
      this.updateQueryBySelectedFilterAndCondition(this.mainFilter, condition);
    },

    updateSelectedFilter(filterObject) {
      this.updateFieldsInWidgetPreferences({ mainFilter: filterObject || {} });
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
