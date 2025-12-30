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
      <span>{{ item.duration | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
        <c-db-export-btn :id="item._id" flapping-rule />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <alarm-status-rules-list-expand-item :rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

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
  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.id'), value: '_id' },
      { text: t('common.name'), value: 'name' },
      { text: t('common.duration'), value: 'duration', sortable: false },
      { text: t('common.priority'), value: 'priority' },
      props.flapping && { text: t('common.frequencyLimit'), value: 'freq_limit' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ].filter(Boolean));

    return {
      headers,
    };
  },
};
</script>
