<template>
  <c-modal-wrapper close>
    <template #title>
      {{ $t('modals.remediationInstructionApproval.title') }}
    </template>

    <template #text>
      <v-fade-transition>
        <c-layout v-if="pending || !approval" justify="center">
          <v-progress-circular color="primary" indeterminate />
        </c-layout>
        <c-layout v-else column>
          <remediation-instruction-approval-alert
            :user-name="userName"
            :comment="approval.approval.comment"
          >
            <remediation-instruction-approval-form
              :disabled="submitting"
              :submitting="submitting"
              @approve="approve"
              @dismiss="dismiss"
            />
          </remediation-instruction-approval-alert>

          <remediation-instruction-approval-tabs
            :original="approval.original"
            :updated="approval.updated"
          />
        </c-layout>
      </v-fade-transition>
    </template>

    <template #actions>
      <v-btn :disabled="submitting" variant="text" @click="close">
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';

import { MODALS } from '@/constants';

import { useInnerModal } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';

import RemediationInstructionApprovalAlert from '@/components/other/remediation/instructions/partials/approval-alert.vue';
import RemediationInstructionApprovalForm from '@/components/other/remediation/instructions/partials/approval-form.vue';
import RemediationInstructionApprovalTabs from '@/components/other/remediation/instructions/partials/approval-tabs.vue';

defineOptions({ name: MODALS.remediationInstructionApproval });

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { config, close } = useInnerModal(props);
const {
  fetchRemediationInstructionApprovalWithoutStore,
  updateRemediationInstructionApproval,
} = useRemediationInstruction();

const approval = ref(null);
const instructionId = computed(() => config.value.remediationInstructionId);
const userName = computed(() => approval.value?.approval?.requested_by?.name);

const { pending, handler: fetchItem } = usePendingHandler(async () => {
  approval.value = await fetchRemediationInstructionApprovalWithoutStore({
    id: instructionId.value,
  });
}, true);

const { pending: submitting, handler: submit } = usePendingHandler(async (data) => {
  await updateRemediationInstructionApproval({
    id: instructionId.value,
    data,
  });

  await config.value.afterSubmit?.();

  close();
});

const approve = () => submit({ approve: true });
const dismiss = comment =>
  submit({
    approve: false,
    comment,
  });

onMounted(fetchItem);
</script>
