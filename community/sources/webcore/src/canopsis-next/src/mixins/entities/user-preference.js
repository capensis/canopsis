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
      getUserPreferenceByWidgetId: 'getItemByWidgetId',
    }),

    userPreference() {
      return this.getUserPreferenceByWidgetId(this.widget._id);
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
  },
};
