<template>
  <c-layout v-if="item || !isStopped" gap="1">
    <c-action-btn
      v-if="item"
      :tooltip="t('jobs.actions.editJob')"
      icon="edit"
      @click="edit"
    />

    <template v-if="!isStopped && isRunning">
      <c-action-btn
        :tooltip="t('jobs.actions.pauseJob')"
        :disabled-button="pausePending"
        :loading="pausePending"
        icon="pause"
        @click="pause"
      />
      <c-action-btn
        :tooltip="t('jobs.actions.repeatJob')"
        :disabled-button="repeatPending"
        :loading="repeatPending"
        icon="refresh"
        @click="repeat"
      />
    </template>
    <c-action-btn
      v-else-if="!isStopped"
      :tooltip="t('jobs.actions.startJob')"
      :disabled-button="playPending"
      :loading="playPending"
      icon="play_arrow"
      @click="play"
    />
  </c-layout>
</template>

<script>
import { keyBy } from 'lodash';
import { computed } from 'vue';

import { JOB_ACTION_REFETCH_DELAY, JOB_STATUS, MODALS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { promisedWait } from '@/helpers/async';

import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useTicketStatusJob } from '@/hooks/store/modules/ticket-status-job';

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
  },

  emits: ['clear:items'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { callActionWithPopup } = useCallActionWithPopup(false, 'info');
    const {
      updateTicketStatusJob,
      playTicketStatusJobs,
      pauseTicketStatusJobs,
      fetchTicketStatusJobsListWithPreviousParams,
    } = useTicketStatusJob();
    const preparedItems = computed(() => (props.item ? [props.item] : props.items));
    const preparedItemsById = computed(() => keyBy(preparedItems.value, '_id'));
    const isRunning = computed(() =>
      preparedItems.value.some(item => item.status === JOB_STATUS.running),
    );
    const isStopped = computed(() =>
      preparedItems.value.every(item => item.status === JOB_STATUS.stopped),
    );

    const showTicketStatusJobActionPopup = (items = [], textKey, popupType = 'info') => {
      const firstItem = preparedItemsById.value[items[0]?.item?._id] ?? {};
      const { rule_name: ruleName, ticket_id: ticketNumber } = firstItem;

      popups[popupType]({
        text: `<span class="font-weight-regular">${t(textKey, items.length, {
          ruleName,
          ticketNumber,
          count: items.length,
        })}</span>`,
      });
    };

    const callTicketStatusJobActionWithPopups = async ({
      action,
      successTextKey,
      errorTextKey,
    }) => {
      try {
        const result = await action();

        await promisedWait(JOB_ACTION_REFETCH_DELAY);

        const { successItems, failedItems } = result.reduce(
          (acc, item) => {
            acc[item.error ? 'failedItems' : 'successItems'].push(item);

            return acc;
          },
          { successItems: [], failedItems: [] },
        );

        if (failedItems.length) {
          showTicketStatusJobActionPopup(failedItems, errorTextKey, 'error');
        }

        if (successItems.length) {
          showTicketStatusJobActionPopup(successItems, successTextKey, 'info');
        }
      } catch (error) {
        console.error(error);

        popups.error({ text: t('errors.default') });
      } finally {
        fetchTicketStatusJobsListWithPreviousParams();
        emit('clear:items');
      }
    };

    const { pending: playPending, handler: play } = usePendingHandler(() =>
      callTicketStatusJobActionWithPopups({
        action: () => playTicketStatusJobs({ data: pickIds(preparedItems.value) }),
        successTextKey: 'jobs.popups.restarted',
        errorTextKey: 'jobs.popups.restartFailed',
      }),
    );
    const { pending: repeatPending, handler: repeat } = usePendingHandler(() =>
      callTicketStatusJobActionWithPopups({
        action: () => playTicketStatusJobs({ data: pickIds(preparedItems.value) }),
        successTextKey: 'jobs.popups.repeated',
        errorTextKey: 'jobs.popups.repeatFailed',
      }),
    );
    const { pending: pausePending, handler: pause } = usePendingHandler(() =>
      callTicketStatusJobActionWithPopups({
        action: () => pauseTicketStatusJobs({ data: pickIds(preparedItems.value) }),
        successTextKey: 'jobs.popups.paused',
        errorTextKey: 'jobs.popups.pauseFailed',
      }),
    );
    const edit = () => {
      modals.show({
        name: MODALS.createTicketStatusJob,
        config: {
          ticketStatusJob: props.item,
          action: ticketStatusJob =>
            callActionWithPopup(
              () =>
                updateTicketStatusJob({
                  id: props.item._id,
                  data: ticketStatusJob,
                }),
              fetchTicketStatusJobsListWithPreviousParams,
              t('jobs.popups.updated'),
            ),
        },
      });
    };

    return {
      t,
      isRunning,
      isStopped,
      playPending,
      repeatPending,
      pausePending,
      play,
      repeat,
      pause,
      edit,
    };
  },
};
</script>
