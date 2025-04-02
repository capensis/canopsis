<template>
  <c-advanced-data-table
    :headers="headers"
    :items="rules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    class="events-record-events-list"
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar="">
      <v-btn
        color="primary"
        fab
        @click.stop="showCreateExternalDataTableRuleItem"
      >
        <v-icon>add</v-icon>
      </v-btn>
    </template>
    <template #timestamp="{ item }">
      {{ item.timestamp | date }}
    </template>
    <template #actions="{ item }">
      <c-action-btn
        v-if="updatable"
        type="edit"
        @click="$emit('edit', item)"
      />
      <c-action-btn
        v-if="removable"
        :disabled="!item.deletable"
        type="delete"
        @click="remove(item)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    rules: {
      type: Array,
      default: () => [],
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
        sortable: false,
      },
      {
        text: t('common.description'),
        value: 'description',
        sortable: false,
      },
      {
        text: t('common.database'),
        value: 'database',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    /**
     * REMOVE
     */
    const remove = externalDataTable => emit('remove', externalDataTable);

    /**
     * OPTIONS
     */
    const updateOptions = options => emit('update:options', options);

    const showCreateExternalDataTableRuleItem = () => {};

    return {
      headers,

      remove,

      updateOptions,
      showCreateExternalDataTableRuleItem,
    };
  },
};
</script>
