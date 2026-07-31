<template>
  <c-advanced-data-table
    :headers="headers"
    :items="commentTemplates"
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

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn
          v-if="editable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
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
    commentTemplates: {
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
    editable: {
      type: Boolean,
      default: false,
    },
    deletable: {
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
        value: 'name',
        sortable: true,
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
        sortable: true,
      },
      {
        text: t('common.lastModifiedBy'),
        value: 'author.display_name',
        sortable: true,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      formatDate: convertDateToString,
      headers,
    };
  },
};
</script>
