<template>
  <c-action-btn
    v-if="icon"
    :tooltip="t('eventsRecord.export')"
    :loading="downloading"
    icon="file_download"
    @click="download"
  />
  <v-btn
    v-else
    :loading="downloading"
    color="primary"
    variant="outlined"
    @click="download"
  >
    <v-icon class="mr-2">
      file_download
    </v-icon>
    <span>{{ t('common.exportToJson') }}</span>
  </v-btn>
</template>

<script>
import { computed, toRef } from 'vue';

import { useI18n } from '@/hooks/i18n';

import { useEventsRecordExportJson } from '../hooks/export-json';

export default {
  emits: ['exported'],

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
    eventPattern: {
      type: Array,
      required: false,
    },
  },

  setup(props, { emit }) {
    const { t } = useI18n();
    const eventIds = computed(() => (props.eventId ? [props.eventId] : props.eventIds));

    const { downloading, exportJson } = useEventsRecordExportJson({
      eventIds,
      eventPattern: toRef(props, 'eventPattern'),
      eventsRecordId: toRef(props, 'eventsRecordId'),
    });
    const download = async () => {
      await exportJson();
      emit('exported');
    };

    return {
      t,
      downloading,
      download,
    };
  },
};
</script>
