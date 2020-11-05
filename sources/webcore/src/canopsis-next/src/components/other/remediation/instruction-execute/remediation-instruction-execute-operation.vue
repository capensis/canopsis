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
      template(v-if="isStartedOperation")
        text-editor-blurred(
          :value="operation.name",
          :label="$t('common.description')"
        )
        v-layout(row, justify-end)
          v-btn.accent(@click="$emit('previous')") {{ $t('common.previous') }}
          v-btn.primary.mr-0(@click="$emit('next')") {{ $t('common.next') }}
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
    isStartedOperation() {
      return !!this.operation.started_at;
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
