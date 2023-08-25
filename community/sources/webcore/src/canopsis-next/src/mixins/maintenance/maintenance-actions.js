import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

const { mapGetters: mapInfoGetters, mapActions: mapInfoActions } = createNamespacedHelpers('info');
const { mapActions: mapAuthActions } = createNamespacedHelpers('auth');

export const maintenanceActionsMixin = {
  computed: {
    ...mapInfoGetters(['maintenance']),
  },
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

    showConfirmationLeaveMaintenanceMode() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          title: this.$t('modals.confirmationLeaveMaintenance.title'),
          text: this.$t('modals.confirmationLeaveMaintenance.text'),
          action: this.disableMaintenanceMode,
        },
      });
    },

    showToggleMaintenanceModeModal() {
      return this.maintenance
        ? this.showConfirmationLeaveMaintenanceMode()
        : this.showCreateMaintenanceModeModal();
    },
  },
};
