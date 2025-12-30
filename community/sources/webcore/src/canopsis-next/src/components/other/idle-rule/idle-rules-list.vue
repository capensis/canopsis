<template>
  <c-advanced-data-table
    :headers="headers"
    :items="idleRules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable || updatable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        idle-rule
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #type="{ item }">
      {{ $t(`idleRules.types.${item.type}`) }}
    </template>
    <template #operation.type="{ item }">
      {{ item | get('operation.type', '-') }}
    </template>
    <template #duration="{ item }">
      <span>{{ item.duration | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
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
        <c-db-export-btn :id="item._id" idle-rule />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <idle-rules-list-expand-item :idle-rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import IdleRulesListExpandItem from './partials/idle-rules-list-expand-item.vue';

export default {
  components: {
    IdleRulesListExpandItem,
  },
  props: {
    idleRules: {
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
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.type'), value: 'type' },
      { text: t('common.enabled'), value: 'enabled', sortable: false },
      { text: t('common.action'), value: 'operation.type', sortable: false },
      { text: t('idleRules.timeAwaiting'), value: 'duration', sortable: false },
      { text: t('common.priority'), value: 'priority' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    return {
      headers,
    };
  },
};
</script>
