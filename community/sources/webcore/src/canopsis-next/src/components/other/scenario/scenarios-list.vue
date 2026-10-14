<template>
  <c-advanced-data-table
    :headers="headers"
    :items="scenarios"
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
        scenario
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #delay="{ item }">
      <span>{{ formatDuration(item.delay) }}</span>
    </template>

    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>

    <template #enabled="{ item }">
      <c-help-icon
        v-if="hasDeprecatedTrigger(item)"
        :text="$t('scenario.errors.deprecatedTriggerExist')"
        color="error"
        icon="error"
        top
      />
      <c-enabled v-else :value="item.enabled" />
    </template>

    <template #created="{ item }">
      {{ formatDate(item.created) }}
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-db-export-btn :id="item._id" scenario />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <scenarios-list-expand-item :scenario="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';
import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

import { useI18n } from '@/hooks/i18n';

import ScenariosListExpandItem from './partials/scenarios-list-expand-item.vue';

export default {
  components: {
    ScenariosListExpandItem,
  },

  props: {
    scenarios: {
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

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.id'), value: '_id' },
      { text: t('common.name'), value: 'name' },
      { text: t('common.delay'), value: 'delay', sortable: false },
      { text: t('common.priority'), value: 'priority' },
      { text: t('common.enabled'), value: 'enabled' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.updated'), value: 'updated' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const hasDeprecatedTrigger = item => item.triggers.some(({ type }) => isDeprecatedTrigger(type));

    return {
      headers,
      formatDate: convertDateToString,
      formatDuration: convertDurationToString,
      hasDeprecatedTrigger,
    };
  },
};
</script>
