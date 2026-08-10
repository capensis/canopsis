<template>
  <c-advanced-data-table
    :headers="headers"
    :items="llms"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        llm
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #name="{ item }">
      {{ item.name }}
    </template>

    <template #type="{ item }">
      {{ item.type ?? item.model_type }}
    </template>

    <template #thinking_level="{ item }">
      <span v-if="item.thinking_level && te(`llm.thinkingLevels.${item.thinking_level}`)">
        {{ t(`llm.thinkingLevels.${item.thinking_level}`) }}
      </span>
      <span v-else>
        {{ item.thinking_level }}
      </span>
    </template>

    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>

    <template #default="{ item }">
      <c-enabled :value="item.default" />
    </template>

    <template #last_used="{ item }">
      {{ formatDate(item.last_used) }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
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

export default {
  props: {
    llms: {
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
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['edit', 'remove', 'refresh', 'update:options'],

  setup() {
    const { t, te } = useI18n();

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('llm.modelType'), value: 'type' },
      { text: t('llm.model'), value: 'model' },
      { text: t('llm.thinkingLevel'), value: 'thinking_level', sortable: false },
      { text: t('common.enabled'), value: 'enabled', sortable: false },
      { text: t('llm.isDefaultModel'), value: 'default', sortable: false },
      { text: t('llm.lastUsedDate'), value: 'last_used' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const formatDate = value => convertDateToString(value, 'long', '-');

    return {
      t,
      te,
      headers,

      formatDate,
    };
  },
};
</script>
