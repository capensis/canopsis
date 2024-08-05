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
      <v-icon
        v-show="item.isResending"
        class="blinking"
        color="blue darken-3"
      >
        play_arrow
      </v-icon>
    </template>
    <template #t="{ item }">
      {{ item.t | date }}
    </template>
    <template #actions="{ item }">
      <v-layout v-if="item.isRecording">
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
        <c-action-btn
          :tooltip="$t('eventsRecord.export')"
          :loading="downloadingsById[item._id]"
          icon="file_download"
          @click="exportJson(item)"
        />
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

import { useExportJson } from './hooks/export-json';

export default {
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

    const inProgress = computed(() => props.eventsRecords.status === 0);
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
        value: 'c',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    /**
     * STORE
     */
    const { downloadingsById, exportJson: exportJsonMethod } = useExportJson();

    const stop = () => emit('stop');
    const show = eventsRecord => emit('show', eventsRecord);
    const remove = eventsRecord => emit('remove', eventsRecord._id);
    const exportJson = eventsRecord => exportJsonMethod(eventsRecord._id);
    const updateOptions = options => emit('update:options', options);

    return {
      inProgress,
      headers,
      downloadingsById,

      stop,
      show,
      exportJson,
      remove,
      updateOptions,
    };
  },
};
</script>
