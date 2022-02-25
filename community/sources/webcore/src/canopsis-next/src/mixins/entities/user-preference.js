import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('userPreference');

export const entitiesUserPreferenceMixin = {
  props: {
    localWidget: {
      type: Boolean,
      default: false,
    },
  },
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
      fetchUserPreferenceItem: 'fetchItem',
      fetchUserPreferenceWithoutStore: 'fetchItemWithoutStore',
      updateUserPreference: 'update',
    }),

    fetchUserPreference(data) {
      if (this.localWidget) {
        return Promise.resolve();
      }

      return this.fetchUserPreferenceItem(data);
    },

    updateContentInUserPreference(content = {}) {
      if (this.localWidget) {
        return Promise.resolve();
      }

      return this.updateUserPreference({
        data: {
          ...this.userPreference,

          content: {
            ...this.userPreference.content,
            ...content,
          },
        },
      });
    },

    /**
     * Send requests to create userPreference by widgetsIdsMappings
     *
     * @param {{oldId: string, newId: string}[]} widgetsIdsMappings
     * @returns {Promise}
     */
    copyUserPreferencesByWidgetsIdsMappings(widgetsIdsMappings) {
      if (this.localWidget) {
        return Promise.resolve();
      }

      return Promise.all(widgetsIdsMappings.map(async ({ oldId, newId }) => {
        const userPreference = await this.fetchUserPreferenceWithoutStore({ id: oldId });

        if (!userPreference) {
          return Promise.resolve();
        }

        return this.updateUserPreference({
          data: {
            ...userPreference,

            widget: newId,
          },
        });
      }));
    },
  },
};
