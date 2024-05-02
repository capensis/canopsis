<template>
  <c-advanced-data-table
    :headers="headers"
    :items="linkRules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
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
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item._id)"
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
    linkRules: {
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
    duplicable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.enabled'), value: 'enabled', sortable: false },
      { text: t('common.lastModifiedOn'), value: 'updated' },
      { text: t('common.lastModifiedBy'), value: 'author.display_name' },
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
