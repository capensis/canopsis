import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

const { mapGetters: mapInfoGetters, mapActions: mapInfoActions } = createNamespacedHelpers('info');

export const maintenanceActionsMixin = {
  computed: {
    ...mapInfoGetters(['maintenance']),
  },
  methods: {
    ...mapInfoActions(['updateMaintenanceMode']),

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
          action: this.enableMaintenanceMode,
          warningText: this.$t('maintenance.logoutWarning'),
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
