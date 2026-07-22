<template>
  <c-advanced-data-table
    :headers="headers"
    :items="eventsRecords"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #loader="{ item }">
      <v-tooltip v-if="item.is_recording" location="right">
        <template #activator="{ props: tooltipProps }">
          <v-progress-circular v-bind="tooltipProps" color="primary" width="3" indeterminate />
        </template>
        <span>{{ item.progress }}</span>
      </v-tooltip>
      <v-icon v-else-if="item.is_resending" class="c-advanced-data-table__blink" color="primary">
        play_arrow
      </v-icon>
    </template>

    <template #t="{ item }">
      {{ formatDate(item.t) }}
    </template>

    <template #e="{ item }">
      <span v-if="item.is_recording">{{ t('common.inProgress') }}</span>
      <span v-else>{{ formatDate(item.e, '-') }}</span>
    </template>

    <template #d="{ item }">
      <span v-if="item.is_recording">{{ t('common.inProgress') }}</span>
      <span v-else>{{ formatDuration(item.d) }}</span>
    </template>

    <template #count="{ item }"> {{ item.count }}{{ item.is_recording ? '...' : '' }} </template>

    <template #author="{ item }">
      {{ item.author || '-' }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
        <c-action-btn
          :tooltip="
            t(item.is_recording ? 'eventsRecord.viewPattern' : 'eventsRecord.viewEventsAndPattern')
          "
          type="view"
          @click="$emit('show', item)"
        />
        <c-action-btn
          v-if="item.is_recording"
          :tooltip="t('eventsRecord.stop')"
          icon="stop"
          color="primary"
          @click="$emit('stop', item._id)"
        />
        <template v-else>
          <events-record-download-btn :events-record-id="item._id" icon />
          <c-action-btn type="delete" @click="$emit('remove', item._id)" />
        </template>
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue'

import { convertDateToString } from '@/helpers/date/date'
import { convertDurationToString } from '@/helpers/date/duration'

import { useI18n } from '@/hooks/i18n'

import EventsRecordDownloadBtn from './partials/events-record-download-btn.vue'

export default {
  components: {
    EventsRecordDownloadBtn,
  },

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

  emits: ['update:options', 'stop', 'show', 'remove'],

  setup() {
    const { t } = useI18n()

    const headers = computed(() => [
      { text: '', value: 'loader', sortable: false },
      { text: t('eventsRecord.recordingStart'), value: 't', sortable: false },
      { text: t('eventsRecord.recordingEnd'), value: 'e', sortable: false },
      { text: t('eventsRecord.duration'), value: 'd', sortable: false },
      { text: t('eventsRecord.eventsCount'), value: 'count', sortable: false },
      { text: t('eventsRecord.launchedBy'), value: 'author', sortable: false },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ])

    const formatDate = (date, defaultValue = '') =>
      convertDateToString(date, undefined, defaultValue)
    const formatDuration = (duration) => convertDurationToString(duration)

    return {
      t,
      headers,
      formatDate,
      formatDuration,
    }
  },
}
</script>
