import { useNotificationSettingsStore } from '@/stores/modules/notification-settings';

/**
 * Notification settings store facade.
 *
 * @returns {Object}
 */
export const useNotificationSettings = () => {
  const notificationSettingsStore = useNotificationSettingsStore();

  return {
    fetchNotificationSettingsWithoutStore: notificationSettingsStore.fetchItemWithoutStore,
    updateNotificationSettings: notificationSettingsStore.update,
  };
};
