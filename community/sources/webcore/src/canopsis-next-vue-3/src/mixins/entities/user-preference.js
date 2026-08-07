import { mapActions, mapStores } from 'pinia';

import { useUserPreferenceStore } from '@/stores/modules/user-preference';

export const entitiesUserPreferenceMixin = {
  props: {
    localWidget: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    ...mapStores(useUserPreferenceStore),

    userPreference() {
      return this.userPreferenceStore.getItemByWidgetId(this.widget._id);
    },
  },

  methods: {
    ...mapActions(useUserPreferenceStore, {
      fetchUserPreferenceItem: 'fetchItem',
      fetchUserPreferenceWithoutStore: 'fetchItemWithoutStore',
      updateUserPreference: 'update',
      updateLocalUserPreference: 'updateLocal',
    }),

    fetchUserPreference(data) {
      if (this.localWidget) {
        return Promise.resolve();
      }

      return this.fetchUserPreferenceItem(data);
    },

    updateContentInUserPreference(content = {}) {
      const method = this.localWidget ? this.updateLocalUserPreference : this.updateUserPreference;

      return method({
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
