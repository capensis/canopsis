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
      {{ formatDate(item.created) }}
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item)" />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

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

  emits: ['edit', 'remove', 'refresh', 'update:options'],

  setup() {
    const { t } = useI18n();
    const formatDate = value => convertDateToString(value);

    const headers = computed(() => [
      {
        text: t('common.connectorName'),
        value: 'name',
      },
      {
        text: t('common.enabled'),
        value: 'enabled',
      },
      {
        text: t('common.created'),
        value: 'created',
      },
      {
        text: t('common.updated'),
        value: 'updated',
      },
      {
        text: t('common.author'),
        value: 'author.display_name',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      headers,
      formatDate,
    };
  },
};
</script>
