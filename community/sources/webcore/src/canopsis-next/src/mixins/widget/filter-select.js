import { getMainFilter } from '@/helpers/filter';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

export const widgetFilterSelectMixin = {
  mixins: [entitiesWidgetMixin],
  computed: {
    mainFilter() {
      return getMainFilter(this.widget, this.userPreference);
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

    async updateSelectedFilter(mainFilter = null) {
      await this.updateFieldsInWidgetPreferences({ mainFilter, mainFilterUpdatedAt: Date.now() });
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
