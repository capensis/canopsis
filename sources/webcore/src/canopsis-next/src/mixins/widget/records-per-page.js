export default {
  methods: {
    updateRecordsPerPage(limit) {
      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        itemsPerPage: limit,
      });

      this.query = { ...this.query, limit };
    },
  },
};
