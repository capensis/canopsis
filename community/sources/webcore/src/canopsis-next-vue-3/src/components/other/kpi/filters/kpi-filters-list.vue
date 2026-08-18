<template>
  <c-advanced-data-table
    :options="options"
    :items="filters"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    search
    advanced-pagination
    hide-actions
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #created="{ item }">
      {{ convertDateToString(item.created) }}
    </template>
    <template #updated="{ item }">
      {{ convertDateToString(item.updated) }}
    </template>
    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn v-if="updatable" type="edit" @click.stop="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click.stop="$emit('duplicate', item)" />
        <c-action-btn v-if="removable" type="delete" @click.stop="$emit('remove', item._id)" />
      </c-layout>
    </template>
    <template #expand="{ item }">
      <kpi-filters-expand-item :filter="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import KpiFiltersExpandItem from './partials/kpi-filters-expand-item.vue';

export default {
  components: { KpiFiltersExpandItem },

  props: {
    filters: {
      type: Array,
      required: true,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    pending: {
      type: Boolean,
      default: false,
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

  emits: ['update:options', 'edit', 'duplicate', 'remove'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.lastModifiedOn'), value: 'updated' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    return {
      headers,
      convertDateToString,
    };
  },
};
</script>
