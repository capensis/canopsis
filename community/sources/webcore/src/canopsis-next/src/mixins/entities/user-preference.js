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

    /**
     * Send requests to create userPreference by widgetsIdsMappings
     *
     * @param {Array.<{oldId: string, newId: string}>} widgetsIdsMappings
     * @returns {Promise.<*[]>}
     */
    copyUserPreferencesByWidgetsIdsMappings(widgetsIdsMappings) {
      return Promise.all(widgetsIdsMappings.map(async ({ oldId, newId }) => {
        const userPreference = await this.fetchUserPreferenceByWidgetIdWithoutStore({ widgetId: oldId });

        if (!userPreference) {
          return Promise.resolve();
        }

        const newUserPreference = generateUserPreferenceByWidgetAndUser({
          _id: newId,
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
