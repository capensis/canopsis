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
    <template #resending="{ item }">
      <v-tooltip right>
        <template #activator="{ on }">
          <v-icon
            v-show="item.is_resending"
            class="blinking"
            color="blue darken-3"
            v-on="on"
          >
            play_arrow
          </v-icon>
        </template>
        <span>{{ $t('eventsRecord.resendingInProgress') }}</span>
      </v-tooltip>
    </template>
    <template #t="{ item }">
      {{ item.t | date }}
    </template>
    <template #actions="{ item }">
      <v-layout v-if="item.is_recording">
        <c-action-btn
          :tooltip="$t('eventsRecord.stop')"
          color="blue darken-3"
          icon="stop"
          @click="stop"
        />
      </v-layout>
      <v-layout v-else>
        <c-action-btn
          :tooltip="$t('eventsRecord.viewEvents')"
          icon="pageview"
          color="#6A6A6A"
          @click="show(item)"
        />
        <events-record-download-btn :events-record-id="item._id" icon />
        <c-action-btn
          type="delete"
          @click="remove(item)"
        />
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
        value: 'resending',
        sortable: false,
      },
      {
        text: t('eventsRecord.recorded'),
        value: 't',
        sortable: false,
      },
      {
        text: t('eventsRecord.eventsCount'),
        value: 'count',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const stop = () => emit('stop');
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
