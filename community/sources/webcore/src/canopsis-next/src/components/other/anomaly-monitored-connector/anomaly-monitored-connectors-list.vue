<template>
  <c-advanced-data-table
    :headers="headers"
    :items="anomalyMonitoredConnectors"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="updatable || removable"
    item-key="id"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        item-id="id"
        anomaly-monitored-connector
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item)"
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
    anomalyMonitoredConnectors: {
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
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.connectorName'),
        value: 'name',
        sortable: false,
      },
      {
        text: t('common.enabled'),
        value: 'enabled',
        sortable: false,
      },
      {
        text: t('common.created'),
        value: 'created',
        sortable: false,
      },
      {
        text: t('common.updated'),
        value: 'updated',
        sortable: false,
      },
      {
        text: t('common.author'),
        value: 'author.display_name',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      headers,
    };
  },
};
</script>
