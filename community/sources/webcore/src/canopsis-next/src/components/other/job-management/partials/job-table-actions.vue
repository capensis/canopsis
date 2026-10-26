<template>
  <div v-if="isRunning">
    <c-action-btn
      :tooltip="t('jobs.actions.stopJob')"
      :disabled-button="stopPending"
      :loading="stopPending"
      icon="stop"
      @click="stop"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import {
  JOB_ACTION_REFETCH_DELAY,
  JOB_EXECUTION_STATUS,
  JOBS_TABS,
  MODALS,
} from '@/constants';

import { promisedWait } from '@/helpers/async';
import {
  prepareJobsStopRequest,
  splitJobsBulkResponse,
} from '@/helpers/entities/job/list';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';

import { useJobManagement } from '../hooks/job-management';

export default {
  props: {
    item: {
      type: Object,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    jobType: {
      type: String,
      required: true,
    },
  },

  emits: ['clear:items'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { stopJobs, fetchJobsListWithPreviousParams } = useJobManagement(props.jobType);
    const preparedItems = computed(() => (props.item ? [props.item] : props.items));
    const isRunning = computed(
      () => preparedItems.value.length > 0 && preparedItems.value.every(
        item => item.status === JOB_EXECUTION_STATUS.running,
      ),
    );
    const withStopPendingJobs = computed(() => props.jobType !== JOBS_TABS.authToken);
    const hasPendingJobs = computed(
      () => withStopPendingJobs.value && preparedItems.value.some(item => item.pending_jobs > 0),
    );

    const showResults = (result = []) => {
      const { successItems, failedItems } = splitJobsBulkResponse(result);

      if (failedItems.length) {
        popups.error({ text: t('jobs.popups.stopFailed', failedItems.length) });
      }

      if (successItems.length) {
        const firstItem = preparedItems.value.find(
          item => item._id === successItems[0]?.item?._id,
        );

        popups.info({
          text: t('jobs.popups.stopped', successItems.length, {
            count: successItems.length,
            ruleName: firstItem?.rule_name ?? '',
          }),
        });
      }
    };

    const { pending: stopPending, handler: executeStop } = usePendingHandler(
      async (stopPendingJobsById = {}) => {
        try {
          const data = prepareJobsStopRequest(
            preparedItems.value,
            stopPendingJobsById,
            withStopPendingJobs.value,
          );
          const result = await stopJobs({ data });

          showResults(result);
          await promisedWait(JOB_ACTION_REFETCH_DELAY);
        } catch (error) {
          console.error(error);
          popups.error({ text: t('errors.default') });
        } finally {
          fetchJobsListWithPreviousParams();
          emit('clear:items');
        }
      },
    );

    const stop = () => {
      if (!hasPendingJobs.value) {
        executeStop();

        return;
      }

      modals.show({
        name: MODALS.stopJobs,
        config: {
          items: preparedItems.value,
          jobType: props.jobType,
          action: executeStop,
        },
      });
    };

    return {
      t,
      isRunning,
      stopPending,
      stop,
    };
  },
};
</script>
