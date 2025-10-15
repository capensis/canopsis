import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useAuth } from '@/hooks/auth';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';

/**
 * Provides actions for managing remediation instructions, including modals for duplication,
 * editing, approval, and removal.
 *
 * @param {Function} [refresh=() => {}] - Callback to refresh data after actions.
 * @returns {Object} Actions for remediation instructions:
 *   - showDuplicateRemediationInstructionModal
 *   - showEditRemediationInstructionModal
 *   - showApproveRemediationInstructionModal
 *   - showRemoveRemediationInstructionModal
 *   - showRemoveSelectedRemediationInstructionModal
 */
export const useRemediationInstructionsActions = (refresh = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();
  const { currentUser } = useAuth();

  const {
    createRemediationInstruction,
    updateRemediationInstruction,
    removeRemediationInstruction,
  } = useRemdeitionInstruction();

  /**
   * Opens a modal to duplicate a remediation instruction.
   *
   * @param {Object} [remediationInstruction={}] - The remediation instruction to duplicate.
   * @returns {Promise<void>} Promise resolving when the modal action completes.
   */
  const showDuplicateRemediationInstructionModal = (remediationInstruction = {}) => modals.show({
    name: MODALS.createRemediationInstruction,
    config: {
      remediationInstruction: omit(remediationInstruction, ['_id']),
      title: t('modals.createRemediationInstruction.duplicate.title'),
      action: async (instruction) => {
        await createRemediationInstruction({ data: instruction });

        popups.success({
          text: t('modals.createRemediationInstruction.duplicate.popups.success', {
            instructionName: remediationInstruction.name,
          }),
        });

        refresh();
      },
    },
  });

  /**
   * Opens a modal to edit a remediation instruction.
   * Disables editing if the instruction was requested by another user.
   *
   * @param {Object} remediationInstruction - The remediation instruction to edit.
   * @returns {void}
   */
  const showEditRemediationInstructionModal = (remediationInstruction) => {
    const wasRequestedByAnotherUser = !!remediationInstruction.approval
      && !(remediationInstruction.approval.requested_by?._id === currentUser.value._id);

    modals.show({
      name: MODALS.createRemediationInstruction,
      config: {
        remediationInstruction,
        disabled: wasRequestedByAnotherUser,
        title: t('modals.createRemediationInstruction.edit.title'),
        action: async (instruction) => {
          await updateRemediationInstruction({ id: remediationInstruction._id, data: instruction });

          popups.success({
            text: t('modals.createRemediationInstruction.edit.popups.success', {
              instructionName: instruction.name,
            }),
          });

          refresh();
        },
      },
    });
  };

  /**
   * Opens a modal to approve a remediation instruction.
   *
   * @param {Object} [remediationInstruction={}] - The remediation instruction to approve.
   * @returns {Promise<void>} Promise resolving when the modal action completes.
   */
  const showApproveRemediationInstructionModal = (remediationInstruction = {}) => modals.show({
    name: MODALS.remediationInstructionApproval,
    config: {
      remediationInstructionId: remediationInstruction._id,
      afterSubmit: refresh,
    },
  });

  /**
   * Opens a confirmation modal to remove a remediation instruction.
   *
   * @param {Object} [remediationInstruction={}] - The remediation instruction to remove.
   * @returns {Promise<void>} Promise resolving when the modal action completes.
   */
  const showRemoveRemediationInstructionModal = (remediationInstruction = {}) => modals.show({
    name: MODALS.confirmation,
    config: {
      action: async () => {
        await removeRemediationInstruction({ id: remediationInstruction._id });

        refresh();
      },
    },
  });

  /**
   * Opens a confirmation modal to remove multiple selected remediation instructions.
   *
   * @param {Array<Object>} [selected=[]] - The selected remediation instructions to remove.
   * @returns {Promise<void>} Promise resolving when the modal action completes.
   */
  const showRemoveSelectedRemediationInstructionModal = (selected = []) => modals.show({
    name: MODALS.confirmation,
    config: {
      action: async () => {
        await Promise.all(selected.map(({ _id: id }) => removeRemediationInstruction({ id })));

        refresh();
      },
    },
  });

  return {
    showDuplicateRemediationInstructionModal,
    showEditRemediationInstructionModal,
    showApproveRemediationInstructionModal,
    showRemoveRemediationInstructionModal,
    showRemoveSelectedRemediationInstructionModal,
  };
};
