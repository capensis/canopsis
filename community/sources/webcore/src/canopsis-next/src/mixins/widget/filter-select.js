import { FILTER_DEFAULT_VALUES } from '@/constants';

import { prepareMainFilterToQueryFilter, getMainFilterAndCondition } from '@/helpers/filter';

export const widgetFilterSelectMixin = {
  computed: {
    mainFilterAndCondition() {
      return getMainFilterAndCondition(this.widget, this.userPreference);
    },

    mainFilterCondition() {
      const { condition } = this.mainFilterAndCondition;

      return condition;
    },

    mainFilter() {
      const { mainFilter } = this.mainFilterAndCondition;

      return mainFilter;
    },

    viewFilters() {
      return this.userPreference.content.viewFilters || [];
    },

    widgetViewFilters() {
      const { mainFilter, viewFilters } = this.widget.parameters;

      if (!this.hasAccessToListFilters) {
        return mainFilter ? [mainFilter] : [];
      }

      return viewFilters || [];
    },
  },
  methods: {
    updateFieldsInWidgetPreferences(fields = {}) {
      if (this.hasAccessToUserFilter) {
        return this.updateContentInUserPreference({
          ...fields,
        });
      }

      return Promise.resolve();
    },

    async updateFilters(viewFilters, mainFilter = this.mainFilter) {
      await this.updateFieldsInWidgetPreferences({ viewFilters, mainFilter });
      this.updateQueryBySelectedFilterAndCondition(mainFilter, this.mainFilterCondition);
    },

    async updateSelectedCondition(condition = FILTER_DEFAULT_VALUES.condition) {
      await this.updateFieldsInWidgetPreferences({ mainFilterCondition: condition });
      this.updateQueryBySelectedFilterAndCondition(this.mainFilter, condition);
    },

    async updateSelectedFilter(filterObject) {
      await this.updateFieldsInWidgetPreferences({ mainFilter: filterObject || {}, mainFilterUpdatedAt: Date.now() });
      this.updateQueryBySelectedFilterAndCondition(filterObject || {}, this.mainFilterCondition);
    },

    updateQueryBySelectedFilterAndCondition(filter, condition) {
      this.query = {
        ...this.query,

        page: 1,
        filter: prepareMainFilterToQueryFilter(filter, condition),
      };
    },
  },
};
