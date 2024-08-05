<template>
  <c-advanced-data-table
    :headers="headers"
    :items="events"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    class="events-record-events-list"
    expand
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <v-expand-transition>
        <v-layout v-if="resending" class="gap-4" align-center>
          <span class="font-italic">Resending in progress</span>
          <v-progress-circular
            color="primary"
            width="3"
            indeterminate
          />
          <c-action-btn
            tooltip="Stop event resending"
            icon="stop"
            color="blue darken-3"
            top
            @click="stopResending"
          />
        </v-layout>
      </v-expand-transition>
    </template>
    <template #mass-actions="{ selected }">
      <c-action-btn
        :tooltip="$tc('eventsRecord.resendEvents', selected.length)"
        :disabled="resending"
        :loading="massResending"
        icon="play_arrow"
        color="#134A9F"
        @click="resendEvents(selected)"
      />
      <c-action-btn
        :tooltip="$t('eventsRecord.export')"
        icon="file_download"
        @click="exportJsonSelected(selected)"
      />
      <c-action-btn
        type="delete"
        @click="removeSelected(selected)"
      />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #actions="{ item }">
      <c-action-btn
        :tooltip="$tc('eventsRecord.resendEvents', 1)"
        :disabled="resending"
        :loading="resendingByIds[item._id]"
        icon="play_arrow"
        color="blue darken-3"
        @click="resendEvent(item)"
      />
      <c-action-btn
        :tooltip="$t('eventsRecord.export')"
        :disabled="downloading"
        icon="file_download"
        @click="exportJson(item)"
      />
      <c-action-btn
        type="delete"
        @click="remove(item)"
      />
    </template>
    <template #expand="{ item }">
      <div class="secondary pa-3">
        <v-card>
          <v-card-text>
            <c-json-treeview :json="item | json" />
          </v-card-text>
        </v-card>
      </div>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, ref, set, unref } from 'vue';

import { MODALS, TIME_UNITS } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePolling } from '@/hooks/polling';
import { useEventsRecord } from '@/hooks/store/modules/events-record';

const useMassActionPending = (action) => {
  const pendingByIds = ref({});
  const massPending = ref(false);

  const setPending = (itemsIds, value) => {
    itemsIds.forEach(id => set(pendingByIds.value, id, value));

    if (itemsIds.length > 1) {
      massPending.value = value;
    }
  };

  const preparedAction = async (itemsIds = [], ...rest) => {
    try {
      setPending(itemsIds, true);

      await action(itemsIds, ...rest);
    } finally {
      setPending(itemsIds, false);
    }
  };

  const pending = computed(() => massPending.value || Object.values(pendingByIds.value).some(value => value));

  return {
    pending,
    pendingByIds,
    massPending,

    action: preparedAction,
  };
};

export const usePlayback = ({ eventsRecordId }) => {
  const {
    playbackEventsRecordEvents,
  } = useEventsRecord();

  const startHandler = (eventIds, delay) => (
    playbackEventsRecordEvents({ id: unref(eventsRecordId), data: { delay, event_ids: eventIds } })
  );
  const processHandler = () => {};
  const endHandler = () => {};

  const {
    poll,
    cancel,
  } = usePolling({
    startHandler,
    processHandler,
    endHandler,
  });

  const {
    pending: resending,
    pendingByIds: resendingByIds,
    massPending: massResending,

    action: startResending,
  } = useMassActionPending((eventIds, delay) => poll(eventIds, delay));

  const stopResending = () => cancel();

  return {
    resending,
    resendingByIds,
    massResending,

    startResending,
    stopResending,
  };
};

export default {
  props: {
    eventsRecordId: {
      type: String,
      required: true,
    },
    events: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t, tc } = useI18n();
    const modals = useModals();

    const {
      playbackEventsRecordEvents,
    } = useEventsRecord();

    const {
      pending: resending,
      pendingByIds: resendingByIds,
      massPending: massResending,

      action: resend,
    } = useMassActionPending((eventIds, delay) => (
      playbackEventsRecordEvents({ id: props.eventsRecordId, data: { delay, event_ids: eventIds } })
    ));

    const showResendEventsModalWindow = (eventIds = []) => modals.show({
      name: MODALS.duration,
      config: {
        title: tc('eventsRecord.resendEvents', 1),
        label: t('eventsRecord.delayBetweenEvents'),
        units: [
          { value: TIME_UNITS.millisecond, text: 'common.times.millisecond' },
          { value: TIME_UNITS.second, text: 'common.times.second' },
        ],
        action: delay => resend(eventIds, delay),
      },
    });

    const resendEvent = event => showResendEventsModalWindow([event._id]);
    const resendEvents = events => showResendEventsModalWindow(mapIds(events));

    const stopResending = () => {};

    const headers = computed(() => [
      {
        text: t('common.timestamp'),
        value: 'created',
        sortable: false,
      },
      {
        text: t('common.eventType'),
        value: 'event_type',
        sortable: false,
      },
      {
        text: t('common.sourceType'),
        value: 'source_type',
        sortable: false,
      },
      {
        text: t('common.connector'),
        value: 'connector',
        sortable: false,
      },
      {
        text: t('common.connectorName'),
        value: 'connector_name',
        sortable: false,
      },
      {
        text: t('common.component'),
        value: 'component',
        sortable: false,
      },
      {
        text: t('common.resource'),
        value: 'resource',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const remove = event => emit('remove', event._id);
    const removeSelected = selected => emit('remove:selected', selected);

    const exportJsonSelected = selected => emit('export', selected);
    const exportJson = event => exportJsonSelected([event._id]);

    const updateOptions = options => emit('update:options', options);

    return {
      headers,

      resending,
      resendingByIds,
      massResending,
      resendEvent,
      resendEvents,
      stopResending,

      remove,
      removeSelected,

      exportJson,
      exportJsonSelected,

      updateOptions,
    };
  },
};
</script>

<style lang="scss">
.events-record-events-list {
  table {
    table-layout: fixed;

    .v-data-table-header th {
      &:nth-child(1) {
        width: 56px !important;
      }

      &:nth-child(2) {
        width: 32px !important;
      }

      &:last-child {
        width: 170px !important;
      }
    }
  }
}
</style>
