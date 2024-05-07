<template>
  <c-advanced-data-table
    :headers="headers"
    :items="eventsRecording"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #actions="{ item }">
      <v-layout v-if="!item.inProgress">
        <c-action-btn
          :tooltip="$t('eventsRecording.viewEvents')"
          icon="pageview"
          color="#6A6A6A"
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
        text: t('eventsRecording.recorded'),
        value: 'created',
        sortable: false,
      },
      {
        text: t('eventsRecording.eventsCount'),
        value: 'count',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const show = eventsRecording => emit('show', eventsRecording);
    const exportJson = eventsRecording => emit('export', eventsRecording);
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
