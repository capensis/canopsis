<template>
  <c-advanced-data-table
    :items="icons"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #icon="{ item }">
      <custom-icon :content="item.content" />
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click.stop="$emit('edit', item)" />
        <c-action-btn v-if="removable" type="delete" @click.stop="$emit('remove', item)" />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import CustomIcon from './partials/custom-icon.vue';

export default {
  components: { CustomIcon },

  props: {
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    icons: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
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

  emits: ['edit', 'remove', 'update:options'],

  setup() {
    const { t } = useI18n();
    const formatDate = (value) => convertDateToString(value);

    const headers = computed(() => [
      {
        text: t('common.icon', 1),
        value: 'icon',
        sortable: false,
      },
      {
        text: t('common.title'),
        value: 'title',
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
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
