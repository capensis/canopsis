import { useNotificationStore } from '@/stores/modules/notification';

export const useNotifications = () => {
  const notificationStore = useNotificationStore();

  return {
    fetchNotificationsListWithoutStore: notificationStore.fetchNotificationsListWithoutStore,
  };
};
