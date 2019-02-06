import omit from 'lodash/omit';
import { createNamespacedHelpers } from 'vuex';

import { generateUserPreferenceByWidgetAndUser } from '@/helpers/entities';

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
      fetchUserPreferenceByWidgetIdWithoutStore: 'fetchItemByWidgetIdWithoutStore',
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

    async copyUserPreferencesForWidgets(widgetsIdsMap) {
      const oldWidgetsIds = Object.keys(widgetsIdsMap);

      const userPreferences = await Promise.all(oldWidgetsIds.map(widgetId => (
        this.fetchUserPreferenceByWidgetIdWithoutStore({ widgetId })
      )));

      return Promise.all(userPreferences.map((userPreference) => {
        if (!userPreference) {
          return Promise.resolve();
        }

        const newWidgetId = widgetsIdsMap[userPreference.widget_id];
        const newUserPreference = generateUserPreferenceByWidgetAndUser({
          _id: newWidgetId,
        }, this.currentUser);

        return this.createUserPreference({
          userPreference: {
            ...newUserPreference,
            ...omit(userPreference, ['_id', 'widget_id']),
          },
        });
      }));
    },
  },
};
