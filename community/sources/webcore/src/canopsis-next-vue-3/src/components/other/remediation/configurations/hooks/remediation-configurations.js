import { omit } from 'lodash';

import { MODALS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useRemediationConfiguration } from '@/hooks/store/modules/remediation-configuration';

/**
 * Provides modal actions for remediation configurations.
 *
 * @param {Function} [refresh=() => {}] - Callback to refresh data after actions.
 * @returns {Object}
 */
export const useRemediationConfigurationsActions = (refresh = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();

  const {
    createRemediationConfiguration,
    updateRemediationConfiguration,
    removeRemediationConfiguration,
  } = useRemediationConfiguration();

  const showEditRemediationConfigurationModal = remediationConfiguration =>
    modals.show({
      name: MODALS.createRemediationConfiguration,
      config: {
        title: t('modals.createRemediationConfiguration.edit.title'),
        remediationConfiguration,
        action: async (configuration) => {
          await updateRemediationConfiguration({
            id: remediationConfiguration._id,
            data: configuration,
          });

          popups.success({
            text: t('modals.createRemediationConfiguration.edit.popups.success', {
              configurationName: configuration.name,
            }),
          });

          await refresh();
        },
      },
    });

  const showRemoveRemediationConfigurationModal = remediationConfiguration =>
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeRemediationConfiguration({ id: remediationConfiguration._id });
          await refresh();
        },
      },
    });

  const showDuplicateRemediationConfigurationModal = remediationConfiguration =>
    modals.show({
      name: MODALS.createRemediationConfiguration,
      config: {
        title: t('modals.createRemediationConfiguration.duplicate.title'),
        remediationConfiguration: omit(remediationConfiguration, ['_id']),
        action: async (configuration) => {
          await createRemediationConfiguration({ data: configuration });

          popups.success({
            text: t('modals.createRemediationConfiguration.duplicate.popups.success', {
              configurationName: remediationConfiguration.name,
            }),
          });

          await refresh();
        },
      },
    });

  return {
    showEditRemediationConfigurationModal,
    showRemoveRemediationConfigurationModal,
    showDuplicateRemediationConfigurationModal,
  };
};
