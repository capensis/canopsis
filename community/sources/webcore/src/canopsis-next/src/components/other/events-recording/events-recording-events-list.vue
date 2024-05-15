<template>
  <c-advanced-data-table
    :headers="headers"
    :items="events"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        :tooltip="$tc('eventsRecording.resendEvents', selected.length)"
        icon="play_arrow"
        color="#134A9F"
        @click="resendSelected(selected)"
      />
      <c-action-btn
        :tooltip="$t('eventsRecording.export')"
        icon="file_download"
        @click="exportJsonSelected(selected)"
      />
      <c-action-btn
        type="delete"
        @click="removeSelected(selected)"
      />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #actions="{ item }">
      <c-action-btn
        :tooltip="$tc('eventsRecording.resendEvents', 1)"
        icon="play_arrow"
        color="blue darken-3"
        @click="resend(item)"
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
    </template>
    <template #expand="{ item }">
      <div class="secondary pa-3">
        <v-card>
          <v-card-text>
            <c-json-treeview :json="item | json" />
          </v-card-text>
        </v-card>
      </div>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    events: {
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

    const resend = event => emit('resend', event);
    const exportJson = event => emit('export', event);
    const remove = event => emit('remove', event._id);
    const resendSelected = selected => emit('resend:selected', selected);
    const exportJsonSelected = selected => emit('export:selected', selected);
    const removeSelected = selected => emit('remove:selected', selected);
    const updateOptions = options => emit('update:options', options);

    return {
      headers,

      resend,
      exportJson,
      remove,
      resendSelected,
      exportJsonSelected,
      removeSelected,
      updateOptions,
    };
  },
};
</script>
