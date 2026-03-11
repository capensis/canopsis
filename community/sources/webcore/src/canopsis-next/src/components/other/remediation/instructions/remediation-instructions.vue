<template>
  <v-card-text>
    <remediation-instructions-list
      :remediation-instructions="remediationInstructions"
      :pending="remediationInstructionsPending"
      :total-items="remediationInstructionsMeta.total_count"
      :options="options"
      :updatable="hasUpdateAnyRemediationInstructionAccess"
      :removable="hasDeleteAnyRemediationInstructionAccess"
      :duplicable="hasCreateAnyRemediationInstructionAccess"
      @remove-selected="showRemoveSelectedRemediationInstructionModal"
      @duplicate="showDuplicateRemediationInstructionModal"
      @remove="showRemoveRemediationInstructionModal"
      @approve="showApproveRemediationInstructionModal"
      @edit="showEditRemediationInstructionModal"
      @update:options="updateOptions"
    />
  </v-card-text>
</template>

<script>
import { onMounted } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useFetchListWithOptions } from '@/hooks/query/shared';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useCRUDPermissions } from '@/hooks/auth';

import { useRemediationInstructionsActions } from './hooks/remediation-instructions';
import RemediationInstructionsList from './remediation-instructions-list.vue';

export default {
  components: { RemediationInstructionsList },
  setup() {
    const {
      remediationInstructions,
      remediationInstructionsMeta,
      remediationInstructionsPending,
      fetchRemediationInstructionsList,
    } = useRemdeitionInstruction();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useFetchListWithOptions({
      fetchListHandler: ({ params }) => fetchRemediationInstructionsList({
        params: {
          ...params,
          with_flags: true,
          with_month_executions: true,
        },
      }),
    });

    const {
      hasCreateAccess: hasCreateAnyRemediationInstructionAccess,
      hasUpdateAccess: hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAccess: hasDeleteAnyRemediationInstructionAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.remediationInstruction);

    const {
      showDuplicateRemediationInstructionModal,
      showEditRemediationInstructionModal,
      showApproveRemediationInstructionModal,
      showRemoveRemediationInstructionModal,
      showRemoveSelectedRemediationInstructionModal,
    } = useRemediationInstructionsActions(fetchList);

    onMounted(fetchList);

    return {
      remediationInstructions,
      remediationInstructionsMeta,
      remediationInstructionsPending,
      options,
      updateOptions,
      hasCreateAnyRemediationInstructionAccess,
      hasUpdateAnyRemediationInstructionAccess,
      hasDeleteAnyRemediationInstructionAccess,
      showDuplicateRemediationInstructionModal,
      showEditRemediationInstructionModal,
      showApproveRemediationInstructionModal,
      showRemoveRemediationInstructionModal,
      showRemoveSelectedRemediationInstructionModal,
    };
  },
};
</script>
