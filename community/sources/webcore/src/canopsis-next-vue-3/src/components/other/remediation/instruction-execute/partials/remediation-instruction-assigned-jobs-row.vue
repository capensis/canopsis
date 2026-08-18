<template>
  <tr>
    <td class="pa-0">
      <c-layout align="center">
        <v-btn
          v-if="executable"
          :disabled="isRunningJob || isFailedJob"
          color="primary"
          size="small"
          rounded
          block
          @click="$emit('execute-job', job)"
        >
          {{ job.name }}
        </v-btn>
        <span
          v-else
          class="text-body-1"
        >
          {{ job.name }}
        </span>
        <v-tooltip
          v-if="isFailedJob || isCompletedJob"
          :disabled="!hasStatusMessage"
          max-width="400"
          location="end"
        >
          <template #activator="{ props: activatorProps }">
            <v-btn
              v-bind="activatorProps"
              :loading="outputPending"
              class="mr-1"
              size="small"
              icon
              @click="toggleExpanded"
            >
              <v-icon :color="isFailedJob ? 'error' : 'success'">
                {{ statusIcon }}
              </v-icon>
            </v-btn>
          </template>
          <div v-if="job.fail_reason">
            <span>{{ t('remediation.instructionExecute.jobs.failedReason') }}:&nbsp;</span>
            <span v-html="job.fail_reason" class="pre-wrap" />
          </div>
        </v-tooltip>
        <v-tooltip
          v-if="cancelable && isRunningJob && hasJobsInQueue"
          location="top"
        >
          <template #activator="{ props: activatorProps }">
            <v-btn
              v-bind="activatorProps"
              color="error"
              class="ml-2"
              size="small"
              rounded
              block
              @click="$emit('cancel-job-execution', job)"
            >
              {{ t('common.cancel') }}
            </v-btn>
          </template>
          <span>{{ queueNumberTooltip }}</span>
        </v-tooltip>
      </c-layout>
    </td>
    <td
      v-if="rowMessage"
      colspan="3"
    >
      <div class="text-error text-center">
        {{ rowMessage }}
      </div>
    </td>
    <template v-else>
      <progress-cell
        :pending="isRunningJob && !job.started_at"
        class="text-center"
      >
        <span v-if="!isCancelledJob">
          {{ convertDateToString(job.started_at, 'long', '-') }}
        </span>
      </progress-cell>
      <progress-cell
        :pending="shownLaunchedPendingJob"
        class="text-center"
      >
        <span>{{ convertDateToString(job.launched_at, 'long', '-') }}</span>
      </progress-cell>
      <progress-cell
        :pending="shownCompletedPendingJob"
        class="text-center"
      >
        <span>{{ convertDateToString(job.completed_at, 'long', '-') }}</span>
      </progress-cell>
    </template>
  </tr>
</template>

<script>
import { computed, ref } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { isJobExecutionCancelled, isJobExecutionRunning } from '@/helpers/entities/remediation/job/form';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useRemediationJobExecution } from '@/hooks/store/modules/remediation-job-execution';

import ProgressCell from '@/components/common/table/progress-cell.vue';

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

  emits: ['execute-job', 'cancel-job-execution', 'update:expanded'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const popups = usePopups();
    const { fetchRemediationJobExecutionOutput } = useRemediationJobExecution();
    const outputPending = ref(false);

    const isRunningJob = computed(() => isJobExecutionRunning(props.job));
    const isCancelledJob = computed(() => isJobExecutionCancelled(props.job));
    const isStartedJob = computed(() => !!props.job.started_at);
    const isLaunchedJob = computed(() => !!props.job.launched_at);
    const isCompletedJob = computed(() => !!props.job.completed_at);
    const isFailedJob = computed(() => !!props.job.fail_reason);
    const shownLaunchedPendingJob = computed(() =>
      !isCancelledJob.value
      && !isFailedJob.value
      && !isLaunchedJob.value
      && isStartedJob.value,
    );
    const shownCompletedPendingJob = computed(() =>
      !isCancelledJob.value
      && !isFailedJob.value
      && !isCompletedJob.value
      && isStartedJob.value
      && isLaunchedJob.value,
    );
    const hasStatusMessage = computed(() => props.job.fail_reason);
    const hasJobsInQueue = computed(() => props.job.queue_number > 0);
    const statusIcon = computed(() => (isFailedJob.value ? 'warning' : 'check'));
    const queueNumberTooltip = computed(() => t('remediation.instructionExecute.queueNumber', {
      number: props.job.queue_number,
      name: props.job.name,
    }));
    const rowMessage = computed(() => {
      if (props.job.queue_number > 0) {
        return t('remediation.instructionExecute.queueNumber', {
          number: props.job.queue_number,
          name: props.job.name,
        });
      }

      if (isCancelledJob.value) {
        return t('remediation.instructionExecute.jobs.stopped');
      }

      return '';
    });

    const toggleExpanded = async () => {
      if (!props.expanded) {
        try {
          outputPending.value = true;

          await fetchRemediationJobExecutionOutput({ id: props.job._id });
        } catch (err) {
          console.error(err);

          popups.error({ text: err.message || err.description || t('errors.default') });
        } finally {
          outputPending.value = false;
        }
      }

      emit('update:expanded', !props.expanded);
    };

    return {
      t,
      outputPending,
      isRunningJob,
      isCancelledJob,
      isCompletedJob,
      isFailedJob,
      shownLaunchedPendingJob,
      shownCompletedPendingJob,
      hasStatusMessage,
      hasJobsInQueue,
      statusIcon,
      queueNumberTooltip,
      rowMessage,
      convertDateToString,
      toggleExpanded,
    };
  },
};
</script>
