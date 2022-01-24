import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('userPreference');

export const entitiesUserPreferenceMixin = {
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
      fetchUserPreference: 'fetchItem',
      fetchUserPreferenceWithoutStore: 'fetchItemWithoutStore',
      updateUserPreference: 'update',
    }),

    updateContentInUserPreference(content = {}) {
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
