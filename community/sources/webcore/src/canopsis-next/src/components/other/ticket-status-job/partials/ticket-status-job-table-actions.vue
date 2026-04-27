<template>
  <v-layout>
    <c-action-btn
      v-if="item"
      :tooltip="$t('jobs.actions.editJob')"
      icon="edit"
      @click="edit"
    />
    <template v-if="isRunning">
      <c-action-btn
        :tooltip="$t('jobs.actions.pauseJob')"
        :disabled="pausePending"
        :loading="pausePending"
        icon="pause"
        @click="pause"
      />
      <c-action-btn
        :tooltip="$t('jobs.actions.repeatJob')"
        :disabled="repeatPending"
        :loading="repeatPending"
        icon="refresh"
        @click="repeat"
      />
    </template>
    <c-action-btn
      v-else
      :tooltip="$t('jobs.actions.startJob')"
      :disabled="playPending"
      :loading="playPending"
      icon="play_arrow"
      @click="play"
    />
  </v-layout>
</template>

<script>
import { keyBy } from 'lodash';
import { computed } from 'vue';

import { MODALS, JOB_STATUS, JOB_ACTION_REFETCH_DELAY } from '@/constants';

import { pickIds } from '@/helpers/array';
import { promisedWait } from '@/helpers/async';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useCallActionWithPopup } from '@/hooks/actions/call';
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
  setup(props) {
    const { t, tc } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const {
      updateTicketStatusJob,
      playTicketStatusJob,
      pauseTicketStatusJob,
      fetchTicketStatusJobsListWithPreviousParams,
    } = useTicketStatusJob();
    const { callActionWithPopup } = useCallActionWithPopup(false, 'info');

    const preparedItems = computed(() => (props.item ? [props.item] : props.items));
    const preparedItemsById = computed(() => keyBy(preparedItems.value, '_id'));
    const isRunning = computed(() => preparedItems.value.some(item => item.status === JOB_STATUS.running));

    /**
     * Displays a popup with a ticket status job action message.
     *
     * @param {Array} [items=[]] - Result items (used for count and first item's rule/ticket info)
     * @param {string} textKey - i18n key for the message (supports pluralization via count)
     * @param {'info'|'error'} [popupType='info'] - Popup type to display
     */
    const showTicketStatusJobActionPopup = (items = [], textKey, popupType = 'info') => {
      const { rule_name: ruleName, ticket_id: ticketNumber } = preparedItemsById.value[items[0]?.item?._id];

      popups[popupType]({
        text: `<span class="font-weight-regular">${tc(textKey, items.length, { ruleName, ticketNumber, count: items.length })}</span>`,
      });
    };

    /**
     * Executes a ticket status job action and shows success/error popups based on the result.
     * Splits the action result into successful and failed items, displays appropriate messages, and refreshes the list.
     *
     * @param {Object} options - Configuration options
     * @param {Function} options.action - Async function returning result items (each may have `error` property)
     * @param {string} options.successTextKey - i18n key for the success message (supports pluralization)
     * @param {string} options.errorTextKey - i18n key for the error message (supports pluralization)
     */
    const callTicketStatusJobActionWithPopups = async ({ action, successTextKey, errorTextKey }) => {
      try {
        const result = await action();

        /**
         * Wait for the refetch delay to ensure the list is updated before showing the popups and refetching
         * because of backend logic.
         */
        await promisedWait(JOB_ACTION_REFETCH_DELAY);

        const { successItems, failedItems } = result.reduce((acc, item) => {
          acc[item.error ? 'failedItems' : 'successItems'].push(item);

          return acc;
        }, { successItems: [], failedItems: [] });

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
      }
    };

    /**
     * Starts or restarts the ticket status job(s).
     */
    const { pending: playPending, handler: play } = usePendingHandler(() => callTicketStatusJobActionWithPopups({
      action: () => playTicketStatusJob({ data: pickIds(preparedItems.value) }),
      successTextKey: 'jobs.popups.restarted',
      errorTextKey: 'jobs.popups.restartFailed',
    }));

    /**
     * Repeats the ticket status job(s) (when running).
     */
    const { pending: repeatPending, handler: repeat } = usePendingHandler(() => callTicketStatusJobActionWithPopups({
      action: () => playTicketStatusJob({ data: pickIds(preparedItems.value) }),
      successTextKey: 'jobs.popups.repeated',
      errorTextKey: 'jobs.popups.repeatFailed',
    }));

    /**
     * Pauses the ticket status job(s).
     */
    const { pending: pausePending, handler: pause } = usePendingHandler(() => callTicketStatusJobActionWithPopups({
      action: () => pauseTicketStatusJob({ data: pickIds(preparedItems.value) }),
      successTextKey: 'jobs.popups.paused',
      errorTextKey: 'jobs.popups.pauseFailed',
    }));

    /**
     * Opens the edit modal for a ticket status job.
     *
     * @param {Object} item - The ticket status job to edit
     */
    const edit = () => modals.show({
      name: MODALS.createTicketStatusJob,
      config: {
        ticketStatusJob: props.item,
        action: newTicketStatusJob => callActionWithPopup(
          () => updateTicketStatusJob({ id: props.item._id, data: newTicketStatusJob }),
          fetchTicketStatusJobsListWithPreviousParams,
          t('jobs.popups.updated'),
        ),
      },
    });

    return {
      preparedItems,
      isRunning,
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
