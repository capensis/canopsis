<template>
  <c-action-btn
    v-if="icon"
    :tooltip="$t('eventsRecord.export')"
    :loading="downloading"
    icon="file_download"
    @click="exportJson"
  />
  <v-btn
    v-else
    :loading="downloading"
    color="primary"
    @click="exportJson"
  >
    <v-icon class="mr-2" color="white">
      file_download
    </v-icon>
    <span>{{ $t('common.exportToJson') }}</span>
  </v-btn>
</template>
<script>
import { computed, toRef } from 'vue';

import { useEventsRecordExportJson } from '../hooks/export-json';

export default {
  props: {
    eventsRecordId: {
      type: String,
      default: '',
    },
    eventId: {
      type: String,
      default: '',
    },
    eventIds: {
      type: Array,
      default: () => [],
    },
    icon: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const eventIds = computed(() => (props.eventId ? [props.eventId] : props.eventIds));

    const { downloading, exportJson } = useEventsRecordExportJson({
      eventIds,
      eventsRecordId: toRef(props, 'eventsRecordId'),
    });

    return {
      downloading,
      exportJson,
    };
  },
};
</script>
