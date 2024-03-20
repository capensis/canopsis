import { mapIds } from '@/helpers/array';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

export const widgetFilterSelectMixin = {
  mixins: [entitiesWidgetMixin],
  computed: {
    widgetFilters() {
      return this.widget.filters ?? [];
    },

    widgetFiltersIds() {
      return mapIds(this.widgetFilters);
    },

    userPreferencesFilters() {
      return this.userPreference.filters ?? [];
    },

    userPreferencesFiltersIds() {
      return mapIds(this.userPreferencesFilters);
    },

    hasMainFilter() {
      const { mainFilter } = this.userPreference.content;

      return this.widgetFiltersIds.includes(mainFilter)
        || this.userPreferencesFiltersIds.includes(mainFilter);
    },

    mainFilter() {
      const { mainFilter } = this.userPreference.content;

      return mainFilter && this.hasMainFilter
        ? mainFilter
        : undefined;
    },

    lockedFilter() {
      return this.widget.parameters.mainFilter;
    },
  },
  methods: {
    updateFieldsInWidgetPreferences(fields = {}) {
      if (this.hasAccessToUserFilter) {
        return this.updateContentInUserPreference(fields);
      }

      return Promise.resolve();
    },

    async updateSelectedFilter(mainFilter = null) {
      await this.updateFieldsInWidgetPreferences({ mainFilter });
      this.updateQueryBySelectedFilter(mainFilter);
    },

    updateQueryBySelectedFilter(filter) {
      this.query = {
        ...this.query,

        filter,
        page: 1,
      };
    },
  },
};
