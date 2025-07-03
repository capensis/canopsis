<template>
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
</template>

<script>
import { toRef } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';

import {
  useRemediationInstructionsActions,
} from '@/components/other/remediation/instructions/hooks/remediation-instructions';

import RemediationInstructionsList from '@/components/other/remediation/instructions/remediation-instructions-list.vue';

import { useNotificationActiveId } from './hooks/notifications';

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
    activeId: {
      type: String,
      required: false,
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

    useNotificationActiveId({
      activeId: toRef(props, 'activeId'),
      items: toRef(props, 'items'),
      action: showApproveRemediationInstructionModal,
    });

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
