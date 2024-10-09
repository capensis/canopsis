<template>
  <v-layout class="gap-2" column>
    <v-layout class="gap-3 my-4" justify-center align-center>
      <span class="text-subtitle-2">{{ $t('modals.eventsRecord.subtitle', { count }) }}</span>
      <c-action-btn
        :tooltip="$t('modals.eventsRecord.buttonTooltip')"
        type="delete"
        @click="remove"
      />
    </v-layout>
    <v-layout class="gap-3 my-4">
      <v-btn
        color="primary"
        @click="applyEventFilter"
      >
        {{ $t('eventsRecord.applyEventFilter') }}
      </v-btn>
      <events-record-download-btn :events-record-id="eventsRecordId" />
    </v-layout>
  </v-layout>
</template>
<script>
import EventsRecordDownloadBtn from './partials/events-record-download-btn.vue';

export default {
  components: { EventsRecordDownloadBtn },
  props: {
    eventsRecordId: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      default: 0,
    },
  },
  setup(props, { emit }) {
    const remove = () => emit('remove');
    const applyEventFilter = () => emit('apply:filter');

    return {
      remove,
      applyEventFilter,
    };
  },
};
</script>
