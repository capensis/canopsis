<template>
  <c-advanced-data-table
    :headers="headers"
    :items="externalDataTables"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #type="{ item }">
      {{ $t(`externalData.tableTypes.${item.type}`) }}
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
    <template #expand="{ item }">
      <external-data-tables-list-expand-panel :external-data-table="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import ExternalDataTablesListExpandPanel from './partials/external-data-tables-list-expand-panel.vue';

export default {
  components: {
    ExternalDataTablesListExpandPanel,
  },
  props: {
    externalDataTables: {
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
      },
      {
        text: t('common.description'),
        value: 'description',
      },
      {
        text: t('common.database'),
        value: 'type',
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

    const getHeaderSlotName = header => `header.${header.value}`;

    return {
      headers,

      remove,
      updateOptions,

      getHeaderSlotName,
    };
  },
};
</script>
