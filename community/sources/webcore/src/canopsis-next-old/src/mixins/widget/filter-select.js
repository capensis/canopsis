import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

export const widgetFilterSelectMixin = {
  mixins: [entitiesWidgetMixin],
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
