<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      v-layout(justify-center)
        v-avatar.white--text(color="primary", size="32") {{ operationNumber }}
    v-flex(xs11)
      v-layout
        v-text-field(
          :value="operation.name",
          :label="$t('common.description')",
          readonly,
          hide-details,
          box
        )
      remediation-instruction-status(
        :completed-at="operation.completed_at",
        :time-to-complete="operation.time_to_complete",
        :failed-at="operation.failed_at",
        :started-at="operation.started_at"
      )
      v-expand-transition(mode="out-in")
        text-editor-blurred(
          v-if="!isCompletedOperation",
          :value="operation.name",
          :label="$t('common.description')"
        )
</template>

<script>
import TextEditorBlurred from '@/components/other/text-editor/text-editor-blurred.vue';

import RemediationInstructionStatus from './partials/remediation-instruction-status.vue';

export default {
  components: { TextEditorBlurred, RemediationInstructionStatus },
  props: {
    operation: {
      type: Object,
      required: true,
    },
    operationNumber: {
      type: [Number, String],
      required: true,
    },
  },
  computed: {
    isCompletedOperation() {
      return !!this.operation.completed_at;
    },
  },
};
</script>

<style lang="scss">
.step-line {
  flex: 1;
  width: 2px;
}
</style>
