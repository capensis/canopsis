<template>
  <v-card-text>
    <div class="text-center pa-4">
      <h3 class="mb-4">
        {{ $t('notifications.tabs.instructionsToApprove') }}
      </h3>
      <remediation-instructions-list
        :remediation-instructions="items"
        :pending="pending"
        :total-items="meta.total_count"
        :options="options"
        :updatable="hasUpdateAnyRemediationInstructionAccess"
        :removable="hasDeleteAnyRemediationInstructionAccess"
        :duplicable="hasCreateAnyRemediationInstructionAccess"
        @update:options="$emit('update:options', $event)"
        @remove-selected="showRemoveSelectedRemediationInstructionModal"
        @duplicate="showDuplicateRemediationInstructionModal"
        @remove="showRemoveRemediationInstructionModal"
        @approve="showApproveRemediationInstructionModal"
        @edit="showEditRemediationInstructionModal"
      />
    </div>
  </v-card-text>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useCRUDPermissions } from '@/hooks/auth';

import RemediationInstructionsList from '@/components/other/remediation/instructions/remediation-instructions-list.vue';

export default {
  components: {
    RemediationInstructionsList,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    meta: {
      type: Object,
      default: () => {},
    },
    options: {
      type: Object,
      default: () => {},
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const { updateInstructionApproval } = useRemdeitionInstruction();
    const {
      hasCreateAccess: hasCreateAnyRemediationInstructionAccess,
      hasUpdateAccess: hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAccess: hasDeleteAnyRemediationInstructionAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.remediationInstruction);

    const showRemoveSelectedRemediationInstructionModal = (instructions) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.removeSelectedInstructions.title'),
          text: t('modals.removeSelectedInstructions.text', { count: instructions.length }),
        },
      });
    };

    const showRemoveRemediationInstructionModal = (instruction) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.removeInstruction.title'),
          text: t('modals.removeInstruction.text', { name: instruction.name }),
          action: async () => {},
        },
      });
    };

    const showDuplicateRemediationInstructionModal = (instruction) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.duplicateInstruction.title'),
          text: t('modals.duplicateInstruction.text', { name: instruction.name }),
          action: async () => {},
        },
      });
    };

    const showEditRemediationInstructionModal = (instruction) => {
      modals.show({
        name: MODALS.remediationInstruction,
        config: {
          title: t('modals.editInstruction.title'),
          instruction,
        },
      });
    };

    const showApproveRemediationInstructionModal = (instruction) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.approveInstruction.title'),
          text: t('modals.approveInstruction.text', { name: instruction.name }),
          action: async () => {
            await updateInstructionApproval({
              id: instruction._id,
              data: { approved: true },
            });
            emit('refresh');
          },
        },
      });
    };

    return {
      hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAnyRemediationInstructionAccess,
      hasCreateAnyRemediationInstructionAccess,

      showApproveRemediationInstructionModal, // TODO: move this functions to hook from remediation-instructions.vue
      showRemoveSelectedRemediationInstructionModal,
      showRemoveRemediationInstructionModal,
      showDuplicateRemediationInstructionModal,
      showEditRemediationInstructionModal,
    };
  },
};
</script>
