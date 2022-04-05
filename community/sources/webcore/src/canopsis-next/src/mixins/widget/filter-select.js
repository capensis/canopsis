import { FILTER_DEFAULT_VALUES } from '@/constants';

import { prepareMainFilterToQueryFilter, getMainFilterAndCondition } from '@/helpers/filter';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

export const widgetFilterSelectMixin = {
  mixins: [entitiesWidgetMixin],
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

    // TODO: remove
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

    async updateSelectedCondition(condition = FILTER_DEFAULT_VALUES.condition) {
      await this.updateFieldsInWidgetPreferences({ mainFilterCondition: condition });
      this.updateQueryBySelectedFilterAndCondition(this.mainFilter, condition);
    },

    // TODO: remove
    async updateSelectedFilter(mainFilter = null) {
      await this.updateFieldsInWidgetPreferences({ mainFilter, mainFilterUpdatedAt: Date.now() });
      this.updateQueryBySelectedFilterAndCondition(mainFilter, this.mainFilterCondition);
    },

    updateQueryBySelectedFilterAndCondition(filter, condition) {
      this.query = {
        ...this.query,
        ...prepareMainFilterToQueryFilter(filter, condition),

        page: 1,
      };
    },
  },
};
