<template>
  <c-advanced-data-table
    :headers="headers"
    :items="eventsRecords"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #loader="{ item }">
      <v-layout class="gap-2">
        <v-tooltip v-if="item.is_recording" right>
          <template #activator="{ on }">
            <v-progress-circular
              size="24"
              width="2"
              indeterminate
              v-on="on"
            />
          </template>
          <span>{{ $t('eventsRecord.inProgress') }}</span>
        </v-tooltip>
        <v-tooltip v-if="item.is_resending" right>
          <template #activator="{ on }">
            <v-icon
              class="blinking"
              v-on="on"
            >
              play_arrow
            </v-icon>
          </template>
          <span>{{ $t('eventsRecord.resendingInProgress') }}</span>
        </v-tooltip>
      </v-layout>
    </template>
    <template #t="{ item }">
      {{ item.t | date }}
    </template>
    <template #recordingEnd="{ item }">
      <span v-if="item.is_recording">
        {{ $t('common.inProgress') }}
      </span>
      <span v-else>{{ item.t_end | date('long', '-') }}</span>
    </template>
    <template #duration="{ item }">
      <span v-if="item.is_recording">
        {{ $t('common.inProgress') }}
      </span>
      <span v-else>
        {{ item.duration | duration }}
      </span>
    </template>
    <template #count="{ item }">
      <span>{{ item.count }}</span>
      <span v-if="item.is_recording">...</span>
    </template>
    <template #author="{ item }">
      {{ item.author || '—' }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="item.is_recording"
          :tooltip="$t('eventsRecord.stop')"
          color="blue darken-3"
          icon="stop"
          @click="stop(item._id)"
        />
        <c-action-btn
          :tooltip="item.is_recording ? $t('eventsRecord.viewPattern') : $t('eventsRecord.viewEventsAndPattern')"
          icon="pageview"
          color="#6A6A6A"
          @click="show(item)"
        />
        <template v-if="!item.is_recording">
          <events-record-download-btn :events-record-id="item._id" icon />
          <c-action-btn
            type="delete"
            @click="remove(item)"
          />
        </template>
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import EventsRecordDownloadBtn from './partials/events-record-download-btn.vue';

export default {
  components: { EventsRecordDownloadBtn },
  props: {
    eventsRecords: {
      type: Array,
      required: true,
    },
    pending: {
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
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: '',
        value: 'loader',
        sortable: false,
      },
      {
        text: t('eventsRecord.recordingStart'),
        value: 't',
        sortable: false,
      },
      {
        text: t('eventsRecord.recordingEnd'),
        value: 'recordingEnd',
        sortable: false,
      },
      {
        text: t('eventsRecord.duration'),
        value: 'duration',
        sortable: false,
      },
      {
        text: t('eventsRecord.eventsCount'),
        value: 'count',
        sortable: false,
      },
      {
        text: t('eventsRecord.launchedBy'),
        value: 'author',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const stop = item => emit('stop', item);
    const show = eventsRecord => emit('show', eventsRecord);
    const remove = eventsRecord => emit('remove', eventsRecord._id);
    const updateOptions = options => emit('update:options', options);

    return {
      headers,

      stop,
      show,
      remove,
      updateOptions,
    };
  },
};
</script>
