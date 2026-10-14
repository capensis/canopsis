<template>
  <c-advanced-data-table
    :headers="headers"
    :items="widgetTemplates"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #type="{ item }">
      {{ t(`widgetTemplate.types.${item.type}`) }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
        <c-action-btn v-if="updatable" type="edit" @click.stop="$emit('edit', item)" />
        <c-action-btn v-if="removable" type="delete" @click.stop="$emit('remove', item._id)" />
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
    widgetTemplates: {
      type: Array,
      required: true,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
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
  },

  emits: ['update:options', 'edit', 'remove'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'title',
        sortable: false,
      },
      {
        text: t('common.type'),
        value: 'type',
        sortable: false,
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
        sortable: false,
      },
      {
        text: t('common.lastModifiedBy'),
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
      t,
      headers,
      formatDate: convertDateToString,
    };
  },
};
</script>
