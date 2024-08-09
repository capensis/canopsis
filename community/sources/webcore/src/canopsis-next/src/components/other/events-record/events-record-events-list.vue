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
    @update:options="updateOptions"
  >
    <template #toolbar>
      <v-expand-transition>
        <v-layout v-if="resending" class="gap-4" align-center>
          <span class="font-italic">{{ $t('eventsRecord.resendingInProgress') }}</span>
          <v-progress-circular
            color="primary"
            width="3"
            indeterminate
          />
          <c-action-btn
            :tooltip="$t('eventsRecord.stopResending')"
            icon="stop"
            color="blue darken-3"
            top
            @click="stopResending"
          />
        </v-layout>
      </v-expand-transition>
    </template>
    <template #mass-actions="{ selected, selectedKeys, count }">
      <c-action-btn
        :tooltip="$tc('eventsRecord.resendEvents', count)"
        :disabled="resendingDisabled || resending"
        icon="play_arrow"
        color="#134A9F"
        @click="startResending(selectedKeys)"
      />
      <events-record-download-btn :events-record-id="eventsRecordId" :event-ids="selectedKeys" icon />
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
        :tooltip="$tc('eventsRecord.resendEvents', 1)"
        :disabled="resendingDisabled || resending"
        icon="play_arrow"
        color="blue darken-3"
        @click="startResending([item._id])"
      />
      <events-record-download-btn :events-record-id="eventsRecordId" :event-id="item._id" icon />
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

import EventsRecordDownloadBtn from '@/components/other/events-record/partials/events-record-download-btn.vue';

export default {
  components: { EventsRecordDownloadBtn },
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
    downloading: {
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

    /**
     * RESEND
     */
    const startResending = eventsIds => emit('start:resending', eventsIds);
    const stopResending = () => emit('stop:resending');

    /**
     * REMOVE
     */
    const remove = event => emit('remove', event);
    const removeSelected = eventsIds => emit('remove:selected', eventsIds);

    /**
     * OPTIONS
     */
    const updateOptions = options => emit('update:options', options);

    return {
      headers,

      startResending,
      stopResending,

      remove,
      removeSelected,

      updateOptions,
    };
  },
};
</script>

<style lang="scss">
.events-record-events-list {
  table {
    table-layout: fixed;

    .v-data-table-header th {
      &:nth-child(1) {
        width: 56px !important;
      }

      &:nth-child(2) {
        width: 32px !important;
      }

      &:last-child {
        width: 170px !important;
      }
    }
  }
}
</style>
