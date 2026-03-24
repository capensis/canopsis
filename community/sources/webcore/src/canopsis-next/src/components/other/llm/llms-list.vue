<template>
  <c-advanced-data-table
    :headers="headers"
    :items="llms"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    expand
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-show="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #name="{ item }">
      {{ item.name }}
    </template>
    <template #type="{ item }">
      {{ item.type ?? item.model_type }}
    </template>
    <template #model="{ item }">
      {{ item.model }}
    </template>
    <template #thinking_level="{ item }">
      <span v-if="item.thinking_level && $te(`llm.thinkingLevels.${item.thinking_level}`)">
        {{ $t(`llm.thinkingLevels.${item.thinking_level}`) }}
      </span>
      <span v-else>{{ item.thinking_level }}</span>
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #last_used="{ item }">
      {{ item.last_used | date('long', '-') }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <llms-list-expand-panel :llm="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import LlmsListExpandPanel from './partials/llms-list-expand-panel.vue';

export default {
  components: {
    LlmsListExpandPanel,
  },
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
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('llm.modelType'), value: 'type' },
      { text: t('llm.model'), value: 'model' },
      { text: t('llm.thinkingLevel'), value: 'thinking_level', sortable: false },
      { text: t('common.enabled'), value: 'enabled', sortable: false },
      { text: t('llm.lastUsedDate'), value: 'last_used' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    return {
      headers,
    };
  },
};
</script>
