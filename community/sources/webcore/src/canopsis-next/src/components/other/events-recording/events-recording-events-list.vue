<template>
  <c-advanced-data-table
    :headers="headers"
    :items="eventsRecording"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #actions="{ item }">
      <v-layout v-if="!item.inProgress">
        <c-action-btn
          :tooltip="$t('eventsRecording.resendEvent')"
          icon="play"
          color="#134A9F"
          @click="show(item)"
        />
        <c-action-btn
          :tooltip="$t('eventsRecording.export')"
          icon="file_download"
          @click="exportJson(item)"
        />
        <c-action-btn
          type="delete"
          @click="remove(item)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <div class="secondary pa-3">
        <c-json-treeview :json="item" />
      </div>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    eventsRecording: {
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
    const inProgress = computed(() => props.eventsRecording.status === 0);
    const headers = computed(() => [
      {
        text: t('common.timestamp'),
        value: 'created',
        sortable: false,
      },
      {
        text: t('common.eventType'),
        value: 'event_type',
        sortable: false,
      },
      {
        text: t('common.sourceType'),
        value: 'source_type',
        sortable: false,
      },
      {
        text: t('common.connector'),
        value: 'connector',
        sortable: false,
      },
      {
        text: t('common.connectorName'),
        value: 'connector_name',
        sortable: false,
      },
      {
        text: t('common.component'),
        value: 'component',
        sortable: false,
      },
      {
        text: t('common.resource'),
        value: 'resource',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const show = eventsRecording => emit('show', eventsRecording._id);
    const exportJson = eventsRecording => emit('export', eventsRecording._id);
    const remove = eventsRecording => emit('remove', eventsRecording._id);
    const updateOptions = options => emit('update:options', options);

    return {
      inProgress,
      headers,

      show,
      exportJson,
      remove,
      updateOptions,
    };
  },
};
</script>
