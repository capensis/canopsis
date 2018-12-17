import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('userPreference');

/**
 * @mixin Helpers for the userPreference entity
 */
export default {
  computed: {
    ...mapGetters({
      getUserPreferenceByWidget: 'getItemByWidget',
    }),

    userPreference() {
      return this.getUserPreferenceByWidget(this.widget);
    },
  },
  methods: {
    ...mapActions({
      fetchUserPreferencesList: 'fetchList',
      fetchUserPreferenceByWidgetId: 'fetchItemByWidgetId',
      createUserPreference: 'create',
    }),

    updateWidgetPreferencesInUserPreference(widgetPreferences = {}) {
      return this.createUserPreference({
        userPreference: {
          ...this.userPreference,
          widget_preferences: widgetPreferences,
        },
      });
    },
  },
};
