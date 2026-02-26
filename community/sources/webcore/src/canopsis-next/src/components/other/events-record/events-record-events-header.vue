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
    <v-layout class="gap-3 my-4" justify-space-between align-center>
      <v-layout class="gap-2" align-center>
        <v-btn
          v-if="hasFilterApplied"
          color="secondary"
          @click="applyEventFilter"
        >
          <v-icon class="mr-2">
            filter_list
          </v-icon>
          {{ $t('eventsRecord.eventFilterApplied') }}
        </v-btn>
        <v-btn
          v-else
          color="secondary"
          outlined
          @click="applyEventFilter"
        >
          {{ $t('eventsRecord.applyEventFilter') }}
        </v-btn>

        <v-btn
          v-if="hasFilterApplied"
          color="secondary"
          outlined
          @click="resetFilter"
        >
          {{ $t('common.reset') }}
        </v-btn>
      </v-layout>
      <v-spacer />
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
    hasFilterApplied: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const remove = () => emit('remove');
    const applyEventFilter = () => emit('apply:filter');
    const resetFilter = () => emit('reset:filter');

    return {
      remove,
      applyEventFilter,
      resetFilter,
    };
  },
};
</script>
