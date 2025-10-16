import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';

/**
 * Hook for rating remediation instructions by a modal dialog.
 *
 * @param {Function} [refresh=() => {}] - Callback to refresh data after rating. Defaults to a no-op.
 * @returns {{ showRateInstructionModal: (instruction: Object) => Promise<void> }}
 */
export const useRemediationInstructionStatsRate = (refresh = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const { rateRemediationInstruction } = useRemdeitionInstruction();

  /**
   * Shows the modal for rating a remediation instruction.
   *
   * @param {Object} [instruction = {}] - The remediation instruction object.
   * @param {string} [instruction._id] - The unique identifier of the instruction.
   * @param {string} [instruction.name] - The name of the instruction.
   * @returns {Promise<void>} Resolves after the modal action and refresh are complete.
   */
  const showRateInstructionModal = (instruction = {}) => modals.show({
    name: MODALS.rate,
    config: {
      title: t('modals.rateInstruction.title', { name: instruction.name }),
      text: t('modals.rateInstruction.text'),
      action: async (data) => {
        await rateRemediationInstruction({ id: instruction._id, data });

        return refresh();
      },
    },
  });

  return {
    showRateInstructionModal,
  };
};
