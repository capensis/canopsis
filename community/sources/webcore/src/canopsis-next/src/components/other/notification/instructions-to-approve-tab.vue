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
import { USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';

import {
  useRemediationInstructionsActions,
} from '@/components/other/remediation/instructions/hooks/remediation-instructions';

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
    const {
      hasCreateAccess: hasCreateAnyRemediationInstructionAccess,
      hasUpdateAccess: hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAccess: hasDeleteAnyRemediationInstructionAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.remediationInstruction);

    const refresh = () => emit('refresh');

    const {
      showDuplicateRemediationInstructionModal,
      showEditRemediationInstructionModal,
      showApproveRemediationInstructionModal,
      showRemoveRemediationInstructionModal,
      showRemoveSelectedRemediationInstructionModal,
    } = useRemediationInstructionsActions(refresh);

    return {
      hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAnyRemediationInstructionAccess,
      hasCreateAnyRemediationInstructionAccess,

      showApproveRemediationInstructionModal,
      showRemoveSelectedRemediationInstructionModal,
      showRemoveRemediationInstructionModal,
      showDuplicateRemediationInstructionModal,
      showEditRemediationInstructionModal,
    };
  },
};
</script>
