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
    @update:options="$emit('update:options', $event)"
  >
    <template #toolbar>
      <v-expand-transition>
        <c-layout v-if="resending" class="gap-4" align="center">
          <span class="font-italic">
            {{ t('eventsRecord.resendingInProgress') }}
          </span>
          <v-progress-circular color="primary" width="3" indeterminate />
          <c-action-btn
            :tooltip="t('eventsRecord.stopResending')"
            icon="stop"
            color="primary"
            @click="$emit('stop:resending')"
          />
        </c-layout>
      </v-expand-transition>
    </template>

    <template #mass-actions="{ selected, selectedKeys, count }">
      <c-layout align="center">
        <c-action-btn
          :tooltip="t('eventsRecord.resendEvents', count)"
          :disabled="resendingDisabled || resending"
          icon="play_arrow"
          color="primary"
          @click="$emit('start:resending', selectedKeys)"
        />
        <events-record-download-btn
          :events-record-id="eventsRecordId"
          :event-ids="selectedKeys"
          icon
        />
        <c-action-btn type="delete" @click="$emit('remove:selected', selected)" />
      </c-layout>
    </template>

    <template #timestamp="{ item }">
      {{ formatDate(item.timestamp) }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
        <c-action-btn
          :tooltip="t('eventsRecord.resendEvents', 1)"
          :disabled="resendingDisabled || resending"
          icon="play_arrow"
          color="primary"
          @click="$emit('start:resending', [item._id])"
        />
        <c-copy-btn
          :value="stringify(item)"
          :tooltip="t('eventsRecord.copyEvent')"
          @success="showCopySuccessPopup"
          @error="showCopyErrorPopup"
        />
        <events-record-download-btn :events-record-id="eventsRecordId" :event-id="item._id" icon />
        <c-action-btn type="delete" @click="$emit('remove', item)" />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <events-record-events-list-expand-panel :event="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue'

import { convertDateToString } from '@/helpers/date/date'

import { useI18n } from '@/hooks/i18n'
import { usePopups } from '@/hooks/popups'

import EventsRecordEventsListExpandPanel from './partials/events-record-events-list-expand-panel.vue'
import EventsRecordDownloadBtn from './partials/events-record-download-btn.vue'

export default {
  components: {
    EventsRecordEventsListExpandPanel,
    EventsRecordDownloadBtn,
  },

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
    resending: {
      type: Boolean,
      default: false,
    },
    resendingDisabled: {
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

  emits: ['update:options', 'remove', 'remove:selected', 'start:resending', 'stop:resending'],

  setup() {
    const { t } = useI18n()
    const popups = usePopups()

    const headers = computed(() => [
      { text: t('common.timestamp'), value: 'timestamp', sortable: false },
      { text: t('common.eventType'), value: 'event_type', sortable: false },
      { text: t('common.sourceType'), value: 'source_type', sortable: false },
      { text: t('common.connector'), value: 'connector', sortable: false },
      { text: t('common.connectorName'), value: 'connector_name', sortable: false },
      { text: t('common.component'), value: 'component', sortable: false },
      { text: t('common.resource'), value: 'resource', sortable: false },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ])

    const formatDate = (date) => convertDateToString(date)
    const stringify = (item) => JSON.stringify(item)
    const showCopySuccessPopup = () => popups.success({ text: t('eventsRecord.eventCopied') })
    const showCopyErrorPopup = () => popups.error({ text: t('errors.default') })

    return {
      t,
      headers,
      formatDate,
      stringify,
      showCopySuccessPopup,
      showCopyErrorPopup,
    }
  },
}
</script>

<style lang="scss">
.events-record-events-list {
  table {
    table-layout: fixed;
  }
}
</style>
