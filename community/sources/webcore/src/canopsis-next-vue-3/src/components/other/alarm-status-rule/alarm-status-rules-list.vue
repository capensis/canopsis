<template>
  <c-advanced-data-table
    :headers="headers"
    :items="rules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable || updatable"
    expand
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
        :flapping-rule="flapping"
        :resolve-rule="!flapping"
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #duration="{ item }">
      <span>{{ durationToString(item.duration) }}</span>
    </template>

    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>

    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <alarm-status-rules-list-expand-item :rule="item" :flapping="flapping" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

import AlarmStatusRulesListExpandItem from './partials/alarm-status-rules-list-expand-item.vue';

export default {
  components: {
    AlarmStatusRulesListExpandItem,
  },

  props: {
    rules: {
      type: Array,
      required: true,
    },
    flapping: {
      type: Boolean,
      default: false,
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

  emits: ['update:options', 'edit', 'duplicate', 'remove', 'refresh'],

  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => [
        { text: t('common.id'), value: '_id' },
        { text: t('common.name'), value: 'name' },
        { text: t('common.duration'), value: 'duration', sortable: false },
        { text: t('common.priority'), value: 'priority' },
        { text: t('common.enabled'), value: 'enabled', sortable: false },
        props.flapping && { text: t('common.frequencyLimit'), value: 'freq_limit' },
        { text: t('common.author'), value: 'author.display_name' },
        { text: t('common.actionsLabel'), value: 'actions', sortable: false },
      ].filter(Boolean),
    );

    const durationToString = duration => convertDurationToString(duration);

    return {
      headers,
      durationToString,
    };
  },
};
</script>
