<template>
  <c-layout class="gap-3 my-4" justify="space-between" align="center">
    <c-layout class="gap-2" align="center">
      <v-btn v-if="hasFilterApplied" color="secondary" @click="$emit('apply:filter')">
        <v-icon class="mr-2"> filter_list </v-icon>
        {{ t('eventsRecord.eventFilterApplied') }}
      </v-btn>
      <v-btn v-else color="secondary" variant="outlined" @click="$emit('apply:filter')">
        {{ t('eventsRecord.applyEventFilter') }}
      </v-btn>

      <v-btn
        v-if="hasFilterApplied"
        color="secondary"
        variant="outlined"
        @click="$emit('reset:filter')"
      >
        {{ t('common.reset') }}
      </v-btn>
    </c-layout>

    <v-spacer />

    <events-record-download-btn
      :events-record-id="eventsRecordId"
      :event-pattern="query.event_pattern"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue'

import { useI18n } from '@/hooks/i18n'

import EventsRecordDownloadBtn from './partials/events-record-download-btn.vue'

export default {
  components: { EventsRecordDownloadBtn },

  props: {
    eventsRecordId: {
      type: String,
      required: true,
    },
    query: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['apply:filter', 'reset:filter'],

  setup(props) {
    const { t } = useI18n()
    const hasFilterApplied = computed(() => props.query.event_pattern !== undefined)

    return {
      t,
      hasFilterApplied,
    }
  },
}
</script>
