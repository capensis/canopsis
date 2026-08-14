<template>
  <v-sheet class="remediation-instruction-assigned-jobs-expand-panel px-3 py-2">
    <div v-if="isFailedJob">
      {{ t('remediation.instructionExecute.jobs.failedReason') }}:&nbsp;
      <span v-html="job.fail_reason" class="pre-wrap" />
    </div>
    <div>
      {{ t('common.output') }}:&nbsp;
      <span v-html="output" class="pre-wrap" />
    </div>
  </v-sheet>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useRemediationJobExecution } from '@/hooks/store/modules/remediation-job-execution';

export default {
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { getRemediationJobExecutionOutputById } = useRemediationJobExecution();

    const output = computed(() => getRemediationJobExecutionOutputById.value(props.job._id));
    const isFailedJob = computed(() => !!props.job.fail_reason);

    return {
      t,
      output,
      isFailedJob,
    };
  },
};
</script>

<style lang="scss">
.v-sheet.remediation-instruction-assigned-jobs-expand-panel {
  background: rgb(var(--v-theme-surface-variant));
}
</style>
