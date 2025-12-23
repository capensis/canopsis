<template>
  <c-advanced-data-table
    :headers="headers"
    :items="declareTicketRules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable || updatable"
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
        declare-ticket
        @clear:items="clearSelected"
      />
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
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
        <c-db-export-btn :id="item._id" declare-ticket />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    declareTicketRules: {
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
      { text: t('common.systemName'), value: 'system_name' },
      { text: t('common.enabled'), value: 'enabled' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.lastModifiedOn'), value: 'updated' },
      { text: t('common.lastModifiedBy'), value: 'author.display_name' },
    ]);

    return {
      headers,
    };
  },
};
</script>
