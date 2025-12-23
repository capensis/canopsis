<template>
  <c-advanced-data-table
    :headers="headers"
    :items="scenarios"
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
        scenario
        @clear:items="clearSelected"
      />
    </template>
    <template #headerCell="{ header }">
      <span class="pre-line header-text">{{ header.text }}</span>
    </template>
    <template #delay="{ item }">
      <span>{{ item.delay | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #enabled="{ item }">
      <c-help-icon
        v-if="hasDeprecatedTrigger(item)"
        :text="$t('scenario.errors.deprecatedTriggerExist')"
        color="error"
        icon="error"
        top
      />
      <c-enabled
        v-else
        :value="item.enabled"
      />
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
        <c-db-export-btn :id="item._id" scenario />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <scenarios-list-expand-item :scenario="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

import { useI18n } from '@/hooks/i18n';

import ScenariosListExpandItem from './partials/scenarios-list-expand-item.vue';

export default {
  components: {
    ScenariosListExpandItem,
  },
  props: {
    scenarios: {
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
      { text: t('common.id'), value: '_id' },
      { text: t('common.name'), value: 'name' },
      { text: t('common.delay'), value: 'delay', sortable: false },
      { text: t('common.priority'), value: 'priority' },
      { text: t('common.enabled'), value: 'enabled' },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.updated'), value: 'updated' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const hasDeprecatedTrigger = item => item.triggers.some(({ type }) => isDeprecatedTrigger(type));

    return {
      headers,
      hasDeprecatedTrigger,
    };
  },
};
</script>
