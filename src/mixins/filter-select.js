import isEmpty from 'lodash/isEmpty';

export default {
  computed: {
    isMainFilterMultiple() {
      return this.userPreference.widget_preferences.isMainFilterMultiple;
    },
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
    updateSelectedFilter(value) {
      this.createUserPreference({
        userPreference: {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            mainFilter: value || {},
          },
        },
      });

      if (value && value.filter) {
        this.query = { ...this.query, filter: value.filter };
      } else {
        this.query = { ...this.query, filter: undefined };
      }
    },

    updateQueryFilter(filter) {
      this.query = { ...this.query, filter };
    },
  },
};
