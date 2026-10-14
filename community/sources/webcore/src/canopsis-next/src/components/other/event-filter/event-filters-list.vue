<template>
  <c-advanced-data-table
    ref="advancedDataTable"
    :items="eventFilters"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable || updatable"
    advanced-pagination
    search
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        event-filter
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #type="{ item }">
      {{ $t(`eventFilter.types.${item.type}`) }}
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled ?? true" />
    </template>
    <template #unread_failures_count="{ item }">
      <c-circle-badge v-if="item.unread_failures_count" color="error" small>
        {{ item.unread_failures_count }}
      </c-circle-badge>
      <template v-else>
        -
      </template>
    </template>
    <template #created="{ item }">
      {{ dateToString(item.created) }}
    </template>
    <template #updated="{ item }">
      {{ dateToString(item.updated) }}
    </template>
    <template #calendar="{ item }">
      <v-icon>{{ isCalendarRule(item) ? 'check' : 'clear' }}</v-icon>
    </template>
    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-db-export-btn :id="item._id" event-filter />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>
    <template #expand="{ item }">
      <event-filters-list-expand-panel
        ref="expandPanel"
        :event-filter="item"
        @refresh="$emit('refresh')"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import EventFiltersListExpandPanel from './partials/event-filters-list-expand-panel.vue';

export default {
  components: {
    EventFiltersListExpandPanel,
  },
  emits: ['update:options', 'refresh', 'edit', 'duplicate', 'remove'],
  props: {
    eventFilters: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.id'), value: '_id' },
      { text: t('common.type'), value: 'type' },
      { text: t('common.priority'), value: 'priority' },
      { text: t('common.enabled'), value: 'enabled' },
      { text: t('common.error', 2), value: 'unread_failures_count' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.updated'), value: 'updated' },
      { text: t('common.calendar'), value: 'calendar', sortable: false },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    /**
     * Check if event filter item has calendar rules configured
     *
     * @param {Object} item - Event filter item
     * @param {*} [item.start] - Start date/time for calendar rule
     * @param {*} [item.rrule] - Recurrence rule (RRULE)
     * @param {Array} [item.exdates] - Array of excluded dates
     * @param {Array} [item.exceptions] - Array of exception dates
     * @return {boolean} True if item has any calendar-related properties configured
     */
    const isCalendarRule = item => item.start || item.rrule || item.exdates?.length || item.exceptions?.length;
    const dateToString = value => convertDateToString(value);

    return {
      headers,
      isCalendarRule,
      dateToString,
    };
  },
};
</script>
