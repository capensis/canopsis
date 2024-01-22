<template>
  <tr>
    <td class="pa-0">
      <v-layout align-center>
        <v-btn
          class="primary"
          v-if="executable"
          :disabled="isRunningJob || isFailedJob"
          rounded
          small
          block
          @click="$emit('execute-job', job)"
        >
          {{ job.name }}
        </v-btn>
        <span
          class="text-body-1"
          v-else
        >
          {{ job.name }}
        </span>
        <v-tooltip
          v-if="isFailedJob || isCompletedJob"
          :disabled="!hasStatusMessage"
          max-width="400"
          right
        >
          <template #activator="{ on }">
            <v-btn
              class="mr-1"
              v-on="on"
              :loading="outputPending"
              icon
              small
              @click="toggleExpanded"
            >
              <v-icon :color="isFailedJob ? 'error' : 'success'">
                {{ statusIcon }}
              </v-icon>
            </v-btn>
          </template>
          <div v-if="job.fail_reason">
            <span>{{ $t('remediation.instructionExecute.jobs.failedReason') }}:&nbsp;</span>
            <span
              class="pre-wrap"
              v-html="job.fail_reason"
            />
          </div>
        </v-tooltip>
        <v-tooltip
          v-if="cancelable && isRunningJob && hasJobsInQueue"
          top
        >
          <template #activator="{ on }">
            <v-btn
              class="error ml-2"
              v-on="on"
              rounded
              small
              block
              @click="$emit('cancel-job-execution', job)"
            >
              {{ $t('common.cancel') }}
            </v-btn>
          </template>
          <span>{{ queueNumberTooltip }}</span>
        </v-tooltip>
      </v-layout>
    </td>
    <td
      v-if="rowMessage"
      colspan="3"
    >
      <div class="error--text text-center">
        {{ rowMessage }}
      </div>
    </td>
    <template v-else>
      <progress-cell
        class="text-center"
        :pending="isRunningJob && !job.started_at"
      >
        <span v-if="!isCancelledJob">{{ job.started_at | date('long', '-') }}</span>
      </progress-cell>
      <progress-cell
        class="text-center"
        :pending="shownLaunchedPendingJob"
      >
        <span>{{ job.launched_at | date('long', '-') }}</span>
      </progress-cell>
      <progress-cell
        class="text-center"
        :pending="shownCompletedPendingJob"
      >
        <span>{{ job.completed_at | date('long', '-') }}</span>
      </progress-cell>
    </template>
  </tr>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { isJobExecutionCancelled, isJobExecutionRunning } from '@/helpers/entities/remediation/job/form';

import ProgressCell from '@/components/common/table/progress-cell.vue';

const { mapActions } = createNamespacedHelpers('remediationJobExecution');

export default {
  components: { ProgressCell },
  props: {
    job: {
      type: Object,
      required: true,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
    executable: {
      type: Boolean,
      default: false,
    },
    cancelable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      outputPending: false,
    };
  },
  computed: {
    isRunningJob() {
      return isJobExecutionRunning(this.job);
    },

    isCancelledJob() {
      return isJobExecutionCancelled(this.job);
    },

    isStartedJob() {
      return !!this.job.started_at;
    },

    isLaunchedJob() {
      return !!this.job.launched_at;
    },

    isCompletedJob() {
      return !!this.job.completed_at;
    },

    isFailedJob() {
      return !!this.job.fail_reason;
    },

    shownLaunchedPendingJob() {
      return !this.isCancelledJob
        && !this.isFailedJob
        && !this.isLaunchedJob
        && this.isStartedJob;
    },

    shownCompletedPendingJob() {
      return !this.isCancelledJob
        && !this.isFailedJob
        && !this.isCompletedJob
        && this.isStartedJob
        && this.isLaunchedJob;
    },

    hasStatusMessage() {
      return this.job.fail_reason;
    },

    hasJobsInQueue() {
      return this.job.queue_number > 0;
    },

    statusIcon() {
      return this.isFailedJob ? 'warning' : 'check';
    },

    queueNumberTooltip() {
      return this.$t('remediation.instructionExecute.queueNumber', {
        number: this.job.queue_number,
        name: this.job.name,
      });
    },

    rowMessage() {
      if (this.job.queue_number > 0) {
        return this.$t('remediation.instructionExecute.queueNumber', {
          number: this.job.queue_number,
          name: this.job.name,
        });
      }

      if (this.isCancelledJob) {
        return this.$t('remediation.instructionExecute.jobs.stopped');
      }

      return '';
    },
  },
  methods: {
    ...mapActions(['fetchOutput']),

    async toggleExpanded() {
      if (!this.expanded) {
        try {
          this.outputPending = true;

          await this.fetchOutput({ id: this.job._id });
        } catch (err) {
          console.error(err);

          this.$popups.error({ text: err.message || err.description || this.$t('errors.default') });

          console.warn(err);
        } finally {
          this.outputPending = false;
        }
      }

      this.$emit('update:expanded', !this.expanded);
    },
  },
};
</script>
