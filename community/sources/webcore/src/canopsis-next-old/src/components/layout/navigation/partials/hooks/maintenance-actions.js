import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useInfo } from '@/hooks/store/modules/info';

/**
 * Hook for maintenance mode actions
 * Provides methods to enable, disable, and toggle maintenance mode
 *
 * @returns {Object} An object containing maintenance-related methods and state
 */
export const useMaintenanceActions = () => {
  const { t } = useI18n();
  const modals = useModals();
  const { maintenance, updateMaintenanceMode, fetchAppInfo } = useInfo();

  /**
   * Enable maintenance mode
   *
   * @param {Object} form - Form data for maintenance mode
   */
  const enableMaintenanceMode = async (form) => {
    await updateMaintenanceMode({
      data: {
        ...form,
        enabled: true,
      },
    });

    await fetchAppInfo();
  };

  /**
   * Disable maintenance mode
   */
  const disableMaintenanceMode = async () => {
    await updateMaintenanceMode({
      data: { enabled: false },
    });

    await fetchAppInfo();
  };

  /**
   * Show modal to create maintenance mode
   */
  const showCreateMaintenanceModeModal = () => {
    modals.show({
      name: MODALS.createMaintenance,
      config: {
        action: enableMaintenanceMode,
        warningText: t('maintenance.logoutWarning'),
      },
    });
  };

  /**
   * Show confirmation modal to leave maintenance mode
   */
  const showConfirmationLeaveMaintenanceMode = () => {
    modals.show({
      name: MODALS.confirmation,
      config: {
        title: t('modals.confirmationLeaveMaintenance.title'),
        text: t('modals.confirmationLeaveMaintenance.text'),
        action: disableMaintenanceMode,
      },
    });
  };

  /**
   * Toggle maintenance mode (show appropriate modal)
   */
  const showToggleMaintenanceModeModal = () => {
    if (maintenance.value) {
      showConfirmationLeaveMaintenanceMode();
    } else {
      showCreateMaintenanceModeModal();
    }
  };

  return {
    maintenance,
    enableMaintenanceMode,
    disableMaintenanceMode,
    showCreateMaintenanceModeModal,
    showConfirmationLeaveMaintenanceMode,
    showToggleMaintenanceModeModal,
  };
};
