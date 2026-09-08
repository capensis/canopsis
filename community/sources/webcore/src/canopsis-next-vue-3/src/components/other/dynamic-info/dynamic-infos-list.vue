<template>
  <c-advanced-data-table
    :options="options"
    :items="dynamicInfos"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    :select-all="removable || updatable"
    :advanced-search-attributes="advancedSearchAttributes"
    advanced-search
    advanced-pagination
    hide-actions
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        dynamic-info
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #created="{ item }">
      {{ formatDate(item.created) }}
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-db-export-btn :id="item._id" dynamic-info />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <dynamic-infos-list-expand-item :info="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import { useDynamicInfoAdvancedSearchAttributes } from '@/components/common/search/hooks/advanced-search';

import DynamicInfosListExpandItem from './partials/dynamic-infos-expand-item.vue';

export default {
  components: {
    DynamicInfosListExpandItem,
  },

  props: {
    dynamicInfos: {
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
    const { attributes: advancedSearchAttributes } = useDynamicInfoAdvancedSearchAttributes();

    const headers = computed(() => [
      { text: t('common.id'), value: '_id' },
      { text: t('common.name'), value: 'name' },
      { text: t('common.description'), value: 'description', sortable: false },
      { text: t('common.enabled'), value: 'enabled', sortable: false },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.updated'), value: 'updated' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const formatDate = convertDateToString;

    return {
      advancedSearchAttributes,
      headers,
      formatDate,
    };
  },
};
</script>
