<template lang="pug">
  v-sheet.px-3.py-2(color="grey lighten-3")
    div(v-if="isFailedJob") {{ $t('remediation.instructionExecute.jobs.failedReason') }}:&nbsp;
      span.pre-wrap(v-html="job.fail_reason")
    div {{ $t('remediation.instructionExecute.jobs.output') }}:&nbsp;
      span.pre-line(v-html="output")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('remediationJobExecution');

export default {
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...mapGetters(['getOutputById']),

    output() {
      return this.getOutputById(this.job._id);
    },

    isFailedJob() {
      return !!this.job.fail_reason;
    },
  },
};
</script>
