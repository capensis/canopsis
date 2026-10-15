import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('notificationSettings');

export const entitiesNotificationSettingsMixin = {
  methods: {
    ...mapActions({
      fetchNotificationSettingsWithoutStore: 'fetchItemWithoutStore',
      updateNotificationSettings: 'update',
    }),
  },
};
