import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('notifications');

export const entitiesNotificationsMixin = {
  methods: {
    ...mapActions({
      fetchNotificationsListWithoutStore: 'fetchListWithoutStore',
      fetchEventFilterFailuresWithoutStore: 'fetchEventFilterFailuresWithoutStore',
      fetchInstructionsToRateWithoutStore: 'fetchInstructionsToRateWithoutStore',
      fetchInstructionsToApproveWithoutStore: 'fetchInstructionsToApproveWithoutStore',
      markNotificationAsRead: 'markAsRead',
    }),
  },
};
