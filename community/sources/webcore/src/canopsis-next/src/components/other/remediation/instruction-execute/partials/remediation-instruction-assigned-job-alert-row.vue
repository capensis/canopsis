<template lang="pug">
  tr
    td.jobs-alert-cell.pa-0(colspan="4")
      v-expand-transition
        remediation-instruction-assigned-job-alert(
          v-if="shownAlert",
          @skip="$emit('skip', job)",
          @await="hideAlert"
        )
</template>

<script>
import { REMEDIATION_JOB_EXECUTION_STATUSES } from '@/constants';

import RemediationInstructionAssignedJobAlert from './remediation-instruction-assigned-job-alert.vue';

export default {
  inject: ['$system'],
  components: { RemediationInstructionAssignedJobAlert },
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      shownAlert: !!this.job.started_at,
      timer: null,
    };
  },
  watch: {
    'job.started_at': function startedAtWatcher(value) {
      if (value && !this.job.launched_at) {
        this.startTimer();
      } else if (this.shownAlert) {
        this.hideAlert();
      }
    },

    'job.launched_at': function launchedAtWatcher() {
      this.hideAlert();
    },

    'job.status': function statusWatcher(status) {
      if (status !== REMEDIATION_JOB_EXECUTION_STATUSES.running) {
        this.hideAlert();
      }
    },
  },
  beforeDestroy() {
    this.stopTimer();
  },
  methods: {
    startTimer() {
      this.timer = setTimeout(() => {
        this.shownAlert = true;
      }, this.$system.jobExecutorFetchTimeoutSeconds * 1000);
    },

    stopTimer() {
      if (this.timer) {
        clearTimeout(this.timer);

        this.timer = null;
      }
    },

    hideAlert() {
      this.stopTimer();
      this.shownAlert = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .jobs-alert-cell {
    height: auto;
  }
</style>
