import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

const { mapActions: mapInfoActions } = createNamespacedHelpers('info');
const { mapActions: mapAuthActions } = createNamespacedHelpers('auth');

export const maintenanceActionsMixin = {
  methods: {
    ...mapInfoActions(['updateMaintenanceMode']),
    ...mapAuthActions(['logout']),

    enableMaintenanceMode(form) {
      return this.updateMaintenanceMode({
        data: {
          ...form,
          enabled: true,
        },
      });
    },

    disableMaintenanceMode() {
      return this.updateMaintenanceMode({
        data: { enabled: false },
      });
    },

    showCreateMaintenanceModeModal() {
      this.$modals.show({
        name: MODALS.createMaintenance,
        config: {
          cancelTimer: 5,
          action: this.enableMaintenanceMode,
          cancel: this.disableMaintenanceMode,
          logout: this.logout,
        },
      });
    },
  },
};
